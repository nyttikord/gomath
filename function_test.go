package gomath

import "testing"

func TestNewFunction(t *testing.T) {
	f, n, err := NewFunction("x -> x^2")
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Errorf("got %d, want 1", n)
	}
	result, err := f(map[string]string{"x": "5"})
	if err != nil {
		t.Fatal(err)
	}
	if result.String() != "25" {
		t.Errorf("got %s, want 25", result.String())
	}

	f, n, err = NewFunction("x, y -> x^y")
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Errorf("got %d, want 2", n)
	}
	result, err = f(map[string]string{"x": "5", "y": "2"})
	if err != nil {
		t.Fatal(err)
	}
	if result.String() != "25" {
		t.Errorf("got %s, want 25", result.String())
	}
}
