package objects

import (
	"lib/es"
	"lib/utils"
	"log"
	"net/http"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	size := utils.GetSizeFromHeader(r.Header)
	status, err := storeObject(r.Body, hash, size)
	if err != nil {
		log.Println(err)
		w.WriteHeader(status)
		return
	}
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	err = es.AddVersion(name, hash, size)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
