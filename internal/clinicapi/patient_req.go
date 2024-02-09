package clinicapi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Vadim992/clinicAPI/pkg/database/postgres"
	"net/http"
)

func (c *ClinicAPI) getPatients(w http.ResponseWriter, r *http.Request) {

	patients, err := c.DB.GetPatients()

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

	//TODO: make one func ???

	// Validate patient

	if !checkEmail(patient.Email) {
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

	//TODO: make one func ???

	// Validate patient

	if !checkEmail(patient.Email) {
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

	// TODO: make one func ???

	// Validate patient
	if !checkEmail(patient.Email) {
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
