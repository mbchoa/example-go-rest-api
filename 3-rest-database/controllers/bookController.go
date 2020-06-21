package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mbchoa/example-go-rest-api/3-rest-database/models"

	"github.com/gorilla/mux"
)

// HandleCreateBook handles request for creating a new book record in the database
func (server *Server) HandleCreateBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: HandleCreate")

	newBook := models.Book{}
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Unable to add new book to library.",
			"error":   err.Error(),
		})
		return
	}

	err = newBook.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Unable to add new book to library.",
			"error":   err.Error(),
		})
		return
	}

	_, err = newBook.SaveBook(server.DB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Unable to add new book to library.",
			"error":   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

// HandleGetAllBooks handles returning all books in the database
func (server *Server) HandleGetAllBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: HandleGetAllBooks")

	book := models.Book{}
	books, err := book.GetAllBooks(server.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Unable to fetch all books from library.",
			"error":   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*books)
}

// HandleGetBookByID handles returning a single book by its uniquue ID
func (server *Server) HandleGetBookByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: HandleGetBookById")

	params := mux.Vars(r)

	// Check if the book ID is valid
	bookID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Unable to find book.",
			"error":   err.Error(),
		})
		return
	}

	// Get book record by ID
	book := models.Book{}
	_, err = book.GetBookByID(server.DB, bookID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Unable to find book.",
			"error":   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

// HandleUpdateBookByID handles modifying an existing book in the database
func (server *Server) HandleUpdateBookByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: HandleUpdate")

	id := mux.Vars(r)["id"]

	// Check if the book ID is valid
	bookID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Unable to update book.",
			"error":   err.Error(),
		})
		return
	}

	// Check if the book exists
	bookToUpdate := models.Book{}
	_, err = bookToUpdate.GetBookByID(server.DB, bookID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Unable to find book to update.",
			"error":   err.Error(),
		})
		return
	}

	// Read updated data from request body
	err = json.NewDecoder(r.Body).Decode(&bookToUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Unable to update book.",
			"error":   err.Error(),
		})
		return
	}

	// Updates book record in database
	_, err = bookToUpdate.UpdateBook(server.DB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Failed to update book record.",
			"error":   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookToUpdate)
}

// HandleDeleteBookByID handles removing a book record from the database
func (server *Server) HandleDeleteBookByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Check if the book ID is valid
	bookID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Unable to delete book.",
			"error":   err.Error(),
		})
		return
	}

	// Check if the book exists
	bookToDelete := models.Book{}
	_, err = bookToDelete.GetBookByID(server.DB, bookID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Unable to find book to delete.",
			"error":   err.Error(),
		})
		return
	}

	// Soft deletes book record in database
	_, err = bookToDelete.DeleteBook(server.DB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Failed to delete book record.",
			"error":   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]uint64{"id": bookID})
}
