package gomath

type Options struct {
	Decimal   bool
	Precision int
}

type statement interface {
	// Eval the statement
	Eval(*Options) (string, error)
	// getExpr returns the expression of the statement
	getExpr() expression
}

type calculationStatement struct {
	Expression expression
}

func (p *calculationStatement) Eval(opt *Options) (string, error) {
	f, err := p.Expression.Eval()
	if err != nil {
		return "", err
	}
	if opt.Decimal {
		return f.Approx(opt.Precision), nil
	}
	return f.String(), nil
}

func (p *calculationStatement) getExpr() expression {
	return p.Expression
}

type latexStatement struct {
	Expression expression
}

func (l *latexStatement) Eval(opt *Options) (string, error) {
	s, _, err := l.Expression.RenderLatex()
	if err != nil {
		return "", err
	}
	return s, nil
}

func (l *latexStatement) getExpr() expression {
	return l.Expression
}
