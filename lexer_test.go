package gomath

import (
	"errors"
	"testing"
)

func TestLexerAdd(t *testing.T) {
	lexr, err := lex("1+2")
	if err != nil {
		t.Fatal(err)
	}
	fn := func(l []*lexer) {
		if len(l) != 3 {
			t.Errorf("lexer has wrong length, got %d, excepted %d", len(lexr), 3)
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
		t.Fatal(err)
	}
}
