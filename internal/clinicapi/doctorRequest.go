package clinicapi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Vadim992/clinicAPI/pkg/database/postgres"
	"github.com/Vadim992/clinicAPI/pkg/validator/validate"
	"net/http"
	"strings"
)

func (c *ClinicAPI) getDoctors(w http.ResponseWriter, r *http.Request) {
	var pageData PageData
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&pageData); err != nil {
		c.serveErr(w, err)
		return
	}

	if !validateNumPage(pageData.Page) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if pageData.DoctorsFilter.SpecializationFilter && pageData.DoctorsFilter.FirstNameFilter {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	offset := (pageData.Page - 1) * pageSize

	var filter string

	switch true {
	case pageData.DoctorsFilter.SpecializationFilter:
		filter = "specialization"
	case pageData.DoctorsFilter.FirstNameFilter:
		filter = "firstname"
	}

	doctors, err := c.DB.GetDoctors(offset, pageSize, filter)

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

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&doctor); err != nil {
		c.serveErr(w, err)
		return
	}

	doctor.Email = strings.ToLower(doctor.Email)

	//TODO: make one func ???

	// Validate doctor
	if !checkStructs(doctor) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if !validate.ValidateEmail(doctor.Email) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	err := c.DB.CheckEmailPatient(doctor.Email)

	if !errors.Is(err, sql.ErrNoRows) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	// make one func for validation ???

	if err := c.DB.InsertDoctor(doctor); err != nil {
		c.serveErr(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *ClinicAPI) putDoctor(w http.ResponseWriter, r *http.Request, id int) {

	decoder := json.NewDecoder(r.Body)

	var doctor postgres.Doctor
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&doctor); err != nil {
		c.serveErr(w, err)
		return
	}
	doctor.Email = strings.ToLower(doctor.Email)

	//TODO: make one func ???

	// Validate doctor
	if !checkStructs(doctor) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if !validate.ValidateEmail(doctor.Email) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	err := c.DB.CheckEmailPatient(doctor.Email)

	if !errors.Is(err, sql.ErrNoRows) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	// make one func for validation ???

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

	decoder := json.NewDecoder(r.Body)

	var doctor postgres.Doctor
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&doctor); err != nil {
		c.serveErr(w, err)
		return
	}

	doctor.Email = strings.ToLower(doctor.Email)

	//TODO: make one func ???

	// Validate doctor

	req := patchStructs(doctor)

	if req == "" {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if doctor.Email != "" {

		if !validate.ValidateEmail(doctor.Email) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}

		err := c.DB.CheckEmailPatient(doctor.Email)

		if !errors.Is(err, sql.ErrNoRows) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}

	}

	// make one func ???

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
