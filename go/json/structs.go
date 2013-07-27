package json

import (
	"encoding/json"
	"os"
	"io"
)

type Category struct {
	Name  string `json:"category"`
	Level int    `json:level`
	A100  string `json:"A100"`
	Q100  string `json:"Q100"`
	A200  string `json:"A200"`
	Q200  string `json:"Q200"`
	A300  string `json:"A300"`
	Q300  string `json:"Q300"`
	A400  string `json:"A400"`
	Q400  string `json:"Q400"`
	A500  string `json:"A500"`
	Q500  string `json:"Q500"`
}

func LoadCategory(fn string) (category *Category, err error) {
	if file, e := os.Open(fn); e != nil {
		return nil, e
	} else {
		category = new(Category)
		decoder := json.NewDecoder(file)
		if err = decoder.Decode(category); err != nil {
			return
		}
	}
	return
}

func CategoryForPlayer(category *Category, writer io.Writer)(err error) {
	catPlayer := *category
	catPlayer.Q100 = ""
	catPlayer.Q200 = ""
	catPlayer.Q300 = ""
	catPlayer.Q400 = ""
	catPlayer.Q500 = ""

	encoder := json.NewEncoder(writer)
	if err = encoder.Encode(catPlayer); err != nil {
		return
	}
	return
}


