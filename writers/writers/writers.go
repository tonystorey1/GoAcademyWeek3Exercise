package writers

import "net/http"

func WriteResponse(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
}

func WriteResponseWithMessage(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	if message != "200" {
		_, _ = w.Write([]byte(message))
	}
	w.WriteHeader(code)
}
