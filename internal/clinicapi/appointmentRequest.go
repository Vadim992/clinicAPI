package clinicapi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Vadim992/clinicAPI/internal/customErr/recordsErr"
	"github.com/Vadim992/clinicAPI/pkg/database/postgres"
	"io"
	"net/http"
)

func (c *ClinicAPI) getAppointments(w http.ResponseWriter, r *http.Request) {
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

	offset := (pageData.Page - 1) * pageData.PageSize

	// In future we can filter datas

	var filter string

	appointments, err := c.DB.GetAppointments(offset, pageData.PageSize, filter)

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

	offset := (pageData.Page - 1) * pageData.PageSize

	appointments, err := c.DB.GetAppointmentsId(offset, pageData.PageSize, id)

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

	if err := decode(r, &record); err != nil {
		if errors.Is(err, io.EOF) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}
		c.serveErr(w, err)
		return
	}

	if !putCheckStructs(record) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if err := c.validateRecordPOST(record); err != nil {

		switch true {
		case errors.Is(err, recordsErr.InsertErr):

			c.clientErr(w, http.StatusBadRequest)
			return

		case errors.Is(err, recordsErr.IdErr):

			c.clientErr(w, http.StatusBadRequest)
			return

		case errors.Is(err, recordsErr.TimeErr):

			c.clientErr(w, http.StatusBadRequest)
			return

		case errors.Is(err, recordsErr.DoctorIdErr):

			c.clientErr(w, http.StatusBadRequest)
			return

		case errors.Is(err, recordsErr.PatientIdErr):

			c.clientErr(w, http.StatusBadRequest)
			return

		default:
			c.serveErr(w, err)
			return
		}

	}

	if err := c.DB.InsertAppointment(record); err != nil {
		c.serveErr(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)

}

// TODO: PUT, PATCH, DELETE for appointments

func (c *ClinicAPI) putAppointment(w http.ResponseWriter, r *http.Request, id int) {

	var appointmentData AppointmentsData

	if err := decode(r, &appointmentData); err != nil {
		if errors.Is(err, io.EOF) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}
		c.serveErr(w, err)
		return
	}

	if !putCheckStructs(appointmentData.Record) {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	if err := c.validateRecordPUT(id, appointmentData); err != nil {

		switch true {
		case errors.Is(err, recordsErr.IdErr):

			c.clientErr(w, http.StatusBadRequest)
			return

		case errors.Is(err, recordsErr.TimeErr):

			c.clientErr(w, http.StatusBadRequest)
			return

		case errors.Is(err, recordsErr.DoctorIdErr):

			c.clientErr(w, http.StatusBadRequest)
			return

		case errors.Is(err, recordsErr.PatientIdErr):

			c.clientErr(w, http.StatusBadRequest)
			return

		case errors.Is(err, sql.ErrNoRows):

			c.clientErr(w, http.StatusBadRequest)
			return

		default:
			c.serveErr(w, err)
			return
		}
	}

	start := appointmentData.StartTime
	record := appointmentData.Record

	if err := c.DB.UpdateAppointmentAll(id, start, record); err != nil {
		c.serveErr(w, err)
		return
	}

}

func (c *ClinicAPI) patchAppointment(w http.ResponseWriter, r *http.Request, id int) {
	var appointmentData AppointmentsData

	if err := decode(r, &appointmentData); err != nil {
		if errors.Is(err, io.EOF) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}
		c.serveErr(w, err)
		return
	}

	if err := c.validateRecordPATCH(id, appointmentData); err != nil {

		switch true {

		case errors.Is(err, recordsErr.IdErr):

			c.clientErr(w, http.StatusBadRequest)
			return

		case errors.Is(err, recordsErr.TimeErr):

			c.clientErr(w, http.StatusBadRequest)
			return

		case errors.Is(err, recordsErr.DoctorIdErr):

			c.clientErr(w, http.StatusBadRequest)
			return

		case errors.Is(err, recordsErr.PatientIdErr):

			c.clientErr(w, http.StatusBadRequest)
			return

		case errors.Is(err, sql.ErrNoRows):

			c.clientErr(w, http.StatusBadRequest)
			return

		default:
			c.serveErr(w, err)
			return
		}
	}
	record := appointmentData.Record
	req := patchCheckStructs(record)

	if req == "" {
		c.clientErr(w, http.StatusBadRequest)
		return
	}

	start := appointmentData.StartTime

	if err := c.DB.UpdateAppointment(id, start, req); err != nil {
		c.serveErr(w, err)
		return
	}
}

func (c *ClinicAPI) deleteAppointment(w http.ResponseWriter, r *http.Request, id int) {

	var appointmentData AppointmentsData

	if err := decode(r, &appointmentData); err != nil {
		if errors.Is(err, io.EOF) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}
		c.serveErr(w, err)
		return
	}

	if err := c.DB.DeleteAppointment(id, appointmentData.StartTime); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.clientErr(w, http.StatusBadRequest)
			return
		}

		c.serveErr(w, err)
		return
	}
}
