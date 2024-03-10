package dtoerr

import "errors"

var (
	PageErr = errors.New("incorrect number of page")
	//FilterErr = errors.New("cannot use phone filter and email filter at the same time")
)
