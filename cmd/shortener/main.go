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

func HeadFunction(w http.ResponseWriter, r *http.Request) {

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
			if bytes.Equal(val, respBody) {
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

		//отвечаем клиенту в виде сокращенной ссылки и статус кодом
		w.WriteHeader(201)
		_, err = w.Write([]byte(fmt.Sprintf("http://localhost:8080/%d", idCount)))
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Fatal("Ошибка: ", err)
			return
		}

	} else if r.Method == http.MethodGet {

		idValue := strings.TrimPrefix(r.URL.Path, "/")
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

		w.Header().Set("Location", string(URLStorage[intIDPart]))
		w.WriteHeader(http.StatusTemporaryRedirect)

	} else {
		http.Error(w, "Ошибка запроса!", 400)
	}

}

func main() {

	http.HandleFunc("/", HeadFunction)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
