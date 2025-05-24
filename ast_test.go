package gomath

import (
	"encoding/json"
	"testing"
)

func TestAstSum(t *testing.T) {
	lexr, err := lex("1+2") // useless to test another sum
	if err != nil {
		t.Fatal(err)
	}
	tree, err := astParse(lexr, "return")
	if err != nil {
		t.Fatal(err)
	}
	if tree.Type != "return" {
		t.Errorf("got type %s; want return", tree.Type)
	}
	val, err := tree.Body.Eval(&Options{})
	if err != nil {
		t.Fatal(err)
	}
	if val != "3" {
		t.Errorf("got %s; want 3", val)
	}
	if t.Failed() {
		printAst(t, tree)
	}
}

func TestAstSub(t *testing.T) {
	lexr, err := lex("1-2") // useless to test another sum
	if err != nil {
		t.Fatal(err)
	}
	tree, err := astParse(lexr, "return")
	if err != nil {
		t.Fatal(err)
	}
	if tree.Type != "return" {
		t.Errorf("got type %s; want return", tree.Type)
	}
	val, err := tree.Body.Eval(&Options{})
	if err != nil {
		t.Fatal(err)
	}
	if val != "-1" {
		t.Errorf("got %s; want -1", val)
	}
	if t.Failed() {
		printAst(t, tree)
	}
}

func printAst(t *testing.T, tree *ast) {
	m, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(m))
}
