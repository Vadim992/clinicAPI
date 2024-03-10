package patientserr

import "errors"

var (
	PhoneErr = errors.New("invalid patient's phone number")
)
