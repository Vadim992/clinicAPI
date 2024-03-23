package postgres

import (
	"encoding/json"
	"net/http"
	"time"
)

// "Email" fields contains in map in /internal/helpers/dbhelpers/dbchecker.go
// in func validateStruct

// "PatientId" fields contains in map in /internal/helpers/dbhelpers/dbchecker.go
// in func validateStruct

type Record struct {
	DoctorId   *int       `json:"doctorId"`
	PatientId  *int       `json:"patientId"` // may be null in DB
	Time_Start *time.Time `json:"time_start"`
	Time_End   *time.Time `json:"time_end"`
}

func (rec *Record) DecodeFromJSON(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(rec); err != nil {
		return err
	}
	return nil
}

type Appointment struct {
	DoctorFirstName      string    `json:"doctorFirstName"`
	DoctorLastName       string    `json:"doctorLastName"`
	DoctorSpecialization string    `json:"doctorSpecialization"`
	TimeStart            time.Time `json:"timeStart"`
}

type AppointmentsCount struct {
	DoctorId             int    `json:"doctorId"`
	DoctorFirstName      string `json:"doctorFirstName"`
	DoctorLastName       string `json:"doctorLastName"`
	DoctorSpecialization string `json:"doctorSpecialization"`
	Appointments         int    `json:"count"`
}
