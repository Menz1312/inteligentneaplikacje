package fuzzy

import (
	"math"
)

type FuzzySet struct {
	x1 float64
	x2 float64
	x3 float64
	x4 float64
}

func Agregate(values ...float64) float64 {
	t := 1.0
	for _, v := range values {
		t *= v
	}
	return t
}

func (f *FuzzySet) Calculate(value float64) float64 {
	if value < f.x1 || value > f.x4 {
		return 0.0
	} else if value >= f.x2 && value <= f.x3 {
		return 1.0
	} else if value < f.x2 {
		val := math.Pi * 0.5 * (f.x2 - value) / (f.x2 - f.x1)
		return math.Cos(val)
	} else { // value > f.x3
		val := math.Pi * 0.5 * (value - f.x3) / (f.x4 - f.x3)
		return math.Cos(val)
	}
}

func NewFuzzySet(min, max, x1, x2, x3, x4 int, log bool) FuzzySet {
	fmin := float64(min)
	fmax := float64(max)
	a := 0.0001 * float64(x1)
	b := 0.0001 * float64(x2)
	c := 0.0001 * float64(x3)
	d := 0.0001 * float64(x4)
	if log {
		a = (math.Pow(10.0, a) - 1.0) / 9.0
		b = (math.Pow(10.0, b) - 1.0) / 9.0
		c = (math.Pow(10.0, c) - 1.0) / 9.0
		d = (math.Pow(10.0, d) - 1.0) / 9.0
	}
	a = fmin + (fmax-fmin)*a
	b = fmin + (fmax-fmin)*b
	c = fmin + (fmax-fmin)*c
	d = fmin + (fmax-fmin)*d
	return FuzzySet{a, b, c, d}
}

func Null() *FuzzySet {
	return &FuzzySet{
		x1: -math.MaxFloat64,
		x2: -math.MaxFloat64,
		x3: math.MaxFloat64,
		x4: math.MaxFloat64,
	}
}
