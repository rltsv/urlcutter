package urlgiver

import (
	"github.com/rltsv/internal/cuttool"
	"net/http"
	"strconv"
	"strings"
)

func GetOrigURL(w http.ResponseWriter, r *http.Request) {
	idValue := r.URL.Path
	if idValue == "" {
		http.Error(w, "Ничего нет", http.StatusBadRequest)
		return
	}
	_, idPart, _ := strings.Cut(idValue, "GET/")
	intIdPart, err := strconv.Atoi(idPart)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	if _, ok := cuttool.UrlStorage[intIdPart]; !ok {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Content-Location", cuttool.UrlStorage[intIdPart])
	w.WriteHeader(http.StatusTemporaryRedirect)

}
