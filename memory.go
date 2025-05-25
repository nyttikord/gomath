package gomath

import (
	"errors"
	"fmt"
	m "math"
)

var (
	ErrUnknownVariable = errors.New("unknown variable")
)

var (
	predefinedVariables = map[string]*fraction{}
	predefinedFunctions = map[string]*mathFunction{}
)

func init() {
	add := func(n string, v float64) {
		f, err := floatToFraction(v)
		if err != nil {
			panic(err)
		}
		predefinedVariables[n] = f
	}
	add("pi", m.Pi)
	add("e", m.E)
	add("phi", m.Phi)
}

type mathFunction struct {
	Definition Space
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
