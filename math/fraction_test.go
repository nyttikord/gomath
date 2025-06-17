package math

import (
	"errors"
	"testing"
)

func TestFraction_Is(t *testing.T) {
	t.Log("testing simple is")
	f := NewFraction(3, 4)
	if !f.Is(f) {
		t.Errorf("%s is not %s", f, f)
	}
	a := NewFraction(6, 8)
	if !f.Is(a) {
		t.Errorf("%s is not %s", f, a)
	}
	b := NewFraction(9, 8)
	if f.Is(b) {
		t.Errorf("%s is %s", f, b)
	}

	t.Log("testing negative denominator")
	f = NewFraction(6, -5)
	expected := NewFraction(-6, 5)
	if !f.Is(expected) {
		t.Errorf("got %s; want %s", f, expected)
	}

	t.Log("testing double negative Fraction")
	f = NewFraction(-6, -5)
	expected = NewFraction(6, 5)
	if !f.Is(expected) {
		t.Errorf("got %s; want %s", f, expected)
	}
}

func TestFractionComparison(t *testing.T) {
	t.Log("testing equal Fraction")
	f := NewFraction(5, 3)
	t.Log("testing smaller or equal")
	if !f.SmallerOrEqualThan(f) {
		t.Errorf("fractions should be equal")
	}
	if !NewFraction(5, 3).Is(f) {
		t.Fatalf("fractions 5/3 and %s should be equal", f)
	}
	t.Log("testing smaller")
	if f.SmallerThan(f) {
		t.Errorf("fractions should be equal")
	}
	t.Log("testing greater or equal")
	if !f.GreaterOrEqualThan(f) {
		t.Errorf("fractions should be equal")
	}
	t.Log("testing greater")
	if f.GreaterThan(f) {
		t.Errorf("fractions should be equal")
	}

	t.Log("testing unequal fractions")
	// f < b
	b := NewFraction(7, 4)
	t.Log("testing smaller or equal")
	if !f.SmallerOrEqualThan(b) {
		t.Errorf("%s should be smaller or equal than %s", b, f)
	}
	t.Log("testing smaller")
	if !f.SmallerThan(b) {
		t.Errorf("%s should be smaller than %s", b, f)
	}
	t.Log("testing greater or equal")
	if f.GreaterOrEqualThan(b) {
		t.Errorf("%s should be greater or equal than %s", f, b)
	}
	t.Log("testing greater")
	if b.GreaterThan(b) {
		t.Errorf("%s should be greater than %s", f, b)
	}
}

func TestFraction_Add(t *testing.T) {
	a := NewFraction(5, 6)
	b := NewFraction(8, 3)
	expected := NewFraction(7, 2)
	res := a.Add(b)
	if !res.Is(expected) {
		t.Errorf("got %s; want %s", res.String(), expected.String())
	}
}

func TestFraction_Sub(t *testing.T) {
	a := NewFraction(5, 6)
	b := NewFraction(8, 3)
	expected := NewFraction(-11, 6)
	res := a.Sub(b)
	if !res.Is(expected) {
		t.Errorf("got %s; want %s", res.String(), expected.String())
	}
}

func TestFraction_Mul(t *testing.T) {
	a := NewFraction(5, 6)
	b := NewFraction(8, 3)
	expected := NewFraction(20, 9)
	res := a.Mul(b)
	if !res.Is(expected) {
		t.Errorf("got %s; want %s", res.String(), expected.String())
	}
}

func TestFraction_Neg(t *testing.T) {
	f := NewFraction(5, 6)
	expected := NewFraction(-5, 6)
	res := f.Neg()
	if !res.Is(expected) {
		t.Errorf("got %s; want %s", res.String(), expected.String())
	}
}

func TestFraction_Div(t *testing.T) {
	t.Log("testing division")
	a := NewFraction(5, 6)
	b := NewFraction(8, 3)
	expected := NewFraction(5, 16)
	res, err := a.Div(b)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Is(expected) {
		t.Errorf("got %s; want %s", res.String(), expected.String())
	}

	t.Log("testing division by null Fraction")
	a = OneFraction
	b = NullFraction
	_, err = a.Div(b)
	if !errors.Is(err, ErrIllegalOperation) {
		t.Errorf("expected illegal operation error, not %s", err)
	}
}

func TestFraction_Inv(t *testing.T) {
	t.Log("testing invert")
	a := NewFraction(5, 6)
	expected := NewFraction(6, 5)
	res, err := a.Inv()
	if err != nil {
		t.Fatal(err)
	}
	if !res.Is(expected) {
		t.Errorf("got %s; want %s", res.String(), expected.String())
	}

	t.Log("testing invert a null Fraction")
	a = NullFraction
	_, err = a.Inv()
	if !errors.Is(err, ErrIllegalOperation) {
		t.Errorf("expected illegal operation error, not %s", err)
	}
}

func TestFraction_Approx(t *testing.T) {
	expected := "3.1415"
	f := NewFraction(6283, 2000)

	t.Log("testing exact precision of value")
	res := f.Approx(4)
	if res != expected {
		t.Errorf("got %s; want %s", res, expected)
	}

	t.Log("testing higher precision than exact value")
	res = f.Approx(10)
	if res != expected {
		t.Errorf("got %s; want %s", res, expected)
	}

	t.Log("testing lower precision than exact value")
	expected = "3.14"
	res = f.Approx(2)
	if res != expected {
		t.Errorf("got %s; want %s", res, expected)
	}

	t.Log("testing integer Fraction")
	f = IntToFraction(357)
	res = f.Approx(10)
	expected = "357"
	if res != expected {
		t.Errorf("got %s; want %s", res, expected)
	}

	t.Log("testing precision of 0")
	res = f.Approx(0)
	if res != expected {
		t.Errorf("got %s; want %s", res, expected)
	}
}

func TestFraction_CanBeRepresentedExactly(t *testing.T) {
	f := NewFraction(1, 16) // Equal exactly to 0.0625

	if !f.CanBeRepresentedExactly(5) {
		t.Errorf("1/16 with precision of 5 digits should be exact")
	}
	if !f.CanBeRepresentedExactly(4) {
		t.Errorf("1/16 with precision of 4 digits should be exact")
	}
	if f.CanBeRepresentedExactly(3) {
		t.Errorf("1/16 with precision of 3 digits should not be exact")
	}
	if f.CanBeRepresentedExactly(0) {
		t.Errorf("1/16 with precision of 0 digits should not be exact")
	}

	f = NewFraction(5, 1) // Integer Fraction
	if !f.CanBeRepresentedExactly(0) {
		t.Errorf("5/1 should be exact no matter the precision")
	}
	if !f.CanBeRepresentedExactly(10) {
		t.Errorf("5/1 should be exact no matter the precision")
	}
}
