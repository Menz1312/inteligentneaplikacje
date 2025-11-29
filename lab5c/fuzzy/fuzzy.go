package fuzzy

import (
	"math"
)

type FuzzySet struct {
	x, k float64
}

func (f *FuzzySet) Calculate(x float64) float64 {
	return math.Exp(f.k * (f.x - x) * (f.x - x))
}

type FuzzyRule struct {
	sets   []*FuzzySet
	output float64
}

func (f *FuzzyRule) Activation(x []float64) float64 {
	act := 1.0
	for i := range f.sets {
		act *= f.sets[i].Calculate(x[i])
	}
	return act
}

type FuzzySystem struct {
	rules []*FuzzyRule
}

func (f *FuzzySystem) Predict(x []float64) float64 {
	top, bot := 0.0, 0.0
	for i := range f.rules {
		act := f.rules[i].Activation(x)
		top += act * f.rules[i].output
		bot += act
	}
	if bot == 0.0 {
		return -0.1 * math.MaxFloat64
	}
	return top / bot
}

func newFuzzySystemTSK(x, mins, maxs []float64, inputs, sets, rules int) *FuzzySystem {
	// create fuzzy sets
	fsets := make([][]*FuzzySet, inputs)
	for i := range inputs {
		fsets[i] = make([]*FuzzySet, sets)
		for j := range fsets[i] {
			px := x[i*2*sets+j*2]
			px = mins[i] + (maxs[i]-mins[i])*px
			po := x[i*2*sets+j*2+1]
			po = 0.10 + 0.20*po
			po = mins[i] + (maxs[i]-mins[i])*po
			po = -0.5 / (po * po)
			fsets[i][j] = &FuzzySet{
				x: px,
				k: po,
			}
		}
	}
	// create fuzzy rules
	start := sets * inputs * 2
	frules := make([]*FuzzyRule, rules)
	for i := range rules {
		rsets := make([]*FuzzySet, inputs)
		for j := range inputs {
			best := fsets[j][0]
			vx := x[start+i*(inputs+1)+j]
			vx = mins[j] + (maxs[j]-mins[j])*vx
			bestfit := best.Calculate(vx)
			for k := 1; k < sets; k++ {
				fit := fsets[j][k].Calculate(vx)
				if fit > bestfit {
					bestfit = fit
					best = fsets[j][k]
				}
			}
			rsets[j] = best
		}
		ov := x[start+i*(inputs+1)+inputs]
		ov = -0.05 + 1.1*ov
		ov = mins[inputs] + (maxs[inputs]-mins[inputs])*ov
		frules[i] = &FuzzyRule{
			sets:   rsets,
			output: ov,
		}
	}
	// return system
	return &FuzzySystem{
		rules: frules,
	}
}
