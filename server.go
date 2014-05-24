package main

import (
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	dbinit()
	defer dbclose()
	http.HandleFunc("/api/v1/iot", iotHandler)
	log.Fatal(http.ListenAndServe(":7654", nil))
}
