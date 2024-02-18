package clinicapi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Vadim992/clinicAPI/internal/customErr/recordsErr"
	"github.com/Vadim992/clinicAPI/pkg/database/postgres"
	"github.com/Vadim992/clinicAPI/pkg/validator/validate"
	"net/http"
)

func (c *ClinicAPI) getAppointments(w http.ResponseWriter, r *http.Request) {
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

	offset := (pageData.Page - 1) * pageSize

	// In future we can filter datas

	var filter string

	appointments, err := c.DB.GetAppointments(offset, pageSize, filter)

	if err != nil {
		c.serveErr(w, err)
		return
	}

	b, err := json.Marshal(appointments)

	if err != nil {
		c.serveErr(w, err)
		return
	}

	w.Write(b)
}

func (c *ClinicAPI) getAppointmentsId(w http.ResponseWriter, r *http.Request, id int) {
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

	offset := (pageData.Page - 1) * pageSize

	appointments, err := c.DB.GetAppointmentsId(offset, pageSize, id)

	if err != nil {

		if !errors.Is(err, sql.ErrNoRows) {
			c.notFound(w)
			return
		}

		c.serveErr(w, err)
		return
	}

	b, err := json.Marshal(appointments)

	if err != nil {
		c.serveErr(w, err)
		return
	}

	w.Write(b)
}

func (c *ClinicAPI) postAppointment(w http.ResponseWriter, r *http.Request) {

	var record postgres.Record
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&record); err != nil {
		c.serveErr(w, err)
		return
	}

	if !validate.ValidateId(record.DoctorId) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if err := c.DB.ValidateRecord(record); err != nil {

		switch true {
		case errors.Is(err, recordsErr.TimeErr):
			c.clientErr(w, http.StatusBadRequest)
		case errors.Is(err, recordsErr.DoctorIdErr):
			c.clientErr(w, http.StatusBadRequest)
		default:
			c.serveErr(w, err)
		}
		return

	}

	if err := c.DB.InsertAppointment(record); err != nil {
		c.serveErr(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)

}

// TODO: PUT, PATCH, DELETE for appointments

func (c *ClinicAPI) putAppointment(w http.ResponseWriter, r *http.Request) {

}

func (c *ClinicAPI) patchAppointment(w http.ResponseWriter, r *http.Request) {

}

func (c *ClinicAPI) deleteAppointment(w http.ResponseWriter, r *http.Request) {

}
