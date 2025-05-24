package gomath

import (
	"errors"
	"testing"
)

func TestLexerSum(t *testing.T) {
	lexr, err := lex("1+2")
	if err != nil {
		t.Fatal(err)
	}
	fn := func(l []*lexer) {
		if len(l) != 3 {
			t.Errorf("lexer has wrong length, got %d, excepted %d", len(lexr), 3)
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
	lexr, err = lex("1 + 2")
	if err != nil {
		t.Fatal(err)
	}
	fn(lexr)
	t.Log("Testing 1+ 2")
	lexr, err = lex("1+ 2")
	if err != nil {
		t.Fatal(err)
	}
	fn(lexr)
	t.Log("Testing error 1 +2")
	lexr, err = lex("1 +2")
	if !errors.Is(err, SameTypeFollowErr) {
		printLex(t, lexr)
		t.Fatal(err)
	}
}

func TestLexerSub(t *testing.T) {
	lexr, err := lex("1-2")
	if err != nil {
		t.Fatal(err)
	}
	fn := func(l []*lexer) {
		if len(l) != 3 {
			t.Errorf("lexer has wrong length, got %d, excepted %d", len(lexr), 3)
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
	lexr, err = lex("1 - 2")
	if err != nil {
		t.Fatal(err)
	}
	fn(lexr)
	t.Log("Testing 1- 2")
	lexr, err = lex("1- 2")
	if err != nil {
		t.Fatal(err)
	}
	fn(lexr)
	t.Log("Testing error 1 -2")
	lexr, err = lex("1 -2")
	if !errors.Is(err, SameTypeFollowErr) {
		printLex(t, lexr)
		t.Fatal(err)
	}
}

func TestLexerUnary(t *testing.T) {
	lexr, err := lex("+1")
	if err != nil {
		t.Fatal(err)
	}
	fn := func(l []*lexer, op string) {
		if len(l) != 2 {
			t.Errorf("lexer has wrong length, got %d, excepted %d", len(lexr), 2)
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
	fn(lexr, "+")
	lexr, err = lex("-1")
	if err != nil {
		t.Fatal(err)
	}
	fn(lexr, "-")
}

func printLex(t *testing.T, lexr []*lexer) {
	s := ""
	for _, l := range lexr {
		s += l.String() + " "
	}
	t.Log(s[:len(s)-1])
}
