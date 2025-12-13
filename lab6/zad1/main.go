package main

import (
	"fmt"
	"sync"
	"time"
)

var sg sync.WaitGroup
var mx sync.Mutex        // Mutex
var liczbyPierwsze []int // Globalna tablica

func pierwszeMutex(od, do int) {
	defer sg.Done()

szukaj:
	for liczba := od; liczba < do; liczba++ {
		if liczba < 2 {
			continue
		}
		for dzielnik := 2; dzielnik < liczba; dzielnik++ {
			if liczba%dzielnik == 0 {
				continue szukaj
			}
		}

		// zapis do zmiennej globalnej
		mx.Lock()
		liczbyPierwsze = append(liczbyPierwsze, liczba)
		mx.Unlock()
	}
}

func main() {
	liczbyPierwsze = make([]int, 0, 20000) // wstepna alokacja dla optymalizacji

	fmt.Println("Start obliczeń z Mutexem...")
	start := time.Now()

	sg.Add(4)
	// Uruchomienie watków
	go pierwszeMutex(0, 80000)
	go pierwszeMutex(80000, 130000)
	go pierwszeMutex(130000, 170000)
	go pierwszeMutex(170000, 200000)

	sg.Wait()
	czas := time.Since(start)

	fmt.Printf("Czas obliczeń (Mutex): %v\n", czas)
	fmt.Printf("Znaleziono %d liczb pierwszych.\n", len(liczbyPierwsze))
}
