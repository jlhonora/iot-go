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
	renderResponse(w, "error", []byte("Bad Auth"))
	return
}

/**
 * Checks the auth for request `r`. If there's
 * an error then it prints the corresponding message
 * in `w`.
 * @see `badAuth` and `checkAuth`
 * @returns `true` if correctly authenticated
 */
func authOk(w http.ResponseWriter, r *http.Request) bool {
	if !checkAuth(r, r.Method) {
		badAuth(w)
		return false
	}
	return true
}
