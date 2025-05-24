package gomath

import (
	"encoding/json"
	"testing"
)

func TestAstSum(t *testing.T) {
	genericTest(t, "1+2", "3")
}

func TestAstSub(t *testing.T) {
	genericTest(t, "1-2", "-1")
}

func TestAstAddUnary(t *testing.T) {
	genericTest(t, "1+-2", "-1")
}

func TestAstMult(t *testing.T) {
	genericTest(t, "2*3", "6")
}

func TestAstMultUnary(t *testing.T) {
	genericTest(t, "2*-3", "-6")
}

func TestAstDiv(t *testing.T) {
	genericTest(t, "2/3", "2/3")
}

func TestAstDivUnary(t *testing.T) {
	genericTest(t, "2/-3", "-2/3")
}

func TestAstDivDecimal(t *testing.T) {
	lexr, err := lex("1/10")
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
	val, err := tree.Body.Eval(&Options{true})
	if err != nil {
		t.Fatal(err)
	}
	if val != "0.1" {
		t.Errorf("got %s; want %s", val, "0.1")
	}
	if t.Failed() {
		printAst(t, tree)
	}
}

func genericTest(t *testing.T, exp string, excepted string) {
	lexr, err := lex(exp)
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
	if val != excepted {
		t.Errorf("got %s; want %s", val, excepted)
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
