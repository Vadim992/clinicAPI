package main

import (
	"net/http"
	"path"
)

// TODO: finish handlers

func (c *clinicAPI) handleClient(w http.ResponseWriter, r *http.Request) {

	pathUrl := r.URL.Path
	pathUrl = path.Clean(pathUrl)[1:]

	switch r.Method {

	case http.MethodGet:

		c.getClient(w, r, pathUrl)

	case http.MethodPost:

		c.postClient(w, r, pathUrl)

	case http.MethodPut:

		c.putClient(w, r, pathUrl)

	case http.MethodPatch:

		c.patchClient(w, r, pathUrl)

	case http.MethodDelete:

		c.deleteClient(w, pathUrl)

	default:
		c.clientErr(w, http.StatusMethodNotAllowed)
	}
}

func (c *clinicAPI) handleDoctor(w http.ResponseWriter, r *http.Request) {

	pathUrl := r.URL.Path
	pathUrl = path.Clean(pathUrl)[1:]

	switch r.Method {

	case http.MethodGet:

		c.getDoctor(w, r, pathUrl)

	case http.MethodPost:

		c.postDoctor(w, r, pathUrl)

	case http.MethodPut:

		c.putDoctor(w, r, pathUrl)

	case http.MethodPatch:

		c.patchDoctor(w, r, pathUrl)

	case http.MethodDelete:

		c.deleteDoctor(w, pathUrl)

	default:

		c.clientErr(w, http.StatusMethodNotAllowed)

	}
}
