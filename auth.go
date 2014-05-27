package main

import (
	"log"
	"net/http"
)

func checkAuth(r * http.Request) bool {
	key := r.Header.Get("X-SESSION-KEY")
	ua := r.Header.Get("User-Agent")
	log.Printf("User agent: %s\n", ua)
	return checkKey(key)
}

func badAuth(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	return
}
