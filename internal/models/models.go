package models

import (
	"fmt"
)

var (
	UnknownErrorMessage       = "unknown error"
	NotFoundErrorMessage      = "not found error"
	AlreadyExistsErrorMessage = "already exists error"
)

func NotFoundError(err error) error {
	return fmt.Errorf(NotFoundErrorMessage, err)
}

func UnknownError(err error) error {
	return fmt.Errorf(UnknownErrorMessage, err)
}
func AlreadyExistsError(err error) error {
	return fmt.Errorf(AlreadyExistsErrorMessage, err)
}
