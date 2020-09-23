package temp

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == http.MethodHead {
		head(w, r)
		return
	} else if method == http.MethodPut {
		put(w, r)
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
