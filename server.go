package main

import (
	auth "github.com/abbot/go-http-auth"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func getRealm() string {
	return "127.0.0.1"
}

func main() {
	dbinit()
	defer dbclose()
	authenticator := auth.NewBasicAuthenticator(getRealm(), Secret)
	http.HandleFunc("/api/v1/iot", authenticator.Wrap(iotHandler))
	log.Fatal(http.ListenAndServe(":7654", nil))
}
