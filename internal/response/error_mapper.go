package response

import (
	"errors"
	"net/http"

	"github.com/mjhddev/ticket-booking-api/internal/errs"
)

func StatusCode(err error) int {
	switch {
	case errors.Is(err, errs.ErrEmailAlreadyExists):
		return http.StatusConflict

	case errors.Is(err, errs.ErrInvalidCredentials):
		return http.StatusUnauthorized

	default:
		return http.StatusInternalServerError
	}
}
