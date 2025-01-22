package interpreter

import (
	"errors"
	"github.com/anhgelus/gomath/lexer"
	"math"
)

var (
	UnknownOperationErr = errors.New("unknown operation")
)

type Expression func(l []*lexer.Lexer, i *int) (Statement, error)

type Statement interface {
	eval() (float64, error)
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
	Value float64
}

func (b *BinaryOperation) eval() (float64, error) {
	lb, err := b.Left.eval()
	if err != nil {
		return 0, err
	}
	lr, err := b.Right.eval()
	if err != nil {
		return 0, err
	}
	switch b.Operator {
	case "+":
		return lb + lr, nil
	case "-":
		return lb - lr, nil
	case "*":
		return lb * lr, nil
	case "/":
		return lb / lr, nil
	case "^":
		return math.Pow(lb, lr), nil
	default:
		return 0, errors.Join(UnknownOperationErr, errors.New("operation "+b.Operator+" is not supported"))
	}
}

func (b *UnaryOperation) eval() (float64, error) {
	lb, err := b.Left.eval()
	if err != nil {
		return 0, err
	}
	switch b.Operator {
	case "+":
		return lb, nil
	case "-":
		return -lb, nil
	default:
		return 0, errors.Join(UnknownOperationErr, errors.New("operation "+b.Operator+" is not supported"))
	}
}

func (l *Literal) eval() (float64, error) {
	return l.Value, nil
}
