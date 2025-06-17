package gomath

import (
	"errors"
	ast2 "github.com/nyttikord/gomath/ast"
	"github.com/nyttikord/gomath/lexer"
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
	Approx(int) string
	// LaTeX returns the LaTeX representation of the expression leading to the Result
	LaTeX() (string, error)
	// IsExact returns true if the fraction can be exactly represented by a string
	IsExact(int) bool
}

type res struct {
	ast    *ast2.Ast
	result *ast2.statementResult
}

func (r *res) String() string {
	return r.result.String()
}

func (r *res) Approx(precision int) string {
	f := r.result.Fraction()
	if f == nil {
		panic(ErrInvalidResult)
	}
	return f.Approx(precision)
}

func (r *res) IsExact(precision int) bool {
	f := r.result.Fraction()
	if f == nil {
		panic(ErrInvalidResult)
	}

	return f.CanBeRepresentedExactly(precision)
}

func (r *res) LaTeX() (string, error) {
	err := r.ast.ChangeType(ast2.astTypeLatex)
	if err != nil {
		return "", err
	}
	result, err := r.ast.Body.Eval(&ast2.Options{})
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

// Parse the given expression and return the Result obtained
func Parse(expression string) (Result, error) {
	tree, err := parseAst(expression, ast2.astTypeCalculation)
	if err != nil {
		return nil, err
	}
	r, err := tree.Body.Eval(&ast2.Options{Decimal: false})
	if err != nil {
		return nil, err
	}
	return &res{ast: tree, result: r}, nil
}

// ParseAndCalculate an expression with given Options
func ParseAndCalculate(expression string, opt *ast2.Options) (string, error) {
	tree, err := parseAst(expression, ast2.astTypeCalculation)
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
func ParseAndConvertToLaTeX(expression string, opt *ast2.Options) (string, error) {
	tree, err := parseAst(expression, ast2.astTypeLatex)
	if err != nil {
		return "", err
	}
	result, err := tree.Body.Eval(opt)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

func parseAst(expression string, tpe ast2.astType) (*ast2.Ast, error) {
	lexed, err := lexer.Lex(expression)
	if err != nil {
		return nil, err
	}
	return ast2.Parse(lexed, tpe)
}
