package interpreter

import (
	"errors"
	"github.com/anhgelus/gomath/lexer"
)

var (
	UnknownOperationErr = errors.New("unknown operation")
)

type Expression func(l []*lexer.Lexer, i *int) (Statement, error)

type Statement interface {
	Eval() (*Fraction, error)
}

type BinaryOperation struct {
	Operator    string
	Left, Right Statement
}

type UnaryOperation struct {
	Operator string
	Left     Statement
}

type Literal struct {
	Value *Fraction
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
