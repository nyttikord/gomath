package gomath

import (
	"errors"
	"fmt"
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
}

type operator string
type separator string

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
	ID string
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
		return lf.Pow(lr)
	default:
		return nil, errors.Join(ErrUnknownOperation, errors.New("operation "+string(b.Operator)+" is not supported"))
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
	default:
		return nil, errors.Join(ErrUnknownOperation, errors.New("operation "+string(b.Operator)+" is not supported"))
	}
}

func (l *literalExp) Eval() (*fraction, error) {
	return l.Value, nil
}

func (v *predefinedVariable) Eval() (*fraction, error) {
	val, ok := predefinedVariables[v.ID]
	if !ok {
		return nil, errors.Join(ErrUnknownVariable, fmt.Errorf("undefined variable \\%s", v.ID))
	}
	return val, nil
}

func (f *predefinedFunction) Eval() (*fraction, error) {
	fn, ok := predefinedFunctions[f.ID]
	if !ok {
		return nil, errors.Join(ErrUnknownVariable, fmt.Errorf("undefined variable \\%s", f.ID))
	}
	val, err := f.exp.Eval()
	if err != nil {
		return nil, err
	}
	return fn.Eval(val)
}

func (r *relation) Eval(f *fraction) *fraction {
	return (*r)(f)
}
