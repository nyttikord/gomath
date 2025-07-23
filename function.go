package gomath

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidFunction     = errors.New("invalid function")
	ErrInvalidFunctionCall = errors.New("invalid function call")
)

type Function func(map[string]string) (Result, error)

func NewFunction(s string) (Function, error) {
	splits := strings.Split(s, "->")
	if len(splits) != 2 {
		return nil, errors.Join(ErrInvalidFunction, errors.New("a function is defined by 'args -> expression'"))
	}
	before := splits[0]
	expression := splits[1]
	var params []string
	for _, p := range strings.Split(before, ",") {
		params = append(params, strings.TrimSpace(p))
	}
	return func(args map[string]string) (Result, error) {
		if len(params) != len(args) {
			return nil, errors.Join(ErrInvalidFunctionCall, errors.New("not all parameters have been defined"))
		}
		cp := expression
		for _, p := range params {
			v, ok := args[p]
			if !ok {
				return nil, errors.Join(ErrInvalidFunctionCall, fmt.Errorf("missing argument for %s", p))
			}
			cp = strings.ReplaceAll(cp, p, v)
		}
		return Parse(cp)
	}, nil
}
