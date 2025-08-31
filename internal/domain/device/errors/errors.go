package errors

import "fmt"

type DeviceErrorType string

const (
	ErrNotFound DeviceErrorType = "not_found"
	ErrInvalid  DeviceErrorType = "invalid"
	ErrInternal DeviceErrorType = "internal"
)

type DeviceError struct {
	Type    DeviceErrorType
	Message string
	Err     error
}

func (e *DeviceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewDeviceError(t DeviceErrorType, msg string, err error) *DeviceError {
	return &DeviceError{
		Type:    t,
		Message: msg,
		Err:     err,
	}
}
