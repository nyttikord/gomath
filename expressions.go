package gomath

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

var (
	// ErrUnknownOperation is thrown when GoMath doesn't know the operation used
	ErrUnknownOperation = errors.New("unknown operation")
	// ErrNumberNotInSpace is thrown when the number is not in the definition space
	ErrNumberNotInSpace = errors.New("number is not in the definition space")
)

type expressionFunc func(l []*lexer, i *int) (expression, error)

type expression interface {
	// Eval the expression
	Eval() (*fraction, error)
	// RenderLatex the expression
	RenderLatex() (string, priority, error)
}

type operator string
type separator string

type priority uint8

const (
	termPriority    priority = 0
	factorPriority  priority = 1
	expPriority     priority = 2
	unaryPriority   priority = 3
	literalPriority priority = 4
)

type binaryOperation struct {
	Operator    operator
	Left, Right expression
}

type unaryOperation struct {
	Operator   operator
	Expression expression
}

type literalExp struct {
	Value *fraction
}

type variable struct {
	ID        string
	OmitSlash bool
}

type function struct {
	ID  string
	exp expression
}

type predefinedVariable variable
type predefinedFunction function

type relation func(*fraction) *fraction

func (b *binaryOperation) Eval() (*fraction, error) {
	chanLf := make(chan *fraction)
	chanLr := make(chan *fraction)
	go func() {
		lf, err := b.Left.Eval()
		chanLf <- lf
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		lr, err := b.Right.Eval()
		chanLr <- lr
		if err != nil {
			panic(err)
		}
	}()
	lf := <-chanLf
	lr := <-chanLr
	close(chanLf)
	close(chanLr)
	switch b.Operator {
	case "+":
		return lf.Add(lr), nil
	case "-":
		return lf.Sub(lr), nil
	case "*":
		return lf.Mul(lr), nil
	case "/":
		return lf.Div(lr)
	case "^":
		return lf.Exp(lr)
	default:
		return nil, errors.Join(ErrUnknownOperation, errors.New("operation "+string(b.Operator)+" is not supported"))
	}
}

func (b *binaryOperation) RenderLatex() (string, priority, error) {
	chanLf := make(chan string)
	chanLr := make(chan string)
	chanLfp := make(chan priority)
	chanLrp := make(chan priority)
	go func() {
		lf, p, err := b.Left.RenderLatex()
		chanLf <- lf
		chanLfp <- p
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		lr, p, err := b.Right.RenderLatex()
		chanLr <- lr
		chanLrp <- p
		if err != nil {
			panic(err)
		}
	}()
	lf := <-chanLf
	lr := <-chanLr
	lfp := <-chanLfp
	lrp := <-chanLrp
	close(chanLf)
	close(chanLr)
	close(chanLfp)
	close(chanLrp)
	switch b.Operator {
	case "+":
		return fmt.Sprintf("%s + %s", lf, lr), termPriority, nil
	case "-":
		return fmt.Sprintf("%s - %s", lf, lr), termPriority, nil
	case "*":
		if strings.Contains(lf, " ") && lfp < factorPriority {
			lf = `\left(` + lf + `\right)`
		}
		if strings.Contains(lr, " ") && lrp < factorPriority {
			lr = `\left(` + lr + `\right)`
		}
		return fmt.Sprintf(`%s \times %s`, lf, lr), factorPriority, nil
	case "/":
		return fmt.Sprintf(`\frac{%s}{%s}`, lf, lr), factorPriority, nil
	case "^":
		var s string
		if strings.Contains(lf, " ") && lfp < expPriority {
			s += `\left(` + lf + `\right)`
		} else {
			s += lf
		}
		s += "^"
		if len(lr) > 1 {
			s += "{" + lr + "}"
		} else {
			s += lr
		}
		return s, expPriority, nil
	default:
		return "", 0, errors.Join(ErrUnknownOperation, errors.New("operation "+string(b.Operator)+" is not supported"))
	}
}

func (b *unaryOperation) Eval() (*fraction, error) {
	lb, err := b.Expression.Eval()
	if err != nil {
		return nil, err
	}
	switch b.Operator {
	case "+":
		return lb, nil
	case "-":
		return lb.Mul(intToFraction(-1)), nil
	case "!":
		var i *big.Int
		if i, err = lb.Int(); err != nil || i.Cmp(nullBigInt) < 0 {
			return nil, errors.Join(ErrNumberNotInSpace, errors.New("operation "+string(b.Operator)+" is not supported for non positive integer"))
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
	default:
		return nil, errors.Join(ErrUnknownOperation, errors.New("operation "+string(b.Operator)+" is not supported"))
	}
}

func (b *unaryOperation) RenderLatex() (string, priority, error) {
	s, p, err := b.Expression.RenderLatex()
	if err != nil {
		return "", unaryPriority, err
	}
	if strings.Contains(s, " ") && p < unaryPriority {
		s = `\left(` + s + `\left)`
	}
	if b.Operator == "!" {
		return fmt.Sprintf("%s!", s), unaryPriority, nil
	}
	return fmt.Sprintf("%s%s", b.Operator, s), unaryPriority, nil
}

func (l *literalExp) Eval() (*fraction, error) {
	return l.Value, nil
}

func (l *literalExp) RenderLatex() (string, priority, error) {
	return l.Value.String(), literalPriority, nil
}

func (v *predefinedVariable) Eval() (*fraction, error) {
	val, ok := predefinedVariables[v.ID]
	if !ok {
		return nil, errors.Join(genErrUnknownVariable(v.ID), fmt.Errorf("undefined variable %s", v.ID))
	}
	return val.Val, nil
}

func (v *predefinedVariable) RenderLatex() (string, priority, error) {
	_, ok := predefinedVariables[v.ID]
	if !ok {
		return "", literalPriority, errors.Join(genErrUnknownVariable(v.ID), fmt.Errorf("undefined variable %s", v.ID))
	}
	if v.OmitSlash {
		return v.ID, literalPriority, nil
	}
	return `\` + v.ID, literalPriority, nil
}

func (f *predefinedFunction) Eval() (*fraction, error) {
	fn, ok := predefinedFunctions[f.ID]
	if !ok {
		return nil, errors.Join(genErrUnknownVariable(f.ID), fmt.Errorf("undefined variable %s", f.ID))
	}
	val, err := f.exp.Eval()
	if err != nil {
		return nil, err
	}
	return fn.Eval(val)
}

func (f *predefinedFunction) RenderLatex() (string, priority, error) {
	_, ok := predefinedFunctions[f.ID]
	if !ok {
		return "", literalPriority, errors.Join(genErrUnknownVariable(f.ID), fmt.Errorf("undefined variable %s", f.ID))
	}
	val, _, err := f.exp.RenderLatex()
	if err != nil {
		return "", literalPriority, err
	}
	return fmt.Sprintf(`\%s\left(%s\right)`, f.ID, val), literalPriority, nil
}

func (r *relation) Eval(f *fraction) *fraction {
	return (*r)(f)
}
