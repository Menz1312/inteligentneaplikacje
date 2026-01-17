package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

type Cuboid struct {
	Name string
	L    float64
	W    float64
	H    float64
}

func getCuboids() []Cuboid {
	return []Cuboid{
		{"CuboidA", 1, 2, 3},
		{"CuboidB", 4, 5, 1},
		{"CuboidC", 2, 3, 4},
		{"CuboidD", 5, 1, 2},
		{"CuboidE", 3, 4, 5},
	}
}

func getCuboidById(cuboids []Cuboid, id int) Cuboid {
	return cuboids[id]
}

func randCuboid(cuboids []Cuboid) Cuboid {
	return getCuboidById(cuboids, rand.IntN(5))
}

// Wyświetlanie wszystkich obiektów HTML
func indexFunc(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("pages/index.html")
	tmpl.Execute(w, getCuboids())
}

// Wyświetlanie losowego obiektu HTML
func randFunc(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("pages/rand.html")
	tmpl.Execute(w, randCuboid(getCuboids()))
}

// Wyświetlanie wybranego id w formacie json
func jsonIdFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cuboidId, _ := strconv.Atoi(r.PathValue("id"))
	cuboid := getCuboidById(getCuboids(), cuboidId)
	data, _ := json.Marshal(cuboid)
	w.Write(data)
}

// Wyświetlanie wszystkich obiektów json
func jsonFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cuboids := getCuboids()
	data, _ := json.Marshal(cuboids)
	w.Write(data)
}

// Wielowątkowość
var ch = make(chan string)

func write(n int, m string) {
	for i := 0; i < n; i++ {
		ch <- m
		time.Sleep(time.Second * time.Duration(1))
	}
}

func main() {
	fmt.Println(getCuboids())

	// //Serwer HTTP
	// http.HandleFunc("/", indexFunc)
	// http.HandleFunc("/rand/", randFunc)
	// http.HandleFunc("/json/{id}", jsonIdFunc)
	// http.HandleFunc("/json/", jsonFunc)
	// http.ListenAndServe("localhost:8080", nil)

	//Wielowątkowość
	go write(10, "Hi")
	go write(8, "Hello")
	for m := range ch {
		fmt.Println(m)
	}

}
