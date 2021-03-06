package objects

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == http.MethodGet {
		get(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
