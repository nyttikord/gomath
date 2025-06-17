package gomath

import (
	"errors"
	"fmt"
	"math/big"
)

type operator interface {
	Eval() (*fraction, error)
	RenderLatex() (string, priority, error)
}

type unaryOperator interface {
	IsSingle() bool
}

type addition struct {
	Left, Right expression
	isSub       bool
}

type negation struct {
	Left     expression
	isSingle bool
}

type multiplication struct {
	Left, Right expression
	isDiv       bool
}

type inversion struct {
	Left     expression
	isSingle bool
}

type exponential struct {
	Left, Right expression
}

type factorial struct {
	Left     expression
	isSingle bool
}

func (a *addition) Eval() (*fraction, error) {
	cf := make(chan *fraction)
	cr := make(chan *fraction)
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
	return fmt.Sprintf("%s%s%s", lf, op, lr), termPriority, nil
}

func (n *negation) Eval() (*fraction, error) {
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
		return s, unaryPriority, nil
	}
	return fmt.Sprintf("%s%s", "-", s), unaryPriority, nil
}

func (n *negation) IsSingle() bool {
	return n.isSingle
}

func (m *multiplication) Eval() (*fraction, error) {
	cf := make(chan *fraction)
	cr := make(chan *fraction)
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
	if m.isDiv {
		return fmt.Sprintf(`\frac{%s}{%s}`, lf, lr), factorPriority, nil
	}
	lf = handleLatexParenthesis(lf, pf, factorPriority)
	lr = handleLatexParenthesis(lr, pr, factorPriority)
	return fmt.Sprintf(`%s \times %s`, lf, lr), factorPriority, nil
}

func (i *inversion) Eval() (*fraction, error) {
	lf, err := i.Left.Eval()
	if err != nil {
		return nil, err
	}
	return lf.Inv()
}

func (i *inversion) RenderLatex() (string, priority, error) {
	s, p, err := i.Left.RenderLatex()
	if err != nil {
		return "", 0, err
	}
	s = handleLatexParenthesis(s, p, unaryPriority)
	if !i.IsSingle() {
		return s, unaryPriority, nil
	}
	return fmt.Sprintf(`\frac{1}{%s}`, s), factorPriority, nil
}

func (i *inversion) IsSingle() bool {
	return i.isSingle
}

func (e *exponential) Eval() (*fraction, error) {
	cf := make(chan *fraction)
	cr := make(chan *fraction)
	getLeftRight(cf, cr, e.Left, e.Right)
	lf := <-cf
	lr := <-cr
	return lf.Exp(lr)
}

func (e *exponential) RenderLatex() (string, priority, error) {
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

func (f *factorial) Eval() (*fraction, error) {
	lf, err := f.Left.Eval()
	if err != nil {
		return nil, err
	}
	var i *big.Int
	if i, err = lf.Int(); err != nil || i.Cmp(nullBigInt) < 0 {
		return nil, errors.Join(ErrNumberNotInSpace, errors.New("factorial is not supported for non positive integer"))
	}
	if !i.IsInt64() {
		return nil, errors.Join(ErrInvalidExpression, fmt.Errorf("number %s is too big", i))
	}
	ii := i.Int64()
	res := ii
	ii--
	for ii > 1 {
		res *= ii
		ii--
	}
	return intToFraction(res), nil
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

func getLeftRight(cl, cr chan<- *fraction, left, right expression) {
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

func getLatexLeftRight(cl, cr chan<- string, cpl, cpr chan<- priority, left, right expression) {
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
