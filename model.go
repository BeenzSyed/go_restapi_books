package main

import (
	"database/sql"
	"fmt"
	"log"
)

//Book Struct (Model)
type Book struct {
	ID     int     `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

//Author Struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func (b *Book) getBook(db *sql.DB) error {
	return db.QueryRow("SELECT isbn, title, author FROM books WHERE id=$1",
		b.ID).Scan(&b.Isbn, &b.Title, &b.Author)
}

func (b *Book) updateBook(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE books SET isbn=$1, title=$2, author= $3 WHERE id=$4",
			b.Isbn, b.Title, b.Author, b.ID)
	return err
}

func (b *Book) deleteBook(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM books WHERE id=$1", b.ID)

	return err
}

func (b *Book) createBook(db *sql.DB) error {
	fmt.Printf("%s %s %s", b.Isbn, b.Title, *b.Author)
	_, err := db.Exec("INSERT INTO books (isbn, title, author) VALUES ($1, $2, $3);", b.Isbn, b.Title, b.Author)
	return err
	//return errors.New("Not implemented")
}

func getBooks(db *sql.DB, start, count int) ([]Book, error) {
	rows, err := db.Query("SELECT * FROM books;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	books := []Book{}

	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Isbn, &b.Title, &b.Author); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("%d %s %s %s", id, isbn, title, author)
		books = append(books, b)
	}
	return books, err
}
