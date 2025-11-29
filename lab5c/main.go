package main

import (
	"fmt"
	"labfs/apts"
	"labfs/data"
	"labfs/fuzzy"
)

func main() {
	// wczytanie danych apartamentów
	apts := apts.LoadData("data/apts.txt")

	// konwersja danych na struktury z wejściami i wyjściami
	samples, mins, maxs, inputs := data.ToSamples(apts)

	// utworzenie obiektu creator pozwalającego tworzyć systemy rozmyte (5 zbiorów na wejście, 12 reguł rozmytych)
	creator := fuzzy.TSKConstructor(mins, maxs, inputs, 5, 12)

	// utworzenie systemu rozmytego o zupełnie losowych parametrach
	system := creator.CreateRandom()

	// test systemu na pierwszych 20 apartamentach
	for i := range samples[:20] {
		sysout := system.Predict(samples[i].Inputs)
		datout := samples[i].Output
		fmt.Printf("S#%d - data price: %0.0f, system price: %0.0f\n", i+1, datout, sysout)
	}
}
