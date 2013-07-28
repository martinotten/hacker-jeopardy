package statemachine

// this contains the communication with the main screen ui
import "jeopardy/json"
import "code.google.com/p/go.net/websocket"


import (
	"fmt"
	"os"
	J "encoding/json"
)

type WebsocketHandler struct {
	conn    *websocket.Conn
	encoder *J.Encoder
}


func (ws * WebsocketHandler)SetSocket (con *websocket.Conn) {
	ws.conn = con
	ws.encoder = J.NewEncoder(con)
}

func (ws * WebsocketHandler) SendGameState(state *json.GameState) {
	println(ws.conn)	
	if err := ws.encoder.Encode(state); err != nil {
		fmt.Printf("%s\n", err)
		panic(err.Error())
	}

	encoder := J.NewEncoder(os.Stdin)
	encoder.Encode(state)
}



// 
func (ws * WebsocketHandler) UISendGame(state []*json.Category) {
	if err := ws.encoder.Encode(state); err != nil {
		fmt.Printf("%s", err)
		panic(err.Error())
	}
}

