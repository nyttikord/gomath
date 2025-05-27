package gomath

type Options struct {
	Decimal   bool
	Precision int
}

type statement interface {
	// Eval the statement
	Eval(*Options) (string, error)
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

type latexStatement struct {
	Expression expression
}

func (l *latexStatement) Eval(opt *Options) (string, error) {
	s, err := l.Expression.RenderLatex()
	if err != nil {
		return "", err
	}
	return s, nil
}
