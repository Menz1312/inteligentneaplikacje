package web

import (
	"fmt"
	"labmobile/fuzzy"
	"strconv"
	"strings"
)

type Data struct {
	FSet          map[string]fuzzy.FuzzySet
	SortingMethod int
}

func (d *Data) GetSet(key string) *fuzzy.FuzzySet {
	set, ok := d.FSet[key]
	if ok {
		return &set
	}
	keys := []string{}
	for k := range d.FSet {
		keys = append(keys, k)
	}
	fmt.Printf("can't find key '%s' in fuzzy sets, possible keys: %+v\n", key, keys)
	return fuzzy.Null()
}

func ParseData(query string) Data {
	sortingmethod := 1
	fuzzysets := map[string]fuzzy.FuzzySet{}
	for _, v := range strings.Split(query, ";") {
		data := strings.Split(v, ",")
		switch data[0] {
		case "0": // fuzzy set
			if len(data) != 9 {
				fmt.Println("Wrong block data", data)
				continue
			}
			name := data[1]
			min, _ := strconv.Atoi(data[2])
			max, _ := strconv.Atoi(data[3])
			log, _ := strconv.Atoi(data[4])
			x1, _ := strconv.Atoi(data[5])
			x2, _ := strconv.Atoi(data[6])
			x3, _ := strconv.Atoi(data[7])
			x4, _ := strconv.Atoi(data[8])
			fuzzysets[name] = fuzzy.NewFuzzySet(min, max, x1, x2, x3, x4, log == 1)
		default: // sorting method
			method, err := strconv.Atoi(data[0])
			if err != nil {
				fmt.Println("Wrong sorting method", method)
			} else {
				sortingmethod = method
			}
		}
	}

	return Data{
		FSet:          fuzzysets,
		SortingMethod: sortingmethod,
	}
}
