package gomath

import "testing"

func TestRealInterval_Contains(t *testing.T) {
	t.Log("testing inclusive bounds")
	set := realInterval{
		LowerBound: &intervalBound{
			Value:        oneFraction.Neg(),
			IncludeValue: true,
			Infinite:     false,
		},
		UpperBound: &intervalBound{
			Value:        oneFraction,
			IncludeValue: true,
			Infinite:     false,
		},
		CustomName: "",
	}

	if !set.Contains(nullFraction) {
		t.Errorf("0 should be in [-1, 1]")
	}
	if !set.Contains(oneFraction) {
		t.Errorf("1 should be in [-1, 1]")
	}

	t.Log("testing exclusive bounds")
	set = realInterval{
		LowerBound: &intervalBound{
			Value:        oneFraction.Mul(intToFraction(-1)),
			IncludeValue: false,
			Infinite:     false,
		},
		UpperBound: &intervalBound{
			Value:        oneFraction,
			IncludeValue: false,
			Infinite:     false,
		},
		CustomName: "",
	}
	if !set.Contains(nullFraction) {
		t.Errorf("0 should be in ]-1, 1[")
	}
	if set.Contains(oneFraction) {
		t.Errorf("1 should not be in ]-1, 1[")
	}

	t.Log("testing infinite bounds")
	set = realInterval{
		LowerBound: &intervalBound{
			Value:        oneFraction,
			Infinite:     false,
			IncludeValue: true,
		},
		UpperBound: &intervalBound{
			Infinite: true,
			Positive: true,
		},
		CustomName: "",
	}
	if set.Contains(nullFraction) {
		t.Errorf("0 should not be in [1, +inf[")
	}
	if !set.Contains(oneFraction) {
		t.Errorf("1 should be in [1, +inf[")
	}
	if !set.Contains(intToFraction(2)) {
		t.Errorf("2 should be in [1, +inf[")
	}
}
