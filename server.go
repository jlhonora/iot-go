package main

import (
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	dbinit()
	defer dbclose()
	http.Handle("/api/v1/iot/", getIotRouter())
	http.Handle("/api/v1/iot", getIotRouter())
	log.Fatal(http.ListenAndServe(":7654", nil))
}
