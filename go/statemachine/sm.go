package statemachine

type State int

const (
	IDLE State = iota    // beginning of the game (expected input: admin starts game)
	NEW_GAME             // admin started game (expected input: player1 name)
	PLAYER1              // player1 (expected input: player2 name)
	PLAYER2              // player2 (expected input: player3 name)
	START_GAME					 // all names know, broadcast board
	PICK_PLAYER          // server picks player (expected input: question picked by player)
	QUESTION_PICKED      // question is picked, display, start timer (expected input: buzzer 1,2 or 3)
	ANSWER_QUESTION      // first buzzer pressed ansers, start timer (expected input: correct, incorrect, timer_expired)
	CORRECT              // adjust player score -> CHK_GAME_OVER
	WRONG                // addjust player score -> QUESTION_PICKED / CHK_GAME_OVER
	LAST_PLAYER          // have all players (unsuccessfully) tried to anser 
	CHK_GAME_OVER        // check if any questions are left on the board -> pick player / GAME_OVER
	GAME_OVER
)



const (
	// events can come from buzzer (high byte 00) or admin (high byte FF)
	START_GAME = 0xFF00
)
type Event struct {
	Id int
	Data string
}

type Player struct {
	Name String
	Score int
}

type Game struct {
	GameState State
	Player1 Player
	Player2 Player
	Player3 Player
	// Buzzer 1
	// Buzzer 2
	// Buzzer 3
	// UI
}

func (g * Game) SwitchState (e Event) {
	switch (g.GameState) {
		case IDLE:
			if e.Id ==  START_GAME {
				g.GameState = NEW_GAME
			}

		case NEW_GAME:
			// alert buzzer
			// alert ui

		case PLAYER1:
		case PLAYER2:
		case START_GAME:
		case PICK_PLAYER:
		case QUESTION_PICKED:
		case ANSWER_QUESTION:
		case CORRECT:
		case WRONG:
		case LAST_PLAYER:
		case CHK_GAME_OVER:
		case GAME_OVER:
	}
}
