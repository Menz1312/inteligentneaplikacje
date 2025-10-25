package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Rate struct {
	Movie int     `json:"movieid"`
	Rate  float64 `json:"rate"`
}

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Rates []Rate `json:"rates"`
}

type Movie struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Oryginal string `json:"oryginal"`
	Year     int    `json:"year"`
	Genre    string `json:"genre"`
}

func GetUser(host string, uid int) *User {
	userData := User{}
	client := http.Client{Timeout: 5 * time.Second}
	url := host + "user/" + strconv.Itoa(uid)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Błąd przy tworzeniu żądania: %v\n", err)
		return nil
	}

	execute, err := client.Do(request)
	if err != nil {
		fmt.Printf("Błąd przy wykonywaniu żądania dla użytkownika %d: %v\n", uid, err)
		return nil
	}
	defer execute.Body.Close()

	body, err := io.ReadAll(execute.Body)
	if err != nil {
		fmt.Printf("Błąd przy czytaniu odpowiedzi: %v\n", err)
		return nil
	}

	err = json.Unmarshal(body, &userData)
	if err != nil {
		fmt.Printf("Błąd przy parsowaniu JSON: %v\n", err)
		return nil
	}

	return &userData
}

// Funkcja do pobierania liczby użytkowników
func GetUserCount(host string) int {
	client := http.Client{Timeout: 5 * time.Second}
	url := host + "usercount/"
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Błąd przy tworzeniu żądania: %v\n", err)
		return 0
	}

	execute, err := client.Do(request)
	if err != nil {
		fmt.Printf("Błąd przy wykonywaniu żądania: %v\n", err)
		return 0
	}
	defer execute.Body.Close()

	body, err := io.ReadAll(execute.Body)
	if err != nil {
		fmt.Printf("Błąd przy czytaniu odpowiedzi: %v\n", err)
		return 0
	}

	var count int
	err = json.Unmarshal(body, &count)
	if err != nil {
		fmt.Printf("Błąd przy parsowaniu JSON: %v\n", err)
		return 0
	}

	return count
}

// Funkcja do pobierania danych wszystkich filmów
func GetMovies(host string) []Movie {
	client := http.Client{Timeout: 5 * time.Second}
	url := host + "movies/"
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Błąd przy tworzeniu żądania: %v\n", err)
		return nil
	}

	execute, err := client.Do(request)
	if err != nil {
		fmt.Printf("Błąd przy wykonywaniu żądania: %v\n", err)
		return nil
	}
	defer execute.Body.Close()

	body, err := io.ReadAll(execute.Body)
	if err != nil {
		fmt.Printf("Błąd przy czytaniu odpowiedzi: %v\n", err)
		return nil
	}

	var movies []Movie
	err = json.Unmarshal(body, &movies)
	if err != nil {
		fmt.Printf("Błąd przy parsowaniu JSON: %v\n", err)
		return nil
	}

	return movies
}

// Funkcja obliczająca różnicę między dwoma użytkownikami
func ObliczRoznice(userA, userB *User) (float64, int) {
	// Krok 1: Znajdź wspólne ocenione obiekty (ocena != 0)
	wspolneOceny := make(map[int]struct {
		rateA float64
		rateB float64
	})

	// Tworzymy mapę ocen użytkownika A (pomijając oceny 0)
	ocenyA := make(map[int]float64)
	for _, rate := range userA.Rates {
		if rate.Rate != 0 {
			ocenyA[rate.Movie] = rate.Rate
		}
	}

	// Tworzymy mapę ocen użytkownika B (pomijając oceny 0)
	ocenyB := make(map[int]float64)
	for _, rate := range userB.Rates {
		if rate.Rate != 0 {
			ocenyB[rate.Movie] = rate.Rate
		}
	}

	// Znajdujemy wspólne filmy ocenione przez obu użytkowników
	for movieID, rateA := range ocenyA {
		if rateB, exists := ocenyB[movieID]; exists {
			wspolneOceny[movieID] = struct {
				rateA float64
				rateB float64
			}{rateA, rateB}
		}
	}

	// Obliczamy C - liczba takich samych obiektów ocenionych przez obu użytkowników
	C := len(wspolneOceny)

	// Jeśli nie ma wspólnych filmów, zwracamy maksymalną różnicę
	if C == 0 {
		return 1.0, 0
	}

	// Obliczamy D - liczba wszystkich ocenionych obiektów przez obu użytkowników
	D := len(ocenyA) + len(ocenyB) - C // Unikamy podwójnego liczenia wspólnych

	// Obliczamy A
	A := float64(C) / float64(D)

	// Obliczamy B - różnica w ocenie tych samych obiektów
	var sumaRoznic float64
	for _, oceny := range wspolneOceny {
		roznica := math.Abs(oceny.rateA - oceny.rateB)
		sumaRoznic += roznica
	}

	B := sumaRoznic / float64(C)

	// Obliczamy finalną różnicę
	roznica := (0.5 + A) * B

	return roznica, C
}

// Funkcja do wyświetlania informacji o użytkowniku
func WyswietlUzytkownika(user *User, movies []Movie) {
	fmt.Printf("=== Użytkownik: %s (ID: %d) ===\n", user.Name, user.Id)

	// Liczba niezerowych ocen
	niezeroweOceny := 0
	for _, rate := range user.Rates {
		if rate.Rate != 0 {
			niezeroweOceny++
		}
	}
	fmt.Printf("Liczba ocenionych filmów: %d\n", niezeroweOceny)

	// Wyświetl kilka ocenionych filmów
	// fmt.Println("Przykładowe ocenione filmy:")
	// displayed := 0
	// for _, rate := range user.Rates {
	// 	if rate.Rate != 0 && displayed < 5 {
	// 		for _, movie := range movies {
	// 			if movie.Id == rate.Movie {
	// 				fmt.Printf("  - %s (%d): %.1f/10\n", movie.Title, movie.Year, rate.Rate)
	// 				displayed++
	// 				break
	// 			}
	// 		}
	// 	}
	// }
	fmt.Println()
}

// Funkcja do wyświetlania wysokich ocen użytkownika
func WyswietlWysokieOceny(user *User, movies []Movie, progOceny float64) {
	fmt.Printf("=== Filmy z wysokimi ocenami użytkownika %s (ID: %d) ===\n", user.Name, user.Id)
	fmt.Printf("Oceny >= %.1f:\n", progOceny)

	licznik := 0
	for _, rate := range user.Rates {
		if rate.Rate >= progOceny {
			for _, movie := range movies {
				if movie.Id == rate.Movie {
					fmt.Printf("  - %s (%d) | Ocena: %.1f/10 | Gatunek: %s\n",
						movie.Title, movie.Year, rate.Rate, movie.Genre)
					licznik++
					break
				}
			}
		}
	}

	if licznik == 0 {
		fmt.Println("  Brak filmów z tak wysokimi ocenami")
	}
	fmt.Printf("Znaleziono %d filmów z wysokimi ocenami\n\n", licznik)
}

// Funkcja znajdująca najbardziej podobnego użytkownika (z co najmniej jednym wspólnym filmem)
func ZnajdzNajbardziejPodobnegoUzytkownika(host string, myID int, userCount int, movies []Movie) (int, *User, float64, int, map[int]float64) {
	fmt.Println("\n=== WYSZUKIWANIE NAJBARDZIEJ PODOBNEGO UŻYTKOWNIKA ===")
	fmt.Println("Szukam użytkownika z co najmniej jednym wspólnym ocenionym filmem...")

	myUser := GetUser(host, myID)
	if myUser == nil {
		fmt.Printf("Błąd: Nie udało się pobrać danych użytkownika %d\n", myID)
		return -1, nil, 0, 0, nil
	}

	najmniejszaRoznica := 1.0 // Startujemy z maksymalną różnicą
	najlepszyUzytkownikID := -1
	var najlepszyUzytkownik *User
	najwiecejWspolnych := 0
	znalezionoKandydatow := 0

	// Mapa do przechowywania wyników wszystkich użytkowników z wspólnymi filmami
	wynikiWszystkich := make(map[int]float64)

	// Sprawdź wszystkich użytkowników (oprócz siebie)
	for i := 0; i < userCount; i++ {
		if i != myID {
			user := GetUser(host, i)
			if user != nil {
				roznica, liczbaWspolnych := ObliczRoznice(myUser, user)

				// Rozważamy tylko użytkowników z co najmniej jednym wspólnym filmem
				if liczbaWspolnych > 0 {
					// Zapisz wynik tego użytkownika
					wynikiWszystkich[i] = roznica
					znalezionoKandydatow++

					// Wyświetl wynik dla tego użytkownika
					fmt.Printf("  Użytkownik %s (ID: %d) - różnica: %.4f, wspólne filmy: %d\n",
						user.Name, i, roznica, liczbaWspolnych)

					// Szukamy najmniejszej różnicy
					if roznica < najmniejszaRoznica {
						najmniejszaRoznica = roznica
						najlepszyUzytkownikID = i
						najlepszyUzytkownik = user
						najwiecejWspolnych = liczbaWspolnych
					}
				}
			}
		}
		// Progress indicator dla dużych zbiorów
		if i%10 == 0 && i > 0 {
			fmt.Printf("  Przetworzono %d/%d użytkowników, znaleziono %d kandydatów\n", i, userCount, znalezionoKandydatow)
		}
	}

	if znalezionoKandydatow == 0 {
		fmt.Println("Nie znaleziono żadnego użytkownika z wspólnymi ocenionymi filmami")
		return -1, nil, 0, 0, nil
	}

	fmt.Printf("\nZnaleziono %d użytkowników z wspólnymi filmami\n", znalezionoKandydatow)
	return najlepszyUzytkownikID, najlepszyUzytkownik, najmniejszaRoznica, najwiecejWspolnych, wynikiWszystkich
}

// Funkcja wyświetlająca rekomendacje - filmy których nie oglądałem z wysokimi ocenami podobnego użytkownika
func WyswietlRekomendacje(myUser *User, podobnyUser *User, movies []Movie, progOceny float64) {
	fmt.Printf("=== REKOMENDACJE - filmy których nie oglądałem z wysokimi ocenami użytkownika %s ===\n", podobnyUser.Name)
	fmt.Printf("Oceny >= %.1f:\n", progOceny)

	// Tworzymy mapę filmów które ja oceniłem (nawet 0 oznacza że widziałem)
	mojeFilmy := make(map[int]bool)
	for _, rate := range myUser.Rates {
		mojeFilmy[rate.Movie] = true
	}

	// Szukamy filmów które podobny użytkownik ocenił wysoko, a ja nie oglądałem
	licznik := 0
	for _, rate := range podobnyUser.Rates {
		if rate.Rate >= progOceny && !mojeFilmy[rate.Movie] {
			for _, movie := range movies {
				if movie.Id == rate.Movie {
					fmt.Printf("  - %s (%d) | Ocena: %.1f/10 | Gatunek: %s\n",
						movie.Title, movie.Year, rate.Rate, movie.Genre)
					licznik++
					break
				}
			}
		}
	}

	if licznik == 0 {
		fmt.Printf("  Brak rekomendacji - wszystkie filmy z wysokimi ocenami użytkownika %s zostały już przez Ciebie ocenione\n", podobnyUser.Name)
	} else {
		fmt.Printf("Znaleziono %d rekomendowanych filmów\n\n", licznik)
	}
}

// Funkcja zapisująca raport do pliku
func ZapiszRaport(user *User, movies []Movie, progOceny float64, myID int, roznica float64, myUser *User, wspolneFilmy int, rekomendacje []Movie) {
	file, err := os.Create("raport_rekomendacji.txt")
	if err != nil {
		fmt.Printf("Błąd przy tworzeniu pliku raportu: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "=== RAPORT SYSTEMU REKOMENDACJI ===\n")
	fmt.Fprintf(file, "Wygenerowano: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "Mój ID: %d\n", myID)
	fmt.Fprintf(file, "Moja nazwa: %s\n", myUser.Name)
	fmt.Fprintf(file, "Najbardziej podobny użytkownik: %s (ID: %d)\n", user.Name, user.Id)
	fmt.Fprintf(file, "RÓŻNICA: %.4f\n", roznica)
	fmt.Fprintf(file, "Liczba wspólnych ocenionych filmów: %d\n\n", wspolneFilmy)

	fmt.Fprintf(file, "=== REKOMENDOWANE FILMY (oceny >= %.1f) ===\n", progOceny)

	if len(rekomendacje) == 0 {
		fmt.Fprintf(file, "Brak rekomendowanych filmów\n")
	} else {
		for _, movie := range rekomendacje {
			fmt.Fprintf(file, "%-40s | Rok: %-4d | Gatunek: %s\n",
				movie.Title, movie.Year, movie.Genre)
		}
		fmt.Fprintf(file, "\nŁącznie znaleziono %d rekomendowanych filmów\n", len(rekomendacje))
	}

	fmt.Printf("Raport zapisany do pliku: raport_rekomendacji.txt\n")
}

// Funkcja do pobierania rekomendowanych filmów
func PobierzRekomendacje(myUser *User, podobnyUser *User, movies []Movie, progOceny float64) []Movie {
	// Tworzymy mapę filmów które ja oceniłem (nawet 0 oznacza że widziałem)
	mojeFilmy := make(map[int]bool)
	for _, rate := range myUser.Rates {
		mojeFilmy[rate.Movie] = true
	}

	// Szukamy filmów które podobny użytkownik ocenił wysoko, a ja nie oglądałem
	var rekomendacje []Movie
	for _, rate := range podobnyUser.Rates {
		if rate.Rate >= progOceny && !mojeFilmy[rate.Movie] {
			for _, movie := range movies {
				if movie.Id == rate.Movie {
					rekomendacje = append(rekomendacje, movie)
					break
				}
			}
		}
	}

	return rekomendacje
}

func main() {
	// Adres serwera - zgodnie z podanym przez Ciebie
	host := "http://172.16.10.115:8080/"
	myID := 7

	fmt.Println("=== System rekomendacji ===")
	fmt.Printf("Łączenie z serwerem: %s\n", host)
	fmt.Printf("Mój ID: %d\n\n", myID)

	// Pobierz dane wszystkich filmów
	fmt.Println("1. Pobieranie danych filmów...")
	movies := GetMovies(host)
	if movies == nil {
		fmt.Println("Błąd: Nie udało się pobrać danych filmów")
		return
	}
	fmt.Printf("   Pobrano dane %d filmów\n", len(movies))

	// Pobierz liczbę użytkowników
	fmt.Println("\n2. Pobieranie liczby użytkowników...")
	userCount := GetUserCount(host)
	fmt.Printf("   Liczba użytkowników w systemie: %d\n", userCount)

	if userCount < 2 {
		fmt.Println("Błąd: Zbyt mało użytkowników w systemie")
		return
	}

	// Pobierz moje dane
	fmt.Println("\n3. Pobieranie moich danych...")
	myUser := GetUser(host, myID)
	if myUser == nil {
		fmt.Printf("Błąd: Nie udało się pobrać danych użytkownika %d\n", myID)
		return
	}

	// Wyświetl informacje o mnie
	fmt.Println("\n" + strings.Repeat("=", 60))
	WyswietlUzytkownika(myUser, movies)

	// Znajdź najbardziej podobnego użytkownika (z co najmniej jednym wspólnym filmem)
	najlepszyID, najlepszyUser, najmniejszaRoznica, wspolneFilmy, _ := ZnajdzNajbardziejPodobnegoUzytkownika(host, myID, userCount, movies)

	if najlepszyUser != nil {
		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Printf("=== NAJBARDZIEJ PODOBNY UŻYTKOWNIK ===\n")
		fmt.Printf("Użytkownik: %s (ID: %d)\n", najlepszyUser.Name, najlepszyID)
		fmt.Printf("RÓŻNICA: %.4f\n", najmniejszaRoznica)
		fmt.Printf("Liczba wspólnych ocenionych filmów: %d\n", wspolneFilmy)

		// Wyświetl informacje o najbardziej podobnym użytkowniku
		WyswietlUzytkownika(najlepszyUser, movies)

		// Wyświetl wspólne filmy z ocenami
		fmt.Println("=== WSPÓLNE FILMY Z OCENAMI ===")
		WyswietlWspolneFilmy(myUser, najlepszyUser, movies)

		// Wyświetl wysokie oceny najbardziej podobnego użytkownika
		progOceny := 8.0 // Próg wysokiej oceny
		WyswietlWysokieOceny(najlepszyUser, movies, progOceny)

		// Wyświetl rekomendacje - filmy których nie oglądałem
		WyswietlRekomendacje(myUser, najlepszyUser, movies, progOceny)

		// Pobierz rekomendacje do zapisu w raporcie
		// rekomendacje := PobierzRekomendacje(myUser, najlepszyUser, movies, progOceny)

		// Zapisz raport do pliku
		// ZapiszRaport(najlepszyUser, movies, progOceny, myID, najmniejszaRoznica, myUser, wspolneFilmy, rekomendacje)

	} else {
		fmt.Println("\nNie znaleziono użytkownika z wspólnymi ocenionymi filmami")
		fmt.Println("Spróbuj zmienić próg ocen lub sprawdź czy masz ocenione jakieś filmy")
	}
}

// Funkcja do wyświetlania wspólnych filmów z ocenami obu użytkowników
func WyswietlWspolneFilmy(userA, userB *User, movies []Movie) {
	// Tworzymy mapy ocen obu użytkowników (pomijając oceny 0)
	ocenyA := make(map[int]float64)
	for _, rate := range userA.Rates {
		if rate.Rate != 0 {
			ocenyA[rate.Movie] = rate.Rate
		}
	}

	ocenyB := make(map[int]float64)
	for _, rate := range userB.Rates {
		if rate.Rate != 0 {
			ocenyB[rate.Movie] = rate.Rate
		}
	}

	// Znajdujemy wspólne filmy
	wspolneFilmy := 0
	for movieID, rateA := range ocenyA {
		if rateB, exists := ocenyB[movieID]; exists {
			// Znajdź informacje o filmie
			for _, movie := range movies {
				if movie.Id == movieID {
					roznica := math.Abs(rateA - rateB)
					fmt.Printf("  - %s: Ja=%.1f, %s=%.1f (różnica=%.1f)\n",
						movie.Title, rateA, userB.Name, rateB, roznica)
					wspolneFilmy++
					break
				}
			}
		}
	}

	if wspolneFilmy == 0 {
		fmt.Println("  Brak wspólnych ocenionych filmów")
	} else {
		fmt.Printf("Łącznie wspólnych filmów: %d\n", wspolneFilmy)
	}
	fmt.Println()
}
