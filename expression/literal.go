package expression

import (
	"errors"
	"fmt"
	"github.com/nyttikord/gomath/math"
)

type Literal interface {
	Eval() (*math.Fraction, error)
	RenderLatex() (string, priority, error)
}

type predefinedVariable variable

type predefinedFunction function

func (v *predefinedVariable) Eval() (*math.Fraction, error) {
	val, ok := predefinedVariables[v.ID]
	if !ok {
		return nil, errors.Join(GenErrUnknownVariable(v.ID), fmt.Errorf("undefined variable %s", v.ID))
	}
	return val.Val, nil
}

func (v *predefinedVariable) RenderLatex() (string, priority, error) {
	_, ok := predefinedVariables[v.ID]
	if !ok {
		return "", literalPriority, errors.Join(GenErrUnknownVariable(v.ID), fmt.Errorf("undefined variable %s", v.ID))
	}
	if v.OmitSlash {
		return v.ID, literalPriority, nil
	}
	return `\` + v.ID, literalPriority, nil
}

func (f *predefinedFunction) Eval() (*math.Fraction, error) {
	fn, ok := predefinedFunctions[f.ID]
	if !ok {
		return nil, errors.Join(GenErrUnknownVariable(f.ID), fmt.Errorf("undefined variable %s", f.ID))
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
		return "", literalPriority, errors.Join(GenErrUnknownVariable(f.ID), fmt.Errorf("undefined variable %s", f.ID))
	}
	val, _, err := f.exp.RenderLatex()
	if err != nil {
		return "", literalPriority, err
	}
	return fmt.Sprintf(`\%s\left(%s\right)`, f.ID, val), literalPriority, nil
}

func LiteralVariable(id string) Literal {
	v := predefinedVariables[id]
	return &predefinedVariable{id, v.OmitSlash}
}

func LiteralFunction(id string, exp Expression) Literal {
	return &predefinedFunction{id, exp}
}
