package expression

import (
	"errors"
	"fmt"
	"github.com/nyttikord/gomath/math"
	"math/big"
)

type Operator interface {
	Eval() (*math.Fraction, error)
	RenderLatex() (string, priority, error)
}

type UnaryOperator interface {
	Eval() (*math.Fraction, error)
	RenderLatex() (string, priority, error)
	IsSingle() bool
}

type addition struct {
	Left, Right Expression
	isSub       bool
}

type negation struct {
	Left     Expression
	isSingle bool
}

type multiplication struct {
	Left, Right Expression
}

type division struct {
	Left, Right Expression
}

type pow struct {
	Left, Right Expression
}

type factorial struct {
	Left     Expression
	isSingle bool
}

func (a *addition) Eval() (*math.Fraction, error) {
	cf := make(chan *math.Fraction)
	cr := make(chan *math.Fraction)
	getLeftRight(cf, cr, a.Left, a.Right)
	lf := <-cf
	lr := <-cr
	return lf.Add(lr), nil
}

func (a *addition) RenderLatex() (string, priority, error) {
	cf := make(chan string)
	cr := make(chan string)
	cpl := make(chan priority)
	cpr := make(chan priority)
	getLatexLeftRight(cf, cr, cpl, cpr, a.Left, a.Right)
	lf := <-cf
	lr := <-cr
	pf := <-cpl
	pr := <-cpr
	lf = handleLatexParenthesis(lf, pf, termPriority)
	lr = handleLatexParenthesis(lr, pr, termPriority)
	op := "+"
	if a.isSub {
		op = "-"
	}
	return fmt.Sprintf("%s %s %s", lf, op, lr), termPriority, nil
}

func (n *negation) Eval() (*math.Fraction, error) {
	lf, err := n.Left.Eval()
	if err != nil {
		return nil, err
	}
	return lf.Neg(), nil
}

func (n *negation) RenderLatex() (string, priority, error) {
	s, p, err := n.Left.RenderLatex()
	if err != nil {
		return "", unaryPriority, err
	}
	s = handleLatexParenthesis(s, p, unaryPriority)
	if !n.IsSingle() {
		return s, literalPriority, nil
	}
	return fmt.Sprintf("%s%s", "-", s), unaryPriority, nil
}

func (n *negation) IsSingle() bool {
	return n.isSingle
}

func (m *multiplication) Eval() (*math.Fraction, error) {
	cf := make(chan *math.Fraction)
	cr := make(chan *math.Fraction)
	getLeftRight(cf, cr, m.Left, m.Right)
	lf := <-cf
	lr := <-cr
	return lf.Mul(lr), nil
}

func (m *multiplication) RenderLatex() (string, priority, error) {
	cf := make(chan string)
	cr := make(chan string)
	cpl := make(chan priority)
	cpr := make(chan priority)
	getLatexLeftRight(cf, cr, cpl, cpr, m.Left, m.Right)
	lf := <-cf
	lr := <-cr
	pf := <-cpl
	pr := <-cpr
	lf = handleLatexParenthesis(lf, pf, factorPriority)
	lr = handleLatexParenthesis(lr, pr, factorPriority)
	return fmt.Sprintf(`%s \times %s`, lf, lr), factorPriority, nil
}

func (m *division) Eval() (*math.Fraction, error) {
	cf := make(chan *math.Fraction)
	cr := make(chan *math.Fraction)
	getLeftRight(cf, cr, m.Left, m.Right)
	lf := <-cf
	lr := <-cr
	return lf.Div(lr)
}

func (m *division) RenderLatex() (string, priority, error) {
	cf := make(chan string)
	cr := make(chan string)
	cpl := make(chan priority)
	cpr := make(chan priority)
	getLatexLeftRight(cf, cr, cpl, cpr, m.Left, m.Right)
	lf := <-cf
	lr := <-cr
	return fmt.Sprintf(`\frac{%s}{%s}`, lf, lr), factorPriority, nil
}

func (e *pow) Eval() (*math.Fraction, error) {
	cf := make(chan *math.Fraction)
	cr := make(chan *math.Fraction)
	getLeftRight(cf, cr, e.Left, e.Right)
	lf := <-cf
	lr := <-cr
	return lf.Exp(lr)
}

func (e *pow) RenderLatex() (string, priority, error) {
	cf := make(chan string)
	cr := make(chan string)
	cpl := make(chan priority)
	cpr := make(chan priority)
	getLatexLeftRight(cf, cr, cpl, cpr, e.Left, e.Right)
	lf := <-cf
	lr := <-cr
	pf := <-cpl
	_ = <-cpr
	s := handleLatexParenthesis(lf, pf, expPriority) + "^"
	if len(lr) > 1 {
		s += "{" + lr + "}"
	} else {
		s += lr
	}
	return s, expPriority, nil
}

func (f *factorial) Eval() (*math.Fraction, error) {
	lf, err := f.Left.Eval()
	if err != nil {
		return nil, err
	}
	var i *big.Int
	if i, err = lf.Int(); err != nil || i.Cmp(math.NullBigInt) < 0 {
		return nil, errors.Join(ErrNumberNotInSpace, errors.New("factorial is not supported for non positive integer"))
	}
	if !i.IsInt64() {
		return nil, errors.Join(ErrNumberNotInSpace, fmt.Errorf("number %s is too big", i))
	}
	ii := i.Int64()
	res := ii
	ii--
	for ii > 1 {
		res *= ii
		ii--
	}
	return math.IntToFraction(res), nil
}

func (f *factorial) RenderLatex() (string, priority, error) {
	s, p, err := f.Left.RenderLatex()
	if err != nil {
		return "", 0, err
	}
	s = handleLatexParenthesis(s, p, unaryPriority)
	return fmt.Sprintf("%s!", s), unaryPriority, nil
}

func (f *factorial) IsSingle() bool {
	return f.isSingle
}

func Neg(l Expression) UnaryOperator {
	return &negation{l, true}
}

func Add(l Expression, r Expression) Operator {
	return &addition{l, r, false}
}

func Sub(l Expression, r Expression) Operator {
	return &addition{l, &negation{r, false}, true}
}

func Mul(l Expression, r Expression) Operator {
	return &multiplication{l, r}
}

func Div(l Expression, r Expression) Operator {
	return &division{l, r}
}

func Factorial(l Expression) UnaryOperator {
	return &factorial{l, true}
}

func Pow(l Expression, r Expression) Operator {
	return &pow{l, r}
}

func getLeftRight(cl, cr chan<- *math.Fraction, left, right Expression) {
	go func() {
		lf, err := left.Eval()
		cl <- lf
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		lr, err := right.Eval()
		cr <- lr
		if err != nil {
			panic(err)
		}
	}()
}

func getLatexLeftRight(cl, cr chan<- string, cpl, cpr chan<- priority, left, right Expression) {
	go func() {
		lf, p, err := left.RenderLatex()
		cl <- lf
		cpl <- p
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		lr, p, err := right.RenderLatex()
		cr <- lr
		cpr <- p
		if err != nil {
			panic(err)
		}
	}()
}
