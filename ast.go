package main

import (
	"errors"
	"slices"
	"strconv"
)

var (
	UnknownOperationErr = errors.New("unknown operation")

	termOperators   = []string{"+", "-"}
	factorOperators = []string{"*", "/"}
)

type Ast struct {
	Type string
	Body []*Statement
}

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

func parse(lexed [][]*Lexer) (*Ast, error) {
	tree := Ast{Type: "program"}
	for _, l := range lexed {
		i := 0
		stmt, err := termExpression(l, &i)
		if err != nil {
			return nil, err
		}
		tree.Body = append(tree.Body, &stmt)
	}
	return &tree, nil
}

func termExpression(l []*Lexer, i *int) (Statement, error) {
	left, err := factorExpression(l, i)
	if err != nil {
		return nil, err
	}
	for *i < len(l) && slices.Contains(termOperators, l[*i].Value) {
		op := l[*i].Value
		*i++
		right, err := factorExpression(l, i)
		if err != nil {
			return nil, err
		}
		left = &BinaryOperation{
			Operator: op,
			Left:     left,
			Right:    right,
		}
	}
	return left, nil
}

func factorExpression(l []*Lexer, i *int) (Statement, error) {
	left, err := literalExpression(l, i)
	if err != nil {
		return nil, err
	}
	for *i < len(l) && slices.Contains(factorOperators, l[*i].Value) {
		op := l[*i].Value
		*i++
		right, err := literalExpression(l, i)
		if err != nil {
			return nil, err
		}
		left = &BinaryOperation{
			Operator: op,
			Left:     left,
			Right:    right,
		}
	}
	return left, nil
}

func literalExpression(l []*Lexer, i *int) (Statement, error) {
	c := l[*i]
	if c.Type == number {
		v, err := strconv.ParseFloat(c.Value, 64)
		if err != nil {
			return nil, err
		}
		*i++
		return &Literal{Value: v}, nil
	}
	return nil, nil
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
		return lb + lr, nil
	case "*":
		return lb + lr, nil
	case "/":
		return lb + lr, nil
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
