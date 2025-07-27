package ast

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nyttikord/gomath/expression"
	"github.com/nyttikord/gomath/lexer"
	"github.com/nyttikord/gomath/math"
	"slices"
	"strconv"
)

var (
	termOperators   = []string{"+", "-"}
	factorOperators = []string{"*", "/"}
	expOperators    = []string{"^"}

	// ErrUnknownExpression is thrown when GoMath does not know the expression
	ErrUnknownExpression = errors.New("unknown expression")
	// ErrInvalidExpression is thrown when the given expression's syntax is invalid
	ErrInvalidExpression = errors.New("invalid expression")
	// ErrUnknownAstType is thrown when GoMath does not know the given Type
	ErrUnknownAstType = errors.New("unknown Ast type")
)

type Type uint

const (
	TypeCalculation Type = 0
	TypeLatex       Type = 1
)

type Ast struct {
	Type Type
	Body statement
}

type expressionFunc func(l []*lexer.Lexer, i *int) (expression.Expression, error)

func (a *Ast) ChangeType(tpe Type) error {
	a.Type = tpe
	return a.setStatement(a.Body.getExpr())
}

func (a *Ast) String() string {
	m, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return ""
	}
	return string(m)
}

func (a *Ast) setStatement(expr expression.Expression) error {
	switch a.Type {
	case TypeCalculation:
		a.Body = &calculationStatement{Expression: expr}
	case TypeLatex:
		a.Body = &latexStatement{Expression: expr}
	default:
		return ErrUnknownAstType
	}
	return nil
}

// Parse the given lexer and returns an Ast
func Parse(lexed []*lexer.Lexer, tpe Type) (*Ast, error) {
	tree := &Ast{Type: tpe}
	i := 0
	exp, err := termExpression(lexed, &i)
	if err != nil {
		return nil, err
	}
	if i < len(lexed) {
		return nil, errors.Join(ErrUnknownExpression, fmt.Errorf("cannot parse expression %s", lexed[i]))
	}
	return tree, tree.setStatement(exp) // works because tree is a pointer
}

func termExpression(l []*lexer.Lexer, i *int) (expression.Expression, error) {
	return binExpression(termOperators, omitParenthesisExpression, l, i)
}

func omitParenthesisExpression(l []*lexer.Lexer, i *int) (expression.Expression, error) {
	return omitExpression(factorExpression, func(l *lexer.Lexer) bool {
		return l.Type == lexer.Separator && l.Value == "("
	}, l, i)
}

func factorExpression(l []*lexer.Lexer, i *int) (expression.Expression, error) {
	return binExpression(factorOperators, omitLiteralExpression, l, i)
}

func omitLiteralExpression(l []*lexer.Lexer, i *int) (expression.Expression, error) {
	return omitExpression(expExpression, func(l *lexer.Lexer) bool {
		return l.Type == lexer.Literal
	}, l, i)
}

func expExpression(l []*lexer.Lexer, i *int) (expression.Expression, error) {
	res, err := binExpression(expOperators, literalExpression, l, i)
	if err != nil {
		return nil, err
	}
	if *i == len(l) {
		return res, nil
	}
	if l[*i].Type != lexer.Operator || l[*i].Value != "!" {
		return res, nil
	}
	*i++
	return expression.Factorial(res), nil
}

func binExpression(ops []string, sub expressionFunc, l []*lexer.Lexer, i *int) (expression.Expression, error) {
	left, err := sub(l, i)
	if err != nil {
		return nil, err
	}
	for *i < len(l) && slices.Contains(ops, l[*i].Value) {
		op := l[*i].Value
		*i++
		if *i >= len(l) {
			return nil, ErrInvalidExpression
		}
		right, err := sub(l, i)
		if err != nil {
			return nil, err
		}
		switch op {
		case "+":
			left = expression.Add(left, right)
		case "-":
			left = expression.Sub(left, right)
		case "*":
			left = expression.Mul(left, right)
		case "/":
			left = expression.Div(left, right)
		case "^":
			left = expression.Pow(left, right)
		default:
			return nil, errors.Join(expression.ErrUnknownOperation, fmt.Errorf("unknown operator %s", op))
		}
	}
	return left, nil
}

func omitExpression(sub expressionFunc, cond func(*lexer.Lexer) bool, l []*lexer.Lexer, i *int) (expression.Expression, error) {
	left, err := sub(l, i)
	if err != nil {
		return nil, err
	}
	for *i < len(l) && cond(l[*i]) {
		right, err := sub(l, i)
		if err != nil {
			return nil, err
		}
		left = expression.Mul(left, right)
	}
	return left, nil
}

func literalExpression(l []*lexer.Lexer, i *int) (expression.Expression, error) {
	c := l[*i]
	*i++
	switch c.Type {
	case lexer.Number:
		v, err := strconv.ParseFloat(c.Value, 64)
		if err != nil {
			return nil, err
		}
		f, err := math.FloatToFraction(v)
		if err != nil {
			return nil, err
		}
		return expression.Const(f), nil
	case lexer.Literal:
		if expression.IsPredefinedFunction(c.Value) {
			return predefinedFunction(l, i, c.Value)
		}
		return expression.LiteralExpression(c.Value)
	case lexer.Separator:
		if c.Value == "(" {
			exp, err := termExpression(l, i)
			if err != nil {
				return nil, err
			}
			if *i >= len(l) {
				return nil, errors.Join(ErrInvalidExpression, fmt.Errorf("')' excepted"))
			} else if l[*i].Value != ")" {
				return nil, errors.Join(ErrInvalidExpression, fmt.Errorf("')' excepted, not '%s'", l[*i].Value))
			}
			*i++
			return exp, nil
		}
	case lexer.Operator:
		exp, err := expExpression(l, i)
		if err != nil {
			return nil, err
		}
		switch c.Value {
		case "-":
			exp = expression.Neg(exp)
		case "+":
		default:
			return nil, errors.Join(expression.ErrUnknownOperation, fmt.Errorf("unknown unary operator %s", c.Value))
		}
		return exp, nil
	}
	return nil, errors.Join(ErrUnknownExpression, fmt.Errorf("unknown type %s: excepting a valid literal expression", c))
}

func predefinedFunction(l []*lexer.Lexer, i *int, id string) (expression.Expression, error) {
	exp, err := operatorExpression(l, i)
	if err != nil {
		return nil, err
	}
	return expression.LiteralFunction(id, exp), nil
}

func operatorExpression(l []*lexer.Lexer, i *int) (expression.Expression, error) {
	c := l[*i]
	if c.Type == lexer.Separator && c.Value == "(" {
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
