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

func (ws * WebsocketHandler) UISendNewGame(cats []*json.Category) error {
	if ws.encoder == nil {
		return fmt.Errorf("no websocket")
	}
	if err := ws.encoder.Encode(cats); err != nil {
		return err
	}
	return nil
}

func (ws * WebsocketHandler) UISendPlayer(p json.Player) error{
	if ws.encoder == nil {
		return fmt.Errorf("no websocket")
	}
	if err := ws.encoder.Encode(p); err != nil {
		return err
	}
	return nil
}

func (ws * WebsocketHandler) UISendCurrentPick(p json.Player) error {
	if ws.encoder == nil {
		return fmt.Errorf("no websocket")
	}
	if err := ws.encoder.Encode(p.Number); err != nil {
		return err
	}
	return nil
}
func (ws * WebsocketHandler) UISendCurrentPlayAnswer(a json.Answer) error {
	if ws.encoder == nil {
		return fmt.Errorf("no websocket")
	}
	if err := ws.encoder.Encode(a); err!= nil {
		return err;
	}
	return nil
}

func (ws * WebsocketHandler) UISendCorrect() error {
	if ws.encoder == nil {
		return fmt.Errorf("no websocket")
	}
	if err := ws.encoder.Encode("correct"); err != nil {
		return err
	}
	return nil
}
func (ws * WebsocketHandler) UISendIncorrect() error {
	if ws.encoder == nil {
		return fmt.Errorf("no websocket")
	}
	if err := ws.encoder.Encode("wrong"); err != nil {
		return err
	}
	return nil
}
