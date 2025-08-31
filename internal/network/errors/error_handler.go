package errors

import (
	"net/http"

	"github.com/labstack/echo/v4"
	deviceerrors "github.com/tiagos4ntos/device-manager/internal/domain/device/errors"
)

func Handle(c echo.Context, err error) error {

	c.Logger().Error(err)

	switch e := err.(type) {
	case *ApiError:
		return c.JSON(mapApiErrorsToStatusCode(e.Type), ErrorResponse(e.Message))
	case *deviceerrors.DeviceError:
		return c.JSON(mapDomainErrorsToStatusCode(e.Type), ErrorResponse(e.Message))
	default:
		return c.JSON(http.StatusInternalServerError, ErrorResponse("Internal server error"))
	}
}

func mapDomainErrorsToStatusCode(t deviceerrors.DeviceErrorType) int {
	switch t {
	case deviceerrors.ErrNotFound:
		return http.StatusNotFound
	case deviceerrors.ErrInvalid:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func mapApiErrorsToStatusCode(t ApiErrorType) int {
	switch t {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrInvalid:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func ErrorResponse(msg string) map[string]string {
	return map[string]string{"error": msg}
}
