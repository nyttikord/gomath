package gomath

import "testing"

func TestAstSum(t *testing.T) {
	lexr, err := lex("1+2") // useless to test another sum like 1+2
	if err != nil {
		t.Fatal(err)
	}
	tree, err := astParse(lexr)
	if err != nil {
		t.Fatal(err)
	}

}
