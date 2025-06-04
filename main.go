package gomath

import (
	"errors"
)

var (
	ErrInvalidResult = errors.New("invalid result")
)

// Result represents the result got after the Parse function.
// You can directly get the exact result with String or with fmt.Sprintf("%s", result)
type Result interface {
	// String returns the string representation of the Result.
	// It is the exact result (fraction form)
	String() string
	// Approx returns an approximation of the Result given by String()
	Approx(int) (string, error)
	// LaTeX returns the LaTeX representation of the expression leading to the Result
	LaTeX() (string, error)
}

type res struct {
	*ast
	result *statementResult
}

func (r *res) String() string {
	return r.result.String()
}

func (r *res) Approx(precision int) (string, error) {
	f := r.result.Fraction()
	if f == nil {
		return "", ErrInvalidResult
	}
	return f.Approx(precision), nil
}

func (r *res) LaTeX() (string, error) {
	err := r.ChangeType(astTypeLatex)
	if err != nil {
		return "", err
	}
	result, err := r.Body.Eval(&Options{})
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

// Parse the given expression and return the Result obtained
func Parse(expression string) (Result, error) {
	tree, err := parseAst(expression, astTypeCalculation)
	if err != nil {
		return nil, err
	}
	r, err := tree.Body.Eval(&Options{Decimal: false})
	if err != nil {
		return nil, err
	}
	return &res{ast: tree, result: r}, nil
}

// ParseAndCalculate an expression with given Options
func ParseAndCalculate(expression string, opt *Options) (string, error) {
	tree, err := parseAst(expression, astTypeCalculation)
	if err != nil {
		return "", err
	}
	result, err := tree.Body.Eval(opt)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

// ParseAndConvertToLaTeX an expression with given Options
func ParseAndConvertToLaTeX(expression string, opt *Options) (string, error) {
	tree, err := parseAst(expression, astTypeLatex)
	if err != nil {
		return "", err
	}
	result, err := tree.Body.Eval(opt)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

func parseAst(expression string, tpe astType) (*ast, error) {
	lexed, err := lex(expression)
	if err != nil {
		return nil, err
	}
	return astParse(lexed, tpe)
}
