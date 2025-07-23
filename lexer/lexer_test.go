package lexer

import (
	"testing"
)

func TestLexerLiteral(t *testing.T) {
	lexr, err := Lex("1")
	if err != nil {
		t.Fatal(err)
	}
	if lexr[0].Type != Number || lexr[0].Value != "1" {
		t.Error("expecting number(1), got", lexr[0])
	}

	lexr, err = Lex("1.2")
	if err != nil {
		t.Fatal(err)
	}
	if lexr[0].Type != Number || lexr[0].Value != "1.2" {
		t.Error("expecting number(1.2), got", lexr[0])
	}

	lexr, err = Lex(".5")
	if err != nil {
		t.Fatal(err)
	}
	if lexr[0].Type != Number || lexr[0].Value != "0.5" {
		t.Error("expecting number(0.5), got", lexr[0])
	}
}

func TestLexerSum(t *testing.T) {
	lexr, err := Lex("1+2")
	if err != nil {
		t.Fatal(err)
	}
	fn := func(l []*Lexer) {
		if len(l) != 3 {
			t.Errorf("Lexer has wrong length, got %d, excepted %d", len(lexr), 3)
			printLex(t, lexr)
			return
		}
		first := l[0]
		second := l[1]
		third := l[2]
		if first.Type != Number {
			t.Errorf("got type %s; want Number", first.Type)
		}
		if first.Value != "1" {
			t.Errorf("got %s; want 1", first.Value)
		}
		if second.Type != Operator {
			t.Errorf("got type %s; want Operator", second.Type)
		}
		if second.Value != "+" {
			t.Errorf("got %s; want +", second.Value)
		}
		if third.Type != Number {
			t.Errorf("got type %s; want Number", third.Type)
		}
		if third.Value != "2" {
			t.Errorf("got %s; want 2", third.Value)
		}
		if t.Failed() {
			printLex(t, lexr)
		}
	}
	t.Log("Testing 1+2")
	fn(lexr)
	t.Log("Testing 1 + 2")
	lexr, err = Lex("1 + 2")
	if err != nil {
		t.Fatal(err)
	}
	fn(lexr)
	t.Log("Testing 1+ 2")
	lexr, err = Lex("1+ 2")
	if err != nil {
		t.Fatal(err)
	}
	fn(lexr)
	t.Log("Testing 1 +2")
	lexr, err = Lex("1 +2")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLexerSub(t *testing.T) {
	lexr, err := Lex("1-2")
	if err != nil {
		t.Fatal(err)
	}
	fn := func(l []*Lexer) {
		if len(l) != 3 {
			t.Errorf("Lexer has wrong length, got %d, excepted %d", len(lexr), 3)
			printLex(t, lexr)
			return
		}
		first := l[0]
		second := l[1]
		third := l[2]
		if first.Type != Number {
			t.Errorf("got type %s; want Number", first.Type)
		}
		if first.Value != "1" {
			t.Errorf("got %s; want 1", first.Value)
		}
		if second.Type != Operator {
			t.Errorf("got type %s; want Operator", second.Type)
		}
		if second.Value != "-" {
			t.Errorf("got %s; want -", second.Value)
		}
		if third.Type != Number {
			t.Errorf("got type %s; want Number", third.Type)
		}
		if third.Value != "2" {
			t.Errorf("got %s; want 2", third.Value)
		}
		if t.Failed() {
			printLex(t, lexr)
		}
	}
	t.Log("Testing 1-2")
	fn(lexr)
	t.Log("Testing 1 - 2")
	lexr, err = Lex("1 - 2")
	if err != nil {
		t.Fatal(err)
	}
	fn(lexr)
	t.Log("Testing 1- 2")
	lexr, err = Lex("1- 2")
	if err != nil {
		t.Fatal(err)
	}
	fn(lexr)
	t.Log("Testing 1 -2")
	lexr, err = Lex("1 -2")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLexerUnary(t *testing.T) {
	lexr, err := Lex("+1")
	if err != nil {
		t.Fatal(err)
	}
	fn := func(l []*Lexer, op string) {
		if len(l) != 2 {
			t.Errorf("Lexer has wrong length, got %d, excepted %d", len(lexr), 2)
			printLex(t, lexr)
			return
		}
		first := l[0]
		second := l[1]
		if first.Type != Operator {
			t.Errorf("got type %s; want Operator", first.Type)
		}
		if first.Value != op {
			t.Errorf("got %s; want %s", first.Value, op)
		}
		if second.Type != Number {
			t.Errorf("got type %s; want Number", second.Type)
		}
		if second.Value != "1" {
			t.Errorf("got %s; want 1", second.Value)
		}
		if t.Failed() {
			printLex(t, lexr)
		}
	}
	t.Log("Testing +1")
	fn(lexr, "+")
	t.Log("Testing -1")
	lexr, err = Lex("-1")
	if err != nil {
		t.Fatal(err)
	}
	fn(lexr, "-")
}

func TestLexerComplex(t *testing.T) {
	lexr, err := Lex("2(2+3)^2")
	if err != nil {
		t.Fatal(err)
	}
	if len(lexr) != 8 {
		t.Errorf("Lexer has wrong length, got %d, excepted %d", len(lexr), 8)
	}
	if lexr[0].Type != Number {
		t.Errorf("got type %s; want Number", lexr[0].Type)
	}
	if lexr[0].Value != "2" {
		t.Errorf("got %s; want 2", lexr[0].Value)
	}
	if lexr[1].Type != Separator {
		t.Errorf("got type %s; want Operator", lexr[1].Type)
	}
	if lexr[1].Value != "(" {
		t.Errorf("got %s; want (", lexr[1].Value)
	}
	if lexr[2].Type != Number {
		t.Errorf("got type %s; want Number", lexr[2].Type)
	}
	if lexr[2].Value != "2" {
		t.Errorf("got %s; want 2", lexr[2].Value)
	}
	if lexr[3].Type != Operator {
		t.Errorf("got type %s; want Operator", lexr[3].Type)
	}
	if lexr[3].Value != "+" {
		t.Errorf("got %s; want +", lexr[3].Value)
	}
	if lexr[4].Type != Number {
		t.Errorf("got type %s; want Number", lexr[4].Type)
	}
	if lexr[4].Value != "3" {
		t.Errorf("got %s; want 3", lexr[4].Value)
	}
	if lexr[5].Type != Separator {
		t.Errorf("got type %s; want Operator", lexr[5].Type)
	}
	if lexr[5].Value != ")" {
		t.Errorf("got %s; want )", lexr[5].Value)
	}
	if lexr[6].Type != Operator {
		t.Errorf("got type %s; want Operator", lexr[6].Type)
	}
	if lexr[6].Value != "^" {
		t.Errorf("got %s; want ^", lexr[6].Value)
	}
	if lexr[7].Type != Number {
		t.Errorf("got type %s; want Number", lexr[7].Type)
	}
	if lexr[7].Value != "2" {
		t.Errorf("got %s; want 2", lexr[7].Value)
	}
}

func TestLexer_Word(t *testing.T) {
	lexr, err := Lex("cos sin exp")
	if err != nil {
		t.Fatal(err)
	}
	if len(lexr) != 3 {
		t.Errorf("Lexer has wrong length, got %d, excepted %d", len(lexr), 3)
	}
	if lexr[0].Type != Literal {
		t.Errorf("got type %s; want Literal", lexr[0].Type)
	}
	if lexr[0].Value != "cos" {
		t.Errorf("got %s; want 'cos'", lexr[0].Value)
	}
	if lexr[1].Type != Literal {
		t.Errorf("got type %s; want Literal", lexr[1].Type)
	}
	if lexr[1].Value != "sin" {
		t.Errorf("got %s; want 'sin'", lexr[1].Value)
	}
	if lexr[2].Type != Literal {
		t.Errorf("got type %s; want Literal", lexr[2].Type)
	}
	if lexr[2].Value != "exp" {
		t.Errorf("got %s; want 'exp'", lexr[2].Value)
	}
}

func printLex(t *testing.T, lexr []*Lexer) {
	s := ""
	for _, l := range lexr {
		s += l.String() + " "
	}
	t.Log(s[:len(s)-1])
}
