package comms


type Admin struct {
	Game *statemachine.Game
}


func (a * admin) Prompt (msg string) {
	fmt.Printf("%s", msg)
}


