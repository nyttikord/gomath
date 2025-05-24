package gomath

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var (
	termOperators   = []operator{"+", "-"}
	factorOperators = []operator{"*", "/"}
	expOperators    = []operator{"^"}

	UnknownExpressionErr = errors.New("unknown expression")
	UnknownStatementErr  = errors.New("unknown statement")
	WrongExpressionErr   = errors.New("wrong expression")
)

type astType string

type ast struct {
	Type astType
	Body statement
}

func astParse(lexed []*lexer, tpe astType) (*ast, error) {
	tree := ast{Type: tpe}
	i := 0
	exp, err := termExpression(lexed, &i)
	if err != nil {
		return nil, err
	}
	tree.Body = &returnStatement{Expression: exp}
	return &tree, nil
}

func termExpression(l []*lexer, i *int) (expression, error) {
	return binExpression(termOperators, factorExpression)(l, i)
}

func factorExpression(l []*lexer, i *int) (expression, error) {
	return binExpression(factorOperators, expExpression)(l, i)
}

func expExpression(l []*lexer, i *int) (expression, error) {
	return binExpression(expOperators, literalExpression)(l, i)
}

func binExpression(operators []operator, sub expressionFunc) expressionFunc {
	return func(l []*lexer, i *int) (expression, error) {
		left, err := sub(l, i)
		if err != nil {
			return nil, err
		}
		for *i < len(l) && slices.Contains(operators, operator(l[*i].Value)) {
			op := operator(l[*i].Value)
			*i++
			right, err := sub(l, i)
			if err != nil {
				return nil, err
			}
			left = &binaryOperation{
				Operator: op,
				Left:     left,
				Right:    right,
			}
		}
		return left, nil
	}
}

func operatorExpression(l []*lexer, i *int) (expression, error) {
	c := l[*i]
	if c.Type == Operator && c.Value == "{" {
		*i++
		exp, err := termExpression(l, i)
		if err != nil {
			return nil, err
		}
		if l[*i].Value != "}" {
			return exp, errors.Join(WrongExpressionErr, errors.New("} excepted"))
		}
		*i++
		return exp, nil
	}
	return literalExpression(l, i)
}

func literalExpression(l []*lexer, i *int) (expression, error) {
	c := l[*i]
	*i++
	switch c.Type {
	case Number:
		v, err := strconv.ParseFloat(c.Value, 64)
		if err != nil {
			return nil, err
		}
		f, err := FloatToFraction(v)
		if err != nil {
			return nil, err
		}
		return &literalExp{Value: f}, nil
	case Literal:
		if strings.HasPrefix(c.Value, "\\") {
			return &predefinedVariable{ID: c.Value[1:]}, nil
		} else if *i < len(l) && l[*i].Type == Operator && l[*i].Value == "{" {
			name := c.Value
			*i++
			exp, err := operatorExpression(l, i)
			if err != nil {
				return nil, err
			}
			return &evaluateOperation{Expression: exp, FunctionName: name}, nil
		}
		return &variable{ID: c.Value}, nil
	case Separator:
		if c.Value == "(" {
			exp, err := termExpression(l, i)
			if err != nil {
				return nil, err
			}
			if l[*i].Value != ")" {
				return nil, errors.Join(WrongExpressionErr, errors.New(") excepted"))
			}
			*i++
			return exp, nil
		}
	case Operator:
		exp, err := expExpression(l, i)
		if err != nil {
			return nil, err
		}
		return &unaryOperation{Operator: operator(c.Value), Expression: exp}, nil
	}
	return nil, errors.Join(UnknownExpressionErr, fmt.Errorf(
		"unknown type %s('%s'): excepting a valid literal expression",
		c.Type,
		c.Value,
	))
}
