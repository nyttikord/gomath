package gomath

import "testing"

func TestRealInterval_Contains(t *testing.T) {
	t.Log("testing inclusive bounds")
	set := realInterval{
		LowerBound: &intervalBound{
			Value:        OneFraction.Mul(intToFraction(-1)),
			IncludeValue: true,
			Infinite:     false,
		},
		UpperBound: &intervalBound{
			Value:        OneFraction,
			IncludeValue: true,
			Infinite:     false,
		},
		CustomName: "",
	}

	if !set.Contains(NullFraction) {
		t.Fatalf("0 should be in [-1, 1]")
	}
	if !set.Contains(OneFraction) {
		t.Fatalf("1 should be in [-1, 1]")
	}

	t.Log("testing exclusive bounds")
	set = realInterval{
		LowerBound: &intervalBound{
			Value:        OneFraction.Mul(intToFraction(-1)),
			IncludeValue: false,
			Infinite:     false,
		},
		UpperBound: &intervalBound{
			Value:        OneFraction,
			IncludeValue: false,
			Infinite:     false,
		},
		CustomName: "",
	}
	if !set.Contains(NullFraction) {
		t.Fatalf("0 should be in [-1, 1]")
	}
	if set.Contains(OneFraction) {
		t.Fatalf("1 should not be in ]-1, 1[")
	}

	t.Log("testing infinite bounds")
	set = realInterval{
		LowerBound: &intervalBound{
			Value:        OneFraction,
			Infinite:     false,
			IncludeValue: true,
		},
		UpperBound: &intervalBound{
			Infinite: true,
			Positive: true,
		},
		CustomName: "",
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
