package postgres

import (
	"encoding/json"
	"net/http"
	"time"
)

// "Email" fields contains in map in /internal/helpers/dbhelpers/dbchecker.go
// in func validateStruct
type Patient struct {
	Id           *int    `json:"id"`
	FirstName    *string `json:"firstName"`
	LastName     *string `json:"lastName"`
	Email        *string `json:"email"` // may be null in DB
	Address      *string `json:"address"`
	Phone_Number *string `json:"phone_number"`
}

func (p *Patient) DecodeFromJSON(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(p); err != nil {
		return err
	}
	return nil
}

type Doctor struct {
	Id             *int    `json:"id"`
	FirstName      *string `json:"firstName"`
	LastName       *string `json:"lastName"`
	Specialization *string `json:"specialization"`
	Room           *int    `json:"room"`
	Email          *string `json:"email"`
}

func (d *Doctor) DecodeFromJSON(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(d); err != nil {
		return err
	}
	return nil
}

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
