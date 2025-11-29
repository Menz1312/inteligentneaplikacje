package web

import (
	"encoding/json"
	"fmt"
	"labmobile/fuzzy"
	"net/http"
	"sort"
	"text/template"
)

var phones []*Phone

func indexFunc(w http.ResponseWriter, r *http.Request) {
	tl, _ := template.ParseFiles("pages/index.html")
	tl.Execute(w, nil)
}

func searchFunc(w http.ResponseWriter, r *http.Request) {
	query := r.PathValue("query")
	data := ParseData(query)

	// tworzenie zbiorów rozmytych
	fsetScreen := data.GetSet("screen")
	fsetBattery := data.GetSet("battery")
	fsetWeight := data.GetSet("weight")

	// obliczenie dopasowania rozmytego (value) dla wszystkich telefonów
	for _, phone := range phones {
		f1 := fsetScreen.Calculate(phone.Screen)
		f2 := fsetBattery.Calculate(phone.Battery)
		f3 := fsetWeight.Calculate(phone.Weight)
		phone.Value = fuzzy.Agregate(f1, f2, f3)
	}

	// posortowanie telefonów
	switch data.SortingMethod {
	case 1: // według dopasowania
		sort.Slice(phones, func(a, b int) bool {
			return phones[a].Value > phones[b].Value
		})
	case 2: // według ceny
		sort.Slice(phones, func(a, b int) bool {
			return phones[a].PricePLN > phones[b].PricePLN
		})
	}

	// zamiana wyników na JSON i zwrócenie
	phdata, err := json.Marshal(phones)
	if err != nil {
		http.Error(w, "json parse error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(phdata)
}

func StartServer() {
	phones = LoadData()
	fs := http.FileServer(http.Dir("./files"))
	http.Handle("/files/", http.StripPrefix("/files/", fs))
	http.HandleFunc("/", indexFunc)
	http.HandleFunc("/search/{query}", searchFunc)
	fmt.Println("Serwer startuje na: http://localhost:8080/")
	http.ListenAndServe("localhost:8080", nil)
}
