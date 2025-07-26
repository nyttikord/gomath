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

type expressionFunc func(*lexer.TokenList) (expression.Expression, error)

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
func Parse(tkl *lexer.TokenList, tpe Type) (*Ast, error) {
	tree := &Ast{Type: tpe}
	if !tkl.Next() {
		return nil, ErrInvalidExpression
	}
	exp, err := termExpression(tkl)
	if err != nil {
		return nil, err
	}
	if !tkl.Empty() {
		return nil, errors.Join(ErrUnknownExpression, fmt.Errorf("cannot parse expression %s", tkl.Current()))
	}
	return tree, tree.setStatement(exp) // works because tree is a pointer
}

func termExpression(tkl *lexer.TokenList) (expression.Expression, error) {
	return binExpression(termOperators, omitParenthesisExpression, tkl)
}

func omitParenthesisExpression(tkl *lexer.TokenList) (expression.Expression, error) {
	return omitExpression(factorExpression, func(l *lexer.Lexer) bool {
		return l.Type == lexer.Separator && l.Value == "("
	}, tkl)
}

func factorExpression(tkl *lexer.TokenList) (expression.Expression, error) {
	return binExpression(factorOperators, omitLiteralExpression, tkl)
}

func omitLiteralExpression(tkl *lexer.TokenList) (expression.Expression, error) {
	return omitExpression(expExpression, func(l *lexer.Lexer) bool {
		return l.Type == lexer.Literal
	}, tkl)
}

func expExpression(tkl *lexer.TokenList) (expression.Expression, error) {
	res, err := binExpression(expOperators, literalExpression, tkl)
	if err != nil {
		return nil, err
	}
	if tkl.Empty() {
		return res, nil
	}
	if tkl.Current().Type != lexer.Operator || tkl.Current().Value != "!" {
		return res, nil
	}
	tkl.Next()
	return expression.Factorial(res), nil
}

func binExpression(ops []string, sub expressionFunc, tkl *lexer.TokenList) (expression.Expression, error) {
	left, err := sub(tkl)
	if err != nil {
		return nil, err
	}
	for !tkl.Empty() && slices.Contains(ops, tkl.Current().Value) {
		op := tkl.Current().Value
		if !tkl.Next() {
			return nil, ErrInvalidExpression
		}
		right, err := sub(tkl)
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
			return nil, expression.ErrUnknownOperation
		}
	}
	return left, nil
}

func omitExpression(sub expressionFunc, cond func(*lexer.Lexer) bool, tkl *lexer.TokenList) (expression.Expression, error) {
	left, err := sub(tkl)
	if err != nil {
		return nil, err
	}
	for !tkl.Empty() && cond(tkl.Current()) {
		right, err := sub(tkl)
		if err != nil {
			return nil, err
		}
		left = expression.Mul(left, right)
	}
	return left, nil
}

func literalExpression(tkl *lexer.TokenList) (expression.Expression, error) {
	c := tkl.Current()
	tkl.Next()
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
		return predefinedExpression(tkl, c.Value)
	case lexer.Separator:
		if c.Value == "(" {
			exp, err := termExpression(tkl)
			if err != nil {
				return nil, err
			}
			if tkl.Empty() {
				return nil, errors.Join(ErrInvalidExpression, fmt.Errorf("')' excepted"))
			} else if tkl.Current().Value != ")" {
				return nil, errors.Join(ErrInvalidExpression, fmt.Errorf("')' excepted, not '%s'", tkl.Current().Value))
			}
			tkl.Next()
			return exp, nil
		}
	case lexer.Operator:
		exp, err := expExpression(tkl)
		if err != nil {
			return nil, err
		}
		switch c.Value {
		case "-":
			exp = expression.Neg(exp)
		case "+":
		default:
			return nil, expression.ErrUnknownOperation
		}
		return exp, nil
	}
	return nil, errors.Join(ErrUnknownExpression, fmt.Errorf("unknown type %s: excepting a valid literal expression", c))
}

func predefinedExpression(tkl *lexer.TokenList, id string) (expression.Expression, error) {
	if expression.IsPredefinedVariable(id) {
		return expression.LiteralVariable(id), nil
	}
	if expression.IsPredefinedFunction(id) {
		exp, err := operatorExpression(tkl)
		if err != nil {
			return nil, err
		}
		return expression.LiteralFunction(id, exp), nil
	}
	return nil, expression.GenErrUnknownVariable(id)
}

func operatorExpression(tkl *lexer.TokenList) (expression.Expression, error) {
	c := tkl.Current()
	if c.Type == lexer.Separator && c.Value == "(" {
		if !tkl.Next() {
			return nil, ErrInvalidExpression
		}
		exp, err := termExpression(tkl)
		if err != nil {
			return nil, err
		}
		if tkl.Current().Value != ")" {
			return exp, errors.Join(ErrInvalidExpression, fmt.Errorf(") excepted, not %s", tkl.Current().Value))
		}
		tkl.Next()
		return exp, nil
	}
	return nil, ErrInvalidExpression
}
