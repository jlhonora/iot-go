package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"time"
	"encoding/json"
)

func Secret(user, realm string) string {
	return getPassword(user)
}

func renderOk(w http.ResponseWriter, content map[string]interface{}) {
	renderResponse(w, "ok", content)
}

func renderResponse(w http.ResponseWriter, status string, content map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"status": status}
	if content != nil {
		response = map[string]interface{}{
			"status":   status,
			"response": content,
		}
	}
	enc := json.NewEncoder(w)
	enc.Encode(response)
}

func iotPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Post")
	// receive posted data
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Fprint(w, "IoT failed")
		return
	}
	json, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Println("error in NewJson:", err)
		return
	}
	payload := json.Get("payload")
	payload_str, err := payload.String()
	if err != nil {
		fmt.Println("Couldn't get payload")
		return
	}
	created_at, err := json.Get("created_at").String()
	if err != nil {
		// If the created at time doesn't come with the
		// JSON object then generate it with RFC3339 format
		t := time.Now()
		created_at = t.UTC().Format("2006-01-02T15:04:05Z07:00")
	}
	fmt.Println("Got payload " + payload_str + " at " + created_at)
	renderOk(w, map[string]interface{}{"content": body})
}

func iotGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get")
	renderOk(w, nil)
}

func iotHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling IoT")
	if !checkAuth(r, r.Method) {
		badAuth(w)
		return
	}
	if r.Method == "POST" {
		iotPost(w, r)
	}
	if r.Method == "GET" {
		iotGet(w, r)
	}
}
