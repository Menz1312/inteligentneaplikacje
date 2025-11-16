package main

import (
	"fmt"
	"lab4/dbase"
)

func main() {
	beers := dbase.GetBeersAll()
	for i := range beers {
		beer := beers[i]
		fmt.Println(i, beer.Id, beer.Name, beer.Rateavg)
	}
}
