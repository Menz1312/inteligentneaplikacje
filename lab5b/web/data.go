package web

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func LoadData() []*Phone {
	// baza danych z: https://www.kaggle.com/datasets/abdulmalik1518/mobiles-dataset-2025
	file, err := os.Open("files/mobiles2025.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan() // skip header
	phones := []*Phone{}
	for scanner.Scan() {
		data := strings.Split(strings.ReplaceAll(scanner.Text(), "\"", "'"), ",")
		weight, err1 := strconv.ParseFloat(strings.ReplaceAll(data[2], "g", ""), 64)
		ram, err2 := strconv.ParseFloat(strings.ReplaceAll(data[3], "GB", ""), 64)
		camf := strings.ReplaceAll(data[4], "MP", "")
		camfd := strings.Split(strings.ReplaceAll(camf, " ", ""), "+")
		maxf := 0.0
		for i := range camfd {
			tmp, err := strconv.ParseFloat(camfd[i], 64)
			if err == nil {
				if tmp > maxf {
					maxf = tmp
				}
			} else {
				fmt.Println(camfd, err)
			}
		}
		camb := strings.ReplaceAll(data[5], "MP", "")
		cambd := strings.Split(strings.ReplaceAll(camb, " ", ""), "+")
		maxb := 0.0
		for i := range cambd {
			tmp, err := strconv.ParseFloat(cambd[i], 64)
			if err == nil {
				if tmp > maxb {
					maxb = tmp
				}
			} else {
				fmt.Println(cambd, err)
			}
		}
		battery, err3 := strconv.ParseFloat(strings.ReplaceAll(data[7], "mAh", ""), 64)
		screen, err4 := strconv.ParseFloat(strings.ReplaceAll(data[8], " inches", ""), 64)
		usd, err5 := strconv.ParseFloat(strings.ReplaceAll(data[12], "USD ", ""), 64)
		year, err6 := strconv.ParseFloat(strings.ReplaceAll(data[14], "USD ", ""), 64)
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
			fmt.Println(data)
			fmt.Println("weight", err1)
			fmt.Println("ram", err2)
			fmt.Println("battery", err3)
			fmt.Println("screen", err4)
			fmt.Println("usd", err5)
			fmt.Println("year", err6)
			continue
		} else {
			phones = append(phones, &Phone{
				Company:   data[0],
				Model:     data[1],
				Weight:    weight,
				RAM:       ram,
				CamFront:  maxf,
				CamBack:   maxb,
				Processor: data[6],
				Battery:   battery,
				Screen:    screen,
				PriceUSD:  usd,
				PricePLN:  math.Round(usd * 3.65),
				Year:      year,
			})
		}
	}
	fmt.Printf("Wczytano %d telefonów\n", len(phones))
	return phones
}
