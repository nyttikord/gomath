package gomath

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrUnknownVariable is thrown when GoMath doesn't know the variable used
	ErrUnknownVariable = func(name string) error { return errors.New("unknown variable: " + name) }
)

var (
	predefinedVariables = map[string]*fraction{}
	predefinedFunctions = map[string]*mathFunction{}
)

func init() {
	addVar := func(n string, v float64) {
		f, err := floatToFraction(v)
		if err != nil {
			panic(err)
		}
		predefinedVariables[n] = f
	}
	addVar("pi", math.Pi)
	addVar("e", math.E)
	addVar("phi", math.Phi)

	addFunc := func(n string, f *mathFunction) {
		predefinedFunctions[n] = f
	}
	createMathFunction := func(def space, mathFunc func(float64) float64) *mathFunction {
		var rel relation
		rel = func(f *fraction) *fraction {
			x := f.Float()
			result, err := floatToFraction(mathFunc(x))
			if err != nil {
				panic(err)
			}
			return result
		}

		return &mathFunction{
			Definition: def,
			Relation:   &rel,
		}
	}

	addFunc("exp", createMathFunction(&realSet{}, math.Exp))
	addFunc("sqrt", createMathFunction(&realSet{}, math.Sqrt))
	addFunc("sin", createMathFunction(&realSet{}, math.Sin))
	addFunc("cos", createMathFunction(&realSet{}, math.Cos))

	pi, err := floatToFraction(math.Pi)
	if err != nil {
		panic(err)
	}
	piOverTwo, err := pi.Div(intToFraction(2))
	if err != nil {
		panic(err)
	}
	tanDef := &periodicInterval{
		Interval: &realInterval{
			LowerBound: &intervalBound{
				Value:        piOverTwo.Mul(intToFraction(-1)),
				IncludeValue: false,
				Infinite:     false,
			},
			UpperBound: &intervalBound{
				Value:        piOverTwo,
				IncludeValue: false,
				Infinite:     false,
			},
			CustomName: "",
		},
		Period:     pi,
		CustomName: "] -pi/2 ; pi/2 [ mod pi",
	}
	addFunc("tan", createMathFunction(tanDef, math.Tan))
	addFunc("ln", createMathFunction(&realInterval{
		LowerBound: &intervalBound{
			Value:        NullFraction,
			IncludeValue: false,
			Infinite:     false,
		},
		UpperBound: &intervalBound{
			Infinite: true,
			Positive: true,
		},
		CustomName: `R \ { 0 }`,
	}, math.Log))
	addFunc("log2", createMathFunction(&realInterval{
		LowerBound: &intervalBound{
			Value:        NullFraction,
			IncludeValue: false,
			Infinite:     false,
		},
		UpperBound: &intervalBound{
			Infinite: true,
			Positive: true,
		},
		CustomName: `R \ { 0 }`,
	}, math.Log2))

	log10 := createMathFunction(&realInterval{
		LowerBound: &intervalBound{
			Value:        NullFraction,
			IncludeValue: false,
			Infinite:     false,
		},
		UpperBound: &intervalBound{
			Infinite: true,
			Positive: true,
		},
		CustomName: `R \ { 0 }`,
	}, math.Log10)
	addFunc("log", log10)
	addFunc("log10", log10)
}

type mathFunction struct {
	Definition space
	Relation   *relation
}

func (mf *mathFunction) Eval(f *fraction) (*fraction, error) {
	if !mf.Definition.Contains(f) {
		return nil, errors.Join(ErrNumberNotInSpace, fmt.Errorf("%s is not in %s", f, mf.Definition))
	}
	return mf.Relation.Eval(f), nil
}

func isPredefinedVariable(id string) bool {
	_, ok := predefinedVariables[id]
	return ok
}

func isPredefinedFunction(id string) bool {
	_, ok := predefinedFunctions[id]
	return ok
}
