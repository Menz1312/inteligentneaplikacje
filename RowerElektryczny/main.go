package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

// Zadanie 1
type RowerElektryczny struct {
	Nazwa       string
	Zasieg      float64
	MaxPredkosc float64
	RozmiarKol  float64
	Cena        float64
}

func getRowery() (rowery []RowerElektryczny) {
	r1 := RowerElektryczny{"RowerA", 200.0, 10.0, 15.0, 7000.0}
	r2 := RowerElektryczny{"RowerB", 100.0, 20.0, 27.0, 4000.0}
	r3 := RowerElektryczny{"RowerC", 50.0, 50.0, 26.0, 10000.0}
	r4 := RowerElektryczny{"RowerD", 80.0, 30.0, 18.0, 4000.0}
	r5 := RowerElektryczny{"RowerE", 120.0, 25.0, 27.0, 6000.0}
	return []RowerElektryczny{r1, r2, r3, r4, r5}
}

func getRowerById(rowery []RowerElektryczny, id int) RowerElektryczny {
	return rowery[id]
}

// Zadanie 2
func roznica(r1 RowerElektryczny, r2 RowerElektryczny) float64 {
	return math.Abs(r1.Zasieg-r2.Zasieg)/100 + math.Abs(r1.MaxPredkosc-r2.MaxPredkosc)/30 + math.Abs(r1.RozmiarKol-r2.RozmiarKol)
}

// Zadanie 3
func indexFunc(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("pages/index.html")
	tmpl.Execute(w, getRowery())
}

// Zadanie 4
func rowerFunc(w http.ResponseWriter, r *http.Request) {
	rowerId, _ := strconv.Atoi(r.PathValue("id"))
	tmpl, _ := template.ParseFiles("pages/rower.html")
	tmpl.Execute(w, getRowerById(getRowery(), rowerId))
}

// Zadanie 5
type fuzzy struct {
	x, o float64
}

func (f *fuzzy) calc(x float64) float64 {
	return math.Exp(-(x - f.x) * (x - f.x) / f.o)
}

var f1, f2, f3 = fuzzy{50, 10}, fuzzy{20, 5}, fuzzy{20, 5}

func dopasowanie(r RowerElektryczny) float64 {
	return f1.calc(r.Zasieg) * f2.calc(r.MaxPredkosc) * f3.calc(r.RozmiarKol)
}

func maxDopasowanie(rowery []RowerElektryczny) RowerElektryczny {
	maxRower := rowery[0]
	for i := 1; i < len(rowery); i++ {
		if dopasowanie(rowery[i]) > dopasowanie(maxRower) {
			maxRower = rowery[i]
		}
	}
	return maxRower
}

// Zadanie 6
func najbardziejPodobny(rowerDoPorownania RowerElektryczny, rowery []RowerElektryczny) RowerElektryczny {
	roznicaMin := roznica(rowerDoPorownania, rowery[0])
	rowerPodobny := rowery[0]
	for i := 1; i < len(rowery); i++ {
		if roznica(rowerDoPorownania, rowery[i]) < roznicaMin {
			rowerPodobny = rowery[i]
		}
	}
	return rowerPodobny
}

func main() {
	rowery := getRowery()
	roznica := roznica(getRowerById(rowery, 1), getRowerById(rowery, 2))
	fmt.Println("Różnica: ", roznica)

	fmt.Println("Max dopasowanie: ", maxDopasowanie(rowery))
	fmt.Println("Najbardziej podobny: ", najbardziejPodobny(RowerElektryczny{"Feobike SX4", 130, 45, 27.5, 7000}, rowery))

	http.HandleFunc("/", indexFunc)
	http.HandleFunc("/rower/{id}", rowerFunc)
	http.ListenAndServe("localhost:8080", nil)
}
