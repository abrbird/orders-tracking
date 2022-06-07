package models

import "errors"

var (
	UnknownError       = errors.New("unknown error")
	NotFoundError      = errors.New("not found error")
	AlreadyExistsError = errors.New("already exists error")
)
