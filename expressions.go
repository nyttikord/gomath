package gomath

import (
	"errors"
	"fmt"
)

var (
	ErrUnknownOperation = errors.New("unknown operation")
	ErrLexerNotValid    = errors.New("lexer is not valid")
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

type evaluateOperation struct {
	FunctionName string
	Expression   expression
}

type literalExp struct {
	Value *fraction
}

type variable struct {
	ID string
}

type predefinedVariable variable

type relation string

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
		return lf.Add(lr)
	case "-":
		return lf.Sub(lr)
	case "*":
		return lf.Mul(lr)
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
		return lb.Mul(intToFraction(-1))
	default:
		return nil, errors.Join(ErrUnknownOperation, errors.New("operation "+string(b.Operator)+" is not supported"))
	}
}

func (e *evaluateOperation) Eval() (*fraction, error) {
	f, ok := functions[e.FunctionName]
	if !ok {
		return nil, errors.Join(ErrUnknownFunction, fmt.Errorf("undefined function %s", e.FunctionName))
	}
	return f.Relation.Eval(f.Definition, f.Variable, e.Expression)
}

func (l *literalExp) Eval() (*fraction, error) {
	return l.Value, nil
}

func (v *variable) Eval() (*fraction, error) {
	val, ok := variables[v.ID]
	if !ok {
		return nil, errors.Join(ErrUnknownVariable, fmt.Errorf("undefined variable %s", v.ID))
	}
	return val, nil
}

func (v *predefinedVariable) Eval() (*fraction, error) {
	val, ok := predefinedVariables[v.ID]
	if !ok {
		return nil, errors.Join(ErrUnknownVariable, fmt.Errorf("undefined variable \\%s", v.ID))
	}
	return val, nil
}

func lexToRel(lexers []*lexer) *relation {
	var s relation
	for _, l := range lexers {
		s += relation(l.Value)
	}
	return &s
}

func (r *relation) String() string {
	return string(*r)
}

func (r *relation) Eval(def Space, variable string, val expression) (*fraction, error) {
	lexed, err := lex(r.String())
	if err != nil {
		return nil, err
	}
	if len(lexed) != 1 {
		return nil, ErrLexerNotValid
	}
	var lexr []*lexer
	for _, l := range lexed {
		// replace all x by their value in brackets
		if l.Type == Literal && l.Value == variable {
			fr, err := val.Eval()
			if err != nil {
				return nil, err
			}
			if !def.Contains(fr) {
				return nil, errors.Join(ErrNumberNotInSpace, fmt.Errorf("%s is not in %s", fr.String(), def.String()))
			}
			lexr = append(lexr, &lexer{Type: Separator, Value: "("})
			l.Type = Number
			l.Value = fr.String()
			lexr = append(lexr, l)
			lexr = append(lexr, &lexer{Type: Separator, Value: ")"})
		} else {
			lexr = append(lexr, l)
		}
	}
	i := 0
	exp, err := termExpression(lexed, &i)
	if err != nil {
		return nil, err
	}
	return exp.Eval()
}
