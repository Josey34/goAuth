package handler

import (
	"goauth/dto"

	"github.com/go-playground/validator/v10"
)

func extractFieldErrors(err error) []dto.FieldError {
	var fieldErrors []dto.FieldError
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, ve := range validationErrors {
			fieldErrors = append(fieldErrors, dto.FieldError{
				Field:   ve.Field(),
				Message: ve.Tag() + " validation failed",
			})
		}
	}
	return fieldErrors
}
