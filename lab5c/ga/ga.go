package ga

import (
	"fmt"
	"labfs/data"
	"math/rand"
)

func NewOptimizer[T data.Predicter, U any](prod data.Producer[T], samples []*data.Sample[U]) GA[T, U] {
	return GA[T, U]{
		best:       nil,
		pops:       []*Individual[T, U]{},
		population: 200,
		parameters: prod.Parameters,
		creator:    prod.Create,
		samples:    samples,
	}
}

type GA[T data.Predicter, U any] struct {
	best       *Individual[T, U]
	pops       []*Individual[T, U]
	population int
	parameters int
	creator    func([]float64) T
	samples    []*data.Sample[U]
}

func (ga *GA[T, U]) Optimize(iterations int) T {
	ga.initialize()
	for i := range iterations {
		wsp := 1.0 - float64(i+1)/float64(iterations)
		fmt.Printf("Iteration %d/%d...", i+1, iterations)
		ga.iterate(wsp)
		fmt.Printf("\rIteration %d done, best: %0.2f      \n", i+1, ga.best.fitness)
	}
	return ga.best.system
}

func (ga *GA[T, U]) initialize() {
	ga.pops = make([]*Individual[T, U], ga.population)
	fmt.Printf("Population initialization...")
	for i := range ga.population {
		ind := newIndividualRandom[T, U](ga.parameters)
		ind.system = ga.creator(ind.x)
		ind.evaluate(ga.samples)
		ga.pops[i] = ind
		if ga.best == nil || ind.fitness < ga.best.fitness {
			ga.best = ind
		}
		fmt.Printf("\rIndividual %d/%d - error: %0.2f      ", i+1, ga.population, ind.fitness)
	}
	fmt.Printf("\rInitialization done, best: %0.2f      \n", ga.best.fitness)
}

func (ga *GA[T, U]) iterate(wsp float64) {
	mc := 0.5 + 0.5*wsp
	mr := 0.001 + 0.10*wsp
	gm := 0.1 + 0.2*wsp
	for i := range ga.population {
		// select parents - tournament
		i1, i2 := rand.Intn(ga.population), rand.Intn(ga.population)
		pa, pb := ga.pops[i1], ga.pops[i2]
		for range 2 {
			i3, i4 := rand.Intn(ga.population), rand.Intn(ga.population)
			if ga.pops[i3].fitness < pa.fitness {
				pa = ga.pops[i3]
			}
			if ga.pops[i4].fitness < pb.fitness {
				pb = ga.pops[i4]
			}
		}
		// crossover and mutate
		ind := newIndividualCrossover(pa, pb)
		if rand.Float64() < mc {
			ind.mutate(gm, mr)
		}
		ind.system = ga.creator(ind.x)
		ind.evaluate(ga.samples)
		// update population
		if ind.fitness < ga.pops[i].fitness {
			ga.pops[i] = ind
			if ind.fitness < ga.best.fitness {
				ga.best = ind
			}
		}
		fmt.Printf("\rIndividual %d/%d - error: %0.2f      ", i+1, len(ga.pops), ind.fitness)
	}
}
