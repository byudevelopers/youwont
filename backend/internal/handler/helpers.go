package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"youwont.api/internal/service"
)

func handleError(c *echo.Context, err error) error {
	code := http.StatusInternalServerError
	errCode := "INTERNAL"

	switch {
	case errors.Is(err, service.ErrNotFound):
		code = http.StatusNotFound
		errCode = "NOT_FOUND"
	case errors.Is(err, service.ErrForbidden):
		code = http.StatusForbidden
		errCode = "FORBIDDEN"
	case errors.Is(err, service.ErrAlreadyExists):
		code = http.StatusConflict
		errCode = "ALREADY_EXISTS"
	case errors.Is(err, service.ErrInsufficientPoints):
		code = http.StatusBadRequest
		errCode = "INSUFFICIENT_POINTS"
	case errors.Is(err, service.ErrBetNotOpen):
		code = http.StatusBadRequest
		errCode = "BET_NOT_OPEN"
	case errors.Is(err, service.ErrAlreadyWagered):
		code = http.StatusConflict
		errCode = "ALREADY_WAGERED"
	case errors.Is(err, service.ErrDeciderCannotWager):
		code = http.StatusForbidden
		errCode = "DECIDER_CANNOT_WAGER"
	case errors.Is(err, service.ErrCannotSelfDecide):
		code = http.StatusBadRequest
		errCode = "CANNOT_SELF_DECIDE"
	case errors.Is(err, service.ErrNoOpposingSide):
		code = http.StatusBadRequest
		errCode = "NO_OPPOSING_SIDE"
	}

	return c.JSON(code, map[string]interface{}{
		"error": map[string]string{
			"code":    errCode,
			"message": err.Error(),
		},
	})
}

func badRequest(c *echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"error": map[string]string{
			"code":    "BAD_REQUEST",
			"message": message,
		},
	})
}
