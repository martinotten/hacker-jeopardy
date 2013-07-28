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
				fmt.Println("%d %s", i, cat.Name)
			}
		}

		var i int
		fmt.Scanln(&i)
		for j, ans := range game.Categories[i].Answers {
			if (!ans.Done) {
				fmt.Println("%d %s", j, ans.Answer)
			}
		}
		var j int
		fmt.Scanln(&j)

		event := Event{E_QUESTION_CHOSEN, fmt.Sprintf("%d_%d")}
		game.HandleEvent(event)
	}()
}
