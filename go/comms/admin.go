package comms



import (
	"fmt"
) 

type Admin struct {
}


func (a * Admin) Prompt (msg string) {
	fmt.Printf("%s", msg)
}


