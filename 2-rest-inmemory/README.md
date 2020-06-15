# Code Walkthrough

Let's go through `main.go` line by line and understand the implementation details of this example API server.

## Imports
```go
import (
  "encoding/json"
  "fmt"
  "net/http"

  "github.com/gorilla/mux"
)
```

* `encoding/json`: Standard go library used to encode/decode the JSON data into and out of the application.
* `fmt`: Standard go library used to trace helpful logs to the standard output.
* `net/http`: Standard go library to start the API server and provide an interface for reading requests and writing responses back to the client.
* `github.com/gorilla/mux`: Third-party routing package used for more complex URL path matching that's lacking from the standard go library.

## Database & Model
```go
type Book struct {
  ID     string `json:"id"`
  Author string `json:"author"`
  Title  string `json:"title"`
}

var library = []Book{}
```

* The `Book` struct defined here describes the data model representing a single book stored in the in-memory database.  One thing of note are the use of struct tags (eg: `json: "id"`) appended at the end of each struct field. The struct tags here are used to control how the properties are encoded into JSON object. In this example, each of the struct fields will be encoding using the struct tag property.

      {
        "id": "1",
        "author": "John Smith",
        "title": "Git 101"
      }

* The `library` variable represents our in-memory database. It is implemented as an array comprised of elements of type `Book`.

## CRUD Request Handlers
### **C**REATE
```go
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
```

* Create `newBook` variable to store `Book` data received from the client:
  ```go
  var newBook Book
  err := json.NewDecoder(r.Body).Decode(&newBook)
  ```

* Handle error from decoding client data and return error response:
  ```go
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(map[string]string{"message": "Unable to add new book to library."})
    return
  }
  ```

* Add newly created book into `library` "database":
  ```go
  library = append(library, newBook)
  ```

* Return `200` status code response along with JSON payload containing data of newly added book:
  ```go
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(newBook)
  ```

### **R**EAD
```go
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
```
* Fetch the dynamic URL parameters from the URL path:
  ```go
  params := mux.Vars(r)
  ```

* Search for book by the ID provided in the URL params and return the book data:
  ```go
  for _, book := range library {
		if book.ID == params["id"] {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&book)
			return
		}
  }
  ```

* Return a `404` error status code with an error message indicating book record was not found:
  ```go
  w.WriteHeader(http.StatusNotFound)
  json.NewEncoder(w).Encode(map[string]string{"message": "Unable to find book."})
  ```

### **U**PDATE
```go
func HandleUpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: HandleUpdate")

	id := mux.Vars(r)["id"]

	for idx := range library {
		book := &library[idx]
		if book.ID == id {
			err := json.NewDecoder(r.Body).Decode(book)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"message": "Unable to parse updated book values."})
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"id": id})
			return
		}
	}

  w.WriteHeader(http.StatusConflict)
  json.NewEncoder(w).Encode(map[string]string{"message": "Unable to update book, selected book not found."})
}
```

* The standard `range` function used for slices or arrays will return a *new* reference (in the second tuple output) to the elements within the array. This is important to note since we are directly writing the received updated book values to the book reference by the input ID.
  ```go
  for idx := range library {
    book := &library[idx]
    // We can fetch the correct address of the book by referencing the book by its index
  }
  ```

### **D**ELETE
```go
func HandleDeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: HandleDelete")

	id := mux.Vars(r)["id"]

	for idx := range library {
		book := &library[idx]
		if book.ID == id {
			library = append(library[:idx], library[idx+1:]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	w.WriteHeader(http.StatusConflict)
	json.NewEncoder(w).Encode(map[string]string{"message": "Unable to delete book, selected book not found."})
}
```

* This is the idiomatic way to remove the selected element from the array:
  ```go
  library = append(library[:idx], library[idx+1:]...)
  ```
