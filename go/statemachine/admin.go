package statemachine



import (
	"fmt"
)


type Admin struct {
}


func (a * Admin) Prompt (msg string) {
	fmt.Printf("%s\n", msg)
}

func (a * Admin) StartGame(game * Game) {
	go func() {
		a.Prompt("Hit Enter to start Game")
		var s string
		fmt.Scanln(&s)
		event := Event{E_START_GAME, ""}
		game.HandleEvent(event)
	}()
}

func (a * Admin) GetPlayer1 (game * Game) {
	go func() {
	a.Prompt("Enter Name Player1: ")
	var name string
	fmt.Scanln(&name)
	event := Event{E_PLAYER_ONE, name}
	game.HandleEvent(event)
	}()
}


func (a * Admin) GetPlayer2 (game * Game) {
	go func() {
	a.Prompt("Enter Name Player2: ")
	var name string
	fmt.Scanln(&name)
	event := Event{E_PLAYER_TWO, name}
	game.HandleEvent(event)
	}()
}
func (a * Admin) GetPlayer3 (game * Game) {
	go func() {
	a.Prompt("Enter Name Player3: ")
	var name string
	fmt.Scanln(&name)
	event := Event{E_PLAYER_THREE, name}
	game.HandleEvent(event)
	}()
}

func (a * Admin) ChooseCategory(game * Game) {
	go func() {
		for i, cat := range game.Categories {
			if (!cat.Done()) {
				fmt.Printf("%d %s\n", i, cat.Name)
			} else {
			}
		}
		print("Enter Category: ")

		var i int
		fmt.Scanln(&i)
		println("")
		for j, ans := range game.Categories[i].Answers {
			if (!ans.Done) {
				fmt.Printf("%d %s\n", j, ans.Answer)
			}
		}
		print("Enter Question: ")
		var j int
		fmt.Scanln(&j)
		println("")

		event := Event{E_QUESTION_CHOSEN, fmt.Sprintf("%d_%d", i, j)}
		game.HandleEvent(event)
	}()
}

func (a * Admin) GetBuzzer(game * Game) {
	go func() {
		var buzzer string
		fmt.Scanln(&buzzer)
		var event Event
		switch buzzer[0] {
			case 114:
				event = Event{E_BUZZER_ONE, ""}
			case 103:
				event = Event{E_BUZZER_TWO, ""}
			case 98:
				event = Event{E_BUZZER_THREE, ""}
		}
		game.HandleEvent(event)
	}()
}

func (a * Admin) AnswerCorrect(game * Game) {
	go func() {
		a.Prompt("Answer Correct? y/n")
		var correct string
		fmt.Scanln(&correct)
		var first = correct[0:1]
		switch(first) {
			case "y":
				game.HandleEvent(Event{E_CORRECT, ""})
			case "n":
				game.HandleEvent(Event{E_INCORRECT, ""})
			default:
				a.AnswerCorrect(game)
		}
	}()
}
