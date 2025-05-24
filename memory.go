package gomath

import (
	"errors"
	m "math"
)

var (
	ErrInvalidVariableDeclaration = errors.New("invalid variable declaration")
	ErrInvalidFunctionDeclaration = errors.New("invalid function declaration")
	ErrVariableMustHaveExpression = errors.New("variable must have an expression")
	ErrUnknownVariable            = errors.New("unknown variable")
	ErrUnknownFunction            = errors.New("unknown function")
)

var (
	variables           = map[string]*fraction{}
	predefinedVariables = map[string]*fraction{}

	functions = map[string]*mathFunction{}
)

func init() {
	add := func(n string, v float64) {
		f, err := floatToFraction(v)
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

// newFunction handles the creation of a new mathFunction
func newFunction(l []*lexer, i *int) (*mathFunction, error) {
	if *i+9 >= len(l) {
		return nil, ErrInvalidFunctionDeclaration
	}
	if l[*i+1].Type != Literal && l[*i+1].Value != "in" {
		return nil, ErrInvalidFunctionDeclaration
	}
	if l[*i+2].Type != Literal {
		return nil, ErrInvalidFunctionDeclaration
	}
	if l[*i+3].Type != Separator && l[*i+1].Value != "," {
		return nil, ErrInvalidFunctionDeclaration
	}
	variable := l[*i].Value
	rawDef := l[*i+2].Value
	*i += 4
	if l[*i].Type != Literal {
		return nil, ErrInvalidFunctionDeclaration
	}
	name := l[*i].Value
	*i += 1
	if l[*i].Type != Operator && l[*i].Value != "{" {
		return nil, ErrInvalidFunctionDeclaration
	}
	if l[*i+1].Type != Literal && l[*i].Value != variable {
		return nil, ErrInvalidFunctionDeclaration
	}
	if l[*i+2].Type != Operator && l[*i+2].Value != "}" {
		return nil, ErrInvalidFunctionDeclaration
	}
	if l[*i+3].Type != Operator && l[*i+2].Value != "=" {
		return nil, ErrInvalidFunctionDeclaration
	}
	*i += 4
	rel := lexToRel(l[*i:])
	def, err := parseSpace(rawDef)
	if err != nil {
		return nil, err
	}
	return &mathFunction{Definition: def, Relation: rel, Name: name, Variable: variable}, nil
}

// isInMemory checks if the given id is already used
func isInMemory(id string) bool {
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
