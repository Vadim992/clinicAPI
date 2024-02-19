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

func (c *ClinicAPI) getPatients(w http.ResponseWriter, r *http.Request) {
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

	if pageData.PatientsFilter.PhoneFilter && pageData.PatientsFilter.FirstNameFilter {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	offset := (pageData.Page - 1) * pageData.PageSize

	var filter string

	switch true {
	case pageData.PatientsFilter.PhoneFilter:
		filter = "phone_number"
	case pageData.PatientsFilter.FirstNameFilter:
		filter = "firstname"
	}

	patients, err := c.DB.GetPatients(offset, pageData.PageSize, filter)

	if err != nil {
		c.serveErr(w, err)
		return
	}

	b, err := json.Marshal(patients)

	if err != nil {
		c.serveErr(w, err)
		return
	}
	w.Write(b)
}

func (c *ClinicAPI) getPatientId(w http.ResponseWriter, r *http.Request, id int) {

	patient, err := c.DB.GetPatientId(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.notFound(w)
			return
		}
		c.serveErr(w, err)
		return
	}

	b, err := json.Marshal(patient)

	if err != nil {
		c.serveErr(w, err)
		return
	}
	w.Write(b)
}

func (c *ClinicAPI) postPatient(w http.ResponseWriter, r *http.Request) {

	var patient postgres.Patient

	if err := decode(r, &patient); err != nil {
		if errors.Is(err, io.EOF) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}
		c.serveErr(w, err)
	}
	patient.Email.String = strings.ToLower(patient.Email.String)

	if !c.validatePatientEmailPhone(w, r, &patient) {
		return
	}

	//patient.Phone_Number = formatPhoneNumber(patient.Phone_Number)

	if !putCheckStructs(patient) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if err := c.DB.InsertPatient(patient); err != nil {
		c.serveErr(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *ClinicAPI) putPatient(w http.ResponseWriter, r *http.Request, id int) {

	var patient postgres.Patient

	if err := decode(r, &patient); err != nil {
		if errors.Is(err, io.EOF) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}
		c.serveErr(w, err)
		return
	}

	patient.Email.String = strings.ToLower(patient.Email.String)

	if !c.validatePatientEmailPhone(w, r, &patient) {
		return
	}

	//patient.Phone_Number = formatPhoneNumber(patient.Phone_Number)

	if !putCheckStructs(patient) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if err := c.DB.UpdatePatientAll(id, patient); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			c.notFound(w)
			return
		}

		c.serveErr(w, err)
		return
	}

}

func (c *ClinicAPI) patchPatient(w http.ResponseWriter, r *http.Request, id int) {

	var patient postgres.Patient

	if err := decode(r, &patient); err != nil {
		if errors.Is(err, io.EOF) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}
		c.serveErr(w, err)
		return
	}

	patient.Email.String = strings.ToLower(patient.Email.String)

	if !c.validatePatientEmailPhone(w, r, &patient) {
		return
	}

	//patient.Phone_Number = formatPhoneNumber(patient.Phone_Number)

	req := patchCheckStructs(patient)

	if req == "" {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if err := c.DB.UpdatePatient(id, req); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			c.notFound(w)
			return
		}
		c.serveErr(w, err)
		return
	}

}

func (c *ClinicAPI) deletePatient(w http.ResponseWriter, r *http.Request, id int) {

	if err := c.DB.DeletePatient(id); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			c.notFound(w)
			return
		}

		c.serveErr(w, err)
		return
	}
}
