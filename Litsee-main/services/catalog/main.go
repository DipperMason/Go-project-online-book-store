package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	_ "github.com/lib/pq"
)

type Book struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
	Genre  string `json:"genre"`
	Desc   string `json:"description"`
	File   string `json:"file_url"`
}

var (
	db  *sql.DB
	rdb *redis.Client
	ctx = context.Background()
)

func main() {

	var err error

	db, err = sql.Open(
		"postgres",
		"postgres://user:pass@db:5432/books?sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}

	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	http.HandleFunc("/books", getBooks)
	http.HandleFunc("/books/", getBookByID)

	fs := http.FileServer(http.Dir("./files"))
	http.Handle("/files/", http.StripPrefix("/files/", fs))

	log.Println("catalog started :8081")
	http.ListenAndServe(":8081", nil)
}

func getBooks(w http.ResponseWriter, r *http.Request) {

	genre := r.URL.Query().Get("genre")
	q := r.URL.Query().Get("q")

	if genre == "" && q == "" {
		cached, err := rdb.Get(ctx, "books").Result()
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(cached))
			return
		}
	}

	query := "SELECT id, name, author, genre, description, file_url FROM books WHERE 1=1"

	var args []interface{}
	i := 1

	if genre != "" {
		query += " AND genre=$" + strconv.Itoa(i)
		args = append(args, genre)
		i++
	}

	if q != "" {
		query += " AND name ILIKE $" + strconv.Itoa(i)
		args = append(args, "%"+q+"%")
		i++
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var b Book
		rows.Scan(&b.ID, &b.Name, &b.Author, &b.Genre, &b.Desc, &b.File)
		books = append(books, b)
	}

	data, _ := json.Marshal(books)

	if genre == "" && q == "" {
		rdb.Set(ctx, "books", data, 60*time.Second)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func getBookByID(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/books/")

	var b Book

	err := db.QueryRow(
		"SELECT id, name, author, genre, description, file_url FROM books WHERE id=$1",
		id,
	).Scan(&b.ID, &b.Name, &b.Author, &b.Genre, &b.Desc, &b.File)

	if err != nil {
		http.Error(w, "not found", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)
}
