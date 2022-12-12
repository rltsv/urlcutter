package cuttool

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

var idCount = 0
var URLStorage = make(map[int][]byte)

func MakeURLShorter(w http.ResponseWriter, r *http.Request) {

	respBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatal("Ошибка: ", err)
		return
	}
	//Здесь будет какая-то логика по сокращению ссылки

	urlString := fmt.Sprint(respBody)
	URLLink, err := url.Parse(urlString)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	partOfURL := URLLink.Path

	//Здесь будет какая-то логика по сокращению ссылки

	//Fill map with id and full link
	idCount++
	URLStorage[idCount] = respBody

	//Answer to client with shortened link
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(201)
	_, err = w.Write([]byte(partOfURL))
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatal("Ошибка: ", err)
		return
	}
}
