package apts

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadData(filename string) []Apartment {
	data := make([]Apartment, 0)
	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Errorf("can't find file '%s'", filename))
	}
	defer file.Close()
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		textdata := strings.Split(scan.Text(), "\t")
		if len(textdata) != 11 {
			fmt.Println("Wrong data: ", textdata)
			continue
		}
		apt := Apartment{}
		apt.Year, _ = strconv.Atoi(textdata[0])
		apt.Age, _ = strconv.Atoi(textdata[1])
		apt.Area, _ = strconv.Atoi(textdata[2])
		apt.Floor, _ = strconv.Atoi(textdata[3])
		apt.Parking, _ = strconv.Atoi(textdata[4])
		apt.Bus, _ = strconv.Atoi(textdata[5])
		apt.Metro, _ = strconv.Atoi(textdata[6])
		apt.Location, _ = strconv.Atoi(textdata[7])
		apt.Parks, _ = strconv.Atoi(textdata[8])
		apt.Schools, _ = strconv.Atoi(textdata[9])
		apt.Price, _ = strconv.Atoi(textdata[10])
		data = append(data, apt)
	}
	return data
}
