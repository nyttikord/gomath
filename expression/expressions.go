package expression

import (
	"errors"
	"github.com/nyttikord/gomath/math"
	"strings"
)

var (
	// ErrUnknownOperation is thrown when GoMath doesn't know the operation used
	ErrUnknownOperation = errors.New("unknown operation")
	// ErrNumberNotInSpace is thrown when the number is not in the definition space
	ErrNumberNotInSpace = errors.New("number is not in the definition space")
)

type Expression interface {
	// Eval the Expression
	Eval() (*math.Fraction, error)
	// RenderLatex the Expression
	RenderLatex() (string, priority, error)
}

type priority uint8

const (
	termPriority    priority = 0
	factorPriority  priority = 1
	expPriority     priority = 2
	unaryPriority   priority = 3
	literalPriority priority = 4
)

type constExp struct {
	Value *math.Fraction
}

type variable struct {
	ID        string
	OmitSlash bool
}

type function struct {
	ID  string
	exp Expression
}

func Const(f *math.Fraction) Expression {
	return &constExp{f}
}

func (l *constExp) Eval() (*math.Fraction, error) {
	return l.Value, nil
}

func (l *constExp) RenderLatex() (string, priority, error) {
	return l.Value.String(), literalPriority, nil
}

func handleLatexParenthesis(s string, stringPriority, currentPriority priority) string {
	if strings.Contains(s, " ") && stringPriority < currentPriority {
		s = `\left(` + s + `\right)`
	}
	return s
}
