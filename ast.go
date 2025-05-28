package gomath

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
)

var (
	termOperators   = []operator{"+", "-"}
	factorOperators = []operator{"*", "/"}
	expOperators    = []operator{"^"}

	// ErrUnknownExpression is thrown when GoMath does not know the expression
	ErrUnknownExpression = errors.New("unknown expression")
	// ErrInvalidExpression is thrown when the given expression's syntax is invalid
	ErrInvalidExpression = errors.New("invalid expression")
	// ErrUnknownAstType is thrown when GoMath does not know the given astType
	ErrUnknownAstType = errors.New("unknown ast type")
)

type astType uint

const (
	astTypeCalculation astType = 0
	astTypeLatex       astType = 1
)

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
	switch tpe {
	case astTypeCalculation:
		tree.Body = &calculationStatement{Expression: exp}
	case astTypeLatex:
		tree.Body = &latexStatement{Expression: exp}
	default:
		return nil, ErrUnknownAstType
	}
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
	res, err := binExpression(expOperators, literalExpression)(l, i)
	if err != nil {
		return nil, err
	}
	if *i == len(l) {
		return res, nil
	}
	if l[*i].Type != Operator || l[*i].Value != "!" {
		return res, nil
	}
	*i++
	return &unaryOperation{"!", res}, nil
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
		return predefinedExpression(l, i, c.Value)
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
		"unknown type %s('%s'): excepting a valid literal expression", c.Type, c.Value,
	))
}

func predefinedExpression(l []*lexer, i *int, id string) (expression, error) {
	if isPredefinedVariable(id) {
		v := predefinedVariables[id]
		return &predefinedVariable{id, v.OmitSlash}, nil
	}
	if isPredefinedFunction(id) {
		exp, err := operatorExpression(l, i)
		if err != nil {
			return nil, err
		}
		return &predefinedFunction{id, exp}, nil
	}
	return nil, ErrUnknownVariable(id)
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
