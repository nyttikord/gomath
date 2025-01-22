package interpreter

import (
	"errors"
	"github.com/anhgelus/gomath/lexer"
)

var (
	InvalidVariableDeclarationErr = errors.New("invalid variable declaration")
	VariableMustHaveExpressionErr = errors.New("variable must have an expression")
	UnknownVariableErr            = errors.New("unknown variable")
)

var variables = map[string]*Fraction{}

type Memory struct {
	ID         string
	Expression Expression
}

func NewMemory(l []*lexer.Lexer, i *int) (*Memory, error) {
	if *i+2 >= len(l) {
		return nil, InvalidVariableDeclarationErr
	}
	if l[*i+1].Type != lexer.Operator && l[*i+1].Value != "=" {
		return nil, VariableMustHaveExpressionErr
	}
	id := l[*i].Value
	*i += 2
	return &Memory{ID: id}, nil
}

func IsInMemory(id string) bool {
	_, ok := variables[id]
	return ok
}

func GetValueInMemory(id string) (*Fraction, error) {
	v, ok := variables[id]
	if !ok {
		return nil, UnknownVariableErr
	}
	return v, nil
}

func (v *Memory) Eval(*Options) error {
	f, err := v.Expression.Eval()
	if err != nil {
		return err
	}
	variables[v.ID] = f
	return nil
}
