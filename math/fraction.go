package math

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

type Fraction struct {
	*big.Rat
}

var (
	NullBigInt   = big.NewInt(0)
	OneFraction  = IntToFraction(1)
	NullFraction = IntToFraction(0)
	Pi           *Fraction

	// ErrFractionNotInt is thrown when a non-integer Fraction is converted into an int
	ErrFractionNotInt = errors.New("fraction is not an int")
	// ErrIllegalOperation is thrown when an illegal operation is performed (like dividing by 0)
	ErrIllegalOperation = errors.New("illegal operation")
	// ErrUnsupportedOperation is thrown when an unsupported operation is performed
	ErrUnsupportedOperation = errors.New("unsupported operation")
)

func init() {
	var err error
	Pi, err = FloatToFraction(math.Pi)
	if err != nil {
		panic(err)
	}
}

func NewFraction(a, b int64) *Fraction {
	return &Fraction{big.NewRat(a, b)}
}

// IntToFraction converts an int64 into a Fraction
func IntToFraction(n int64) *Fraction {
	return NewFraction(n, 1)
}

// FloatToFraction converts a float64 into a Fraction
func FloatToFraction(f float64) (*Fraction, error) {
	if f == float64(int64(f)) {
		return IntToFraction(int64(f)), nil
	}
	s := strconv.FormatFloat(f, 'f', -1, 64)
	sp := strings.Split(s, ".")
	i, err := strconv.ParseInt(sp[0]+sp[1], 10, 64)
	if err != nil {
		return nil, err
	}
	return NewFraction(i, int64(math.Pow(10, float64(len(sp[1]))))), nil
}

func (f Fraction) String() string {
	return f.Rat.RatString()
}

func (f Fraction) Is(a *Fraction) bool {
	return f.Rat.Num().Cmp(a.Rat.Num()) == 0 && f.Denom().Cmp(a.Denom()) == 0
}

func (f Fraction) Approx(precision int) string {
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

func (f Fraction) CanBeRepresentedExactly(precision int) bool {
	if f.IsInt() {
		return true
	}

	// Fraction is not an int, and the precision is negative or null, the representation
	if precision <= 0 {
		return false
	}

	rest := big.NewInt(0).Mod(f.Num(), f.Denom())

	for n := 1; n <= precision && rest.Cmp(big.NewInt(0)) == 1; n++ {
		rest.Mod(rest.Mul(rest, big.NewInt(10)), f.Denom())
	}

	// Is exact only if rest is null
	return rest.Cmp(big.NewInt(0)) == 0
}

func (f Fraction) Copy() *Fraction {
	b := IntToFraction(1)
	b.Set(f.Rat)
	return b
}

func (f Fraction) SmallerOrEqualThan(b *Fraction) bool {
	x := big.NewInt(0)
	y := big.NewInt(0)
	// fractions are always simplified
	x.Mul(f.Num(), b.Denom())
	y.Mul(b.Num(), f.Denom())

	return x.Cmp(y) <= 0
}

func (f Fraction) SmallerThan(b *Fraction) bool {
	return f.SmallerOrEqualThan(b) && !f.Is(b)
}

func (f Fraction) GreaterOrEqualThan(b *Fraction) bool {
	return !f.SmallerThan(b)
}

func (f Fraction) GreaterThan(b *Fraction) bool {
	return !f.SmallerOrEqualThan(b)
}

// Add a Fraction
func (f Fraction) Add(a *Fraction) *Fraction {
	c := f.Copy()
	c.Rat.Add(f.Rat, a.Rat)
	return c
}

func (f Fraction) Neg() *Fraction {
	c := f.Copy()
	c.Num().Mul(f.Num(), big.NewInt(-1))
	return c
}

// Sub (subtract) a Fraction
func (f Fraction) Sub(a *Fraction) *Fraction {
	c := f.Copy()
	return c.Add(a.Neg())
}

// Mul (multiply) by Fraction
func (f Fraction) Mul(a *Fraction) *Fraction {
	c := f.Copy()
	c.Rat.Mul(f.Rat, a.Rat)
	return c
}

// Inv (invert) the Fraction
func (f Fraction) Inv() (*Fraction, error) {
	c := f.Copy()
	if c.Num().Cmp(NullBigInt) == 0 {
		return c, errors.Join(ErrIllegalOperation, errors.New("cannot invert a null Fraction"))
	}
	c.Rat.Inv(f.Rat)
	return c, nil
}

// Div (divide) by a Fraction
func (f Fraction) Div(a *Fraction) (*Fraction, error) {
	invA, err := a.Inv() // is a copy
	if err != nil {
		return f.Copy(), errors.Join(err, errors.New("cannot divide by a null Fraction"))
	}
	mul := f.Mul(invA) // is a copy
	return mul, nil
}

// IsInt returns true if the Fraction is an int
func (f Fraction) IsInt() bool {
	return f.Rat.IsInt()
}

// Int converts the Fraction to an int.
// Returns ErrFractionNotInt if the Fraction isn't an int (check before with Fraction.IsInt)
func (f Fraction) Int() (*big.Int, error) {
	if !f.IsInt() {
		return NullBigInt, errors.Join(ErrFractionNotInt, errors.New(f.String()+" is not an int"))
	}
	r := big.Int{}
	return r.Div(f.Num(), f.Denom()), nil
}

// Float converts the Fraction to a float
func (f Fraction) Float() (float64, bool) {
	return f.Float64()
}

// Exp the Fraction by another
func (f Fraction) Exp(a *Fraction) (*Fraction, error) {
	if a.IsInt() {
		n, _ := a.Int()
		fl, _ := f.Float()
		if fl == 0 {
			if n.Cmp(NullBigInt) == 0 {
				return OneFraction, nil
			}
			return NullFraction, nil
		}
		c := f.Copy()
		c.Num().Exp(f.Num(), n, nil)
		c.Denom().Exp(f.Denom(), n, nil)
		return c, nil
	}
	//afl, _ := a.Float()
	//nf, err := FloatToFraction(math.Pow(float64(f.Num().P), afl))
	//if err != nil {
	//	return NullFraction, errors.Join(err, errors.New("cannot convert numerator^a into a Fraction"))
	//}
	//nff, err := FloatToFraction(math.Pow(float64(f.Denominator), afl))
	//if err != nil {
	//	return NullFraction, errors.Join(err, errors.New("cannot convert denominator^a into a Fraction"))
	//}
	//pf, err := nf.Div(nff)
	//if err != nil {
	//	return NullFraction, err
	//}
	return nil, errors.Join(ErrUnsupportedOperation, fmt.Errorf("Fraction.Exp(%s) is not supported because it's not an int", a))
}
