package main

import (
	"net/http"
)

// TODO: finish handlers
func (c *clinicAPI) handleClient(w http.ResponseWriter, r *http.Request) {
	// Transfer data using Querry Parametrs ?

	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
	case http.MethodPatch:
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *clinicAPI) handleDoctor(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
	case http.MethodPatch:
	default:

	}
}
