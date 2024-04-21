package domain

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrBadRequest      = errors.New("Error bad request")
	ErrInternalFailure = errors.New("Error internal failure")
	ErrNotFound        = errors.New("Error not found")
	ErrValidation      = errors.New("Error validation failed")
)

type ValidationErrors struct {
	Errors []string
}

func (e *ValidationErrors) Error() string {
	return strings.Join(e.Errors, ";")
}

func InternalFailure(err error) error {
	return fmt.Errorf("%w: %s", ErrInternalFailure, err.Error())
}

func BadRequestError(err error) error {
	return fmt.Errorf("%w: %s", ErrBadRequest, err.Error())
}

func NotFoundError(message string) error {
	return fmt.Errorf("%w: "+message, ErrNotFound)
}

func ValidationError(err *ValidationErrors) error {
	return fmt.Errorf("%w: %s", ErrValidation, err.Error())
}

func ExtractValidationErrorsToList(err error) []string {
	return strings.Split(strings.Replace(err.Error(), ErrValidation.Error()+": ", "", -1), ";")
}
