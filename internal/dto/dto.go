package dto

import (
	"encoding/json"
	"github.com/Vadim992/clinicAPI/internal/database/postgres"
	"net/http"
	"time"
)

type PageData struct {
	Page                 int            `json:"page"`
	PageSize             int            `json:"pageSize"`
	PatientsFilter       PatientsFilter `json:"patientsFilter"`
	DoctorsFilter        DoctorsFilter  `json:"doctorsFilter"`
	AppointmentsDoctorId []int          `json:"appointmentsDoctorId"`
}

func (p *PageData) DecodeFromJSON(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(p); err != nil {
		return err
	}

	return nil
}

type PatientsFilter struct {
	PhoneFilter     bool `json:"phoneFilter"`
	FirstNameFilter bool `json:"firstNameFilter"`
}

type DoctorsFilter struct {
	SpecializationFilter bool `json:"specializationFilter"`
	FirstNameFilter      bool `json:"firstNameFilter"`
}

//type AppointmentsFilter struct {
//
//}

// AppoitmentData is used for PUT, PATCH (I should know Doctor's ID and Doctor's start time of appointment)
type AppointmentsData struct {
	StartTime time.Time       `json:"startTime"`
	Record    postgres.Record `json:"record"`
}

func (a *AppointmentsData) DecodeFromJSON(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(a); err != nil {
		return err
	}
	return nil
}
