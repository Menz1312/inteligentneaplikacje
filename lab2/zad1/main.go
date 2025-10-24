package main

import (
	"fmt"
	"html/template"
	"lab0201/knn"
	"net/http"
	"strconv"
)

type ItemPageData struct {
	Beer     *knn.Beer
	RecBeers []*knn.Beer
}

// wczytanie bazy danych piw (zmienna globalna)
var recsys *knn.Knn = knn.Initialize()

func piwaFunc(w http.ResponseWriter, r *http.Request) {
	randombeers := recsys.Get10RandomBeers()
	tmpl, err := template.ParseFiles("pages/beer.html")
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		fmt.Printf("template parse error: %+v\r\n", err)
		return
	}
	err = tmpl.Execute(w, randombeers)
	if err != nil {
		http.Error(w, "template exec error", http.StatusInternalServerError)
		fmt.Printf("template executing error: %+v\r\n", err)
		return
	}
}

func rekoFunc(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	for name, value := range r.PostForm {
		// odczyt danych z formularza
		id, er1 := strconv.Atoi(name[4:])
		rate, er2 := strconv.ParseFloat(value[0], 64)
		beer := recsys.GetBeerByID(id)
		// sprawdzenie błędów
		if beer == nil || er1 != nil || er2 != nil {
			fmt.Printf("parse error %d %f %v\r\n", id, rate, beer)
			continue
		}
		// ocena piwa
		beer.Rate = rate
		fmt.Printf("Beer %d ('%s') rated as %f\r\n", id, beer.Name, rate)
	}
	// pobranie rekomendacji
	recbeers := recsys.GetRecommendation()
	tmpl, _ := template.ParseFiles("pages/reko.html")
	tmpl.Execute(w, recbeers)
}

func itemFunc(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	beer := recsys.GetBeerByID(id)
	if beer == nil {
		http.NotFound(w, r)
		return
	}

	// tymczasowo oceń bieżące piwo
	beer.Rate = 5.0

	fmt.Println(beer.Id, beer.Name, beer.Rate)

	// pobranie rekomendacji (na podstawie ocenionego piwa)
	recbeers := recsys.GetClosestBeers(beer, 3)

	// jeśli chcesz, przywróć ocenę na 0, żeby nie zapisywać jej na stałe
	beer.Rate = 0

	data := ItemPageData{
		Beer:     beer,
		RecBeers: recbeers,
	}

	tmpl, _ := template.ParseFiles("pages/item.html")
	tmpl.Execute(w, data)
}

func main() {
	// for rates := 0; rates < 10; {
	// 	// wczytaj losowe piwo z bazy i wyświetl
	// 	beer := recsys.GetRandomBeer()
	// 	beer.DisplayInformation(recsys)
	// 	fmt.Print("Ocen piwo (od 1 do 5) lub 0 aby wybrać następne: ")
	// 	// odczytaj z konsoli ocene użytkownika
	// 	var rateText string
	// 	fmt.Scanln(&rateText)
	// 	rate, error := strconv.ParseFloat(rateText, 64)
	// 	// sprawdź poprawność danych
	// 	if error == nil && rate >= 1 && rate <= 5 {
	// 		beer.Rate = rate
	// 		rates++
	// 	} else {
	// 		fmt.Print("Nieprawidłowa ocena\r\n")
	// 	}
	// }
	// // wygeneruj rekomendacje
	// reco := recsys.GetRecommendation()
	// fmt.Print("\r\n\r\nRekomendowane piwa: \r\n")
	// reco[0].DisplayInformation(recsys)
	// reco[1].DisplayInformation(recsys)
	// reco[2].DisplayInformation(recsys)

	http.HandleFunc("/piwa/", piwaFunc)
	http.HandleFunc("/reko/", rekoFunc)
	http.HandleFunc("/piwo/{id}", itemFunc)
	http.ListenAndServe("localhost:8080", nil)
}
