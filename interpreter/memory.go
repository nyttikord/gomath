package interpreter

import (
	"errors"
	math "github.com/anhgelus/gomath/interpreter/math"
	"github.com/anhgelus/gomath/lexer"
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
	variables           = map[string]*math.Fraction{}
	predefinedVariables = map[string]*math.Fraction{}

	functions = map[string]*Function{}
)

func init() {
	add := func(n string, v float64) {
		f, err := math.FloatToFraction(v)
		if err != nil {
			panic(err)
		}
		predefinedVariables[n] = f
	}
	add("pi", m.Pi)
	add("e", m.E)
	add("phi", m.Phi)
}

type Memory struct {
	ID         string
	Expression Expression
}

type Function struct {
	Definition math.Space
	Relation   *Relation
	Name       string
	Variable   string
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

func NewFunction(l []*lexer.Lexer, i *int) (*Function, error) {
	if *i+9 >= len(l) {
		return nil, InvalidFunctionDeclarationErr
	}
	if l[*i+1].Type != lexer.Literal && l[*i+1].Value != "in" {
		return nil, InvalidFunctionDeclarationErr
	}
	if l[*i+2].Type != lexer.Literal {
		return nil, InvalidFunctionDeclarationErr
	}
	if l[*i+3].Type != lexer.Separator && l[*i+1].Value != "," {
		return nil, InvalidFunctionDeclarationErr
	}
	variable := l[*i].Value
	rawDef := l[*i+2].Value
	*i += 4
	if l[*i].Type != lexer.Literal {
		return nil, InvalidFunctionDeclarationErr
	}
	name := l[*i].Value
	*i += 1
	if l[*i].Type != lexer.Operator && l[*i].Value != "{" {
		return nil, InvalidFunctionDeclarationErr
	}
	if l[*i+1].Type != lexer.Literal && l[*i].Value != variable {
		return nil, InvalidFunctionDeclarationErr
	}
	if l[*i+2].Type != lexer.Operator && l[*i+2].Value != "}" {
		return nil, InvalidFunctionDeclarationErr
	}
	if l[*i+3].Type != lexer.Operator && l[*i+2].Value != "=" {
		return nil, InvalidFunctionDeclarationErr
	}
	*i += 4
	rel := LexToRel(l[*i:])
	def, err := math.ParseSpace(rawDef)
	if err != nil {
		return nil, err
	}
	return &Function{Definition: def, Relation: rel, Name: name, Variable: variable}, nil
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

func (v *Memory) Eval(*Options) error {
	f, err := v.Expression.Eval()
	if err != nil {
		return err
	}
	variables[v.ID] = f
	return nil
}

func (f *Function) Eval(*Options) error {
	functions[f.Name] = f
	return nil
}
