package main

import (
	"fmt"
	"html/template"
	"lab4/db"
	"lab4/dbase"

	// "lab4/kmeans"
	"log"
	"net/http"
	"strconv"
)

var tmpl *template.Template

func beerHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/beer/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Nieprawidłowe ID", http.StatusBadRequest)
		return
	}

	beer := dbase.GetBeer(id)
	if beer.Id == 0 {
		http.Error(w, "Nie znaleziono piwa", http.StatusNotFound)
		return
	}

	var grupa int
	err = db.QueryRow("SELECT grupa FROM grpa WHERE id = ?", id).Scan(&grupa)
	if err != nil {
		// Piwo istnieje, ale nie ma go w tabeli grpa (błąd lub nieuruchomione zad. 1)
		http.Error(w, "Nie znaleziono grupy dla piwa", http.StatusInternalServerError)
		log.Println("Błąd pobierania grupy:", err)
		return
	}

	randoms := dbase.GetBeersFromGroup(grupa)

	beerdata := struct {
		Beer    *dbase.Beer
		Randoms []*dbase.Beer
	}{
		Beer:    beer,
		Randoms: randoms,
	}

	err = tmpl.ExecuteTemplate(w, "beer.html", beerdata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// beers := dbase.GetBeersAll()
	// const liczbaGrup = 10
	// clusters := kmeans.Kmeans(beers, liczbaGrup)

	// fmt.Println("Wynik grupowania (liczba piw w każdej grupie):")
	// for i, cluster := range clusters {
	// 	fmt.Printf("Grupa %d: %d piw\n", i, len(cluster.Beers))
	// }

	// _, err := db.Exec("DELETE FROM grpa")
	// if err != nil {
	// 	log.Fatal("Błąd przy czyszczeniu tabeli 'grpa': ", err)
	// }

	// fmt.Println("\nZapisywanie wyników grupowania do bazy danych...")

	// for _, cluster := range clusters {
	// 	numerGrupy := cluster.Grupa
	// 	for _, beer := range cluster.Beers {
	// 		idPiwa := beer.Id
	// 		_, err := db.Exec("INSERT INTO grpa (id, grupa) VALUES (?, ?)", idPiwa, numerGrupy)
	// 		if err != nil {
	// 			fmt.Printf("Błąd przy wstawianiu (piwo: %d, grupa: %d): %s\n", idPiwa, numerGrupy, err)
	// 		}
	// 	}
	// }

	// fmt.Println("Zakończono zapisywanie.")

	fmt.Println("Serwer uruchomiony na http://localhost:8080")
	fmt.Println("Test: http://localhost:8080/beer/1")

	tmpl = template.Must(template.ParseFiles("pages/beer.html"))

	http.HandleFunc("/beer/", beerHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
