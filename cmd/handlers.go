package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"path"
	"regexp"
	"strconv"
)

// TODO: finish handlers

func (c *clinicAPI) handleClient(w http.ResponseWriter, r *http.Request) {

	pathUrl := r.URL.Path
	pathUrl = path.Clean(pathUrl)[1:]

	switch r.Method {
	case http.MethodGet:
		// TODO: refactor code !!!
		rePathAll := regexp.MustCompile("^client$")
		rePathId := regexp.MustCompile("^client/[0-9]+$")

		if rePathAll.MatchString(pathUrl) {
			clients, err := c.DB.GetAllClients()

			if err != nil {
				c.errLog.Fatal("Error")
			}
			c.infoLog.Printf("Get slice of clients %v", clients)

			b, err := json.Marshal(clients)

			if err != nil {
				c.errLog.Fatal("Error")
			}
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return

		} else if rePathId.MatchString(pathUrl) {

			idStr := path.Base(pathUrl)

			id, err := strconv.Atoi(idStr)

			if err != nil {
				c.errLog.Fatal("cant convert id from string to int")
			}

			client, err := c.DB.GetClient(id)

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					http.NotFound(w, r)
					return
				}
				c.errLog.Fatal("Error")
			}
			c.infoLog.Printf("Get slice of clients %v", client)

			b, err := json.Marshal(client)

			if err != nil {
				c.errLog.Fatal("Error")
			}
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return

		}
		http.NotFound(w, r)

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
