package clinicapi

import (
	"github.com/Vadim992/clinicAPI/pkg/database/postgres"
	"log"
	"time"
)

type ClinicAPI struct {
	ErrLog  *log.Logger
	InfoLog *log.Logger
	DB      *postgres.DB
}

type PageData struct {
	Page           int            `json:"page"`
	PatientsFilter PatientsFilter `json:"patientsFilter"`
	DoctorsFilter  DoctorsFilter  `json:"doctorsFilter"`
}

type PatientsFilter struct {
	PhoneFilter     bool `json:"phoneFilter"`
	FirstNameFilter bool `json:"firstNameFilter"`
}

type DoctorsFilter struct {
	SpecializationFilter bool `json:"specializationFilter"`
	FirstNameFilter      bool `json:"firstNameFilter"`
}

type AppointmentsData struct {
	startTime time.Time `json:"startTime"`
}
