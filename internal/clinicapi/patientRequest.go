package clinicapi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Vadim992/clinicAPI/pkg/database/postgres"
	"net/http"
	"strings"
)

func (c *ClinicAPI) getPatients(w http.ResponseWriter, r *http.Request) {
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

	if pageData.PatientsFilter.PhoneFilter && pageData.PatientsFilter.FirstNameFilter {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	offset := (pageData.Page - 1) * pageSize

	var filter string

	switch true {
	case pageData.PatientsFilter.PhoneFilter:
		filter = "phone_number"
	case pageData.PatientsFilter.FirstNameFilter:
		filter = "firstname"
	}

	patients, err := c.DB.GetPatients(offset, pageSize, filter)

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

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&patient); err != nil {
		c.serveErr(w, err)
	}
	patient.Email.String = strings.ToLower(patient.Email.String)

	//TODO: make one func ???

	// Validate patient

	if !checkEmail(patient.Email) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if !checkPhoneNum(patient.PhoneNumber) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if !checkStructs(patient) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if patient.Email.Valid {

		err := c.DB.CheckEmailPatient(patient.Email.String)

		if !errors.Is(err, sql.ErrNoRows) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}

	}

	if patient.PhoneNumber != "" {

		err := c.DB.CheckPhoneNumberPatient(patient.PhoneNumber)

		if !errors.Is(err, sql.ErrNoRows) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}
	}
	// make one func ???

	if err := c.DB.InsertPatient(patient); err != nil {
		c.serveErr(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *ClinicAPI) putPatient(w http.ResponseWriter, r *http.Request, id int) {

	decoder := json.NewDecoder(r.Body)

	var patient postgres.Patient
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&patient); err != nil {
		c.serveErr(w, err)
		return
	}
	patient.Email.String = strings.ToLower(patient.Email.String)

	//TODO: make one func ???

	// Validate patient

	if !checkEmail(patient.Email) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if !checkPhoneNum(patient.PhoneNumber) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if !checkStructs(patient) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if patient.Email.Valid {

		err := c.DB.CheckEmailPatient(patient.Email.String)

		if !errors.Is(err, sql.ErrNoRows) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}

	}

	if patient.PhoneNumber != "" {

		err := c.DB.CheckPhoneNumberPatient(patient.PhoneNumber)

		if !errors.Is(err, sql.ErrNoRows) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}
	}
	// make one func ???

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

	decoder := json.NewDecoder(r.Body)

	var patient postgres.Patient
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&patient); err != nil {
		c.serveErr(w, err)
		return
	}

	patient.Email.String = strings.ToLower(patient.Email.String)

	// TODO: make one func ???

	// Validate patient
	if !checkEmail(patient.Email) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if patient.PhoneNumber != "" && !checkPhoneNum(patient.PhoneNumber) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	req := patchStructs(patient)

	if req == "" {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if patient.Email.Valid {

		err := c.DB.CheckEmailPatient(patient.Email.String)

		if !errors.Is(err, sql.ErrNoRows) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}

	}

	if patient.PhoneNumber != "" {

		err := c.DB.CheckPhoneNumberPatient(patient.PhoneNumber)

		if !errors.Is(err, sql.ErrNoRows) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}
	}

	//make one func ???

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
