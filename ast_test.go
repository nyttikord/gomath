package gomath

import (
	"errors"
	"testing"
)

func TestAstErrors(t *testing.T) {
	genericTestAstError := func(exp string, exceptedErr error) {
		lexr, err := lex(exp)
		if err != nil {
			t.Fatal(err)
		}
		tree, err := astParse(lexr, astTypeCalculation)
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
