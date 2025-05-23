package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func ContextMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "TraceID ", CreateUuid())
		newReq := r.WithContext(ctx)
		h.ServeHTTP(w, newReq)
	})
}

func LogMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s %s", r.Header.Get("TraceID"), r.Method, r.RequestURI)
		h.ServeHTTP(w, r)
	})
}

func CompileMiddleware(h http.HandlerFunc, m []Middleware) http.HandlerFunc {
	if len(m) < 1 {
		return h
	}

	wrapped := h

	// loop in reverse to preserve middleware order
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}

	return wrapped
}

func CreateUuid() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}
