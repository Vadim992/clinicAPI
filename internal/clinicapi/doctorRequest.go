package clinicapi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Vadim992/clinicAPI/pkg/database/postgres"
	"io"
	"net/http"
	"strings"
)

func (c *ClinicAPI) getDoctors(w http.ResponseWriter, r *http.Request) {
	var pageData PageData

	if err := decode(r, &pageData); err != nil {
		if errors.Is(err, io.EOF) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}
		c.serveErr(w, err)
		return
	}

	if !validateNumPage(pageData.Page, pageData.PageSize) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if pageData.DoctorsFilter.SpecializationFilter && pageData.DoctorsFilter.FirstNameFilter {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	offset := (pageData.Page - 1) * pageData.PageSize

	var filter string

	switch true {
	case pageData.DoctorsFilter.SpecializationFilter:
		filter = "specialization"
	case pageData.DoctorsFilter.FirstNameFilter:
		filter = "firstname"
	}

	doctors, err := c.DB.GetDoctors(offset, pageData.PageSize, filter)

	if err != nil {
		c.serveErr(w, err)
		return
	}

	b, err := json.Marshal(doctors)

	if err != nil {
		c.serveErr(w, err)
		return
	}

	w.Write(b)

}

func (c *ClinicAPI) getDoctorId(w http.ResponseWriter, r *http.Request, id int) {

	doctor, err := c.DB.GetDoctorId(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.notFound(w)
			return
		}
		c.serveErr(w, err)
		return
	}

	b, err := json.Marshal(doctor)

	if err != nil {
		c.serveErr(w, err)
		return
	}

	w.Write(b)

}

func (c *ClinicAPI) postDoctor(w http.ResponseWriter, r *http.Request) {

	var doctor postgres.Doctor

	if err := decode(r, &doctor); err != nil {
		if errors.Is(err, io.EOF) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}
		c.serveErr(w, err)
		return
	}

	doctor.Email = strings.ToLower(doctor.Email)

	if !c.validateDoctor(w, r, doctor) {
		return
	}

	if err := c.DB.InsertDoctor(doctor); err != nil {
		c.serveErr(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *ClinicAPI) putDoctor(w http.ResponseWriter, r *http.Request, id int) {

	var doctor postgres.Doctor

	if err := decode(r, &doctor); err != nil {
		if errors.Is(err, io.EOF) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}
		c.serveErr(w, err)
		return
	}

	doctor.Email = strings.ToLower(doctor.Email)

	if !c.validateDoctor(w, r, doctor) {
		return
	}

	if err := c.DB.UpdateDoctorAll(id, doctor); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			c.notFound(w)
			return
		}

		c.serveErr(w, err)
		return
	}
}

func (c *ClinicAPI) patchDoctor(w http.ResponseWriter, r *http.Request, id int) {

	var doctor postgres.Doctor

	if err := decode(r, &doctor); err != nil {
		if errors.Is(err, io.EOF) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}
		c.serveErr(w, err)
		return
	}

	doctor.Email = strings.ToLower(doctor.Email)

	if !c.validateDoctor(w, r, doctor) {
		return
	}

	req := patchCheckStructs(doctor)

	if req == "" {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if err := c.DB.UpdateDoctor(id, req); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			c.notFound(w)
			return
		}

		c.serveErr(w, err)
		return
	}

}

func (c *ClinicAPI) deleteDoctor(w http.ResponseWriter, r *http.Request, id int) {

	if err := c.DB.DeleteDoctor(id); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			c.notFound(w)
			return
		}

		c.serveErr(w, err)
		return
	}
}
