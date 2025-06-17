package ast

import (
	"errors"
	"github.com/nyttikord/gomath/lexer"
	"testing"
)

func TestAstErrors(t *testing.T) {
	genericTestAstError := func(exp string, exceptedErr error) {
		lexr, err := lexer.Lex(exp)
		if err != nil {
			t.Fatal(err)
		}
		tree, err := Parse(lexr, TypeCalculation)
		if err == nil {
			t.Errorf("expected error %s", exceptedErr)
		} else if !errors.Is(err, exceptedErr) {
			t.Errorf("got %v; want %v", err, exceptedErr)
		}
		if t.Failed() {
			t.Log(lexr)
			t.Log(tree)
		}
	}
	genericTestAstError("1+1)", ErrUnknownExpression)
	genericTestAstError("(1+1", ErrInvalidExpression)
	genericTestAstError("1+1+", ErrInvalidExpression)
	//genericTestAstError("1×1+1", ErrUnknownVariable) // will be valid when omission between number and literal is added
}
