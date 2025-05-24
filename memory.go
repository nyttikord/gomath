package gomath

import (
	"errors"
	m "math"
)

var (
	InvalidVariableDeclarationErr = errors.New("invalid variable declaration")
	InvalidFunctionDeclarationErr = errors.New("invalid function declaration")
	VariableMustHaveExpressionErr = errors.New("variable must have an expression")
	UnknownVariableErr            = errors.New("unknown variable")
	UnknownFunctionErr            = errors.New("unknown function")
)

var (
	variables           = map[string]*Fraction{}
	predefinedVariables = map[string]*Fraction{}

	functions = map[string]*mathFunction{}
)

func init() {
	add := func(n string, v float64) {
		f, err := FloatToFraction(v)
		if err != nil {
			panic(err)
		}
		predefinedVariables[n] = f
	}
	add("pi", m.Pi)
	add("e", m.E)
	add("phi", m.Phi)
}

type mathFunction struct {
	Definition Space
	Relation   *relation
	Name       string
	Variable   string
}

func NewFunction(l []*lexer, i *int) (*mathFunction, error) {
	if *i+9 >= len(l) {
		return nil, InvalidFunctionDeclarationErr
	}
	if l[*i+1].Type != Literal && l[*i+1].Value != "in" {
		return nil, InvalidFunctionDeclarationErr
	}
	if l[*i+2].Type != Literal {
		return nil, InvalidFunctionDeclarationErr
	}
	if l[*i+3].Type != Separator && l[*i+1].Value != "," {
		return nil, InvalidFunctionDeclarationErr
	}
	variable := l[*i].Value
	rawDef := l[*i+2].Value
	*i += 4
	if l[*i].Type != Literal {
		return nil, InvalidFunctionDeclarationErr
	}
	name := l[*i].Value
	*i += 1
	if l[*i].Type != Operator && l[*i].Value != "{" {
		return nil, InvalidFunctionDeclarationErr
	}
	if l[*i+1].Type != Literal && l[*i].Value != variable {
		return nil, InvalidFunctionDeclarationErr
	}
	if l[*i+2].Type != Operator && l[*i+2].Value != "}" {
		return nil, InvalidFunctionDeclarationErr
	}
	if l[*i+3].Type != Operator && l[*i+2].Value != "=" {
		return nil, InvalidFunctionDeclarationErr
	}
	*i += 4
	rel := lexToRel(l[*i:])
	def, err := ParseSpace(rawDef)
	if err != nil {
		return nil, err
	}
	return &mathFunction{Definition: def, Relation: rel, Name: name, Variable: variable}, nil
}

func IsInMemory(id string) bool {
	_, ok := variables[id]
	if ok {
		return true
	} else if _, ok = functions[id]; ok {
		return true
	}
	return false
}

func (f *mathFunction) Eval(*Options) error {
	functions[f.Name] = f
	return nil
}
