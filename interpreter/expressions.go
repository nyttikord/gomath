package interpreter

import (
	"errors"
	"github.com/anhgelus/gomath/lexer"
)

var (
	UnknownOperationErr = errors.New("unknown operation")
)

type ExpressionFunc func(l []*lexer.Lexer, i *int) (Expression, error)

type Expression interface {
	Eval() (*Fraction, error)
}

type BinaryOperation struct {
	Operator    string
	Left, Right Expression
}

type UnaryOperation struct {
	Operator string
	Left     Expression
}

type Literal struct {
	Value *Fraction
}

type Variable struct {
	ID string
}

func (b *BinaryOperation) Eval() (*Fraction, error) {
	lb, err := b.Left.Eval()
	if err != nil {
		return NullFraction, err
	}
	lr, err := b.Right.Eval()
	if err != nil {
		return NullFraction, err
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
		return NullFraction, errors.Join(UnknownOperationErr, errors.New("operation "+b.Operator+" is not supported"))
	}
}

func (b *UnaryOperation) Eval() (*Fraction, error) {
	lb, err := b.Left.Eval()
	if err != nil {
		return NullFraction, err
	}
	switch b.Operator {
	case "+":
		return lb, nil
	case "-":
		return lb.Mul(IntToFraction(-1)), nil
	default:
		return NullFraction, errors.Join(UnknownOperationErr, errors.New("operation "+b.Operator+" is not supported"))
	}
}

func (l *Literal) Eval() (*Fraction, error) {
	return l.Value, nil
}

func (v *Variable) Eval() (*Fraction, error) {
	return GetValueInMemory(v.ID)
}
