package main

import (
	"github.com/rltsv/internal/cuttool"
	"github.com/rltsv/internal/urlgiver"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/POST", cuttool.MakeUrlShorter)
	mux.HandleFunc("/GET/", urlgiver.GetOrigURL)

	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}
