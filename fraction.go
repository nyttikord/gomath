package gomath

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type fraction struct {
	Numerator   int64
	Denominator int64
}

var (
	NullFraction = &fraction{Numerator: 0, Denominator: 1}
	OneFraction  = &fraction{Numerator: 1, Denominator: 1}

	ErrFractionNotInt   = errors.New("fraction is not an int")
	ErrIllegalOperation = errors.New("illegal operation")
)

// intToFraction converts an int64 into a fraction
func intToFraction(n int64) *fraction {
	return &fraction{
		Numerator:   n,
		Denominator: 1,
	}
}

// floatToFraction converts a float64 into a fraction
func floatToFraction(f float64) (*fraction, error) {
	if f == float64(int64(f)) {
		return intToFraction(int64(f)), nil
	}
	s := strconv.FormatFloat(f, 'f', -1, 64)
	sp := strings.Split(s, ".")
	i, err := strconv.ParseInt(sp[0]+sp[1], 10, 64)
	if err != nil {
		return nil, err
	}
	return &fraction{
		Numerator:   i,
		Denominator: int64(math.Pow(10, float64(len(sp[1])))),
	}, nil
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func (f *fraction) String() string {
	if f.Denominator != 1 {
		return fmt.Sprintf("%d/%d", f.Numerator, f.Denominator)
	}
	return fmt.Sprintf("%d", f.Numerator)
}

// Simplify the fraction
func (f *fraction) Simplify() *fraction {
	pgcd := gcd(f.Numerator, f.Denominator)
	if pgcd == 0 {
		return f
	}
	f.Numerator = f.Numerator / pgcd
	f.Denominator = f.Denominator / pgcd
	return f
}

// Add a fraction
func (f *fraction) Add(a *fraction) (*fraction, error) {
	f.Numerator = f.Numerator*a.Denominator + a.Numerator*f.Denominator
	f.Denominator = f.Denominator * a.Denominator
	return f.Simplify(), nil
}

// Sub (subtrack) a fraction
func (f *fraction) Sub(a *fraction) (*fraction, error) {
	f.Numerator = f.Numerator*a.Denominator - a.Numerator*f.Denominator
	f.Denominator = f.Denominator * a.Denominator
	return f.Simplify(), nil
}

// Mul (multiply) by fraction
func (f *fraction) Mul(a *fraction) (*fraction, error) {
	f.Numerator = f.Numerator * a.Numerator
	f.Denominator = f.Denominator * a.Denominator
	return f.Simplify(), nil
}

// Inv (invert) the fraction
func (f *fraction) Inv() (*fraction, error) {
	if f.Numerator == 0 {
		return f, errors.Join(ErrIllegalOperation, errors.New("cannot invert a null fraction"))
	}
	f.Numerator, f.Denominator = f.Denominator, f.Numerator
	return f.Simplify(), nil
}

// Div (divide) by a fraction
func (f *fraction) Div(a *fraction) (*fraction, error) {
	inva, err := a.Inv()
	if err != nil {
		return f, errors.Join(err, errors.New("cannot divide by a null fraction"))
	}
	mul, _ := f.Mul(inva) // avoid checking error because it's always nil for fraction.Mul
	return mul.Simplify(), nil
}

// IsInt returns true if the fraction is an int
func (f *fraction) IsInt() bool {
	return f.Numerator%f.Denominator == 0
}

// Int convers the fraction to an int.
// Returns ErrFractionNotInt if the fraction isn't an int (check before with fraction.IsInt)
func (f *fraction) Int() (int64, error) {
	if !f.IsInt() {
		return 0, errors.Join(ErrFractionNotInt, errors.New(f.String()+" is not an int"))
	}
	return f.Numerator / f.Denominator, nil
}

// Float converts the fraction to a float
func (f *fraction) Float() float64 {
	return float64(f.Numerator) / float64(f.Denominator)
}

// Pow the fraction by another
func (f *fraction) Pow(a *fraction) (*fraction, error) {
	if a.IsInt() {
		n, _ := a.Int()
		if f.Float() == 0 {
			if n == 0 {
				return OneFraction, nil
			}
			return NullFraction, nil
		}
		nf := fraction{
			Numerator:   int64(math.Pow(float64(f.Numerator), float64(n))),
			Denominator: int64(math.Pow(float64(f.Denominator), float64(n))),
		}
		*f = nf
		return f.Simplify(), nil
	}
	afl := a.Float()
	nf, err := floatToFraction(math.Pow(float64(f.Numerator), afl))
	if err != nil {
		return NullFraction, errors.Join(err, errors.New("cannot convert numerator^a into a fraction"))
	}
	nff, err := floatToFraction(math.Pow(float64(f.Denominator), afl))
	if err != nil {
		return NullFraction, errors.Join(err, errors.New("cannot convert denominator^a into a fraction"))
	}
	pf, err := nf.Div(nff)
	if err != nil {
		return NullFraction, err
	}
	*f = *pf
	return f.Simplify(), nil
}
