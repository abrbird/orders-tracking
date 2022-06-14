package models

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	UnknownErrorMessage            = "unknown error"
	NotFoundErrorMessage           = "not found error"
	AlreadyExistsErrorMessage      = "already exists error"
	OrderIsOnIssuingMessage        = "order is on issuing"
	OrderAlreadyIssuedMessage      = "order is already issued"
	OrderIsNotReadyForIssueMessage = "order is not ready for issue"

	RetryErrorMessage       = "please retry"
	MaxAttemptsErrorMessage = "reached max attempts: %d"

	BrokerSendErrorMessage = "can not send message"
)

type retryError struct {
	Inner   error
	Message string
}

func NewRetryError(err error) error {
	return retryError{err, RetryErrorMessage}
}

func (e retryError) Error() string {
	return wrapError(e.Inner, e.Message).Error()
}

func (e retryError) Unwrap() error {
	return e.Inner
}

func wrapError(err error, message string) error {
	if err == nil {
		return errors.New(message)
	}
	return errors.Wrap(err, message)
}

func UnknownError(err error) error {
	return wrapError(err, UnknownErrorMessage)
}
func NotFoundError(err error) error {
	return wrapError(err, NotFoundErrorMessage)
}
func AlreadyExistsError(err error) error {
	return wrapError(err, AlreadyExistsErrorMessage)
}
func OrderIsOnIssuingError(err error) error {
	return wrapError(err, OrderIsOnIssuingMessage)
}
func OrderAlreadyIssuedError(err error) error {
	return wrapError(err, OrderAlreadyIssuedMessage)
}
func OrderIsNotReadyForIssueError(err error) error {
	return wrapError(err, OrderIsNotReadyForIssueMessage)
}
func MaxAttemptsError(err error, maxAttempts int64) error {
	return wrapError(err, fmt.Sprintf(MaxAttemptsErrorMessage, maxAttempts))
}
func BrokerSendError(err error) error {
	return wrapError(err, BrokerSendErrorMessage)
}

var RetryError = NewRetryError(nil)
