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
	eval() (*Fraction, error)
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

func (b *BinaryOperation) eval() (*Fraction, error) {
	lb, err := b.Left.eval()
	if err != nil {
		return NullFraction, err
	}
	lr, err := b.Right.eval()
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

func (b *UnaryOperation) eval() (*Fraction, error) {
	lb, err := b.Left.eval()
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

func (l *Literal) eval() (*Fraction, error) {
	return l.Value, nil
}
