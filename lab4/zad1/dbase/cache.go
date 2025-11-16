package dbase

import (
	"time"
)

var top10 []*Beer = GetBeersTop(10)
var top10time = time.Now()
var pop10 []*Beer = GetBeersPopular(10)
var pop10time = time.Now()

func GetBeersTop10() []*Beer {
	if time.Since(top10time).Minutes() > 10 {
		top10 = GetBeersTop(10)
		top10time = time.Now()
	}
	return top10
}

func GetBeersPopular10() []*Beer {
	if time.Since(pop10time).Minutes() > 10 {
		pop10 = GetBeersTop(10)
		pop10time = time.Now()
	}
	return pop10
}
