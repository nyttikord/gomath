package gomath

import (
	"fmt"
)

type space interface {
	Contains(f *fraction) bool
	String() string
}

type realSet struct{}
type realInterval struct {
	LowerBound         *fraction
	UpperBound         *fraction
	CustomName         string
	ContainsLowerBound bool
	ContainsUpperBound bool
}
type unionSet struct {
	Sets       []space
	CustomName string
}
type periodicInterval struct {
	Interval   *realInterval
	Period     *fraction
	CustomName string
}

func (*realSet) Contains(*fraction) bool {
	return true
}
func (*realSet) String() string {
	return "R"
}

func (i *realInterval) Contains(f *fraction) bool {
	var b1, b2 bool
	if i.ContainsUpperBound {
		b1 = f.SmallerOrEqualThan(i.UpperBound)
	} else {
		b1 = f.SmallerThan(i.UpperBound)
	}
	if i.ContainsLowerBound {
		b2 = f.SmallerOrEqualThan(i.LowerBound)
	} else {
		b2 = f.SmallerThan(i.LowerBound)
	}
	return b1 && b2
}
func (i *realInterval) String() string {
	if i.CustomName != "" {
		return i.CustomName
	}
	return fmt.Sprintf("[%s, %s]", i.LowerBound.String(), i.UpperBound.String())
}

func (s *unionSet) Contains(f *fraction) bool {
	for _, space := range s.Sets {
		if !space.Contains(f) {
			return false
		}
	}
	return true
}
func (s *unionSet) String() string {
	if s.CustomName != "" {
		return s.CustomName
	}
	st := ""
	for i, space := range s.Sets {
		if i < len(s.Sets)-1 {
			st += space.String() + " ∪ "
		} else {
			st += space.String()
		}
	}
	return st
}

func (set *periodicInterval) Contains(f *fraction) bool {
	if set.Interval.Contains(f) {
		return true
	}

	if f.GreaterThan(set.Interval.UpperBound) {
		for f.GreaterThan(set.Interval.UpperBound) {
			f = f.Sub(set.Period)
		}
		return set.Interval.Contains(f)
	}
	// here f is necessarily smaller than the lower bound
	for f.SmallerThan(set.Interval.LowerBound) {
		f = f.Add(set.Period)
	}
	return set.Interval.Contains(f)
}
func (set *periodicInterval) String() string {
	if set.CustomName != "" {
		return set.CustomName
	}
	return set.Interval.String() + " mod " + set.Period.String()
}
