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

type Player struct {
	Name string `json:"name"`
	Number int `json:"number"`
	Score int `json:"score"`
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



