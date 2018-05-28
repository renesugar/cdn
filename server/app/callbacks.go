package app

import (
	"encoding/json"
	"net/http"
)

// Info handles the HTTP request sent on one of the info endpoints
func Info(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Incorrect method for info request"))
		return
	}
	var req SearchRequest
	err := req.ParseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Incorrect list request"))
		return
	}
	req.operation = "info"
	files := Retrieve(req)
	result, _ := json.Marshal(files)
	w.Write([]byte(result))
}

// List handles the HTTP request sent on one of the list endpoints
func List(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Incorrect method for list request"))
		return
	}
	var req SearchRequest
	err := req.ParseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Incorrect list request"))
		return
	}
	req.operation = "list"
	files := Retrieve(req)
	results, _ := json.Marshal(files)
	w.Write([]byte(results))
}
