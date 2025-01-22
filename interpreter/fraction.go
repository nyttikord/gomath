package interpreter

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Fraction struct {
	Numerator   int64
	Denominator int64
}

var (
	NullFraction = &Fraction{Numerator: 0, Denominator: 1}
	OneFraction  = &Fraction{Numerator: 1, Denominator: 1}

	FractionNotIntErr = errors.New("fraction is not an int")
)

func IntToFraction(n int64) *Fraction {
	return &Fraction{
		Numerator:   n,
		Denominator: 1,
	}
}

func FloatToFraction(f float64) (*Fraction, error) {
	if f == float64(int64(f)) {
		return IntToFraction(int64(f)), nil
	}
	s := strconv.FormatFloat(f, 'f', -1, 64)
	sp := strings.Split(s, ".")
	i, err := strconv.ParseInt(sp[0]+sp[1], 10, 64)
	if err != nil {
		return NullFraction, FractionNotIntErr
	}
	return &Fraction{
		Numerator:   i,
		Denominator: int64(math.Pow(10, float64(len(sp[1])))),
	}, nil
}

func (f *Fraction) String() string {
	if f.Denominator != 1 {
		return fmt.Sprintf("%d/%d", f.Numerator, f.Denominator)
	}
	return fmt.Sprintf("%d", f.Numerator)
}

func (f *Fraction) Add(a *Fraction) *Fraction {
	f.Numerator = f.Numerator*a.Denominator + a.Numerator*f.Denominator
	f.Denominator = f.Denominator * a.Denominator
	return f
}

func (f *Fraction) Sub(a *Fraction) *Fraction {
	f.Numerator = f.Numerator*a.Denominator - a.Numerator*f.Denominator
	f.Denominator = f.Denominator * a.Denominator
	return f
}

func (f *Fraction) Mul(a *Fraction) *Fraction {
	f.Numerator = f.Numerator * a.Numerator
	f.Denominator = f.Denominator * a.Denominator
	return f
}

func (f *Fraction) Inv() *Fraction {
	f.Numerator, f.Denominator = f.Denominator, f.Numerator
	return f
}

func (f *Fraction) Div(a *Fraction) *Fraction {
	return f.Mul(a.Inv())
}

func (f *Fraction) IsInt() bool {
	return f.Numerator%f.Denominator == 0
}

func (f *Fraction) Int() (int64, error) {
	if !f.IsInt() {
		return 0, errors.Join(FractionNotIntErr, errors.New(f.String()+" is not an int"))
	}
	return f.Numerator / f.Denominator, nil
}

func (f *Fraction) Float() float64 {
	return float64(f.Numerator) / float64(f.Denominator)
}

func (f *Fraction) Pow(a *Fraction) *Fraction {
	if a.IsInt() {
		n, err := a.Int()
		if err != nil {
			panic(err)
		}
		if f.Float() == 0 {
			if n == 0 {
				return OneFraction
			}
			return NullFraction
		}
		nf := Fraction{
			Numerator:   int64(math.Pow(float64(f.Numerator), float64(n))),
			Denominator: int64(math.Pow(float64(f.Denominator), float64(n))),
		}
		*f = nf
		return f
	}
	fl := f.Float()
	if fl == 0 {
		return NullFraction
	}
	d := math.Pow(fl, a.Float())
	ff, err := FloatToFraction(d)
	if err != nil {
		panic(err)
	}
	*f = *ff
	return f
}
