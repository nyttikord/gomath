package expression

import (
	"github.com/nyttikord/gomath"
	"github.com/nyttikord/gomath/ast"
	"testing"
)

var testOpt = &ast.Options{
	Decimal:   true,
	Precision: 6,
}

func TestMathFunction_Exp(t *testing.T) {
	res, err := gomath.ParseAndCalculate("exp(0)", testOpt)
	expected := "1"
	if err != nil {
		t.Fatal(err)
	}
	if res != expected {
		t.Errorf("got %v; want %v", res, expected)
	}
}

func TestMathFunction_Cos(t *testing.T) {
	res, err := gomath.ParseAndCalculate("cos(0)", testOpt)
	expected := "1"
	if err != nil {
		t.Fatal(err)
	}
	if res != expected {
		t.Errorf("got %v; want %v", res, expected)
	}
}

func TestMathFunction_Sin(t *testing.T) {
	res, err := gomath.ParseAndCalculate("sin(0)", testOpt)
	expected := "0"
	if err != nil {
		t.Fatal(err)
	}
	if res != expected {
		t.Errorf("got %v; want %v", res, expected)
	}
}

func TestMathFunction_Tan(t *testing.T) {
	res, err := gomath.ParseAndCalculate("tan(0)", testOpt)
	expected := "0"
	if err != nil {
		t.Fatal(err)
	}
	if res != expected {
		t.Errorf("got %v; want %v", res, expected)
	}
}
