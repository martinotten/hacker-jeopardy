package json

import (
	"encoding/json"
	"os"
)

type Category struct {
	Name    string `json:"name"`
	Answers []Answer `json:"answers"`
}

type Answer struct {
	Answer string `json:"answer"`
	Question string `json:"question"`
	Value int `json:"value"`
	Done bool`json:"done"`
}

func (c * Category) Done() bool {
	for _, ans := range c.Answers {
		if !ans.Done {
			return true
		}
	}
	return false
}


type Player struct {
	Name string `json:"name"`
	Score int `json:"score"`
	Status string `json:"status"`
}

type GameState struct {
	Categories []*Category `json:"categories"`
	Players    []*Player   `json:"players"`
	Answer     string     `json:"answer"`
}

// messages
// server dir client
// newGame >  
//         < Player1
//         < Player2
//         < Player3
// PlayerX >
// 

func LoadCategories(fn string) (categories []*Category, err error) {
	if file, e := os.Open(fn); e != nil {
		return nil, e
	} else {
		categories = make([]*Category, 1000)
		decoder := json.NewDecoder(file)
		if err = decoder.Decode(&categories); err != nil {
			return
		}
	}
	return
}



