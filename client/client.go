package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	client := &http.Client{
		Timeout: time.Millisecond * 300,
	}

	resp, err := client.Get("http://localhost:8080/cotacao")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	println(string(body))

	file, err := os.Create("./cotacao.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(string(body))
	if err != nil {
		panic(err)
	}
	log.Println("Successfully saved to file")
}
