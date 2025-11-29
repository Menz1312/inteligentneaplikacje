package data

import "math"

type Sampleable interface {
	GetData() ([]float64, float64)
}

type Predicter interface {
	Predict([]float64) float64
}

type Producer[T Predicter] struct {
	Parameters   int
	Create       func([]float64) T
	CreateRandom func() T
}

type Sample[T any] struct {
	Inputs []float64
	Output float64
	Object T
}

func NewSample[T Sampleable](o T) *Sample[T] {
	inputs, output := o.GetData()
	return &Sample[T]{
		Inputs: inputs,
		Output: output,
		Object: o,
	}
}

func Evaluate[T Predicter, U any](system T, samples []*Sample[U]) float64 {
	predicter := system.Predict
	mae := 0.0
	for _, sample := range samples {
		systout := predicter(sample.Inputs)
		dataout := sample.Output
		mae += math.Abs(systout - dataout)
	}
	return mae / float64(len(samples))
}

func ToSamples[T Sampleable](in []T) ([]*Sample[T], []float64, []float64, int) {
	if len(in) == 0 {
		panic("no data in samples")
	}
	inp, _ := in[0].GetData()
	if len(inp) == 0 {
		panic("no inputs in sample")
	}
	inpl := len(inp)
	mins := make([]float64, inpl+1)
	maxs := make([]float64, inpl+1)
	for i := range mins {
		mins[i] = math.MaxFloat64
		maxs[i] = -math.MaxFloat64
	}
	out := make([]*Sample[T], len(in))
	for i, v := range in {
		sample := NewSample(v)
		inputs := sample.Inputs
		output := sample.Output
		for j := range inputs {
			if inputs[j] < mins[j] {
				mins[j] = inputs[j]
			}
			if inputs[j] > maxs[j] {
				maxs[j] = inputs[j]
			}
			if output < mins[inpl] {
				mins[inpl] = output
			}
			if output > maxs[inpl] {
				maxs[inpl] = output
			}
		}
		out[i] = sample
	}
	return out, mins, maxs, inpl
}
