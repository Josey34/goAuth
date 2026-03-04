package errors

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthError struct {
	Message string
}

func (e AuthError) Error() string {
	return e.Message
}

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return e.Message
}

type ConflictError struct {
	Message string
}

func (e ConflictError) Error() string {
	return e.Message
}

func ToHTTPStatus(err error) int {
	var authErr AuthError
	if errors.As(err, &authErr) {
		return http.StatusUnauthorized
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return http.StatusBadRequest
	}

	var valErr ValidationError
	if errors.As(err, &valErr) {
		return http.StatusBadRequest
	}

	var notFound NotFoundError
	if errors.As(err, &notFound) {
		return http.StatusNotFound
	}

	var conflict ConflictError
	if errors.As(err, &conflict) {
		return http.StatusConflict
	}

	return http.StatusInternalServerError
}
