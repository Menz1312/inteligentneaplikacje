package dbase

import (
	"fmt"
	"lab4/db"
)

type Beer struct {
	Id      int
	Abv     float64
	Ibu     float64
	Name    string
	StyleId int
	Rateavg float64
	Ratesum float64
	Ratecnt int
}

func beerFromRow(row rowscan) *Beer {
	var beer Beer
	err := row.Scan(
		&beer.Id, &beer.Abv, &beer.Ibu,
		&beer.Name, &beer.StyleId,
		&beer.Rateavg, &beer.Ratesum, &beer.Ratecnt,
	)
	if err != nil {
		fmt.Println(err)
		return &Beer{
			Id:   0,
			Name: "not found",
		}
	}
	return &beer
}

func GetBeer(id int) *Beer {
	row := db.QueryRow("SELECT * FROM beers WHERE id = ?;", id)
	return beerFromRow(row)
}

func GetBeersTop(cnt int) []*Beer {
	return queryBeers("SELECT * FROM beers ORDER by rateavg DESC LIMIT ?;", cnt)
}

func GetBeersPopular(cnt int) []*Beer {
	return queryBeers("SELECT * FROM beers ORDER by rateavg*(1/(1+EXP(-0.2*ratecnt))) DESC LIMIT ?;", cnt)
}

func GetBeersRandom(cnt int) []*Beer {
	return queryBeers("SELECT * FROM beers ORDER by rand() LIMIT ?;", cnt)
}

func GetBeersForStyle(style int) []*Beer {
	return queryBeers("SELECT * FROM beers WHERE style = ? ORDER by rand() LIMIT 10;", style)
}

func GetBeersAll() []*Beer {
	return queryBeers("SELECT * FROM beers;")
}

func GetBeersFromUser(user int) []*Beer {
	return queryBeers("SELECT b.* FROM beers b LEFT JOIN rates r ON (b.id = r.beer_id) WHERE r.user_id = ? LIMIT 10;", user)
}

func queryBeers(query string, par ...any) []*Beer {
	beers := []*Beer{}
	rows, err := db.Query(query, par...)
	if err != nil {
		fmt.Println(err)
		return beers
	}
	for rows.Next() {
		beer := beerFromRow(rows)
		if beer.Id != 0 {
			beers = append(beers, beer)
		}
	}
	return beers
}
