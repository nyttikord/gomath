package gomath

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

type fraction struct {
	*big.Rat
}

var (
	nullBigInt   = big.NewInt(0)
	oneFraction  = intToFraction(1)
	nullFraction = intToFraction(0)
	pi           *fraction

	// ErrFractionNotInt is thrown when a non-integer fraction is converted into an int
	ErrFractionNotInt = errors.New("fraction is not an int")
	// ErrIllegalOperation is thrown when an illegal operation is performed (like dividing by 0)
	ErrIllegalOperation = errors.New("illegal operation")
	// ErrUnsupportedOperation is thrown when an unsupported operation is performed
	ErrUnsupportedOperation = errors.New("unsupported operation")
)

func init() {
	var err error
	pi, err = floatToFraction(math.Pi)
	if err != nil {
		panic(err)
	}
}

func newFraction(a, b int64) *fraction {
	return &fraction{big.NewRat(a, b)}
}

// intToFraction converts an int64 into a fraction
func intToFraction(n int64) *fraction {
	return newFraction(n, 1)
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
	return newFraction(i, int64(math.Pow(10, float64(len(sp[1]))))), nil
}

func (f fraction) String() string {
	return f.Rat.RatString()
}

func (f fraction) Is(a *fraction) bool {
	return f.Rat.Num().Cmp(a.Rat.Num()) == 0 && f.Denom().Cmp(a.Denom()) == 0
}

func (f fraction) Approx(precision int) string {
	if f.IsInt() {
		n, _ := f.Int()
		return fmt.Sprintf("%d", n)
	}
	s := strings.TrimSuffix(f.Rat.FloatString(precision), "0")
	for strings.HasSuffix(s, "0") {
		s = strings.TrimSuffix(s, "0")
	}
	return s
}

func (f fraction) Copy() *fraction {
	b := intToFraction(1)
	b.Set(f.Rat)
	return b
}

func (f fraction) SmallerOrEqualThan(b *fraction) bool {
	x := big.NewInt(0)
	y := big.NewInt(0)
	// fractions are always simplified
	x.Mul(f.Num(), b.Denom())
	y.Mul(b.Num(), f.Denom())

	return x.Cmp(y) <= 0
}

func (f fraction) SmallerThan(b *fraction) bool {
	return f.SmallerOrEqualThan(b) && !f.Is(b)
}

func (f fraction) GreaterOrEqualThan(b *fraction) bool {
	return !f.SmallerThan(b)
}

func (f fraction) GreaterThan(b *fraction) bool {
	return !f.SmallerOrEqualThan(b)
}

// Add a fraction
func (f fraction) Add(a *fraction) *fraction {
	c := f.Copy()
	c.Rat.Add(f.Rat, a.Rat)
	return c
}

func (f fraction) Neg() *fraction {
	c := f.Copy()
	c.Num().Mul(f.Num(), big.NewInt(-1))
	return c
}

// Sub (subtract) a fraction
func (f fraction) Sub(a *fraction) *fraction {
	c := f.Copy()
	return c.Add(a.Neg())
}

// Mul (multiply) by fraction
func (f fraction) Mul(a *fraction) *fraction {
	c := f.Copy()
	c.Rat.Mul(f.Rat, a.Rat)
	return c
}

// Inv (invert) the fraction
func (f fraction) Inv() (*fraction, error) {
	c := f.Copy()
	if c.Num().Cmp(nullBigInt) == 0 {
		return c, errors.Join(ErrIllegalOperation, errors.New("cannot invert a null fraction"))
	}
	c.Rat.Inv(f.Rat)
	return c, nil
}

// Div (divide) by a fraction
func (f fraction) Div(a *fraction) (*fraction, error) {
	invA, err := a.Inv() // is a copy
	if err != nil {
		return f.Copy(), errors.Join(err, errors.New("cannot divide by a null fraction"))
	}
	mul := f.Mul(invA) // is a copy
	return mul, nil
}

// IsInt returns true if the fraction is an int
func (f fraction) IsInt() bool {
	return f.Rat.IsInt()
}

// Int convers the fraction to an int.
// Returns ErrFractionNotInt if the fraction isn't an int (check before with fraction.IsInt)
func (f fraction) Int() (*big.Int, error) {
	if !f.IsInt() {
		return nullBigInt, errors.Join(ErrFractionNotInt, errors.New(f.String()+" is not an int"))
	}
	r := big.Int{}
	return r.Div(f.Num(), f.Denom()), nil
}

// Float converts the fraction to a float
func (f fraction) Float() (float64, bool) {
	return f.Float64()
}

// Exp the fraction by another
func (f fraction) Exp(a *fraction) (*fraction, error) {
	if a.IsInt() {
		n, _ := a.Int()
		fl, _ := f.Float()
		if fl == 0 {
			if n.Cmp(nullBigInt) == 0 {
				return oneFraction, nil
			}
			return nullFraction, nil
		}
		c := f.Copy()
		c.Num().Exp(f.Num(), n, nil)
		c.Denom().Exp(f.Denom(), n, nil)
		return c, nil
	}
	//afl, _ := a.Float()
	//nf, err := floatToFraction(math.Pow(float64(f.Num().P), afl))
	//if err != nil {
	//	return nullFraction, errors.Join(err, errors.New("cannot convert numerator^a into a fraction"))
	//}
	//nff, err := floatToFraction(math.Pow(float64(f.Denominator), afl))
	//if err != nil {
	//	return nullFraction, errors.Join(err, errors.New("cannot convert denominator^a into a fraction"))
	//}
	//pf, err := nf.Div(nff)
	//if err != nil {
	//	return nullFraction, err
	//}
	return nil, errors.Join(ErrUnsupportedOperation, fmt.Errorf("fraction.Exp(%s) is not supported because it's not an int", a))
}
