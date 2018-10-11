package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"))

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS books
(
id SERIAL,
isbn INT NOT NULL,
title TEXT NOT NULL,
author TEXT NOT NULL, 
PRIMARY KEY (id)
)`

func clearTable() {
	a.DB.Exec("DELETE FROM books")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentBook(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/book/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Book not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Book not found'. Got '%s'", m["error"])
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/books", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestCreateBook(t *testing.T) {
	clearTable()

	payload := []byte(`{"isbn": "12345","title": "test book","author": "Sabeen Syed"}`)

	req, _ := http.NewRequest("POST", "/book", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["title"] != "test book" {
		t.Errorf("Expected book title to be 'test book'. Got '%v'", m["title"])
	}

	if m["isbn"] != "12345" {
		t.Errorf("Expected book isbn to be '12345'. Got '%v'", m["isbn"])
	}

	if m["author"] != "Sabeen Syed" {
		t.Errorf("Expected book author to be 'Sabeen Syed'. Got '%v'", m["author"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 0.0 {
		t.Errorf("Expected book ID to be '0'. Got '%v'", m["id"])
	}
}

func TestGetBook(t *testing.T) {
	clearTable()
	addBooks(1)

	req, _ := http.NewRequest("GET", "/book/35", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	fmt.Println(response.Body)
}

func addBooks(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO books (isbn, title, author) VALUES (12345, 'Book One', 'Sabeen');")
	}
}
