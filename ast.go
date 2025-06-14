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

func (a *ast) ChangeType(tpe astType) error {
	a.Type = tpe
	return a.setStatement(a.Body.getExpr())
}

func (a *ast) String() string {
	m, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return ""
	}
	return string(m)
}

func (a *ast) setStatement(expr expression) error {
	switch a.Type {
	case astTypeCalculation:
		a.Body = &calculationStatement{Expression: expr}
	case astTypeLatex:
		a.Body = &latexStatement{Expression: expr}
	default:
		return ErrUnknownAstType
	}
	return nil
}

// astParse the given lexer and returns an ast
func astParse(lexed []*lexer, tpe astType) (*ast, error) {
	tree := &ast{Type: tpe}
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

func termExpression(l []*lexer, i *int) (expression, error) {
	return binExpression(termOperators, omitParenthesisExpression, l, i)
}

func omitParenthesisExpression(l []*lexer, i *int) (expression, error) {
	return omitExpression(factorExpression, func(l *lexer) bool {
		return l.Type == Separator && l.Value == "("
	}, l, i)
}

func factorExpression(l []*lexer, i *int) (expression, error) {
	return binExpression(factorOperators, omitLiteralExpression, l, i)
}

func omitLiteralExpression(l []*lexer, i *int) (expression, error) {
	return omitExpression(expExpression, func(l *lexer) bool {
		return l.Type == Literal
	}, l, i)
}

func expExpression(l []*lexer, i *int) (expression, error) {
	res, err := binExpression(expOperators, literalExpression, l, i)
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

func binExpression(ops []operator, sub expressionFunc, l []*lexer, i *int) (expression, error) {
	left, err := sub(l, i)
	if err != nil {
		return nil, err
	}
	for *i < len(l) && slices.Contains(ops, operator(l[*i].Value)) {
		op := operator(l[*i].Value)
		*i++
		if *i >= len(l) {
			return nil, ErrInvalidExpression
		}
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

func omitExpression(sub expressionFunc, cond func(*lexer) bool, l []*lexer, i *int) (expression, error) {
	left, err := sub(l, i)
	if err != nil {
		return nil, err
	}
	for *i < len(l) && cond(l[*i]) {
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
			if *i >= len(l) {
				return nil, errors.Join(ErrInvalidExpression, fmt.Errorf("')' excepted"))
			} else if l[*i].Value != ")" {
				return nil, errors.Join(ErrInvalidExpression, fmt.Errorf("')' excepted, not '%s'", l[*i].Value))
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
	return nil, errors.Join(ErrUnknownExpression, fmt.Errorf("unknown type %s: excepting a valid literal expression", c))
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
	return nil, genErrUnknownVariable(id)
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
