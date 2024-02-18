package clinicapi

import (
	"github.com/Vadim992/clinicAPI/pkg/validator/validate"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var pageSize = 3

func (c *ClinicAPI) handlePatients(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:

		c.getPatients(w, r)

	case http.MethodPost:

		c.postPatient(w, r)

	}
}

func (c *ClinicAPI) handlePatientId(w http.ResponseWriter, r *http.Request) {

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.serveErr(w, err)
		return
	}

	if !validate.ValidateId(id) {
		c.notFound(w)
		return
	}

	switch r.Method {
	case http.MethodGet:

		c.getPatientId(w, r, id)

	case http.MethodPut:

		c.putPatient(w, r, id)

	case http.MethodPatch:

		c.patchPatient(w, r, id)

	case http.MethodDelete:

		c.deletePatient(w, r, id)

	}
}

func (c *ClinicAPI) handleDoctors(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:

		c.getDoctors(w, r)

	case http.MethodPost:

		c.postDoctor(w, r)
	}
}

func (c *ClinicAPI) handleDoctorId(w http.ResponseWriter, r *http.Request) {

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.serveErr(w, err)
		return
	}

	if !validate.ValidateId(id) {
		c.notFound(w)
		return
	}

	switch r.Method {

	case http.MethodGet:

		c.getDoctorId(w, r, id)

	case http.MethodPut:

		c.putDoctor(w, r, id)

	case http.MethodPatch:

		c.patchDoctor(w, r, id)

	case http.MethodDelete:

		c.deleteDoctor(w, r, id)

	}
}

func (c *ClinicAPI) handleAppointments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.getAppointments(w, r)
	case http.MethodPost:
		c.postAppointment(w, r)
	}
}

func (c *ClinicAPI) handleAppointmentsId(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.serveErr(w, err)
		return
	}

	if !validate.ValidateId(id) {
		c.notFound(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		c.getAppointmentsId(w, r, id)
	case http.MethodPut:
	case http.MethodPatch:
	case http.MethodDelete:

	}
}
