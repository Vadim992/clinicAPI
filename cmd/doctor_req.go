package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Vadim992/clinicAPI/internal/database/postgres"
	"github.com/Vadim992/clinicAPI/internal/validator/docrotvalid"
	"net/http"
	"path"
	"reflect"
	"strconv"
)

func (c *clinicAPI) getDoctor(w http.ResponseWriter, r *http.Request, pathUrl string) {

	switch {
	case docrotvalid.ValidPathAll(pathUrl):
		doctors, err := c.DB.GetAllDoctors()

		if err != nil {
			c.serveErr(w, err)
		}

		b, err := json.Marshal(doctors)

		if err != nil {
			c.serveErr(w, err)
		}

		w.Write(b)

	case docrotvalid.ValidPathId(pathUrl):

		idStr := path.Base(pathUrl)

		id, err := strconv.Atoi(idStr)

		if err != nil {
			c.serveErr(w, err)
		}

		doctor, err := c.DB.GetDoctor(id)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.notFound(w)
				return
			}
			c.serveErr(w, err)
		}

		b, err := json.Marshal(doctor)

		if err != nil {
			c.serveErr(w, err)
		}

		w.Write(b)
	default:
		c.notFound(w)
	}

}

func (c *clinicAPI) postDoctor(w http.ResponseWriter, r *http.Request, pathUrl string) {

	if !docrotvalid.ValidPathAll(pathUrl) {
		c.notFound(w)
		return
	}

	var doctor postgres.Doctor

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&doctor); err != nil {
		c.serveErr(w, err)
	}

	if err := c.DB.InsertDoctor(&doctor); err != nil {
		c.serveErr(w, err)
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *clinicAPI) putDoctor(w http.ResponseWriter, r *http.Request, pathUrl string) {

	if !docrotvalid.ValidPathId(pathUrl) {
		c.notFound(w)
		return
	}

	idStr := path.Base(pathUrl)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.serveErr(w, err)
	}

	decoder := json.NewDecoder(r.Body)

	var doctor postgres.Doctor
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&doctor); err != nil {
		c.serveErr(w, err)
	}

	if err := c.DB.UpdateDoctorAll(id, &doctor); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.notFound(w)
			return
		}
		c.serveErr(w, err)
	}
}

func (c *clinicAPI) patchDoctor(w http.ResponseWriter, r *http.Request, pathUrl string) {

	if !docrotvalid.ValidPathId(pathUrl) {
		c.notFound(w)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var doctor postgres.Doctor
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&doctor); err != nil {
		c.serveErr(w, err)
	}

	idStr := path.Base(pathUrl)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.serveErr(w, err)
	}

	v := reflect.ValueOf(doctor)
	t := v.Type()

	fieldsMap := make(map[string]string, postgres.CanChangeFieldDoctor)

	for i := 1; i < v.NumField(); i++ {

		if v.Field(i).CanInt() {
			val := v.Field(i).Int()
			if val != 0 {
				fieldsMap[t.Field(i).Name] = strconv.Itoa(int(val))
			}
		} else {
			val := v.Field(i).String()
			if val != "" {
				fieldsMap[t.Field(i).Name] = val
			}
		}
	}

	if err := c.DB.UpdateDoctor(id, fieldsMap); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.notFound(w)
			return
		}
		c.serveErr(w, err)
	}

}

func (c *clinicAPI) deleteDoctor(w http.ResponseWriter, pathUrl string) {

	if !docrotvalid.ValidPathId(pathUrl) {
		c.notFound(w)
		return
	}

	idStr := path.Base(pathUrl)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.serveErr(w, err)
	}

	if err := c.DB.DeleteDoctor(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.notFound(w)
			return
		}

		c.serveErr(w, err)
	}
}
