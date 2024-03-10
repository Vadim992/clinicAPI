package structserr

import "errors"

var (
	EmptyFieldErr          = errors.New("the field is empty")
	InvalidTypeOfStructErr = errors.New("invalid type of struct")
	NoDataForPatchReqErr   = errors.New("no data for parch request")
	EmailErr               = errors.New("invalid email")
)
