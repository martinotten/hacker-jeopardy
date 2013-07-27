package statemachine

import (
	"jeopardy/json"
	"jeopardy/comms"
)

import (
	"fmt"
	"math/rand"
	"strings"
	"strconv"
)
type StateId int

//const (
//	IDLE StateId = iota    // beginning of the game (expected input: admin starts game)
//	NEW_GAME             // admin started game (expected input: player1 name)
//	PLAYER              // player1 (expected input: player2 name)
////	PLAYER2              // player2 (expected input: player3 name)
//	START_GAME					 // all names know, broadcast board
//	PICK_PLAYER          // server picks player (expected input: question picked by player)
//	QUESTION_PICKED      // question is picked, display, start timer (expected input: buzzer 1,2 or 3)
//	ANSWER_QUESTION      // first buzzer pressed ansers, start timer (expected input: correct, incorrect, timer_expired)
//	ADJUST_SCORE
//	CORRECT              // adjust player score -> CHK_GAME_OVER
//	WRONG                // addjust player score -> QUESTION_PICKED / CHK_GAME_OVER
//	LAST_PLAYER          // have all players (unsuccessfully) tried to anser 
//	CHK_GAME_OVER        // check if any questions are left on the board -> pick player / GAME_OVER
//	GAME_OVER
//	DETERMINE_WINNER
//)


const (
	// events can come from buzzer (high byte 00) or admin (high byte FF)
	E_START_GAME = 0xFF00
	E_PLAYER_ONE = 0xFF01
	E_PLAYER_TWO = 0xFF02
	E_PLAYER_THREE = 0xFF03
	E_QUESTION_CHOSEN = 0xFF04
	E_CORRECT = 0xFF05
	E_INCORRECT = 0xFF06

	E_BUZZER_ONE = 0x0000
	E_BUZZER_TWO = 0x0001
	E_BUZZER_THREE = 0x0002

	E_TIMEOUT_NO_ANSWER = 0xF000

)
type Event struct {
	Id int
	Data string
}

type State interface {
	Game() *Game
	EnterState(Event) 
	HandleEvent(Event) State

}
//type Player struct {
//	Name string
//	Score int
//}

//type Question struct {
//	Answer string
//	Question string
//	Value int
//}

type Game struct {
	GameState State
	Players []*json.Player
	CurrentPlayer int // 0 noone ...
	LastCorrectAnswer int // 0 noone, 1 player 1, 2 = player 2...
	CurrentQuestion * json.Answer
	CurrentAttempts string  // keeps track of who has tried to answer the current Question
	QuestionsRemaining int
	Categories []*json.Category

	Admin	*comms.Admin
	UI    *comms.WebsocketHandler
	// Buzzer 1
	// Buzzer 2
	// Buzzer 3
	// UI
}

func NewGame(fn string, admin *comms.Admin)*Game {
	var err error
	game := new(Game)
	game.Admin = admin

	if game.Categories, err = json.LoadCategories(fn); err != nil {
		fmt.Printf("%s\n", err)
		panic(err.Error())
	}
	game.Players = make([]*json.Player, 3)
	game.Players[0] = &json.Player{"1", 0, "default"}
	game.Players[1] = &json.Player{"2", 0, "default"}
	game.Players[2] = &json.Player{"3", 0, "default"}

	game.GameState = new(S_Idle)

	return game
}

// sends the current gamestate to the web client
func (g * Game) SendGameState() {
	state := json.GameState{}
	state.Categories = g.Categories
	state.Players    = g.Players
	if (g.CurrentQuestion != nil) {
		state.Answer = g.CurrentQuestion.Answer
	}
	if (g.UI != nil) {
		g.UI.SendGameState(&state)
	} else {
		g.Admin.Prompt("no websocket!")
	}
}


// Input from Admin
func (g * Game) HandleEvent (e Event) {
	g.GameState = g.GameState.HandleEvent(e)
	g.GameState.EnterState(e)
}

type baseState struct {
	game *Game
}

func (s * baseState) Game()*Game {
	return s.game
}

type S_Idle struct {
	baseState
}

func (s * S_Idle) EnterState(e Event) {return}
func (s * S_Idle) HandleEvent(e Event)State {
	if (e.Id == E_START_GAME) {
		var snew S_NewGame
		snew.game = s.game
		return &snew
	}
	return s;
}

type S_NewGame struct {
	baseState
}
func (s * S_NewGame) EnterState(e Event) {
	s.game.SendGameState()
	// send (something) to buzzer
}
func (s * S_NewGame) HandleEvent(e Event)State {
	if (e.Id == E_PLAYER_ONE) {
		// s.sendBoard(e.Data)
		snew := S_Player{}
		snew.game = s.game
		return &snew
	} else {
		return s
	}
}

type S_Player struct {
	baseState
}
func (s * S_Player) EnterState(e Event) {
	switch (e.Id) {
		case E_PLAYER_ONE:
			s.game.Players[0] = &json.Player{e.Data, 0, "default"}
		case E_PLAYER_TWO:
			s.game.Players[1] = &json.Player{e.Data, 0, "default"}
		case E_PLAYER_THREE:
			s.game.Players[2] = &json.Player{e.Data, 0, "default"}
		default:
			s.game.Admin.Prompt("unexpected event")
			return
	}
	s.game.SendGameState()
}
func (s * S_Player) HandleEvent(e Event)(state State) {
	switch (e.Id) {
		case E_PLAYER_TWO:
			snew := S_Player{}
			snew.game = s.game
			return &snew
		case E_PLAYER_THREE:
			// new game will start. broadcast board.
			snew := S_StartGame{}
			snew.game = s.game
			return  &snew
		default:
			return s
	}
}

type S_StartGame struct {
	baseState
}

func (s * S_StartGame) EnterState(e Event) {
	// set up board. broadcast
	s.game.SendGameState()
	s.game.HandleEvent(e) // advance to next state automatically.
}
func (s * S_StartGame) HandleEvent(e Event)State{
	new_state := new(S_PickPlayer)
	new_state.game = s.game
	return new_state
}

type S_PickPlayer struct {
	baseState
}
func (s * S_PickPlayer) EnterState(e Event) {
	// reset some state
	s.game.CurrentAttempts = ""
	// pick player and broadcast
	if s.game.LastCorrectAnswer != 0 {
		s.game.CurrentPlayer = s.game.LastCorrectAnswer
	} else {
		s.game.CurrentPlayer = (rand.Int() % 3) + 1
	}
	s.game.Players[s.game.CurrentPlayer -1].Status = "active"
	s.game.SendGameState()
}
func (s * S_PickPlayer) HandleEvent(e Event)State {
	if (e.Id == E_QUESTION_CHOSEN) {
		nstate := new(S_QuestionChosen)
		nstate.game = s.game
		return nstate
	}
	return s
}

type S_QuestionChosen struct {
	baseState
}

func(s * S_QuestionChosen) EnterState(e Event) {
		// tell ui question
		// display question to admin
		cat_ques := strings.Split(e.Data, "_")
		var cat int64
		var ques int64
		var err error
		if cat, err = strconv.ParseInt(cat_ques[0], 10, 32); err != nil {
			mes := fmt.Sprintf("error parsing %s : %s", cat_ques[0], err.Error())
			s.game.Admin.Prompt(mes)
		}
		if ques, err = strconv.ParseInt(cat_ques[1], 10, 32); err != nil {
			mes := fmt.Sprintf("error parsing %s : %s", cat_ques[1], err.Error())
			s.game.Admin.Prompt(mes)
		}

		category := s.game.Categories[cat-1]
		answer   := category.Answers[ques-1]
		s.game.CurrentQuestion = &answer

		s.game.Admin.Prompt(answer.Answer)
		s.game.Admin.Prompt(answer.Question)

		s.game.SendGameState()

		// TODO start_timer
}
func(s * S_QuestionChosen) HandleEvent(e Event) State {
	nstate := new(S_AnswerExpected)
	nstate.game = s.game
	switch (e.Id) {
		case E_BUZZER_ONE:
			return nstate
		case E_BUZZER_TWO:
			return nstate
		case E_BUZZER_THREE:
			return nstate
		case E_TIMEOUT_NO_ANSWER:
			nstate2 := new(S_CheckGameOver)	
			nstate2.game = s.game
			return nstate2
		default:
			return s
	}
}

type S_AnswerExpected struct {
	baseState
}

func(s * S_AnswerExpected) EnterState(e Event) {
		
		switch (e.Id) {
		case E_BUZZER_ONE:
			s.game.CurrentPlayer = 1
		case E_BUZZER_TWO:
			s.game.CurrentPlayer = 2
		case E_BUZZER_THREE:
			s.game.CurrentPlayer = 3
		}
		for i, player := range s.game.Players {
			if i+1 == s.game.CurrentPlayer {
				player.Status = "active"
			} else {
				player.Status = "default"
			}
		}
		s.game.SendGameState()
		s.game.Admin.Prompt("Was the given Answer correct?")

		// TODO start_timer
}

func(s * S_AnswerExpected) HandleEvent(e Event) State {
	var nstate State
	switch (e.Id) {
		case E_CORRECT:
			nstate := new(S_Adjust_Score)
			nstate.game = s.Game()
		case E_INCORRECT:
			nstate := new(S_Adjust_Score)
			nstate.game = s.Game()
		case E_TIMEOUT_NO_ANSWER:
			nstate := new(S_CheckGameOver)
			nstate.game = s.Game()
		default:
			return s
	}
	return nstate
}

type S_Adjust_Score struct {
	baseState
}

func(s * S_Adjust_Score) EnterState(e Event) {
		switch (e.Id) {
		case E_CORRECT:
			s.game.Players[s.game.CurrentPlayer-1].Score = s.game.CurrentQuestion.Value
			s.game.LastCorrectAnswer = s.game.CurrentPlayer
		case E_INCORRECT:
			s.game.Players[s.game.CurrentPlayer-1].Score = s.game.CurrentQuestion.Value
		}
		s.game.SendGameState()

		s.game.HandleEvent(e)

}
func(s * S_Adjust_Score) HandleEvent(e Event) State {
	switch (e.Id) {
		case E_CORRECT:
			nstate := new(S_CheckGameOver)
			nstate.game = s.game
			return nstate
		case E_INCORRECT:
			nstate := new(S_CheckLastPlayer)
			nstate.game = s.game
			if (nstate.game.CurrentPlayer == 1) {
				nstate.game.CurrentAttempts += "1"
			}else if (nstate.game.CurrentPlayer == 2) {
				nstate.game.CurrentAttempts += "2"
			} else {
				nstate.game.CurrentAttempts += "3"	
			}
			return nstate
		default:
			return s
	}
}

type S_CheckGameOver struct {
	baseState
}

func(s * S_CheckGameOver) EnterState(e Event) {
	s.HandleEvent(e)
}
func(s * S_CheckGameOver) HandleEvent(e Event) State {
	if (s.game.QuestionsRemaining == 0) { // TODO!!
		nstate := new(S_DetermineWinner)
		nstate.game = s.game
		return nstate
	} else {
		nstate := new(S_PickPlayer)
		nstate.game = s.game
		return nstate
	}
}

type S_CheckLastPlayer struct {
	baseState
}

func(s * S_CheckLastPlayer) EnterState(e Event) {
	s.game.HandleEvent(e)
}
func(s * S_CheckLastPlayer) HandleEvent(e Event) State {
	if (len(s.game.CurrentAttempts) == 3) {
		nstate := new(S_CheckGameOver)
		nstate.game = s.game
		return nstate
	}
	nstate := new(S_QuestionChosen)
	nstate.game = s.game
	return nstate
}

type S_DetermineWinner struct {
	baseState
}

func(s * S_DetermineWinner) EnterState(e Event) {
	// broadcast winner.
}
func(s * S_DetermineWinner) HandleEvent(e Event) State {
	return s
}





