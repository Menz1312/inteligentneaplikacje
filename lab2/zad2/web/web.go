package web

import (
	"encoding/json"
	"fmt"
	"lab2/zad2/knn"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

var recsys *knn.Knn

type Beersandstyle struct {
	Beers []*knn.Beer
	Style *knn.Styles
}

func indexFunc(w http.ResponseWriter, r *http.Request) {
	tenbeers := recsys.Get12RandomBeers()
	senddata := Beersandstyle{tenbeers, recsys.GetStyles()}
	tmpl, _ := template.ParseFiles("pages/index.html")
	tmpl.Execute(w, senddata)
}

func recoFunc(w http.ResponseWriter, r *http.Request) {
	data := strings.Split(r.RequestURI, "/")
	if len(data) < 3 {
		fmt.Fprintf(w, "ERROR - wrong path")
		return
	}
	pairs := strings.Split(data[2], ";")
	if len(pairs) < 6 {
		fmt.Fprintf(w, "ERROR - not enough rated beers")
		return
	}
	for i := 0; i < len(pairs)-1; i += 2 {
		bid, _ := strconv.Atoi(pairs[i])
		brate, _ := strconv.Atoi(pairs[i+1])
		beer := recsys.GetBeerByID(bid)
		beer.Rate = float64(brate)
	}
	reco := recsys.GetRecommendation()
	senddata := Beersandstyle{reco, recsys.GetStyles()}
	tmpl, _ := template.ParseFiles("pages/reco.html")
	tmpl.Execute(w, senddata)
}

func explFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // nagłówek JSON
	// obtain beer ID
	beerid, _ := strconv.Atoi(r.PathValue("id"))
	// obtain beer
	beer := recsys.GetBeerByID(beerid)
	fmt.Println(beer)
	// obtain 3 similiar beers
	if beer != nil {
		beers := recsys.GetSimiliar(beer)
		fmt.Println(beers)
		// zadanie 2 - część 1 - zwrocic beers jako JSON
		data, _ := json.Marshal(beers) // konwersja na JSON
		w.Write(data)
	} else {
		// zadanie 2 - część 2 - zwrocic informacje o zlym ID jako JSON
		data, _ := json.Marshal(map[string]string{"error": "Nieprawidłowe ID"})
		w.Write(data)
	}
}

func StartServer() {
	recsys = knn.Initialize()

	fs := http.FileServer(http.Dir("./data"))
	http.Handle("/files/", http.StripPrefix("/files/", fs))
	http.HandleFunc("/", indexFunc)
	http.HandleFunc("/reco/", recoFunc)
	http.HandleFunc("/expl/{id}", explFunc)
	fmt.Println("Serwer wystartował: http://localhost:8080/")
	http.ListenAndServe("localhost:8080", nil)
}
