package main

import (
	"errors"
	"github.com/anhgelus/gomath/interpreter"
	"github.com/anhgelus/gomath/lexer"
	"slices"
	"strconv"
)

var (
	termOperators   = []string{"+", "-"}
	factorOperators = []string{"*", "/"}
	expOperators    = []string{"^"}

	UnknownExpressionErr = errors.New("unknown expression")
)

type Ast struct {
	Type string
	Body []*interpreter.Statement
}

func parse(lexed [][]*lexer.Lexer) (*Ast, error) {
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

func termExpression(l []*lexer.Lexer, i *int) (interpreter.Statement, error) {
	return binExpression(termOperators, factorExpression)(l, i)
}

func factorExpression(l []*lexer.Lexer, i *int) (interpreter.Statement, error) {
	return binExpression(factorOperators, expExpression)(l, i)
}

func expExpression(l []*lexer.Lexer, i *int) (interpreter.Statement, error) {
	return binExpression(expOperators, literalExpression)(l, i)
}

func binExpression(operators []string, sub interpreter.Expression) interpreter.Expression {
	return func(l []*lexer.Lexer, i *int) (interpreter.Statement, error) {
		left, err := sub(l, i)
		if err != nil {
			return nil, err
		}
		for *i < len(l) && slices.Contains(operators, l[*i].Value) {
			op := l[*i].Value
			*i++
			right, err := sub(l, i)
			if err != nil {
				return nil, err
			}
			left = &interpreter.BinaryOperation{
				Operator: op,
				Left:     left,
				Right:    right,
			}
		}
		return left, nil
	}
}

func literalExpression(l []*lexer.Lexer, i *int) (interpreter.Statement, error) {
	c := l[*i]
	if c.Type == lexer.Number {
		v, err := strconv.ParseFloat(c.Value, 64)
		if err != nil {
			return nil, err
		}
		*i++
		f, err := interpreter.FloatToFraction(v)
		if err != nil {
			return nil, err
		}
		return &interpreter.Literal{Value: f}, nil
	}
	return nil, errors.Join(UnknownExpressionErr, errors.New("unknown type "+c.Type))
}
