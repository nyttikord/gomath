package main

import (
	"errors"
	"fmt"
	"github.com/anhgelus/gomath/interpreter"
	"github.com/anhgelus/gomath/lexer"
	"slices"
	"strconv"
)

var (
	termOperators   = []string{"+", "-"}
	factorOperators = []string{"*", "/"}
	expOperators    = []string{"^"}

	statements = []string{"let"}

	UnknownExpressionErr = errors.New("unknown expression")
	UnknownStatementErr  = errors.New("unknown statement")
)

type Ast struct {
	Type string
	Body []interpreter.Statement
}

func parse(lexed [][]*lexer.Lexer) (*Ast, error) {
	tree := Ast{Type: "program"}
	for _, l := range lexed {
		i := 0
		var stmt interpreter.Statement

		switch l[0].Value {
		case "let":
			i++
			// create new variable
			v, err := interpreter.NewMemory(l, &i)
			if err != nil {
				return nil, err
			}
			// get variable expression
			exp, err := termExpression(l, &i)
			if err != nil {
				return nil, err
			}
			v.Expression = exp
			stmt = v
		default:
			exp, err := termExpression(l, &i)
			if err != nil {
				return nil, err
			}
			stmt = &interpreter.PrintStatement{Expression: exp}
		}
		tree.Body = append(tree.Body, stmt)
	}
	return &tree, nil
}

func termExpression(l []*lexer.Lexer, i *int) (interpreter.Expression, error) {
	return binExpression(termOperators, factorExpression)(l, i)
}

func factorExpression(l []*lexer.Lexer, i *int) (interpreter.Expression, error) {
	return binExpression(factorOperators, expExpression)(l, i)
}

func expExpression(l []*lexer.Lexer, i *int) (interpreter.Expression, error) {
	return binExpression(expOperators, literalExpression)(l, i)
}

func binExpression(operators []string, sub interpreter.ExpressionFunc) interpreter.ExpressionFunc {
	return func(l []*lexer.Lexer, i *int) (interpreter.Expression, error) {
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

func literalExpression(l []*lexer.Lexer, i *int) (interpreter.Expression, error) {
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
	} else if c.Type == lexer.Literal {
		*i++
		return &interpreter.Variable{ID: c.Value}, nil
	}
	return nil, errors.Join(UnknownExpressionErr, fmt.Errorf(
		"unknown type %s('%s'): excepting a valid literal expression",
		c.Type,
		c.Value,
	))
}
