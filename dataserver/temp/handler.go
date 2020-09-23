package temp

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == http.MethodPut {
		put(w, r)
		return
	} else if method == http.MethodPatch {
		patch(w, r)
		return
	} else if method == http.MethodPost {
		post(w, r)
		return
	} else if method == http.MethodDelete {
		del(w, r)
		return
	} else if method == http.MethodHead {
		head(w, r)
		return
	} else if method == http.MethodGet {
		get(w, r)
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
