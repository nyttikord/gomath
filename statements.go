package gomath

type Options struct {
	Decimal   bool
	Precision int
}
type statementResult struct {
	fraction *fraction
	result   string
}

// String gives the natural result of the statement.
func (c *statementResult) String() string {
	return c.result
}

// Fraction gives the computed fraction during the evaluation.
// Is nil if no fraction was computed
func (c *statementResult) Fraction() *fraction {
	return c.fraction
}

type statement interface {
	// Eval the statement
	Eval(*Options) (*statementResult, error)
	// getExpr returns the expression of the statement
	getExpr() expression
}

type calculationStatement struct {
	Expression expression
}

func (p *calculationStatement) Eval(opt *Options) (*statementResult, error) {
	f, err := p.Expression.Eval()
	if err != nil {
		return nil, err
	}
	r := &statementResult{}
	r.fraction = f
	if opt.Decimal {
		r.result = f.Approx(opt.Precision)
		return r, nil
	}
	r.result = f.String()
	return r, nil
}

func (p *calculationStatement) getExpr() expression {
	return p.Expression
}

type latexStatement struct {
	Expression expression
}

func (l *latexStatement) Eval(opt *Options) (*statementResult, error) {
	s, _, err := l.Expression.RenderLatex()
	if err != nil {
		return nil, err
	}
	r := &statementResult{}
	r.result = s
	r.fraction = nil
	return r, nil
}

func (l *latexStatement) getExpr() expression {
	return l.Expression
}
