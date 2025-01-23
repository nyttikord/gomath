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
		fmt.Printf("%f\n", f.Float())
		return nil
	}
	println(f.String())
	return nil
}
