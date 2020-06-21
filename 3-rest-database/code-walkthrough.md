# Code Walkthrough

## Project Structure

```
.
├── controllers		// Package containing business logic to interface with the database
├── middlewares		// Package containing server middlewares for modifying responses to the client
├── models			// Package containing database models
└── main.go			// API server entry point
```

## Model

* `Book` definition
  ```diff
  // Book struct representing our Book DB model
  type Book struct {
  -  ID     string `json:"id"`
  -  Author string `json:"author"`
  -  Title  string `json:"title"`
  +  ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
  +  Author    string    `gorm:"size:100;" json:"author"`
  +  Title     string    `gorm:"size:100;" json:"title"`
  +  CreatedAt time.Time `json:"createdAt"`
  +  UpdatedAt time.Time `json:"updatedAt"`
  }
  ```
  * The `Book` struct has now been updated to utilize `gorm`-specific struct tags. The unique `ID` is an important field that is used by the ORM to determine how to handle operations against the database. Omitting this will cause unexpected behaviors!

        {
          "id": "1",
          "author": "John Smith",
          "title": "Git 101",
          "createdAt": "2020-06-21T00:40:59.035312Z",
          "updatedAt": "2020-06-21T00:51:44.539768Z": 
        }

* `Book.Validate()`
  ```go
  // Validate verifies that the required fields are present
  func (b *Book) Validate() error {
    if b.Title == "" {
      return errors.New("book: missing required title")
    }
    if b.Author == "" {
      return errors.New("book: missing required author")
    }
    return nil
  }
  ```
  * This method is a simple validation check against the values of a `Book` instance. Here, we are ensuring that the `author` and `title` fields are populated.
  * golang does not have a concept of a "class" construct that exists in object-oriented languages, however, the language supports the implementation of "methods" which can be defined on any type.
  * In this particular example, we see that the `Validate` function has a pointer receiver. Pointer receivers allow us to modify the values to which the receiver points to (Book).

* `Book.SaveBook()`
  ```go
  // SaveBook saves the model instance data to the database
  func (b *Book) SaveBook(db *gorm.DB) (*Book, error) {
    err := db.Debug().Create(&b).Error
    if err != nil {
      return nil, err
    }
    return b, nil
  }
  ```
  * Each of the CRUD operations we're implementing will receive a reference to the `db` object which is the ORM representation of the Postgres database. We can run commands against this `db` which will execute SQL queries for us under the hood.
  * The `Debug()` invocation here will print out a useful log containing the SQL query executed:
    ```bash
    (/home/user/go/src/github.com/user/example-go-rest-api/3-rest-database/models/book.go:40)
    [2020-06-20 15:23:18]  [1.22ms]  INSERT INTO "books" ("author","title","created_at","updated_at") VALUES ('J. K. Rowling','Harry Potter and the Sorcerers Stone','2020-06-20 15:23:18','2020-06-20 15:23:18') RETURNING "books"."id"
    [1 rows affected or returned ]
    ```

* `Book.GetAllBooks()`
  ```go
  // GetAllBooks returns a reference to array of first 100 books in the database
  func (b *Book) GetAllBooks(db *gorm.DB) (*[]Book, error) {
    books := []Book{}
    err := db.Debug().Limit(100).Order("id").Find(&books).Error
    if err != nil {
      return nil, err
    }
    return &books, nil
  }
  ```
  * An empty `Book` slice literal is created store the list of books returned from the database.
  * Here we're limiting the number of books returned to us to 100 records and ensuring that the order of the records returned is sorted by the `id` field.

* `Book.GetBookByID()`
  ```go
  // GetBookByID returns a reference to the book given the book ID
  func (b *Book) GetBookByID(db *gorm.DB, bid uint64) (*Book, error) {
    err := db.Debug().First(&b, bid).Error
    if err != nil {
      return nil, err
    }
    return b, nil
  }

  ```
  * Pretty self-explanatory here. The `.First()` method will return the first instance of a book record given the `bid`.

* `Book.UpdateBook()`
  ```go
  // UpdateBook updates the title and author fields in the database and returns a reference to the updated book
  func (b *Book) UpdateBook(db *gorm.DB) (*Book, error) {
    err := db.Debug().Model(b).Updates(Book{Title: b.Title, Author: b.Author}).Error
    if err != nil {
      return nil, err
    }
    return b, nil
  }
  ```
  * Here we're updating the existing model instance via the pointer receiver.

* `Book.DeleteBook()`
  ```go
  // DeleteBook removes the selected book from the database
  func (b *Book) DeleteBook(db *gorm.DB) (uint64, error) {
    err := db.Debug().Delete(b).Error
    if err != nil {
      return 0, err
    }
    return b.ID, nil
  }
  ```
  * Again, pretty straight forward here. The `.Delete()` method will remove the record from the database. One caveat here with the ORM is that this operation is considered a "soft delete" meaning the `DeletedAt` field is set to today's date.

## Controllers
### base.go

The base controller provides the `Server` which contains references to the gorm `DB` and also the gorilla/mux `Router`.  The `Server` type will be used to append all the router handlers. In addition, the `Server` type will also implement the `ConnectDb` and `StartServer` methods which will allow us to connect to the external Postgres service and start the API server on a specific port.

* `Server.initializeRoutes()`
  ```go
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
  ```
  * The base controller here defines the routes used for the API. If we wanted to introduce other API groupings, then we can extract and create multiple `initializeRoutes` for each set of APIs we want to create routes for (eg: authentication routes, other model routes).

* `Server.ConnectDb()`
  ```go
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
  ```
  * At this point, we'll have loaded the environment variables from the `.env` file. Here, we use the standard `os` library to read these environment variables to form the `DBURL` path to connect to our Postgres database:
    ```go
    DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
    ```

  * We use the 3rd party gorm library to open a connection to the Postgres database using the `DBURL`:
    ```go
    db, err := gorm.Open("postgres", DBURL)
    if err != nil {
      fmt.Printf("Cannot connect to Postgres database.")
      log.Fatal("This is the error:", err)
    }
    ```

  * Finally, we saved a reference to the database instance on our `Server` instance, perform a migration to setup the `book` table in Postgres and bind the API request handlers to the `Server` instance.

* `Server.StartServer()`
  ```go
  // StartServer initiates the API server
  func (server *Server) StartServer(addr string) {
    fmt.Printf("Listening to port %s", addr)
    log.Fatal(http.ListenAndServe(addr, server.Router))
  }
  ```
  * After we've made our database connection and setup our request handlers, let's start up the server by opening the provided port to handle incoming requests.

### bookController.go

Each of the CRUD function handlers is defined with a pointer receiver to a `Server` instance. We do this so that we can run the ORM commands against the `DB` reference on the `Server` struct.

* `Server.HandleCreateBook()`
  * Create `newBook` variable to store `Book` data received from the client and create error handler to send error response:
    ```go
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
    ```

  * Verify that the `author` and `title` fields are populated:
    ```go
    err = newBook.Validate()
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      json.NewEncoder(w).Encode(map[string]string{
        "message": "Unable to add new book to library.",
        "error":   err.Error(),
      })
      return
    }
    ```
  * Save the new book data to the database and return the newly added book data in the response:
    ```go
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
    ```

* `Server.HandleGetAllBooks()`
  * Creates a temporary `Book` instance:
    ```go
    book := models.Book{}
    ```
  * Use the temporary `Book` instance to fetch all the book records from the database:
    ```go
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
    ```
  * Notice that the value that's encoded in the response is `*books`. We're dereferencing the `books` variable since `book.GetAllBooks()` is returning a *pointer* to the value.
  * This implementation is a little unusual since we're effectively creating a dummy `Book` instance simply to reach for the `GetAllBooks()` method defined on it.

* `Server.HandleGetBookByID()`
  * Fetch the dynamic URL parameters from the URL path:
    ```go
    params := mux.Vars(r)
    ```

  * A quick sanity check to verify the input ID provided is a valid ID:
    ```go
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
    ```

  * Fetches the book from the database by the provided ID:
    ```go
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
    ```

* `Server.HandleUpdateBookByID()`
  * Get the `id` of the book to update:
    ```go
    id := mux.Vars(r)["id"]
    ```

  * Again, let's sanity check that the provided `id` is valid:
    ```go
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
    ```

  * Before attempting to update the book with the given `id`, let's verify that the book exists:
    ```go
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
    ```

  * Let's read the updated data values received from the request body and write them into the existing book reference:
    ```go
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
    ```
  
  * Finally, let's run the ORM command to update the book record in the database with the updated values received:
    ```go
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
    ```

* `Server.HandleDeleteBookByID()`
  * Get the book `id` that we want to delete:
    ```go
    id := mux.Vars(r)["id"]
    ```

  * Sanity check on if the book `id` provided is valid:
    ```go
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
    ```

  * Check if the book to delete exists:
    ```go
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
    ```

  * Run the ORM command to (soft) delete the book record from the database:
    ```go
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
    ```