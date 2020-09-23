package versions

import (
	"encoding/json"
	"lib/es"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	from := 0
	size := 1000
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	for {
		metas, err := es.SearchAllVersions(name, from, size)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, m := range metas {
			bytes, _ := json.Marshal(m)
			_, _ = w.Write(bytes)
			_, _ = w.Write([]byte("\n"))
		}
		if len(metas) != size {
			return
		}
		from += size
	}

}
