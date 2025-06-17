package gomath

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrUnknownOperation is thrown when GoMath doesn't know the operation used
	ErrUnknownOperation = errors.New("unknown operation")
	// ErrNumberNotInSpace is thrown when the number is not in the definition space
	ErrNumberNotInSpace = errors.New("number is not in the definition space")
)

type expressionFunc func(l []*lexer, i *int) (expression, error)

type expression interface {
	// Eval the expression
	Eval() (*fraction, error)
	// RenderLatex the expression
	RenderLatex() (string, priority, error)
}

type separator string

type priority uint8

const (
	termPriority    priority = 0
	factorPriority  priority = 1
	expPriority     priority = 2
	unaryPriority   priority = 3
	literalPriority priority = 4
)

type literalExp struct {
	Value *fraction
}

type variable struct {
	ID        string
	OmitSlash bool
}

type function struct {
	ID  string
	exp expression
}

type predefinedVariable variable
type predefinedFunction function

type relation func(*fraction) *fraction

func (l *literalExp) Eval() (*fraction, error) {
	return l.Value, nil
}

func (l *literalExp) RenderLatex() (string, priority, error) {
	return l.Value.String(), literalPriority, nil
}

func (v *predefinedVariable) Eval() (*fraction, error) {
	val, ok := predefinedVariables[v.ID]
	if !ok {
		return nil, errors.Join(genErrUnknownVariable(v.ID), fmt.Errorf("undefined variable %s", v.ID))
	}
	return val.Val, nil
}

func (v *predefinedVariable) RenderLatex() (string, priority, error) {
	_, ok := predefinedVariables[v.ID]
	if !ok {
		return "", literalPriority, errors.Join(genErrUnknownVariable(v.ID), fmt.Errorf("undefined variable %s", v.ID))
	}
	if v.OmitSlash {
		return v.ID, literalPriority, nil
	}
	return `\` + v.ID, literalPriority, nil
}

func (f *predefinedFunction) Eval() (*fraction, error) {
	fn, ok := predefinedFunctions[f.ID]
	if !ok {
		return nil, errors.Join(genErrUnknownVariable(f.ID), fmt.Errorf("undefined variable %s", f.ID))
	}
	val, err := f.exp.Eval()
	if err != nil {
		return nil, err
	}
	return fn.Eval(val)
}

func (f *predefinedFunction) RenderLatex() (string, priority, error) {
	_, ok := predefinedFunctions[f.ID]
	if !ok {
		return "", literalPriority, errors.Join(genErrUnknownVariable(f.ID), fmt.Errorf("undefined variable %s", f.ID))
	}
	val, _, err := f.exp.RenderLatex()
	if err != nil {
		return "", literalPriority, err
	}
	return fmt.Sprintf(`\%s\left(%s\right)`, f.ID, val), literalPriority, nil
}

func (r *relation) Eval(f *fraction) *fraction {
	return (*r)(f)
}

func handleLatexParenthesis(s string, stringPriority, currentPriority priority) string {
	if strings.Contains(s, " ") && stringPriority < currentPriority {
		s = `\left(` + s + `\right)`
	}
	return s
}
