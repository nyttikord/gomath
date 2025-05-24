package gomath

import (
	"errors"
	"fmt"
	math2 "github.com/anhgelus/gomath/math"
)

var (
	UnknownOperationErr = errors.New("unknown operation")
	LexerNotValidErr    = errors.New("lexer is not valid")
	NumberNotInSpaceErr = errors.New("number is not in the definition space")
)

type ExpressionFunc func(l []*Lexer, i *int) (Expression, error)

type Expression interface {
	Eval() (*math2.Fraction, error)
}

type BinaryOperation struct {
	Operator    string
	Left, Right Expression
}

type UnaryOperation struct {
	Operator   string
	Expression Expression
}

type EvaluateOperation struct {
	FunctionName string
	Expression   Expression
}

type TLiteral struct {
	Value *math2.Fraction
}

type Variable struct {
	ID string
}

type PredefinedVariable Variable

type Relation string

func (b *BinaryOperation) Eval() (*math2.Fraction, error) {
	lb, err := b.Left.Eval()
	if err != nil {
		return nil, err
	}
	lr, err := b.Right.Eval()
	if err != nil {
		return nil, err
	}
	switch b.Operator {
	case "+":
		return lb.Add(lr), nil
	case "-":
		return lb.Sub(lr), nil
	case "*":
		return lb.Mul(lr), nil
	case "/":
		return lb.Div(lr), nil
	case "^":
		return lb.Pow(lr), nil
	default:
		return nil, errors.Join(UnknownOperationErr, errors.New("operation "+b.Operator+" is not supported"))
	}
}

func (b *UnaryOperation) Eval() (*math2.Fraction, error) {
	lb, err := b.Expression.Eval()
	if err != nil {
		return nil, err
	}
	switch b.Operator {
	case "+":
		return lb, nil
	case "-":
		return lb.Mul(math2.IntToFraction(-1)), nil
	default:
		return nil, errors.Join(UnknownOperationErr, errors.New("operation "+b.Operator+" is not supported"))
	}
}

func (e *EvaluateOperation) Eval() (*math2.Fraction, error) {
	f, ok := functions[e.FunctionName]
	if !ok {
		return nil, errors.Join(UnknownFunctionErr, fmt.Errorf("undefined function %s", e.FunctionName))
	}
	return f.Relation.Eval(f.Definition, f.Variable, e.Expression)
}

func (l *TLiteral) Eval() (*math2.Fraction, error) {
	return l.Value, nil
}

func (v *Variable) Eval() (*math2.Fraction, error) {
	val, ok := variables[v.ID]
	if !ok {
		return nil, errors.Join(UnknownVariableErr, fmt.Errorf("undefined variable %s", v.ID))
	}
	return val, nil
}

func (v *PredefinedVariable) Eval() (*math2.Fraction, error) {
	val, ok := predefinedVariables[v.ID]
	if !ok {
		return nil, errors.Join(UnknownVariableErr, fmt.Errorf("undefined variable \\%s", v.ID))
	}
	return val, nil
}

func LexToRel(lexers []*Lexer) *Relation {
	var s Relation
	for _, l := range lexers {
		s += Relation(l.Value)
	}
	return &s
}

func (r *Relation) String() string {
	return string(*r)
}

func (r *Relation) Eval(def math2.Space, variable string, val Expression) (*math2.Fraction, error) {
	lexed, err := Lex(r.String())
	if err != nil {
		return nil, err
	}
	if len(lexed) != 1 {
		return nil, LexerNotValidErr
	}
	var lex []*Lexer
	for _, l := range lexed {
		// replace all x by their value in brackets
		if l.Type == Literal && l.Value == variable {
			fr, err := val.Eval()
			if err != nil {
				return nil, err
			}
			if !def.Contains(fr) {
				return nil, errors.Join(NumberNotInSpaceErr, fmt.Errorf("%s is not in %s", fr.String(), def.String()))
			}
			lex = append(lex, &Lexer{Type: Separator, Value: "("})
			l.Type = Number
			l.Value = fr.String()
			lex = append(lex, l)
			lex = append(lex, &Lexer{Type: Separator, Value: ")"})
		} else {
			lex = append(lex, l)
		}
	}
	i := 0
	exp, err := termExpression(lexed, &i)
	if err != nil {
		return nil, err
	}
	return exp.Eval()
}
