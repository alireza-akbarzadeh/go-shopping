package controllers

import "github.com/go-playground/validator/v10"

// formatValidationErrors converts validator.ValidationErrors to a map.
func formatValidationErrors(err error) map[string]string {
	errorsMap := make(map[string]string)
	for _, ve := range err.(validator.ValidationErrors) {
		errorsMap[ve.Field()] = ve.Tag()
	}
	return errorsMap
}
