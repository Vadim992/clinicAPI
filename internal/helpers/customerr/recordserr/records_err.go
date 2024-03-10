package recordserr

import "errors"

var (
	TimeErr      = errors.New("start or end time of appointment is not valid")
	DoctorIdErr  = errors.New("this doctor is not exist")
	PatientIdErr = errors.New("this patient is not exist")
	IdErr        = errors.New("zero or negative patient's or doctor's id")
	InsertErr    = errors.New("this data already exists")
	InsertIdErr  = errors.New("dont need id for request") // id is unique and auto increment
)
