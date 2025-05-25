package gomath

<<<<<<< HEAD
type Space interface {
=======
import (
	"fmt"
)

type space interface {
>>>>>>> origin/main
	Contains(f *fraction) bool
	String() string
}

<<<<<<< HEAD
type RealSet struct{}

type RealInterval struct {
	lowerBound IntervalBound
	upperBound IntervalBound
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
=======
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
>>>>>>> origin/main
}

func (*realSet) Contains(*fraction) bool {
	return true
}
func (*realSet) String() string {
	return "R"
}

<<<<<<< HEAD
func (set RealInterval) Contains(f *fraction) bool {
	return f.smallerThanBound(set.upperBound) && f.greaterThanBound(set.lowerBound)
}
func (set RealInterval) String() string {
	if set.customName != "" {
		return set.customName
	}
	s := ""

	if set.lowerBound.infinite {
		s = "] "
		if set.upperBound.positive {
			s += "+inf"
		} else {
			s += "-inf"
		}
	} else {
		if set.lowerBound.includeValue {
			s = "[ "
		} else {
			s = "] "
		}
		s += set.lowerBound.value.String() + " ; "
	}

	if set.upperBound.infinite {
		if set.upperBound.positive {
			s += "+inf"
		} else {
			s += "-inf"
		}
		s += " ["
	} else {
		s += set.upperBound.value.String()
		if set.upperBound.includeValue {
			s += " ]"
		} else {
			s += " ["
		}
	}

	return s
}

func (set UnionSet) Contains(f *fraction) bool {
	for _, space := range set.sets {
=======
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
>>>>>>> origin/main
		if !space.Contains(f) {
			return false
		}
	}
	return true
}
<<<<<<< HEAD
func (set UnionSet) String() string {
	if set.customName != "" {
		return set.customName
=======
func (s *unionSet) String() string {
	if s.CustomName != "" {
		return s.CustomName
>>>>>>> origin/main
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

<<<<<<< HEAD
func (set PeriodicInterval) Contains(f *fraction) bool {
	if set.interval.Contains(f) {
		return true
	}

	if f.strictlyGreaterThanBound(set.interval.upperBound) {
		for f.strictlyGreaterThanBound(set.interval.upperBound) {
			f = f.Sub(set.period)
=======
func (set *periodicInterval) Contains(f *fraction) bool {
	if set.Interval.Contains(f) {
		return true
	}

	if f.GreaterThan(set.Interval.UpperBound) {
		for f.GreaterThan(set.Interval.UpperBound) {
			f = f.Sub(set.Period)
>>>>>>> origin/main
		}
		return set.Interval.Contains(f)
	}
	// here f is necessarily smaller than the lower bound
<<<<<<< HEAD
	for f.strictlySmallerThanBound(set.interval.lowerBound) {
		f = f.Add(set.period)
=======
	for f.SmallerThan(set.Interval.LowerBound) {
		f = f.Add(set.Period)
>>>>>>> origin/main
	}
	return set.Interval.Contains(f)
}
<<<<<<< HEAD
func (set PeriodicInterval) String() string {
	if set.customName != "" {
		return set.customName
=======
func (set *periodicInterval) String() string {
	if set.CustomName != "" {
		return set.CustomName
>>>>>>> origin/main
	}
	return set.Interval.String() + " mod " + set.Period.String()
}

type IntervalBound struct {
	value        *fraction
	includeValue bool
	infinite     bool
	positive     bool
}

func (f fraction) smallerThanBound(b IntervalBound) bool {
	if b.infinite {
		return b.positive
	}
	if b.includeValue {
		return f.SmallerOrEqualThan(b.value)
	}
	return f.SmallerThan(b.value)
}

func (f fraction) greaterThanBound(b IntervalBound) bool {
	if b.infinite {
		return !b.positive
	}
	if b.includeValue {
		return f.GreaterOrEqualThan(b.value)
	}
	return f.GreaterThan(b.value)
}

func (f fraction) strictlySmallerThanBound(b IntervalBound) bool {
	if b.infinite {
		return b.positive
	}
	return f.SmallerThan(b.value)
}

func (f fraction) strictlyGreaterThanBound(b IntervalBound) bool {
	if b.infinite {
		return !b.positive
	}
	return f.GreaterThan(b.value)
}
