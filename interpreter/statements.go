package interpreter

import "fmt"

type Options struct {
	Decimal bool
}

type Statement interface {
	Eval(*Options) error
}

type PrintStatement struct {
	Expression Expression
}

func (p *PrintStatement) Eval(opt *Options) error {
	f, err := p.Expression.Eval()
	if err != nil {
		return err
	}
	if opt.Decimal {
		if f.IsInt() {
			i, _ := f.Int()
			fmt.Printf("%d\n", i)
			return nil
		}
		var i1, i2 int64
		i2 = f.Numerator % f.Denominator
		i1 = (f.Numerator - i2) / f.Denominator
		fmt.Printf("%d.%d\n", i1, i2)
		return nil
	}
	println(f.String())
	return nil
}
