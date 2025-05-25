package gomath

import (
	"testing"
)

func TestEvalSum(t *testing.T) {
	genericTest(t, "1+2", "3")
}

func TestEvalSub(t *testing.T) {
	genericTest(t, "1-2", "-1")
}

func TestEvalAddUnary(t *testing.T) {
	genericTest(t, "1+-2", "-1")
}

func TestEvalMult(t *testing.T) {
	t.Log("Testing 1+2")
	genericTest(t, "2*3", "6")
}

func TestEvalMultUnary(t *testing.T) {
	genericTest(t, "2*-3", "-6")
}

func TestEvalDiv(t *testing.T) {
	genericTest(t, "2/3", "2/3")
}

func TestEvalDivUnary(t *testing.T) {
	genericTest(t, "2/-3", "-2/3")
}

func TestEvalDivDecimal(t *testing.T) {
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
		t.Log(tree)
	}
}

func TestEvalPriority(t *testing.T) {
	t.Log("testing 2*(1+2)")
	genericTest(t, "2*(1+2)", "6")
	t.Log("testing 2*1+2")
	genericTest(t, "2*1+2", "4")
	t.Log("testing 2*(1+2)^2")
	genericTest(t, "2*(1+2)^2", "18")
}

func TestEvalOmitMultSigne(t *testing.T) {
	t.Log("testing 2(3+2)")
	genericTest(t, "2(3+2)", "10")
	t.Log("testing 2^2(3+2)")
	genericTest(t, "2^2(3+2)", "20")
	t.Log("testing 2(3+2)^2")
	genericTest(t, "2(3+2)^2", "50")
}

func TestEvalPrioritySpecialCase(t *testing.T) {
	t.Log("testing -3^2")
	genericTest(t, "-3^2", "-9")
	t.Log("testing 6/2(1+2)")
	genericTest(t, "6/2(1+2)", "9")
	t.Log("testing 3^2^3")
	genericTest(t, "3^2^3", "729")
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
		t.Log(tree)
	}
}
