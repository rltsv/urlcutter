package cuttool

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var idCount = 0
var UrlStorage = make(map[int]string)

func MakeUrlShorter(w http.ResponseWriter, r *http.Request) {

	respBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatal("Ошибка: ", err)
		return
	}
	urlString := fmt.Sprint(respBody)

	//Здесь будет какая-то логика по сокращению ссылки

	//

	idCount++
	UrlStorage[idCount] = urlString

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("Сокращенная строка - " + urlString))
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatal("Ошибка: ", err)
		return
	}

	_, err = fmt.Fprint(os.Stdout, urlString)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
