package knn

import (
	"math/rand"
)

type Knn struct {
	beers *Beers
}

func Initialize() *Knn {
	k := Knn{}
	k.beers = LoadBeers("files/beers.csv")
	return &k
}

func (k *Knn) GetStyles() *Styles {
	return &k.beers.styles
}

func (k *Knn) GetRandomBeer() *Beer {
	index := rand.Intn(len(k.beers.beers))
	return &k.beers.beers[index]
}

func (k *Knn) Get10RandomBeers() []*Beer {
	beers := []*Beer{}
	for range 10 {
		index := rand.Intn(len(k.beers.beers))
		beers = append(beers, &k.beers.beers[index])
	}
	return beers
}

func (k *Knn) GetBeerByID(id int) *Beer {
	for b := range k.beers.beers {
		if k.beers.beers[b].Id == id {
			return &k.beers.beers[b]
		}
	}
	return nil
}

func (k *Knn) GetRecommendation() []*Beer {
	return k.beers.Recomendation()
}

func (k *Knn) GetStyleName(id int) string {
	return k.beers.styles.GetStyleName(id)
}

func (k *Knn) GetClosestBeers(base *Beer, count int) []*Beer {
	return k.beers.GetClosestBeers(base, count)
}
