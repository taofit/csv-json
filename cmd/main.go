package main

import (
	"fmt"
	"html"
	"net/http"
	"encoding/csv"
)

type hierarchyHandler struct{}

func (f *hierarchyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
	fmt.Println(csv.NewReader(r.Body).ReadAll())
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", &hierarchyHandler{})
	err := http.ListenAndServe(":8080", &hierarchyHandler{})
	if (err != nil) {
		fmt.Printf("Error: %v", err.Error())
	}
}
