package dto

import (
	"fmt"
	"time"
)

var validDeviceStatuses = map[string]bool{
	"available": true,
	"in-use":    true,
	"inactive":  true,
}

type CreateDeviceRequest struct {
	Name  string `json:"name" validate:"required"`
	Brand string `json:"brand" validate:"required"`
	State string `json:"state" validate:"required,oneof=available in-use inactive"`
}

func (r CreateDeviceRequest) Validate() error {
	if !validDeviceStatuses[r.State] {
		return fmt.Errorf("invalid device status: %s", r.State)
	}
	return nil
}

type DeviceResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Brand     string     `json:"brand"`
	State     string     `json:"state"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type UpdateDeviceRequest struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Brand string `json:"brand"`
	State string `json:"state"`
}
