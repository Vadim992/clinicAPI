package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Vadim992/clinicAPI/internal/database/postgres"
	"github.com/Vadim992/clinicAPI/internal/validator/clientvalid"
	"net/http"
	"path"
	"reflect"
	"strconv"
)

func (c *clinicAPI) getClient(w http.ResponseWriter, r *http.Request, pathUrl string) {

	switch {
	case clientvalid.ValidPathAll(pathUrl):

		clients, err := c.DB.GetAllClients()

		if err != nil {
			c.serveErr(w, err)
		}

		b, err := json.Marshal(clients)

		if err != nil {
			c.serveErr(w, err)
		}
		w.Write(b)

	case clientvalid.ValidPathId(pathUrl):

		idStr := path.Base(pathUrl)

		id, err := strconv.Atoi(idStr)

		if err != nil {
			c.serveErr(w, err)
		}

		client, err := c.DB.GetClient(id)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.notFound(w)
				return
			}
			c.serveErr(w, err)
		}

		b, err := json.Marshal(client)

		if err != nil {
			c.serveErr(w, err)
		}
		w.Write(b)

	default:
		c.notFound(w)
	}

}

func (c *clinicAPI) postClient(w http.ResponseWriter, r *http.Request, pathUrl string) {

	//rePath := regexp.MustCompile(`^client/[a-zA-Z]+/[a-zA-Z]+/[a-zA-Z]+/([a-zA-Z0-9]+(_?))+$`)

	//if rePath.MatchString(pathUrl) {
	//	fmt.Println("OK")
	//}

	if !clientvalid.ValidPathAll(pathUrl) {
		c.notFound(w)
		return
	}

	var client postgres.Client

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&client); err != nil {
		c.serveErr(w, err)
	}

	if err := c.DB.InsertClient(&client); err != nil {
		c.serveErr(w, err)
	}
}

func (c *clinicAPI) putClient(w http.ResponseWriter, r *http.Request, pathUrl string) {

	//rePath := regexp.MustCompile(`^client/[0-9]+/[a-zA-Z]+/[a-zA-Z]+/[a-zA-Z]+/([a-zA-Z0-9]+(_?))+$`)
	//
	//if rePath.MatchString(pathUrl) {
	//	fmt.Println("OK")
	//}

	if !clientvalid.ValidPathId(pathUrl) {
		c.notFound(w)
		return
	}

	idStr := path.Base(pathUrl)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.serveErr(w, err)
	}

	decoder := json.NewDecoder(r.Body)

	var client postgres.Client
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&client); err != nil {
		c.serveErr(w, err)
	}

	if err := c.DB.UpdateClientAll(id, &client); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.notFound(w)
			return
		}
		c.serveErr(w, err)
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *clinicAPI) patchClient(w http.ResponseWriter, r *http.Request, pathUrl string) {

	if !clientvalid.ValidPathId(pathUrl) {
		c.notFound(w)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var client postgres.Client
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&client); err != nil {
		c.serveErr(w, err)
	}

	idStr := path.Base(pathUrl)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.serveErr(w, err)
	}

	v := reflect.ValueOf(client)
	t := v.Type()

	fieldsMap := make(map[string]string, postgres.CanChangeFieldClient)

	for i := 1; i < v.NumField(); i++ {
		val := v.Field(i).String()

		if val != "" {
			fieldsMap[t.Field(i).Name] = val
		}
	}

	if err := c.DB.UpdateClient(id, fieldsMap); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.notFound(w)
			return
		}
		c.serveErr(w, err)
	}

}

func (c *clinicAPI) deleteClient(w http.ResponseWriter, pathUrl string) {

	if !clientvalid.ValidPathId(pathUrl) {
		c.notFound(w)
		return
	}

	idStr := path.Base(pathUrl)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.serveErr(w, err)
	}

	if err := c.DB.DeleteClient(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.notFound(w)
			return
		}

		c.serveErr(w, err)
	}
}
