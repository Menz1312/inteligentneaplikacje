package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

// Zadanie 1
type Pralka struct {
	Nazwa          string
	Obroty         float64
	ZuzycieEnergii float64
	ZuzycieWody    float64
	Cena           float64
}

type Data interface {
	Get() []float64
}

func (p Pralka) Get() []float64 {
	return []float64{p.Obroty, p.ZuzycieEnergii, p.ZuzycieWody}
}

var pralki = []Pralka{
	{"PralkaA", 300, 100, 20, 1000},
	{"PralkaB", 400, 200, 25, 1200},
	{"PralkaC", 500, 300, 30, 1500},
	{"PralkaD", 600, 250, 30, 1600},
	{"PralkaE", 550, 300, 25, 1700},
}

func (p1 *Pralka) podobienstwo(p2 Pralka) float64 {
	wynik := 0.0
	wynik += math.Abs(p1.Obroty-p2.Obroty) / 600.0
	wynik += math.Abs(p1.ZuzycieEnergii-p2.ZuzycieEnergii) / 300
	wynik += math.Abs(p1.ZuzycieWody-p2.ZuzycieWody) / 30
	return 3.0 - wynik
}

func getPralkaById(pralki []Pralka, id int) Pralka {
	return pralki[id]
}

func itemFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p := getPralkaById(pralki, rand.Intn(5))
	data, _ := json.Marshal(p)
	w.Write(data)
}

func obiektyFunc(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("pages/obiekty.html")
	tmpl.Execute(w, pralki)
}

func showFunc(w http.ResponseWriter, r *http.Request) {
	pralkaId, _ := strconv.Atoi(r.PathValue("id"))
	fmt.Println(pralkaId)
	tmpl, _ := template.ParseFiles("pages/show.html")
	tmpl.Execute(w, getPralkaById(pralki, pralkaId))
}

func main() {
	fmt.Println(pralki)
	fmt.Println("Podobienstwo: ", pralki[0].podobienstwo(pralki[2]))

	http.HandleFunc("/item/", itemFunc)
	http.HandleFunc("/obiekty/", obiektyFunc)
	http.HandleFunc("/show/{id}", showFunc)
	http.ListenAndServe("localhost:8080", nil)
}
