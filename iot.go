package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"time"
	"strconv"
)

func Secret(user, realm string) string {
	return getPassword(user)
}

func getIotRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(false)
	//r.HandleFunc("/api/v1/iot", iotHandler)

	s := r.PathPrefix("/api/v1/iot/").Subrouter()

	// Relative URIs
	s.Path("/").HandlerFunc(iotHandler)
	s.Path("/sensors").HandlerFunc(sensorsHandler)
	s.Path("/sensors/{id:[0-9]+}").HandlerFunc(sensorHandler)
	s.Path("/sensors/{id:[0-9]+}/measurements").HandlerFunc(sensorMeasurementsHandler)
	s.Queries("interval", "(peek|minute|hour|day)")

	return r
}

func renderOk(w http.ResponseWriter, content []byte) {
	renderResponse(w, "ok", content)
}

func renderError(w http.ResponseWriter, content []byte) {
	renderResponse(w, "error", content)
}

func renderResponse(w http.ResponseWriter, status string, content []byte) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"status": status}
	if content != nil {
		var content_iface interface{}
		json.Unmarshal(content, &content_iface)
		response = map[string]interface{}{
			"status":   status,
			"response": content_iface,
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
	renderOk(w, []byte(`{"content": body}`))
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

func sensorsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling sensors")
	if r.Method == "GET" {
		fmt.Println("Get")
		renderOk(w, getSensorsJson(nil))
	} else {
		renderError(w, []byte(`{"error": "Not supported"}`))
	}
}

func sensorMeasurementsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling sensor measurements")
	vars := mux.Vars(r)
	sensor_id, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		renderError(w, nil)
		return
	}
	fmt.Println("Got id: ", sensor_id)
	interval := r.FormValue("interval")
	fmt.Println("Interval: " + interval)
	renderOk(w, getMeasurementsJson(sensor_id, interval))
}

func sensorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling sensor")
	vars := mux.Vars(r)
	sensor_id, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		renderError(w, nil)
		return
	}
	fmt.Println("Got id: ", sensor_id)
	renderOk(w, getSensorJson(sensor_id))
}
