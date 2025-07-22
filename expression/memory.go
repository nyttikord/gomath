package expression

import (
	"errors"
	"fmt"
	m "github.com/nyttikord/gomath/math"
	"math"
)

var (
	// ErrUnknownVariable is thrown when GoMath doesn't know the variable used
	ErrUnknownVariable = errors.New("unknown variable")
)

var (
	predefinedVariables = map[string]*savedVariable{}
	predefinedFunctions = map[string]*mathFunction{}
)

type savedVariable struct {
	Val       *m.Fraction
	OmitSlash bool
}

type relation func(*m.Fraction) *m.Fraction

func init() {
	addVar := func(n string, v float64, omitSlash bool) {
		f, err := m.FloatToFraction(v)
		if err != nil {
			panic(err)
		}
		predefinedVariables[n] = &savedVariable{f, omitSlash}
	}
	predefinedVariables["pi"] = &savedVariable{m.Pi, false}
	addVar("e", math.E, true)
	addVar("phi", math.Phi, false)

	addFunc := func(n string, f *mathFunction) {
		predefinedFunctions[n] = f
	}
	createMathFunction := func(def m.Space, mathFunc func(float64) float64) *mathFunction {
		var rel relation
		rel = func(f *m.Fraction) *m.Fraction {
			x, _ := f.Float()
			result, err := m.FloatToFraction(mathFunc(x))
			if err != nil {
				panic(err)
			}
			return result
		}

		return &mathFunction{
			Definition: def,
			Relation:   rel,
		}
	}

	addFunc("exp", createMathFunction(&m.RealSet{}, math.Exp))
	addFunc("sqrt", createMathFunction(&m.RealSet{}, math.Sqrt))
	addFunc("sin", createMathFunction(&m.RealSet{}, math.Sin))
	addFunc("cos", createMathFunction(&m.RealSet{}, math.Cos))

	piOverTwo, err := m.Pi.Div(m.IntToFraction(2))
	if err != nil {
		panic(err)
	}
	tanDef := &m.PeriodicInterval{
		Interval: &m.RealInterval{
			LowerBound: &m.IntervalBound{
				Value:        piOverTwo.Neg(),
				IncludeValue: false,
				Infinite:     false,
			},
			UpperBound: &m.IntervalBound{
				Value:        piOverTwo,
				IncludeValue: false,
				Infinite:     false,
			},
			CustomName: "",
		},
		Period:     m.Pi,
		CustomName: "] -pi/2 ; pi/2 [ mod pi",
	}
	addFunc("tan", createMathFunction(tanDef, math.Tan))
	addFunc("ln", createMathFunction(m.SpaceRStar, math.Log))
	addFunc("log2", createMathFunction(m.SpaceRStarPositive, math.Log2))

	log10 := createMathFunction(m.SpaceRStarPositive, math.Log10)
	addFunc("log", log10)
	addFunc("log10", log10)
}

type mathFunction struct {
	Definition m.Space
	Relation   relation
}

func (mf *mathFunction) Eval(f *m.Fraction) (*m.Fraction, error) {
	if !mf.Definition.Contains(f) {
		return nil, errors.Join(ErrNumberNotInSpace, fmt.Errorf("%s is not in %s", f, mf.Definition))
	}
	return mf.Relation(f), nil
}

func IsPredefinedVariable(id string) bool {
	_, ok := predefinedVariables[id]
	return ok
}

func IsPredefinedFunction(id string) bool {
	_, ok := predefinedFunctions[id]
	return ok
}

func GenErrUnknownVariable(name string) error {
	return errors.Join(ErrUnknownVariable, fmt.Errorf("unknown %s", name))
}
