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
	predefinedVariables = map[string]*savedVariable{}
	predefinedFunctions = map[string]*mathFunction{}
)

type savedVariable struct {
	Val       *fraction
	OmitSlash bool
}

func init() {
	addVar := func(n string, v float64, omitSlash bool) {
		f, err := floatToFraction(v)
		if err != nil {
			panic(err)
		}
		predefinedVariables[n] = &savedVariable{f, omitSlash}
	}
	predefinedVariables["pi"] = &savedVariable{pi, false}
	addVar("e", math.E, true)
	addVar("phi", math.Phi, false)

	addFunc := func(n string, f *mathFunction) {
		predefinedFunctions[n] = f
	}
	createMathFunction := func(def space, mathFunc func(float64) float64) *mathFunction {
		var rel relation
		rel = func(f *fraction) *fraction {
			x, _ := f.Float()
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

	piOverTwo, err := pi.Div(intToFraction(2))
	if err != nil {
		panic(err)
	}
	tanDef := &periodicInterval{
		Interval: &realInterval{
			LowerBound: &intervalBound{
				Value:        piOverTwo.Neg(),
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
	addFunc("ln", createMathFunction(spaceRStar, math.Log))
	addFunc("log2", createMathFunction(spaceRStarPositive, math.Log2))

	log10 := createMathFunction(spaceRStarPositive, math.Log10)
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
