package main

import (
	"fmt"
	"lab4/db"
	"lab4/dbase"
	"lab4/kmeans"
	"log"
)

func main() {
	beers := dbase.GetBeersAll()
	const liczbaGrup = 10
	clusters := kmeans.Kmeans(beers, liczbaGrup)

	fmt.Println("Wynik grupowania (liczba piw w każdej grupie):")
	for i, cluster := range clusters {
		fmt.Printf("Grupa %d: %d piw\n", i, len(cluster.Beers))
	}

	_, err := db.Exec("DELETE FROM grpa")
	if err != nil {
		log.Fatal("Błąd przy czyszczeniu tabeli 'grpa': ", err)
	}

	fmt.Println("\nZapisywanie wyników grupowania do bazy danych...")

	for _, cluster := range clusters {
		numerGrupy := cluster.Grupa

		for _, beer := range cluster.Beers {
			idPiwa := beer.Id
			_, err := db.Exec("INSERT INTO grpa (id, grupa) VALUES (?, ?)", idPiwa, numerGrupy)
			if err != nil {
				fmt.Printf("Błąd przy wstawianiu (piwo: %d, grupa: %d): %s\n", idPiwa, numerGrupy, err)
			}
		}
	}
	fmt.Println("Zakończono zapisywanie.")
}
