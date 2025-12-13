package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getcountry(ip, api, field string) string {
	m := map[string]any{}
	resp, err := http.Get(api + ip)
	if err != nil {
		return "Error"
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &m)

	if val, ok := m[field].(string); ok {
		return val
	}
	return "Unknown"
}

func main() {
	ip := "212.87.244.0"

	result := make(chan string)

	go func() {
		result <- getcountry(ip, "http://ip-api.com/json/", "country")
	}()

	go func() {
		result <- getcountry(ip, "https://freeipapi.com/api/json/", "countryName")
	}()

	go func() {
		result <- getcountry(ip, "http://ipwho.is/", "country")
	}()

	najszybszy := <-result

	fmt.Println("Najszybsza odpowiedź:", najszybszy)
}
