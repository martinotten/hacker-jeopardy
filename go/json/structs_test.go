package json

import (
	"testing"
	"bytes"
	"strings"
)

const testJsonFn = "../test/colors.json"
const jsonPlayer = `{"category":"Colors","Level":1,"A100":"Smurfs","Q100":"","A200":"Simpsons","Q200":"","A300":"Kermit","Q300":"","A400":"Prince","Q400":"","A500":"EL James","Q500":""}`

func TestJson(t *testing.T) {
	var category *Category
	var err error
	if category, err = LoadCategory(testJsonFn); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if category == nil {
		t.Fatalf("category is nil")
	}

	if category.Name != "Colors" {
		t.Errorf("unexpected answer %s", category.Name)
	}
	if category.Level != 1 {
		t.Errorf("unexpected answer %s", category.Level)
	}
	if category.A100 != "Smurfs" {
		t.Errorf("unexpected answer %s", category.A100)
	}
	if category.Q100 != "Blue" {
		t.Errorf("unexpected answer %s", category.Q100)
	}
	if category.A200 != "Simpsons" {
		t.Errorf("unexpected answer %s", category.A200)
	}
	if category.Q200 != "Yellow" {
		t.Errorf("unexpected answer %s", category.Q200)
	}
	if category.A300 != "Kermit" {
		t.Errorf("unexpected answer %s", category.A300)
	}
	if category.Q300 != "Green" {
		t.Errorf("unexpected answer %s", category.Q300)
	}
	if category.A400 != "Prince" {
		t.Errorf("unexpected answer %s", category.A400)
	}
	if category.Q400 != "Purple" {
		t.Errorf("unexpected answer %s", category.Q400)
	}
	if category.A500 != "EL James" {
		t.Errorf("unexpected answer %s", category.A500)
	}
	if category.Q500 != "Grey" {
		t.Errorf("unexpected answer %s", category.Q500)
	}
}

func TestCatPlayer (t *testing.T) {
	var category *Category
	var err error
	if category, err = LoadCategory(testJsonFn); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if category == nil {
		t.Fatalf("category is nil")
	}

	var buf  bytes.Buffer
	if err:= CategoryForPlayer(category, &buf); err != nil {
		t.Errorf("unexpected err: %s", err)
	}

	if (strings.TrimSpace(buf.String()) != jsonPlayer) {
		t.Errorf (">%s<", jsonPlayer)
		t.Errorf (">%s<", buf.String())
	}
}
