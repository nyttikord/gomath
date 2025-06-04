package gomath

import "testing"

func TestRes_Approx(t *testing.T) {
	r, err := Parse("3/4")
	if err != nil {
		t.Fatal(err)
	}
	excepted := "1"
	got := r.Approx(0)
	if excepted != got {
		t.Errorf("excepted: %s, got: %s", excepted, got)
	}
	excepted = "0.8"
	got = r.Approx(1)
	if excepted != got {
		t.Errorf("excepted: %s, got: %s", excepted, got)
	}
	excepted = "0.75"
	got = r.Approx(2)
	if excepted != got {
		t.Errorf("excepted: %s, got: %s", excepted, got)
	}
}

func TestRes_IsExact(t *testing.T) {
	r, err := Parse("3/4")
	if err != nil {
		t.Fatal(err)
	}
	if !r.IsExact() {
		t.Errorf("excepted: %t, got: %t", true, r.IsExact())
	}
	r, err = Parse("2/3")
	if err != nil {
		t.Fatal(err)
	}
	if r.IsExact() {
		t.Errorf("excepted: %t, got: %t", false, r.IsExact())
	}
}

func TestRes_LaTeX(t *testing.T) {
	r, err := Parse("3/4")
	if err != nil {
		t.Fatal(err)
	}
	got, err := r.LaTeX()
	if err != nil {
		t.Fatal(err)
	}
	excepted, err := ParseAndConvertToLaTeX("3/4", &Options{})
	if err != nil {
		t.Fatal(err)
	}
	if got != excepted {
		t.Errorf("excepted: %s, got: %s", excepted, got)
	}
}
