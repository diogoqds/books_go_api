package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"
)

type Book struct {
	ID     int    `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   string `json:year`
}

func init() {
	gotenv.Load()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal((err))
	}
}

var books []Book
var db *sql.DB

func main() {
	psqlInfo, err := pq.ParseURL(os.Getenv("DATABASE_URL"))
	psqlInfo += " sslmode=disable"
	logFatal(err)

	db, err = sql.Open("postgres", psqlInfo)
	logFatal(err)

	err = db.Ping()
	logFatal(err)

	router := mux.NewRouter()

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var book Book
	books = []Book{}

	rows, err := db.Query("select * from books")
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)

		books = append(books, book)

	}
	json.NewEncoder(w).Encode(books)

}

func getBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	params := mux.Vars(r)

	rows := db.QueryRow("select * from books where id=$1", params["id"])

	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	logFatal(err)
	json.NewEncoder(w).Encode(book)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	var bookID int

	json.NewDecoder(r.Body).Decode(&book)

	err := db.QueryRow("insert into books (title, author, year) values($1, $2, $3) returning id;", book.Title, book.Author, book.Year).Scan(&bookID)
	logFatal(err)
	json.NewEncoder(w).Encode(bookID)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
}

func removeBook(w http.ResponseWriter, r *http.Request) {
}
