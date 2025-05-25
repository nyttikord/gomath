package gomath

type Space interface {
	Contains(f *fraction) bool
	String() string
}

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
}

func (RealSet) Contains(*fraction) bool {
	return true
}
func (RealSet) String() string {
	return "R"
}

func (set RealInterval) Contains(f *fraction) bool {
	return f.smallerThanBound(set.upperBound) && f.greaterThanBound(set.lowerBound)
}
func (set RealInterval) String() string {
	if set.customName != "" {
		return set.customName
	}
	s := ""

	if set.lowerBound.includeValue {
		s = "[ "
	} else {
		s = "] "
	}
	s += set.lowerBound.value.String() + " ; " + set.upperBound.value.String()

	if set.upperBound.includeValue {
		s += " ]"
	} else {
		s += " ["
	}
	return s
}

func (set UnionSet) Contains(f *fraction) bool {
	for _, space := range set.sets {
		if !space.Contains(f) {
			return false
		}
	}
	return true
}
func (set UnionSet) String() string {
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

func (set PeriodicInterval) Contains(f *fraction) bool {
	if set.interval.Contains(f) {
		return true
	}

	if f.StrictlyGreaterThanBound(set.interval.upperBound) {
		for f.StrictlyGreaterThanBound(set.interval.upperBound) {
			f = f.Sub(set.period)
		}
		return set.interval.Contains(f)
	}
	// here f is necessarily smaller than the lower bound
	for f.StrictlySmallerThanBound(set.interval.lowerBound) {
		f = f.Add(set.period)
	}
	return set.interval.Contains(f)
}
func (set PeriodicInterval) String() string {
	if set.customName != "" {
		return set.customName
	}
	return set.interval.String() + " mod " + set.period.String()
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

func (f fraction) StrictlySmallerThanBound(b IntervalBound) bool {
	if b.infinite {
		return b.positive
	}
	return f.SmallerThan(b.value)
}

func (f fraction) StrictlyGreaterThanBound(b IntervalBound) bool {
	if b.infinite {
		return !b.positive
	}
	return f.GreaterThan(b.value)
}
