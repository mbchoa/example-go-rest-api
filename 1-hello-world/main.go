package main

import (
	"io"
	"net/http"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func main() {
	// add request handler for "/" route
	http.HandleFunc("/", HandleRoot)

	// serve and listen on http://localhost:8080
	http.ListenAndServe(":8080", nil)
}
