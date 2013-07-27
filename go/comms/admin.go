package comms


import (
	"jeopardy/statemachine"
)

import (
	"fmt"
) 

type Admin struct {
	Game *statemachine.Game
}


func (a * Admin) Prompt (msg string) {
	fmt.Printf("%s", msg)
}


