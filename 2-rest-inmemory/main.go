package main

import (
	"encoding/json"
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
var library = []Book{}

// HandleCreate callback for adding a new book to the library
func HandleCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: HandleCreate")

	var newBook Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unable to add new book to library."})
		return
	}

	library = append(library, newBook)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

// HandleGetBookByID callback for returning a book from the library
func HandleGetBookByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: HandleGetBookById")

	params := mux.Vars(r)

	for _, book := range library {
		if book.ID == params["id"] {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&book)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"message": "Unable to find book."})
}

// HandleGetAllBooks callback for returning all books from the library
func HandleGetAllBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: HandleGetAllBooks")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(library)
}

// HandleUpdateByID callback for updating a single book from the library
func HandleUpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: HandleUpdate")

	id := mux.Vars(r)["id"]

	for idx := range library {
		book := &library[idx]
		if book.ID == id {
			// write updated book fields received from client into book at index
			err := json.NewDecoder(r.Body).Decode(book)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"message": "Unable to parse updated book values."})
				return
			}

			// send OK response
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"id": id})
			return
		}
	}

	// book not found, send 409 response
	// https://stackoverflow.com/questions/10727699/is-http-404-an-appropriate-response-for-a-put-operation-where-some-linked-resour
	w.WriteHeader(http.StatusConflict)
	json.NewEncoder(w).Encode(map[string]string{"message": "Unable to update book, selected book not found."})
}

// HandleDeleteByID callback for deleting a single book from the library
func HandleDeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: HandleDelete")

	id := mux.Vars(r)["id"]

	for idx := range library {
		book := &library[idx]
		if book.ID == id {
			// re-slice and remove the book element
			// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
			library = append(library[:idx], library[idx+1:]...)

			// send OK response
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	// book to delete not found, send 409 response
	w.WriteHeader(http.StatusConflict)
	json.NewEncoder(w).Encode(map[string]string{"message": "Unable to delete book, selected book not found."})
}

func jsonContentTypeHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()
	r.Use(jsonContentTypeHeaderMiddleware)

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
