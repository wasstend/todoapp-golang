package core_errors

import "errors"

var (
	ErrorNotFound        = errors.New("Not found")
	ErrorInvalidArgument = errors.New("Invalid argument")
	ErrorConflict        = errors.New("Conflict")
)
