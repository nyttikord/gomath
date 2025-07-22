package math

import (
	"testing"
)

func TestRealInterval_Contains(t *testing.T) {
	t.Log("testing inclusive bounds")
	set := RealInterval{
		LowerBound: &IntervalBound{
			Value:        OneFraction.Neg(),
			IncludeValue: true,
			Infinite:     false,
		},
		UpperBound: &IntervalBound{
			Value:        OneFraction,
			IncludeValue: true,
			Infinite:     false,
		},
		CustomName: "",
	}

	if !set.Contains(NullFraction) {
		t.Errorf("0 should be in [-1, 1]")
	}
	if !set.Contains(OneFraction) {
		t.Errorf("1 should be in [-1, 1]")
	}

	t.Log("testing exclusive bounds")
	set = RealInterval{
		LowerBound: &IntervalBound{
			Value:        OneFraction.Mul(IntToFraction(-1)),
			IncludeValue: false,
			Infinite:     false,
		},
		UpperBound: &IntervalBound{
			Value:        OneFraction,
			IncludeValue: false,
			Infinite:     false,
		},
		CustomName: "",
	}
	if !set.Contains(NullFraction) {
		t.Errorf("0 should be in ]-1, 1[")
	}
	if set.Contains(OneFraction) {
		t.Errorf("1 should not be in ]-1, 1[")
	}

	t.Log("testing infinite bounds")
	set = RealInterval{
		LowerBound: &IntervalBound{
			Value:        OneFraction,
			Infinite:     false,
			IncludeValue: true,
		},
		UpperBound: &IntervalBound{
			Infinite: true,
			Positive: true,
		},
		CustomName: "",
	}
	if set.Contains(NullFraction) {
		t.Errorf("0 should not be in [1, +inf[")
	}
	if !set.Contains(OneFraction) {
		t.Errorf("1 should be in [1, +inf[")
	}
	if !set.Contains(IntToFraction(2)) {
		t.Errorf("2 should be in [1, +inf[")
	}
}
