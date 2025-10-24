package knn

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Beer struct {
	Id    int     `json:"id"`   // id piwa
	Abv   float64 `json:"abv"`  // poziom abv
	Ibu   float64 `json:"ibu"`  // goryczka
	Style int     `json:"-"`    // styl
	Name  string  `json:"name"` // nazwa
	Rate  float64 `json:"-"`    // ocena
	Estim float64 `json:"-"`    // estymowana ocena
	Simi  []*Beer `json:"-"`    // podobne piwa
}

type Beers struct {
	beers  []Beer
	styles Styles
}

func (b *Beer) DisplayInformation(knn *Knn) {
	fmt.Print("%--------------------------------------------------%\r\n")
	fmt.Printf("| Nazwa: %s\r\n", b.Name)
	fmt.Printf("| Alkohol: %0.1f\t Goryczka: %0.1f\r\n", b.Abv*100.0, b.Ibu)
	fmt.Printf("| Styl: %s\r\n", knn.GetStyles().GetStyleName(b.Style))
	if b.Estim > 0 {
		fmt.Printf("| Przewidywana ocena: %0.1f\r\n", b.Estim)
	}
	fmt.Print("%--------------------------------------------------%\r\n")
}

func (b1 *Beer) Distance(b2 *Beer) float64 {
	var d float64 = 0
	if b1.Style != b2.Style {
		d += 1.0
	}
	d += math.Abs(b1.Abv-b2.Abv) * 5.0
	d += math.Abs(b1.Ibu-b2.Ibu) * 0.01
	return d
}

func (b *Beer) GetAlcohol() string {
	return strconv.FormatFloat(b.Abv*100.0, 'f', 2, 64)
}

func (b *Beer) GetEstim() string {
	return strconv.FormatFloat(b.Estim, 'f', 2, 64)
}

func (b1 *Beer) GetDistance(b2 *Beer) string {
	return strconv.FormatFloat(b1.Distance(b2), 'f', 2, 64)
}

func LoadBeers(name string) *Beers {
	b := Beers{}
	file, error := os.Open(name)
	if error != nil {
		fmt.Println(error.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stringData := scanner.Text()
		data := strings.Split(stringData, ",")
		id, _ := strconv.Atoi(data[0])
		abv, _ := strconv.ParseFloat(data[1], 64)
		ibu, _ := strconv.ParseFloat(data[2], 64)
		style := data[5]
		name := data[4]
		if len(style) < 2 {
			continue
		}
		if ibu == 0.0 {
			continue
		}
		styleId := b.styles.CheckStyle(style)
		b.beers = append(b.beers, Beer{id, abv, ibu, styleId, name, 0, 0, nil})
	}
	return &b
}

func (b *Beer) EstimateRate(bo []*Beer) {
	d1, d2, d3 := math.MaxFloat64, math.MaxFloat64, math.MaxFloat64
	i1, i2, i3 := 0, 0, 0
	for i := range bo {
		tmp := b.Distance(bo[i])
		if tmp < d1 {
			d3 = d2
			d2 = d1
			d1 = tmp
			i3 = i2
			i2 = i1
			i1 = i
		} else if tmp < d2 {
			d3 = d2
			d2 = tmp
			i3 = i2
			i2 = i
		} else if tmp < d3 {
			d3 = tmp
			i3 = i
		}
	}
	b.Estim = (bo[i1].Rate + bo[i2].Rate + bo[i3].Rate) / 3.0
	b.Simi = []*Beer{bo[i1], bo[i2], bo[i3]}
}

func (b *Beers) Recomendation() []*Beer {
	// find rated beers
	beers := []*Beer{}
	for k := range b.beers {
		if b.beers[k].Rate > 0 {
			beers = append(beers, &b.beers[k])
		}
	}
	// calculate estimated rates
	for k := range b.beers {
		if b.beers[k].Rate == 0 {
			b.beers[k].EstimateRate(beers)
		}
	}
	// find best 3 results
	o1, o2, o3 := 0.0, 0.0, 0.0
	i1, i2, i3 := 0, 0, 0
	for k := range b.beers {
		tmp := b.beers[k].Estim
		if tmp > o1 {
			o3 = o2
			o2 = o1
			o1 = tmp
			i3 = i2
			i2 = i1
			i1 = k
		} else if tmp > o2 {
			o3 = o2
			o2 = tmp
			i3 = i2
			i2 = k
		} else if tmp > o3 {
			o3 = tmp
			i3 = k
		}
	}
	// erase rates
	for k := range b.beers {
		if b.beers[k].Rate > 0 {
			b.beers[k].Rate = 0
		}
	}
	// return results
	var recom []*Beer = []*Beer{&b.beers[i1], &b.beers[i2], &b.beers[i3]}
	return recom
}

func (b *Beer) GetSimiliar(bo []Beer) []Beer {
	d1, d2, d3 := math.MaxFloat64, math.MaxFloat64, math.MaxFloat64
	i1, i2, i3 := 0, 0, 0
	for i := range bo {
		if b.Id == bo[i].Id {
			continue
		}
		tmp := b.Distance(&bo[i])
		if tmp < d1 {
			d3 = d2
			d2 = d1
			d1 = tmp
			i3 = i2
			i2 = i1
			i1 = i
		} else if tmp < d2 {
			d3 = d2
			d2 = tmp
			i3 = i2
			i2 = i
		} else if tmp < d3 {
			d3 = tmp
			i3 = i
		}
	}
	return []Beer{bo[i1], bo[i2], bo[i3]}
}
