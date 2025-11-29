package fuzzy

import (
	"labfs/data"
	"math/rand"
)

func TSKConstructor(mins, maxs []float64, inputs, sets, rules int) data.Producer[*FuzzySystem] {
	return data.Producer[*FuzzySystem]{
		Parameters: inputs*sets*2 + rules*(inputs+1),
		Create: func(x []float64) *FuzzySystem {
			return newFuzzySystemTSK(x, mins, maxs, inputs, sets, rules)
		},
		CreateRandom: func() *FuzzySystem {
			x := make([]float64, inputs*sets*2+rules*(inputs+1))
			for i := range x {
				x[i] = rand.Float64()
			}
			return newFuzzySystemTSK(x, mins, maxs, inputs, sets, rules)
		},
	}
}
