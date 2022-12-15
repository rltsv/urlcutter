package main

import (
	"bytes"
	"fmt"
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

func MakeShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if strings.TrimLeft(r.URL.Path, "/") != "" {
			http.Error(w, "Некорректный запрос.", 400)
			return
		}

		respBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if len(respBody) == 0 {
			http.Error(w, "Мне нечего сокращать, уточните ссылку!", 400)
			return
		}

		for key, val := range URLStorage {
			if bytes.Equal(val, respBody) == true {
				w.WriteHeader(400)
				_, err = w.Write([]byte(fmt.Sprint("Такая ссылка уже есть в базе - ", key)))
				if err != nil {
					http.Error(w, err.Error(), 500)
				}
				return
			}
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
	} else {
		http.Error(w, "Ошибка запроса!", 400)
	}

}

func GiveOriginalLinkToRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idValue := strings.TrimPrefix(r.URL.Path, "/GET/")
		if idValue == "" {
			http.Error(w, "Уточните запрос.", 400)
			return
		}
		intIDPart, err := strconv.Atoi(idValue)
		if err != nil {
			http.Error(w, "Уточните запрос.", 400)
			return
		}

		if len(URLStorage) == 0 {
			http.Error(w, "В хранилище ничего нет!", 400)
			return
		}

		mux.Lock()
		if _, ok := URLStorage[intIDPart]; !ok {
			http.Error(w, "По данному id, ничего нет. Уточните запрос.", 400)
			return
		}
		mux.Unlock()

		w.Header().Set("Content-Location", string(URLStorage[intIDPart]))
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "Ошибка запроса!", 400)
	}

}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", MakeShortLink)
	mux.HandleFunc("/GET/", GiveOriginalLinkToRequest)

	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}
