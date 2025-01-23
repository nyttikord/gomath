package interpreter

import (
	"errors"
	"fmt"
	"github.com/anhgelus/gomath/lexer"
	"github.com/anhgelus/gomath/utils"
	"slices"
	"strconv"
	"strings"
)

var (
	termOperators   = []string{"+", "-"}
	factorOperators = []string{"*", "/"}
	expOperators    = []string{"^"}

	UnknownExpressionErr = errors.New("unknown expression")
	UnknownStatementErr  = errors.New("unknown statement")
	WrongExpressionErr   = errors.New("wrong expression")
)

type Ast struct {
	Type string
	Body []Statement
}

func Parse(lexed [][]*lexer.Lexer) (*Ast, error) {
	tree := Ast{Type: "program"}
	for j, l := range lexed {
		i := 0
		var stmt Statement

		switch l[0].Value {
		case "let":
			i++
			// create new variable
			v, err := NewMemory(l, &i)
			if err != nil {
				return nil, errors.Join(utils.GenErrorLine(j), err)
			}
			// get variable expression
			exp, err := termExpression(l, &i)
			if err != nil {
				return nil, errors.Join(utils.GenErrorLine(j), err)
			}
			v.Expression = exp
			stmt = v
		case "for":
			i++
			// create new function
			f, err := NewFunction(l, &i)
			if err != nil {
				return nil, errors.Join(utils.GenErrorLine(j), err)
			}
			stmt = f
		default:
			exp, err := termExpression(l, &i)
			if err != nil {
				return nil, errors.Join(utils.GenErrorLine(j), err)
			}
			stmt = &PrintStatement{Expression: exp}
		}
		tree.Body = append(tree.Body, stmt)
	}
	return &tree, nil
}

func termExpression(l []*lexer.Lexer, i *int) (Expression, error) {
	return binExpression(termOperators, factorExpression)(l, i)
}

func factorExpression(l []*lexer.Lexer, i *int) (Expression, error) {
	return binExpression(factorOperators, expExpression)(l, i)
}

func expExpression(l []*lexer.Lexer, i *int) (Expression, error) {
	return binExpression(expOperators, literalExpression)(l, i)
}

func binExpression(operators []string, sub ExpressionFunc) ExpressionFunc {
	return func(l []*lexer.Lexer, i *int) (Expression, error) {
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
			left = &BinaryOperation{
				Operator: op,
				Left:     left,
				Right:    right,
			}
		}
		return left, nil
	}
}

func operatorExpression(l []*lexer.Lexer, i *int) (Expression, error) {
	c := l[*i]
	if c.Type == lexer.Operator && c.Value == "{" {
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

func literalExpression(l []*lexer.Lexer, i *int) (Expression, error) {
	c := l[*i]
	*i++
	switch c.Type {
	case lexer.Number:
		v, err := strconv.ParseFloat(c.Value, 64)
		if err != nil {
			return nil, err
		}
		f, err := FloatToFraction(v)
		if err != nil {
			return nil, err
		}
		return &Literal{Value: f}, nil
	case lexer.Literal:
		if strings.HasPrefix(c.Value, "\\") {
			return &PredefinedVariable{ID: c.Value[1:]}, nil
		} else if *i < len(l) && l[*i].Type == lexer.Operator {
			name := c.Value
			*i++
			exp, err := operatorExpression(l, i)
			if err != nil {
				return nil, err
			}
			return &EvaluateOperation{Expression: exp, FunctionName: name}, nil
		}
		return &Variable{ID: c.Value}, nil
	case lexer.Separator:
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
	case lexer.Operator:
		exp, err := expExpression(l, i)
		if err != nil {
			return nil, err
		}
		return &UnaryOperation{Operator: c.Value, Expression: exp}, nil
	}
	return nil, errors.Join(UnknownExpressionErr, fmt.Errorf(
		"unknown type %s('%s'): excepting a valid literal expression",
		c.Type,
		c.Value,
	))
}
