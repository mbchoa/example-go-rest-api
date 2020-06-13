package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Book struct representing our Book DB model
type Book struct {
	ID     string `json:"id"`
	Author string `json:"author"`
	Title  string `json:"title"`
}

// Library represents the in-memory database of books
var Library []Book

// HandleCreate callback for adding a new book to the library
func HandleCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: HandleCreate")
}

// HandleGetBookByID callback for returning a book from the library
func HandleGetBookByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: HandleGetBookById")
}

// HandleGetAllBooks callback for returning all books from the library
func HandleGetAllBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: HandleGetAllBooks")
}

// HandleUpdateByID callback for updating a single book from the library
func HandleUpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: HandleUpdate")
}

// HandleDeleteByID callback for deleting a single book from the library
func HandleDeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: HandleDelete")
}

func main() {
	r := mux.NewRouter()

	// add a new book to library catalog
	r.HandleFunc("/books", HandleCreate).Methods(http.MethodPost)

	// get a single book by ID
	r.HandleFunc("/books/{id}", HandleGetBookByID).Methods(http.MethodGet)

	// get all books
	r.HandleFunc("/books", HandleGetAllBooks).Methods(http.MethodGet)

	// update a book by ID
	r.HandleFunc("/books/{id}", HandleUpdateByID).Methods(http.MethodPut)

	// delete a book by ID
	r.HandleFunc("/books/{id}", HandleDeleteByID).Methods(http.MethodDelete)

	// serve and listen on http://localhost:7070
	http.ListenAndServe(":7070", r)
}
