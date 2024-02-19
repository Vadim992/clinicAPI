package postgres

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type Patient struct {
	Id           int            `json:"id"`
	FirstName    string         `json:"firstName"`
	LastName     string         `json:"lastName"`
	Email        sql.NullString `json:"email"`
	Address      string         `json:"address"`
	Phone_Number string         `json:"phone_number"`
}

func (p *Patient) Decode(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(p); err != nil {
		return err
	}
	return nil
}

type Doctor struct {
	Id             int    `json:"id"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Specialization string `json:"specialization"`
	Room           int    `json:"room"`
	Email          string `json:"email"`
}

func (d *Doctor) Decode(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(d); err != nil {
		return err
	}
	return nil
}

type Record struct {
	DoctorId   int           `json:"doctorId"`
	PatientId  sql.NullInt64 `json:"patientId"`
	Time_Start time.Time     `json:"time_start"`
	Time_End   time.Time     `json:"time_end"`
}

func (rec *Record) Decode(r *http.Request) error {
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
