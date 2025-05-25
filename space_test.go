package gomath

import "testing"

func TestRealInterval_Contains(t *testing.T) {
	t.Log("testing inclusive bounds")
	set := RealInterval{
		lowerBound: IntervalBound{
			value:        OneFraction.Mul(intToFraction(-1)),
			includeValue: true,
			infinite:     false,
		},
		upperBound: IntervalBound{
			value:        OneFraction,
			includeValue: true,
			infinite:     false,
		},
		customName: "",
	}

	if !set.Contains(NullFraction) {
		t.Fatalf("0 should be in [-1, 1]")
	}
	if !set.Contains(OneFraction) {
		t.Fatalf("1 should be in [-1, 1]")
	}

	t.Log("testing exclusive bounds")
	set = RealInterval{
		lowerBound: IntervalBound{
			value:        OneFraction.Mul(intToFraction(-1)),
			includeValue: false,
			infinite:     false,
		},
		upperBound: IntervalBound{
			value:        OneFraction,
			includeValue: false,
			infinite:     false,
		},
		customName: "",
	}
	if !set.Contains(NullFraction) {
		t.Fatalf("0 should be in [-1, 1]")
	}
	if set.Contains(OneFraction) {
		t.Fatalf("1 should not be in ]-1, 1[")
	}

	t.Log("testing infinite bounds")
	set = RealInterval{
		lowerBound: IntervalBound{
			value:        OneFraction,
			infinite:     false,
			includeValue: true,
		},
		upperBound: IntervalBound{
			infinite: true,
			positive: true,
		},
		customName: "",
	}
	if set.Contains(NullFraction) {
		t.Fatalf("0 should not be in [1, +inf[")
	}
	if !set.Contains(OneFraction) {
		t.Fatalf("1 should be in [1, +inf[")
	}
	if !set.Contains(OneFraction.Mul(intToFraction(2))) {
		t.Fatalf("2 should be in [1, +inf[")
	}
}
