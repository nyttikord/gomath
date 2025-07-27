package ast

import (
	"github.com/nyttikord/gomath/expression"
	"github.com/nyttikord/gomath/math"
)

type Options struct {
	Decimal   bool
	Precision int
}
type StatementResult struct {
	fraction *math.Fraction
	result   string
}

// String gives the natural result of the statement.
func (c *StatementResult) String() string {
	return c.result
}

// Fraction gives the computed fraction during the evaluation.
// Is nil if no fraction was computed
func (c *StatementResult) Fraction() *math.Fraction {
	return c.fraction
}

type statement interface {
	// Eval the statement
	Eval(*Options) (*StatementResult, error)
	// getExpr returns the expression of the statement
	getExpr() expression.Expression
}

type calculationStatement struct {
	Expression expression.Expression
}

func (p *calculationStatement) Eval(opt *Options) (*StatementResult, error) {
	f, err := p.Expression.Eval()
	if err != nil {
		return nil, err
	}
	r := &StatementResult{}
	r.fraction = f
	if opt.Decimal {
		r.result = f.Approx(opt.Precision)
		return r, nil
	}
	r.result = f.String()
	return r, nil
}

func (p *calculationStatement) getExpr() expression.Expression {
	return p.Expression
}

type latexStatement struct {
	Expression expression.Expression
}

func (l *latexStatement) Eval(_ *Options) (*StatementResult, error) {
	s, _, err := l.Expression.RenderLatex()
	if err != nil {
		return nil, err
	}
	r := &StatementResult{}
	r.result = s
	r.fraction = nil
	return r, nil
}

func (l *latexStatement) getExpr() expression.Expression {
	return l.Expression
}
