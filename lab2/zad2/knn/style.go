package knn

import "fmt"

type Style struct {
	name string
	cnt  int
}

type Styles struct {
	styles []Style
}

func (s *Styles) CheckStyle(name string) int {
	for i := range s.styles {
		if s.styles[i].name == name {
			s.styles[i].cnt++
			return i
		}
	}
	s.styles = append(s.styles, Style{name, 1})
	return len(s.styles) - 1
}

func (s *Styles) PrintAllStyles() {
	fmt.Println("All styles: ")
	for i, k := range s.styles {
		fmt.Println(i, k.name, k.cnt)
	}
}

func (s *Styles) GetStyleName(id int) string {
	if id >= 0 && id < len(s.styles) {
		return s.styles[id].name
	}
	return "wrong index"
}
