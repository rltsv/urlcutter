package main

import (
	"net/http"
)

func MethodPost(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

	}
}

func main() {

	h1 := http.HandlerFunc(MethodPost)

	http.HandleFunc("/", h1)

	http.ListenAndServe("localhost:8080", nil)
}
