package json

import (
	"testing"
)

const testJsonFn = "../questions/questions.json"
const jsonPlayer = `{"category":"Colors","Level":1,"A100":"Smurfs","Q100":"","A200":"Simpsons","Q200":"","A300":"Kermit","Q300":"","A400":"Prince","Q400":"","A500":"EL James","Q500":""}`

func TestJson(t *testing.T) {
	if categories, err := LoadCategories(testJsonFn); err != nil {
		t.Fatalf("unexpected error: %s", err)
	} else {
		for _, cat := range(categories) {
			t.Error(cat)	
		}
	}

}

