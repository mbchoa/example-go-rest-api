package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mbchoa/example-go-rest-api/3-rest-database/middlewares"
	"github.com/mbchoa/example-go-rest-api/3-rest-database/models"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Server provides the DB reference, the gorilla/mux Router and model controller methods
type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) initializeRoutes() {
	// setup middlewares
	server.Router.Use(middlewares.JSONMiddleware)

	// add a new book to library catalog
	server.Router.HandleFunc("/books", server.HandleCreateBook).Methods(http.MethodPost)

	// get a single book by ID
	server.Router.HandleFunc("/books/{id}", server.HandleGetBookByID).Methods(http.MethodGet)

	// get all books
	server.Router.HandleFunc("/books", server.HandleGetAllBooks).Methods(http.MethodGet)

	// update a book by ID
	server.Router.HandleFunc("/books/{id}", server.HandleUpdateBookByID).Methods(http.MethodPut)

	// delete a book by ID
	server.Router.HandleFunc("/books/{id}", server.HandleDeleteBookByID).Methods(http.MethodDelete)
}

// ConnectDb establishes the connection to the external database
func (server *Server) ConnectDb(dbUser, dbPassword, dbName, dbPort, dbHost string) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
	db, err := gorm.Open("postgres", DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to Postgres database.")
		log.Fatal("This is the error:", err)
	}

	server.DB = db
	server.DB.Debug().AutoMigrate(&models.Book{})
	server.Router = mux.NewRouter()
	server.initializeRoutes()
}

// StartServer initiates the API server
func (server *Server) StartServer(addr string) {
	fmt.Printf("Listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
