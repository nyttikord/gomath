package gomath

import "testing"

var opt = &Options{
	Decimal:   true,
	Precision: 6,
}

func TestMathFunction_Exp(t *testing.T) {
	res, err := Parse("exp(0)", opt)
	expected := "1"
	if err != nil {
		t.Fatal(err)
	}
	if res != expected {
		t.Fatalf("got %v; want %v", res, expected)
	}
}

func TestMathFunction_Cos(t *testing.T) {
	res, err := Parse("cos(0)", opt)
	expected := "1"
	if err != nil {
		t.Fatal(err)
	}
	if res != expected {
		t.Fatalf("got %v; want %v", res, expected)
	}
}

func TestMathFunction_Sin(t *testing.T) {
	res, err := Parse("sin(0)", opt)
	expected := "0"
	if err != nil {
		t.Fatal(err)
	}
	if res != expected {
		t.Fatalf("got %v; want %v", res, expected)
	}
}

func TestMathFunction_Tan(t *testing.T) {
	res, err := Parse("tan(0)", opt)
	expected := "0"
	if err != nil {
		t.Fatal(err)
	}
	if res != expected {
		t.Fatalf("got %v; want %v", res, expected)
	}
}
