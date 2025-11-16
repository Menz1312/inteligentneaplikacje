package kmeans

import (
	"lab4/dbase"
	"math/rand"
)

// globalny numer klastra (grupy)
var CLUSTERID int = 0

// liczba iteracji algorytmu
const ITERATIONS int = 50

// struktura klastra - przeznaczenie tylko dla Beers (uproszczenie)
type Cluster struct {
	Grupa int           // numer klastra
	x     []float64     // współrzędne klastra
	Beers []*dbase.Beer // lista przypisanych piw do klastra
}

// inicjalizacja klastra - wartości od 0 do 20 (ABV) i od 0 do 150 (IBU)
func (c *Cluster) initialize() {
	c.x = []float64{rand.Float64() * 20.0, rand.Float64() * 150.0}
	c.Grupa = CLUSTERID
	CLUSTERID++
}

// dystans pomiędyz klastrem a piwem (z normalizacją)
func (c *Cluster) distance(b *dbase.Beer) float64 {
	distance := 0.0
	abvDiff := (b.Abv - c.x[0]) * 0.02
	ibuDiff := (b.Ibu - c.x[1]) * 0.01
	distance += abvDiff * abvDiff
	distance += ibuDiff * ibuDiff
	return distance
}

// zresetowanie przypisanych piw
func (c *Cluster) reset() {
	c.Beers = make([]*dbase.Beer, 0)
}

// pobranie piw z klastra
func (c *Cluster) GetBeers() []*dbase.Beer {
	return c.Beers
}

// uaktualnienie klastrów
func (c *Cluster) update() {
	// pusty klaster? wylosowanie nowej pozycji
	if len(c.Beers) < 1 {
		c.x = []float64{rand.Float64() * 0.2, rand.Float64() * 150.0}
		return
	}
	// w innym wypadku średnia pozycjia piw
	c.x[0] = 0.0
	c.x[1] = 0.0
	for b := range c.Beers {
		c.x[0] += c.Beers[b].Abv
		c.x[1] += c.Beers[b].Ibu
	}
	div := float64(len(c.Beers))
	c.x[0] /= div
	c.x[1] /= div
}

// funkcji grupująca Kmeans (poprawna implementacja wymaga uogólnienia i rozbicia na mniejsze funkcje)
func Kmeans(beers []*dbase.Beer, count int) []Cluster {
	clusters := make([]Cluster, count)
	// inicjalizacja klastrów (losowe położenie)
	CLUSTERID = 0
	for i := range count {
		clusters[i].initialize()
	}
	// główna pętla algorytmu
	for range ITERATIONS {
		// resetujemy piwa przypisane wcześniej do klastrów
		for i := range count {
			clusters[i].reset()
		}
		// każde piwo przypisujemy do klastra
		for b := range beers {
			// szukamy najbliższego klastra
			gdzie := 0
			dystans := clusters[gdzie].distance(beers[b])
			for i := 1; i < count; i++ {
				tmp := clusters[i].distance(beers[b])
				if tmp < dystans {
					dystans = tmp
					gdzie = i
				}
			}
			// przypisujemy piwo do najbliższego klastra
			clusters[gdzie].Beers = append(clusters[gdzie].Beers, beers[b])
		}
		// uaktualniamy klastry (ich pozycję jako średnią pozycję piw)
		for i := range count {
			clusters[i].update()
		}
	}
	// zwrócenie wyniku
	return clusters
}
