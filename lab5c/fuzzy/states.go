package fuzzy

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (f *FuzzySystem) Save(name string) error {
	file, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	for _, rule := range f.rules {
		for _, set := range rule.sets {
			fmt.Fprintf(file, "%e\t%e\t", set.x, set.k)
		}
		fmt.Fprintf(file, "%e\n", rule.output)
	}
	return nil
}

func Load(name string) (*FuzzySystem, error) {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()
	scan := bufio.NewScanner(file)
	rules := []*FuzzyRule{}
samples:
	for scan.Scan() {
		data := strings.Split(scan.Text(), "\t")
		sets := []*FuzzySet{}
		for i := 0; i < len(data)-1; i += 2 {
			x, e1 := strconv.ParseFloat(data[i], 64)
			k, e2 := strconv.ParseFloat(data[i+1], 64)
			if e1 == nil && e2 == nil {
				sets = append(sets, &FuzzySet{x: x, k: k})
			} else {
				fmt.Println(e1, e2)
				continue samples
			}
		}
		output, e3 := strconv.ParseFloat(data[len(data)-1], 64)
		if e3 != nil {
			fmt.Println(e3)
		} else {
			rules = append(rules, &FuzzyRule{sets: sets, output: output})
		}
	}
	return &FuzzySystem{rules: rules}, nil
}
