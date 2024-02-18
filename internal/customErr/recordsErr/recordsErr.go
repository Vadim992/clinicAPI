package recordsErr

import "errors"

var (
	TimeErr     = errors.New("start or end time of appointment is not valid")
	DoctorIdErr = errors.New("this doctor is not exist")
)
