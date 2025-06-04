package gomath

import (
	"errors"
	"math/big"
	"strconv"
	"strings"
)

var (
	ErrInvalidResult = errors.New("invalid result")
)

type Result interface {
	String() string
	Approx(int) (string, error)
	LaTeX() (string, error)
}

type res struct {
	*ast
	result string
}

func (r *res) String() string {
	return r.result
}

func (r *res) Approx(precision int) (string, error) {
	splits := strings.Split(r.result, "/")
	if len(splits) == 1 {
		return r.result, nil
	} else if len(splits) != 2 {
		return "", ErrInvalidResult
	}
	num, err := strconv.Atoi(splits[0])
	if err != nil {
		return "", err
	}
	den, err := strconv.Atoi(splits[1])
	if err != nil {
		return "", err
	}
	return big.NewRat(int64(num), int64(den)).FloatString(precision), nil
}

func (r *res) LaTeX() (string, error) {
	err := r.ChangeType(astTypeLatex)
	if err != nil {
		return "", err
	}
	return r.Body.Eval(&Options{})
}

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
	return tree.Body.Eval(opt)
}

// ParseAndConvertToLaTeX an expression with given Options
func ParseAndConvertToLaTeX(expression string, opt *Options) (string, error) {
	tree, err := parseAst(expression, astTypeLatex)
	if err != nil {
		return "", err
	}
	return tree.Body.Eval(opt)
}

func parseAst(expression string, tpe astType) (*ast, error) {
	lexed, err := lex(expression)
	if err != nil {
		return nil, err
	}
	return astParse(lexed, tpe)
}
