package gomath

type space interface {
	Contains(f *fraction) bool
	String() string
}

type realSet struct{}

type intervalBound struct {
	Value        *fraction
	IncludeValue bool
	Infinite     bool
	Positive     bool
}

type realInterval struct {
	LowerBound *intervalBound
	UpperBound *intervalBound
	CustomName string
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
	return f.smallerThanBound(i.UpperBound) && f.greaterThanBound(i.LowerBound)
}
func (i *realInterval) String() string {
	if i.CustomName != "" {
		return i.CustomName
	}
	s := ""

	if i.LowerBound.Infinite {
		s = "] "
		if i.UpperBound.Positive {
			s += "+inf"
		} else {
			s += "-inf"
		}
	} else {
		if i.LowerBound.IncludeValue {
			s = "[ "
		} else {
			s = "] "
		}
		s += i.LowerBound.Value.String() + " ; "
	}

	if i.UpperBound.Infinite {
		if i.UpperBound.Positive {
			s += "+inf"
		} else {
			s += "-inf"
		}
		s += " ["
	} else {
		s += i.UpperBound.Value.String()
		if i.UpperBound.IncludeValue {
			s += " ]"
		} else {
			s += " ["
		}
	}

	return s
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

	if f.strictlyGreaterThanBound(set.Interval.UpperBound) {
		for f.strictlyGreaterThanBound(set.Interval.UpperBound) {
			f = f.Sub(set.Period)
		}
		return set.Interval.Contains(f)
	}
	// here f is necessarily smaller than the lower bound
	for f.strictlySmallerThanBound(set.Interval.LowerBound) {
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

func (f *fraction) smallerThanBound(b *intervalBound) bool {
	if b.Infinite {
		return b.Positive
	}
	if b.IncludeValue {
		return f.SmallerOrEqualThan(b.Value)
	}
	return f.SmallerThan(b.Value)
}

func (f *fraction) greaterThanBound(b *intervalBound) bool {
	if b.Infinite {
		return !b.Positive
	}
	if b.IncludeValue {
		return f.GreaterOrEqualThan(b.Value)
	}
	return f.GreaterThan(b.Value)
}

func (f *fraction) strictlySmallerThanBound(b *intervalBound) bool {
	if b.Infinite {
		return b.Infinite
	}
	return f.SmallerThan(b.Value)
}

func (f *fraction) strictlyGreaterThanBound(b *intervalBound) bool {
	if b.Infinite {
		return !b.Positive
	}
	return f.GreaterThan(b.Value)
}
