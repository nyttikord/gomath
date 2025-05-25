package gomath

import (
	"errors"
	"testing"
)

func TestFraction_Simplify(t *testing.T) {
	t.Log("testing positive denominator")
	f := fraction{Numerator: 6, Denominator: 8}.Simplify()
	expected := fraction{Numerator: 3, Denominator: 4}
	if *f != expected {
		t.Errorf("got %s; want %s", f.String(), expected.String())
	}

	t.Log("testing negative denominator")
	f = fraction{Numerator: 6, Denominator: -5}.Simplify()
	expected = fraction{Numerator: -6, Denominator: 5}
	if *f != expected {
		t.Errorf("got %s; want %s", f.String(), expected.String())
	}

	t.Log("testing double negative fraction")
	f = &fraction{Numerator: -6, Denominator: -5}
	f = f.Simplify()
	expected = fraction{Numerator: 6, Denominator: 5}
	if *f != expected {
		t.Errorf("got %s; want %s", f.String(), expected.String())
	}
}

func TestFractionComparison(t *testing.T) {
	t.Log("testing equal fraction")
	f := fraction{Numerator: 5, Denominator: 3}
	t.Log("smaller or equal")
	if !f.SmallerOrEqualThan(&f) {
		t.Errorf("fractions should be equal")
	}
	t.Log("smaller")
	if f.SmallerThan(&f) {
		t.Errorf("fractions should be equal")
	}
	t.Log("greater or equal")
	if !f.GreaterOrEqualThan(&f) {
		t.Errorf("fractions should be equal")
	}
	t.Log("greater")
	if f.GreaterThan(&f) {
		t.Errorf("fractions should be equal")
	}

	t.Log("testing unequal fractions")
	// a < b
	a := f
	b := fraction{Numerator: 7, Denominator: 4}
	t.Log("smaller or equal")
	if !a.SmallerOrEqualThan(&b) {
		t.Errorf("a should be smaller than b")
	}
	t.Log("smaller")
	if !a.SmallerThan(&b) {
		t.Errorf("a should be smaller than b")
	}
	t.Log("greater or equal")
	if a.GreaterOrEqualThan(&b) {
		t.Errorf("a should be smaller than b")
	}
	t.Log("greater")
	if b.GreaterThan(&b) {
		t.Errorf("a should be smaller than b")
	}
}

func TestFraction_Add(t *testing.T) {
	a := fraction{Numerator: 5, Denominator: 6}
	b := fraction{Numerator: 8, Denominator: 3}
	expected := fraction{Numerator: 7, Denominator: 2}
	res := *a.Add(&b)
	if res != expected {
		t.Errorf("got %s; want %s", res.String(), expected.String())
	}
}

func TestFraction_Sub(t *testing.T) {
	a := fraction{Numerator: 5, Denominator: 6}
	b := fraction{Numerator: 8, Denominator: 3}
	expected := fraction{Numerator: -11, Denominator: 6}
	res := *a.Sub(&b)
	if res != expected {
		t.Errorf("got %s; want %s", res.String(), expected.String())
	}
}

func TestFraction_Mul(t *testing.T) {
	a := fraction{Numerator: 5, Denominator: 6}
	b := fraction{Numerator: 8, Denominator: 3}
	expected := fraction{Numerator: 20, Denominator: 9}
	res := *a.Mul(&b)
	if res != expected {
		t.Errorf("got %s; want %s", res.String(), expected.String())
	}
}

func TestFraction_Div(t *testing.T) {
	t.Log("testing division")
	a := fraction{Numerator: 5, Denominator: 6}
	b := fraction{Numerator: 8, Denominator: 3}
	expected := fraction{Numerator: 5, Denominator: 16}
	res, err := a.Div(&b)
	if err != nil {
		t.Fatal(err)
	}
	if *res != expected {
		t.Errorf("got %s; want %s", res.String(), expected.String())
	}

	t.Log("testing division by null fraction")
	a = *OneFraction
	b = *NullFraction
	_, err = a.Div(&b)
	if !errors.Is(err, ErrIllegalOperation) {
		t.Errorf("expected illegal operation error, not %s", err)
	}
}

func TestFraction_Inv(t *testing.T) {
	t.Log("testing division")
	a := fraction{Numerator: 5, Denominator: 6}
	expected := fraction{Numerator: 6, Denominator: 5}
	res, err := a.Inv()
	if err != nil {
		t.Fatal(err)
	}
	if *res != expected {
		t.Errorf("got %s; want %s", res.String(), expected.String())
	}

	t.Log("testing division by null fraction")
	a = *NullFraction
	_, err = a.Inv()
	if !errors.Is(err, ErrIllegalOperation) {
		t.Errorf("expected illegal operation error, not %s", err)
	}
}

func TestFraction_Approx(t *testing.T) {
	expected := "3.1415"
	f := fraction{Numerator: 6283, Denominator: 2000}

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

	t.Log("testing integer fraction")
	f = *intToFraction(357)
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
