package main

import (
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInsertIntoTable(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	sqlStmt := `
    CREATE TABLE IF NOT EXISTS usdPrice (
      id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
      usdPrice DOUBLE NOT NULL
    );`

	_, err = db.Exec(sqlStmt)
	assert.NoError(t, err)

	mockValue := 5.30
	start := time.Now()

	_, err = db.Exec("INSERT INTO usdPrice (usdPrice) VALUES (?)", mockValue)
	assert.NoError(t, err)

	timeTaken := time.Since(start)
	assert.LessOrEqual(t, timeTaken.Milliseconds(), int64(10))

	_, err = db.Exec("INSERT INTO usdPrice (usdPrice) VALUES (?)", nil)
	assert.Error(t, err)

	var result float64
	err = db.QueryRow("SELECT usdPrice FROM usdPrice LIMIT 1").Scan(&result)
	assert.NoError(t, err, "Retrieving inserted value should not fail")

	assert.Equal(t, mockValue, result, "Inserted value should match the retrieved one")
}

const mockAPIResponse = `{"USDBRL": {"bid": "5.42"}}`

func TestFetchUsdToBrl(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	sqlStmt := `
    CREATE TABLE IF NOT EXISTS usdPrice (
      id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
      usdPrice DOUBLE NOT NULL
    );`

	_, err = db.Exec(sqlStmt)
	assert.NoError(t, err)

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockAPIResponse))
	}))
	defer mockServer.Close()

	oldHttpGet := httpGet
	httpGet = func(url string) (*http.Response, error) {
		return http.Get(mockServer.URL)
	}
	defer func() { httpGet = oldHttpGet }()

	req := httptest.NewRequest("GET", "/cotacao", nil)
	rec := httptest.NewRecorder()
	convertDollarToBrl(rec, req, db)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	assert.Contains(t, string(body), "5.42")
}
