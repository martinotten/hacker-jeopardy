package comms

// this contains the communication with the main screen ui
import "jeopardy/json"
import "code.google.com/p/go.net/websocket"


import (
	"fmt"
	J "encoding/json"
)

type WebsocketHandler struct {
	encoder *J.Encoder
}


func (ws * WebsocketHandler)SetSocket (con *websocket.Conn) {
	ws.encoder = J.NewEncoder(con)
}

func (ws * WebsocketHandler) SendGameState(state *json.GameState) {
	if err := ws.encoder.Encode(state); err != nil {
		fmt.Printf("%s", err)
		panic(err.Error())
	}
}



// 
func (ws * WebsocketHandler) UISendGame(state []*json.Category) {
	if err := ws.encoder.Encode(state); err != nil {
		fmt.Printf("%s", err)
		panic(err.Error())
	}
}

