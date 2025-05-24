package gomath

import (
	"errors"
	"fmt"
	math2 "github.com/anhgelus/gomath/math"
)

var (
	UnknownOperationErr = errors.New("unknown operation")
	LexerNotValidErr    = errors.New("lexer is not valid")
	NumberNotInSpaceErr = errors.New("number is not in the definition space")
)

type expressionFunc func(l []*lexer, i *int) (expression, error)

type expression interface {
	Eval() (*math2.Fraction, error)
}

type operator string

type binaryOperation struct {
	Operator    operator
	Left, Right expression
}

type unaryOperation struct {
	Operator   operator
	Expression expression
}

type evaluateOperation struct {
	FunctionName string
	Expression   expression
}

type literalExp struct {
	Value *math2.Fraction
}

type variable struct {
	ID string
}

type predefinedVariable variable

type relation string

func (b *binaryOperation) Eval() (*math2.Fraction, error) {
	lb, err := b.Left.Eval()
	if err != nil {
		return nil, err
	}
	lr, err := b.Right.Eval()
	if err != nil {
		return nil, err
	}
	switch b.Operator {
	case "+":
		return lb.Add(lr), nil
	case "-":
		return lb.Sub(lr), nil
	case "*":
		return lb.Mul(lr), nil
	case "/":
		return lb.Div(lr), nil
	case "^":
		return lb.Pow(lr), nil
	default:
		return nil, errors.Join(UnknownOperationErr, errors.New("operation "+string(b.Operator)+" is not supported"))
	}
}

func (b *unaryOperation) Eval() (*math2.Fraction, error) {
	lb, err := b.Expression.Eval()
	if err != nil {
		return nil, err
	}
	switch b.Operator {
	case "+":
		return lb, nil
	case "-":
		return lb.Mul(math2.IntToFraction(-1)), nil
	default:
		return nil, errors.Join(UnknownOperationErr, errors.New("operation "+string(b.Operator)+" is not supported"))
	}
}

func (e *evaluateOperation) Eval() (*math2.Fraction, error) {
	f, ok := functions[e.FunctionName]
	if !ok {
		return nil, errors.Join(UnknownFunctionErr, fmt.Errorf("undefined function %s", e.FunctionName))
	}
	return f.Relation.Eval(f.Definition, f.Variable, e.Expression)
}

func (l *literalExp) Eval() (*math2.Fraction, error) {
	return l.Value, nil
}

func (v *variable) Eval() (*math2.Fraction, error) {
	val, ok := variables[v.ID]
	if !ok {
		return nil, errors.Join(UnknownVariableErr, fmt.Errorf("undefined variable %s", v.ID))
	}
	return val, nil
}

func (v *predefinedVariable) Eval() (*math2.Fraction, error) {
	val, ok := predefinedVariables[v.ID]
	if !ok {
		return nil, errors.Join(UnknownVariableErr, fmt.Errorf("undefined variable \\%s", v.ID))
	}
	return val, nil
}

func LexToRel(lexers []*lexer) *relation {
	var s relation
	for _, l := range lexers {
		s += relation(l.Value)
	}
	return &s
}

func (r *relation) String() string {
	return string(*r)
}

func (r *relation) Eval(def math2.Space, variable string, val expression) (*math2.Fraction, error) {
	lexed, err := lex(r.String())
	if err != nil {
		return nil, err
	}
	if len(lexed) != 1 {
		return nil, LexerNotValidErr
	}
	var lexr []*lexer
	for _, l := range lexed {
		// replace all x by their value in brackets
		if l.Type == Literal && l.Value == variable {
			fr, err := val.Eval()
			if err != nil {
				return nil, err
			}
			if !def.Contains(fr) {
				return nil, errors.Join(NumberNotInSpaceErr, fmt.Errorf("%s is not in %s", fr.String(), def.String()))
			}
			lexr = append(lexr, &lexer{Type: Separator, Value: "("})
			l.Type = Number
			l.Value = fr.String()
			lexr = append(lexr, l)
			lexr = append(lexr, &lexer{Type: Separator, Value: ")"})
		} else {
			lexr = append(lexr, l)
		}
	}
	i := 0
	exp, err := termExpression(lexed, &i)
	if err != nil {
		return nil, err
	}
	return exp.Eval()
}
