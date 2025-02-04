package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type UsdToBrl struct {
	USDBRL struct {
		Bid string `json:"bid"`
	}
}

var httpGet = http.Get

func main() {
	db, err := sql.Open("sqlite3", "../app.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStmt := `
  CREATE TABLE IF NOT EXISTS usdPrice (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    usdPrice DOUBLE
  );`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("Error creating table: %s", err)
	}
	log.Println("Table created successfully")

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		convertDollarToBrl(w, r, db)
	})

	log.Println("Server started on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func convertDollarToBrl(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	resp, err := httpGet("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	var usdToBrl UsdToBrl
	err = json.Unmarshal(body, &usdToBrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	bid, err := strconv.ParseFloat(usdToBrl.USDBRL.Bid, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	_, err = db.ExecContext(ctx, "INSERT INTO usdPrice (usdPrice) VALUES (?)", bid)
	if err != nil {
		panic(err)
	}
	log.Println("Successfully inserted into database")

	log.Println("Received API Response:", usdToBrl.USDBRL.Bid)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Dolar: " + usdToBrl.USDBRL.Bid))
}
