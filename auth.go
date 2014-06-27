package main

import (
	"log"
	"net/http"
)

func checkAuth(r *http.Request, method string) bool {
	key := r.Header.Get("X-SESSION-KEY")
	ua := r.Header.Get("User-Agent")
	log.Printf("User agent: %s\n", ua)
	return checkKey(key, method)
}

func badAuth(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	renderResponse(w, "error", []byte(`{"response": "Bad Auth"}`))
	return
}
