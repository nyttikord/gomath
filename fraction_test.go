package gomath

import (
	"errors"
	"testing"
)

func TestFraction_Simplify(t *testing.T) {
	t.Log("testing positive denominator")
	f := newFraction(6, 8).Simplify()
	expected := newFraction(3, 4)
	if !f.Is(expected) {
		t.Errorf("got %s; want %s", f, expected)
	}

	t.Log("testing negative denominator")
	f = newFraction(6, -5).Simplify()
	expected = newFraction(-6, 5)
	if !f.Is(expected) {
		t.Errorf("got %s; want %s", f, expected)
	}

	t.Log("testing double negative fraction")
	f = newFraction(-6, -5).Simplify()
	expected = newFraction(6, 5)
	if !f.Is(expected) {
		t.Errorf("got %s; want %s", f, expected)
	}
}

func TestFractionComparison(t *testing.T) {
	t.Log("testing equal fraction")
	f := newFraction(5, 3)
	t.Log("smaller or equal")
	if !f.SmallerOrEqualThan(f) {
		t.Errorf("fractions should be equal")
	}
	t.Log("smaller")
	if f.SmallerThan(f) {
		t.Errorf("fractions should be equal")
	}
	t.Log("greater or equal")
	if !f.GreaterOrEqualThan(f) {
		t.Errorf("fractions should be equal")
	}
	t.Log("greater")
	if f.GreaterThan(f) {
		t.Errorf("fractions should be equal")
	}

	t.Log("testing unequal fractions")
	// a < b
	a := f
	b := newFraction(7, 4)
	t.Log("smaller or equal")
	if !a.SmallerOrEqualThan(b) {
		t.Errorf("%s should be smaller or equal than %s", b, a)
	}
	t.Log("smaller")
	if !a.SmallerThan(b) {
		t.Errorf("%s should be smaller than %s", b, a)
	}
	t.Log("greater or equal")
	if a.GreaterOrEqualThan(b) {
		t.Errorf("%s should be greater or equal than %s", a, b)
	}
	t.Log("greater")
	if b.GreaterThan(b) {
		t.Errorf("%s should be greater than %s", a, b)
	}
}

func TestFraction_Add(t *testing.T) {
	a := newFraction(5, 6)
	b := newFraction(8, 3)
	expected := newFraction(7, 2)
	res := a.Add(b)
	if res.Is(expected) {
		t.Errorf("got %s; want %s", res.String(), expected.String())
	}
}

func TestFraction_Sub(t *testing.T) {
	a := newFraction(5, 6)
	b := newFraction(8, 3)
	expected := newFraction(-11, 6)
	res := a.Sub(b)
	if res.Is(expected) {
		t.Errorf("got %s; want %s", res.String(), expected.String())
	}
}

func TestFraction_Mul(t *testing.T) {
	a := newFraction(5, 6)
	b := newFraction(8, 3)
	expected := newFraction(20, 9)
	res := a.Mul(b)
	if !res.Is(expected) {
		t.Errorf("got %s; want %s", res.String(), expected.String())
	}
}

func TestFraction_Neg(t *testing.T) {
	f := newFraction(5, 6)
	expected := newFraction(-5, 6)
	res := f.Neg()
	if res.Is(expected) {
		t.Errorf("got %s; want %s", res.String(), expected.String())
	}
}

func TestFraction_Div(t *testing.T) {
	t.Log("testing division")
	a := newFraction(5, 6)
	b := newFraction(8, 3)
	expected := newFraction(5, 16)
	res, err := a.Div(b)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Is(expected) {
		t.Errorf("got %s; want %s", res.String(), expected.String())
	}

	t.Log("testing division by null fraction")
	a = oneFraction
	b = nullFraction
	_, err = a.Div(b)
	if !errors.Is(err, ErrIllegalOperation) {
		t.Errorf("expected illegal operation error, not %s", err)
	}
}

func TestFraction_Inv(t *testing.T) {
	t.Log("testing invert")
	a := newFraction(5, 6)
	expected := newFraction(6, 5)
	res, err := a.Inv()
	if err != nil {
		t.Fatal(err)
	}
	if res.Is(expected) {
		t.Errorf("got %s; want %s", res.String(), expected.String())
	}

	t.Log("testing invert a null fraction")
	a = nullFraction
	_, err = a.Inv()
	if !errors.Is(err, ErrIllegalOperation) {
		t.Errorf("expected illegal operation error, not %s", err)
	}
}

func TestFraction_Approx(t *testing.T) {
	expected := "3.1415"
	f := newFraction(6283, 2000)

	t.Log("testing exact precision of value")
	res := f.Approx(4)
	if res != expected {
		t.Errorf("got %s; want %s", res, expected)
	}

	t.Log("testing higher precision than exact value")
	res = f.Approx(10)
	expected = "3.1415000000"
	if res != expected {
		t.Errorf("got %s; want %s", res, expected)
	}

	t.Log("testing lower precision than exact value")
	expected = "3.14"
	res = f.Approx(2)
	if res != expected {
		t.Errorf("got %s; want %s", res, expected)
	}

	t.Log("testing integer fraction")
	f = intToFraction(357)
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
