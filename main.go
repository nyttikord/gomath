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

// ParseAndCalculate an expression with given Options
func ParseAndCalculate(expression string, opt *Options) (string, error) {
	return parseAndEval(expression, opt, astTypeCalculation)
}

// ParseAndConvertToLatex an expression with given Options
func ParseAndConvertToLatex(expression string, opt *Options) (string, error) {
	return parseAndEval(expression, opt, astTypeLatex)
}

func parseAndEval(expression string, opt *Options, tpe astType) (string, error) {
	lexed, err := lex(expression)
	if err != nil {
		return "", err
	}
	p, err := astParse(lexed, tpe)
	if err != nil {
		return "", err
	}
	return p.Body.Eval(opt)
}
