package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var idCount = 0
var URLStorage = make(map[int][]byte)
var mux sync.Mutex

//func MakeShortLink(w http.ResponseWriter, r *http.Request) {
//	respBody, err := io.ReadAll(r.Body)
//	if err != nil {
//		http.Error(w, err.Error(), 500)
//		log.Fatal("Ошибка: ", err)
//		return
//	}
//
//	//инкрементим счетчик айди на один и записываем полученную ссылку по данному айдишнику в хранилище
//	mux.Lock()
//	idCount++
//	URLStorage[idCount] = respBody
//	mux.Unlock()
//
//	//здесь сокращаем нашу ссылку, пока что будем делать просто порядковые номера: 1, 2...
//	var LinkAfterHashFunction = strconv.Itoa(idCount)
//
//	//отвечаем клиенту в виде сокращенной ссылки и статус кодом
//	w.WriteHeader(201)
//	_, err = w.Write([]byte(LinkAfterHashFunction))
//	if err != nil {
//		http.Error(w, err.Error(), 500)
//		log.Fatal("Ошибка: ", err)
//		return
//	}
//}
//
//func GiveOriginalLinkToRequest(w http.ResponseWriter, r *http.Request) {
//	idValue := r.URL.Path
//	if idValue == "" {
//		http.Error(w, "", http.StatusBadRequest)
//		return
//	}
//	_, idPart, _ := strings.Cut(idValue, "GET/")
//	intIDPart, err := strconv.Atoi(idPart)
//	if err != nil {
//		http.Error(w, err.Error(), 500)
//	}
//
//	if _, ok := URLStorage[intIDPart]; !ok {
//		http.Error(w, "", http.StatusBadRequest)
//	}
//
//	w.Header().Set("Content-Location", string(URLStorage[intIDPart]))
//	w.WriteHeader(http.StatusTemporaryRedirect)
//
//}

func MapMethodToFunction(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idValue := strings.TrimLeft(r.URL.Path, "/")
		if idValue == "" {
			http.Error(w, "", 400)
			return
		}
		intIDPart, err := strconv.Atoi(idValue)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if len(URLStorage) == 0 {
			http.Error(w, "По данному id, ничего нет. Уточните запрос.", 400)
			return
		}

		if _, ok := URLStorage[intIDPart]; !ok {
			http.Error(w, err.Error(), 400)
			return
		}

		w.Header().Set("Content-Location", string(URLStorage[intIDPart]))
		w.WriteHeader(http.StatusTemporaryRedirect)
	}

	if r.Method == http.MethodPost {
		respBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Fatal("Ошибка: ", err)
			return
		}

		//инкрементим счетчик айди на один и записываем полученную ссылку по данному айдишнику в хранилище
		mux.Lock()
		idCount++
		URLStorage[idCount] = respBody
		mux.Unlock()

		//здесь сокращаем нашу ссылку, пока что будем делать просто порядковые номера: 1, 2...
		var LinkAfterHashFunction = strconv.Itoa(idCount)

		//отвечаем клиенту в виде сокращенной ссылки и статус кодом
		w.WriteHeader(201)
		_, err = w.Write([]byte(LinkAfterHashFunction))
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Fatal("Ошибка: ", err)
			return
		}
	}
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", MapMethodToFunction)

	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}
