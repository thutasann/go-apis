package main

import (
	"fmt"
	"net/http"
)

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to homepage")
}

func main() {
	http.HandleFunc("/", homepage)
	http.ListenAndServe("localhost:10000", nil)
}
