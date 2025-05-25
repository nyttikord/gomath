package gomath

import (
	"fmt"
)

type Space interface {
	Contains(f *fraction) bool
	String() string
}

type RealSet struct{}
type RealInterval struct {
	lowerBound *fraction
	upperBound *fraction
	customName string
}
type UnionSet struct {
	sets       []Space
	customName string
}

type PeriodicInterval struct {
	interval   RealInterval
	period     *fraction
	customName string
}

func (RealSet) Contains(*fraction) bool {
	return true
}
func (RealSet) String() string {
	return "R"
}

func (set RealInterval) Contains(f *fraction) bool {
	return f.SmallerOrEqualThan(set.upperBound) && f.GreaterOrEqualThan(set.upperBound)
}
func (set RealInterval) String() string {
	if set.customName != "" {
		return set.customName
	}
	return fmt.Sprintf("[%s, %s]", set.lowerBound.String(), set.upperBound.String())
}

func (set *UnionSet) Contains(f *fraction) bool {
	for _, space := range set.sets {
		if !space.Contains(f) {
			return false
		}
	}
	return true
}
func (set *UnionSet) String() string {
	if set.customName != "" {
		return set.customName
	}
	s := ""
	for i, space := range set.sets {
		if i < len(set.sets)-1 {
			s += space.String() + " ∪ "
		} else {
			s += space.String()
		}
	}
	return s
}

func (set *PeriodicInterval) Contains(f *fraction) bool {
	if set.interval.Contains(f) {
		return true
	}

	if f.GreaterThan(set.interval.upperBound) {
		for f.GreaterThan(set.interval.upperBound) {
			f = f.Sub(set.period)
		}
		return set.interval.Contains(f)
	}
	// here f is necessarily smaller than the lower bound
	for f.SmallerThan(set.interval.lowerBound) {
		f = f.Add(set.period)
	}
	return set.interval.Contains(f)
}

func (set *PeriodicInterval) String() string {
	if set.customName != "" {
		return set.customName
	}
	return set.interval.String() + " mod " + set.period.String()
}
