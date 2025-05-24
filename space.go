package gomath

import (
	"errors"
	"fmt"
)

var (
	UnknownSpaceErr = errors.New("unknown space")
)

type Space interface {
	Contains(f *fraction) bool
	String() string
}

type RealSet struct{}

type RationalSet struct{}

type RelativeSet struct{}

type NaturalSet struct{}

func parseSpace(s string) (Space, error) {
	switch s {
	case "R":
		return RealSet{}, nil
	case "Q":
		return RationalSet{}, nil
	case "Z":
		return RelativeSet{}, nil
	case "N":
		return NaturalSet{}, nil
	}
	return nil, errors.Join(UnknownSpaceErr, fmt.Errorf("unknown space %s", s))
}

func (RealSet) Contains(f *fraction) bool {
	return true
}

func (RealSet) String() string {
	return "R"
}

func (RationalSet) Contains(f *fraction) bool {
	return true
}

func (RationalSet) String() string {
	return "Q"
}

func (NaturalSet) String() string {
	return "N"
}

func (NaturalSet) Contains(f *fraction) bool {
	if !f.IsInt() {
		return false
	}
	i, _ := f.Int()
	return i >= 0
}

func (RelativeSet) String() string {
	return "Z"
}

func (RelativeSet) Contains(f *fraction) bool {
	return f.IsInt()
}
