package gomath

import (
	"encoding/json"
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

	// ErrUnknownExpression is thrown when GoMath does not know the expression
	ErrUnknownExpression = errors.New("unknown expression")
	// ErrInvalidExpression is thrown when the given expression's syntax is invalid
	ErrInvalidExpression = errors.New("invalid expression")
)

type astType string

type ast struct {
	Type astType
	Body statement
}

func (a *ast) String() string {
	m, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return ""
	}
	return string(m)
}

// astParse the given lexer and returns an ast
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
	return binExpression(termOperators, omitExpression)(l, i)
}

func omitExpression(l []*lexer, i *int) (expression, error) {
	sub := factorExpression
	left, err := sub(l, i)
	if err != nil {
		return nil, err
	}
	for *i < len(l) && l[*i].Value == "(" {
		right, err := sub(l, i)
		if err != nil {
			return nil, err
		}
		left = &binaryOperation{
			Operator: "*",
			Left:     left,
			Right:    right,
		}
	}
	return left, nil

}

func factorExpression(l []*lexer, i *int) (expression, error) {
	return binExpression(factorOperators, expExpression)(l, i)
}

func expExpression(l []*lexer, i *int) (expression, error) {
	return binExpression(expOperators, literalExpression)(l, i)
}

func binExpression(ops []operator, sub expressionFunc) expressionFunc {
	return func(l []*lexer, i *int) (expression, error) {
		left, err := sub(l, i)
		if err != nil {
			return nil, err
		}
		for *i < len(l) && slices.Contains(ops, operator(l[*i].Value)) {
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

func literalExpression(l []*lexer, i *int) (expression, error) {
	c := l[*i]
	*i++
	switch c.Type {
	case Number:
		v, err := strconv.ParseFloat(c.Value, 64)
		if err != nil {
			return nil, err
		}
		f, err := floatToFraction(v)
		if err != nil {
			return nil, err
		}
		return &literalExp{Value: f}, nil
	case Literal:
		if strings.HasPrefix(c.Value, `\`) {
			return predefinedExpression(l, i, c.Value[1:])
		}
		return nil, errors.Join(ErrUnknownExpression, fmt.Errorf("unknown literal %s", c.Value))
	case Separator:
		if c.Value == "(" {
			exp, err := termExpression(l, i)
			if err != nil {
				return nil, err
			}
			if l[*i].Value != ")" {
				return nil, errors.Join(ErrInvalidExpression, fmt.Errorf(") excepted, not %s", l[*i].Value))
			}
			*i++
			return exp, nil
		}
	case Operator:
		exp, err := expExpression(l, i)
		if err != nil {
			return nil, err
		}
		return &unaryOperation{operator(c.Value), exp}, nil
	}
	return nil, errors.Join(ErrUnknownExpression, fmt.Errorf(
		"unknown type %s('%s'): excepting a valid literal expression",
		c.Type,
		c.Value,
	))
}

func predefinedExpression(l []*lexer, i *int, id string) (expression, error) {
	if isPredefinedVariable(id) {
		return &predefinedVariable{id}, nil
	}
	if isPredefinedFunction(id) {
		exp, err := operatorExpression(l, i)
		if err != nil {
			return nil, err
		}
		return &predefinedFunction{id, exp}, nil
	}
	return nil, ErrUnknownVariable
}

func operatorExpression(l []*lexer, i *int) (expression, error) {
	c := l[*i]
	if c.Type == Separator && c.Value == "(" {
		*i++
		exp, err := termExpression(l, i)
		if err != nil {
			return nil, err
		}
		if l[*i].Value != ")" {
			return exp, errors.Join(ErrInvalidExpression, fmt.Errorf(") excepted, not %s", l[*i].Value))
		}
		*i++
		return exp, nil
	}
	return nil, ErrInvalidExpression
}
