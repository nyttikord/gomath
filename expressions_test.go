package gomath

import (
	"testing"
)

func TestEvalSum(t *testing.T) {
	genericTest(t, "1+2", "3")
}

func TestEvalSub(t *testing.T) {
	genericTest(t, "1-2", "-1")
}

func TestEvalAddUnary(t *testing.T) {
	genericTest(t, "1+-2", "-1")
}

func TestEvalMult(t *testing.T) {
	genericTest(t, "2*3", "6")
}

func TestEvalMultUnary(t *testing.T) {
	genericTest(t, "2*-3", "-6")
}

func TestEvalDiv(t *testing.T) {
	genericTest(t, "2/3", "2/3")
}

func TestEvalDivUnary(t *testing.T) {
	genericTest(t, "2/-3", "-2/3")
}

func TestEvalDivDecimal(t *testing.T) {
	lexr, err := lex("1/10")
	if err != nil {
		t.Fatal(err)
	}
	tree, err := astParse(lexr, astTypeCalculation)
	if err != nil {
		t.Fatal(err)
	}
	if tree.Type != astTypeCalculation {
		t.Errorf("got type %d; want %d", tree.Type, astTypeCalculation)
	}
	val, err := tree.Body.Eval(&Options{true, 3})
	if err != nil {
		t.Fatal(err)
	}
	if val.String() != "0.1" {
		t.Errorf("got %s; want %s", val, "0.1")
	}
	if t.Failed() {
		t.Log(tree)
	}
}

func TestEvalPriority(t *testing.T) {
	t.Log("testing 2*(1+2)")
	genericTest(t, "2*(1+2)", "6")
	t.Log("testing 2*1+2")
	genericTest(t, "2*1+2", "4")
	t.Log("testing 2*(1+2)^2")
	genericTest(t, "2*(1+2)^2", "18")
}

func TestEvalOmitMultSigne(t *testing.T) {
	t.Log("testing 2(3+2)")
	genericTest(t, "2(3+2)", "10")
	t.Log("testing 2^2(3+2)")
	genericTest(t, "2^2(3+2)", "20")
	t.Log("testing 2(3+2)^2")
	genericTest(t, "2(3+2)^2", "50")
	t.Log("testing 2cos(0)")
	genericTest(t, "2cos(0)", "2")
	t.Log("testing 2^2cos(0)")
	genericTest(t, "2^2cos(0)", "4")
}

func TestEvalPrioritySpecialCase(t *testing.T) {
	// check https://en.wikipedia.org/wiki/Order_of_operations#Special_cases
	t.Log("testing -3^2") // must be interpreted as -(3^2)
	genericTest(t, "-3^2", "-9")
	t.Log("testing 6/2(1+2)") // must be interpreted as (6/2)(1+2)
	genericTest(t, "6/2(1+2)", "9")
	t.Log("testing 3^2^3") // must be interpreted as (3^2)^3
	genericTest(t, "3^2^3", "729")
	t.Log("testing 6/2*2") // must be interpreted as 2*(6/2)
	genericTest(t, "6/2*2", "6")
	t.Log("testing 6/2cos(1)") // must be interpreted as 6/(2cos(1))
	res, err := ParseAndCalculate("6/(2cos(1))", &Options{})
	if err != nil {
		t.Fatal(err)
	}
	genericTest(t, "6/2cos(1)", res)
}

func TestEvalFactorial(t *testing.T) {
	genericTest(t, "3!", "6")
	genericTest(t, "3^2!", "362880")
	genericTest(t, "3*2^2!", "72")
	genericTest(t, "(2*2)!", "24")
}

func genericTest(t *testing.T, exp string, expected string) {
	lexr, err := lex(exp)
	if err != nil {
		t.Fatal(err)
	}
	tree, err := astParse(lexr, astTypeCalculation)
	if err != nil {
		t.Fatal(err)
	}
	if tree.Type != astTypeCalculation {
		t.Errorf("got type %d; want %d", tree.Type, astTypeCalculation)
	}
	val, err := tree.Body.Eval(&Options{})
	if err != nil {
		t.Fatal(err)
	}
	if val.String() != expected {
		t.Errorf("got %s; want %s", val, expected)
	}
	if t.Failed() {
		t.Log(tree)
	}
}

func TestEvalLatex(t *testing.T) {
	genericTestRenderLatex(t, "(1+2)/3", `\frac{1 + 2}{3}`)
	genericTestRenderLatex(t, "3/(1+2)", `\frac{3}{1 + 2}`)
	genericTestRenderLatex(t, "cos(2*pi)", `\cos\left(2 \times \pi\right)`)
	genericTestRenderLatex(t, "e^(5+2)", `e^{5 + 2}`)
	genericTestRenderLatex(t, "e^5", `e^5`)
	genericTestRenderLatex(t, "(1+2)^5", `\left(1 + 2\right)^5`)
	genericTestRenderLatex(t, "5*(1+2)^5", `5 \times \left(1 + 2\right)^5`)
	genericTestRenderLatex(t, "5(1+2)^5", `5 \times \left(1 + 2\right)^5`)
	genericTestRenderLatex(t, "(1+2/3)/2", `\frac{1 + \frac{2}{3}}{2}`)
	genericTestRenderLatex(t, "2*(1+2)", `2 \times \left(1 + 2\right)`)
	genericTestRenderLatex(t, "2!", `2!`)
	genericTestRenderLatex(t, "2*2!", `2 \times 2!`)
	genericTestRenderLatex(t, "2!*2!", `2! \times 2!`)
	genericTestRenderLatex(t, "3^2!", `3^2!`)
	genericTestRenderLatex(t, "(3+2)!", `\left(3 + 2\right)!`)
}

func genericTestRenderLatex(t *testing.T, exp string, excepted string) {
	lexr, err := lex(exp)
	if err != nil {
		t.Fatal(err)
	}
	tree, err := astParse(lexr, astTypeLatex)
	if err != nil {
		t.Fatal(err)
	}
	if tree.Type != astTypeLatex {
		t.Errorf("got type %d; want %d", tree.Type, astTypeLatex)
	}
	val, err := tree.Body.Eval(&Options{})
	if err != nil {
		t.Fatal(err)
	}
	if val.String() != excepted {
		t.Errorf("got %s; want %s", val, excepted)
	}
	if t.Failed() {
		t.Log(tree)
	}
}
