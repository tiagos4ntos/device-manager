package errors

import "fmt"

type ApiErrorType string

const (
	ErrNotFound ApiErrorType = "not_found"
	ErrInvalid  ApiErrorType = "invalid"
)

type ApiError struct {
	Type    ApiErrorType
	Message string
	Err     error
}

func (e *ApiError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewApiError(t ApiErrorType, msg string, err error) *ApiError {
	return &ApiError{
		Type:    t,
		Message: msg,
		Err:     err,
	}
}
