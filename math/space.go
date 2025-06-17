package math

type Space interface {
	Contains(f *Fraction) bool
	String() string
}

type RealSet struct{}

type IntervalBound struct {
	Value        *Fraction
	IncludeValue bool
	Infinite     bool
	Positive     bool
}

type RealInterval struct {
	LowerBound *IntervalBound
	UpperBound *IntervalBound
	CustomName string
}

type UnionSet struct {
	Sets       []Space
	CustomName string
}
type PeriodicInterval struct {
	Interval   *RealInterval
	Period     *Fraction
	CustomName string
}

var (
	SpaceRStar = &RealInterval{
		LowerBound: &IntervalBound{
			Value:        NullFraction,
			IncludeValue: false,
			Infinite:     false,
		},
		UpperBound: &IntervalBound{
			Infinite: true,
			Positive: true,
		},
		CustomName: `R \ { 0 }`,
	}
	SpaceRStarPositive = &RealInterval{
		LowerBound: &IntervalBound{
			Value:        NullFraction,
			IncludeValue: false,
			Infinite:     false,
		},
		UpperBound: &IntervalBound{
			Infinite: true,
			Positive: true,
		},
		CustomName: `] 0 ; +inf [`,
	}
)

func (*RealSet) Contains(*Fraction) bool {
	return true
}
func (*RealSet) String() string {
	return "R"
}

func (i *RealInterval) Contains(f *Fraction) bool {
	return f.smallerThanBound(i.UpperBound) && f.greaterThanBound(i.LowerBound)
}
func (i *RealInterval) String() string {
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

func (s *UnionSet) Contains(f *Fraction) bool {
	for _, space := range s.Sets {
		if !space.Contains(f) {
			return false
		}
	}
	return true
}
func (s *UnionSet) String() string {
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

func (set *PeriodicInterval) Contains(f *Fraction) bool {
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
func (set *PeriodicInterval) String() string {
	if set.CustomName != "" {
		return set.CustomName
	}
	return set.Interval.String() + " mod " + set.Period.String()
}

func (f Fraction) smallerThanBound(b *IntervalBound) bool {
	if b.Infinite {
		return b.Positive
	}
	if b.IncludeValue {
		return f.SmallerOrEqualThan(b.Value)
	}
	return f.SmallerThan(b.Value)
}

func (f Fraction) greaterThanBound(b *IntervalBound) bool {
	if b.Infinite {
		return !b.Positive
	}
	if b.IncludeValue {
		return f.GreaterOrEqualThan(b.Value)
	}
	return f.GreaterThan(b.Value)
}

func (f Fraction) strictlySmallerThanBound(b *IntervalBound) bool {
	if b.Infinite {
		return b.Infinite
	}
	return f.SmallerThan(b.Value)
}

func (f Fraction) strictlyGreaterThanBound(b *IntervalBound) bool {
	if b.Infinite {
		return !b.Positive
	}
	return f.GreaterThan(b.Value)
}
