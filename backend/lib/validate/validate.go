package validate

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type HttpValidator struct {
	validator *validator.Validate
}

func (hv *HttpValidator) Validate(i any) error {
	if err := hv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func NewHttpValidator(val *validator.Validate) *HttpValidator {
	return &HttpValidator{validator: val}
}
