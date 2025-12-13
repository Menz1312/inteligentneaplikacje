package main

import (
	"fmt"
	"net/http"
	"time"
)

// globalny kanal
var ch = make(chan string)

func main() {
	// endpointy
	http.HandleFunc("/swiat/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "swiat")
		ch <- "swiat"
	})

	http.HandleFunc("/sport/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "sport")
		ch <- "sport"
	})

	http.HandleFunc("/polityka/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "polityka")
		ch <- "polityka"
	})

	go http.ListenAndServe(":8080", nil)
	fmt.Println("Serwer uruchomiony. Statystyki co 30 sekund...")

	stats := make(map[string]int)

	tick := time.NewTicker(30 * time.Second)

	// endless loop
	for {
		select {
		case msg := <-ch:
			stats[msg]++ // Zwiększamy licznik w mapie
			fmt.Println(time.Now().Format("15:04:05"), "– odwiedzono:", msg)

		case <-tick.C:
			fmt.Println("\n===STATYSTYKI===")
			for k, v := range stats {
				fmt.Printf("Endpoint /%s/: %d odwiedzin\n", k, v)
			}
			fmt.Println("")
		}
	}
}
