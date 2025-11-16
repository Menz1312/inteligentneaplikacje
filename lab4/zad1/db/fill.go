package db

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func checkDatabase(database string) {
	// check if there are all beers
	beerscount := 0
	_, err := Exec(fmt.Sprintf("USE %s;", database))
	if err != nil {
		goto create
	}
	QueryRow("SELECT COUNT(*) FROM beers;").Scan(&beerscount)
	if beerscount == 0 {
		goto create
	}
	fmt.Printf("%d beers in database\r\n", beerscount)
	return
	// if not create database, tables and everything
create:
	file, err := os.Open("db/data.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	query := ""
	instruction := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		query += line
		if line[len(line)-1] == ';' {
			fmt.Printf("\rDB creation: %d", instruction)
			res, err := Exec(query)
			instruction++
			if err != nil {
				fmt.Println(err, res)
			}
			query = ""
		}
	}
	generateRandomRates(2000)
}

func generateRandomRates(count int) {
	for r := range count {
		fmt.Printf("\rRATES insertion: %d", r)
		var bid, uid int
		err1 := QueryRow("SELECT id FROM users ORDER BY RAND() LIMIT 1;").Scan(&uid)
		err2 := QueryRow("SELECT id FROM beers ORDER BY RAND() LIMIT 1;").Scan(&bid)
		if err1 == nil && err2 == nil {
			var rate float64 = 0.5 * float64(rand.Intn(11))
			_, err := Exec("INSERT INTO rates VALUES (NULL, ?, ?, ?)", uid, bid, rate)
			if err != nil {
				Exec("UPDATE beers SET ratesum = ratesum+?, ratecnt = ratecnt+1 WHERE id = ?", rate, bid)
				Exec("UPDATE beers SET rateavg = (ratesum / ratecnt) WHERE id = ?", bid)
			}
		} else {
			fmt.Println(err1, err2)
		}
	}
}
