package gomath

import "fmt"

type Options struct {
	Decimal bool
}

type statement interface {
	Eval(*Options) (string, error)
}

type printStatement struct {
	Expression expression
}

func (p *printStatement) Eval(opt *Options) (string, error) {
	f, err := p.Expression.Eval()
	if err != nil {
		return "", err
	}
	if opt.Decimal {
		if f.IsInt() {
			i, _ := f.Int()
			return fmt.Sprintf("%d\n", i), nil
		}
		if f.Denominator%10 != 0 {
			return fmt.Sprintf("%f\n", f.Float()), nil
		}
		var i1, i2 int64
		i2 = f.Numerator % f.Denominator
		i1 = (f.Numerator - i2) / f.Denominator
		return fmt.Sprintf("%d.%d\n", i1, i2), nil
	}
	return f.String(), nil
}
