package ga

import (
	"labfs/data"
	"math"
	"math/rand"
)

type Individual[T data.Predicter, U any] struct {
	fitness float64
	x       []float64
	system  T
}

func newIndividualRandom[T data.Predicter, U any](parameters int) *Individual[T, U] {
	x := make([]float64, parameters)
	for i := range x {
		x[i] = rand.Float64()
	}
	return &Individual[T, U]{
		fitness: math.MaxFloat64,
		x:       x,
	}
}

func newIndividualCrossover[T data.Predicter, U any](pa *Individual[T, U], pb *Individual[T, U]) *Individual[T, U] {
	parameters := len(pa.x)
	x := make([]float64, parameters)
	for i := range x {
		if rand.Float64() < 0.8 {
			x[i] = pa.x[i] + (pb.x[i]-pa.x[i])*rand.Float64()
		} else if rand.Float64() < 0.5 {
			x[i] = pa.x[i]
		} else {
			x[i] = pb.x[i]
		}
	}
	return &Individual[T, U]{
		fitness: math.MaxFloat64,
		x:       x,
	}
}

func (in *Individual[T, U]) evaluate(samples []*data.Sample[U]) {
	in.fitness = data.Evaluate(in.system, samples)
}

func (in *Individual[T, U]) mutate(gm, mr float64) {
	for i := range in.x {
		if rand.Float64() < gm {
			in.x[i] += (rand.Float64() - 0.5) * 2.0 * mr
			if in.x[i] < 0.0 {
				in.x[i] = 0.0
			} else if in.x[i] > 1.0 {
				in.x[i] = 1.0
			}
		}
	}
}
