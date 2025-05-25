package gomath

import "testing"

func TestFraction_Simplify(t *testing.T) {
	t.Log("testing positive denominator")
	f := fraction{Numerator: 6, Denominator: 8}
	f.Simplify()
	expected := fraction{Numerator: 3, Denominator: 4}
	if f != expected {
		t.Errorf("got %s; want %s", f.String(), expected.String())
	}

	t.Log("testing negative denominator")
	f = fraction{Numerator: 6, Denominator: -5}
	f.Simplify()
	expected = fraction{Numerator: -6, Denominator: 5}
	if f != expected {
		t.Errorf("got %s; want %s", f.String(), expected.String())
	}

	t.Log("testing double negative fraction")
	f = fraction{Numerator: -6, Denominator: -5}
	f.Simplify()
	expected = fraction{Numerator: 6, Denominator: 5}
	if f != expected {
		t.Errorf("got %s; want %s", f.String(), expected.String())
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
	if err == nil {
		t.Errorf("expected error")
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
	if err == nil {
		t.Errorf("expected error")
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
