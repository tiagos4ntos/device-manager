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
	Name  string `json:"name" validate:"required" example:"Moto G100"`
	Brand string `json:"brand" validate:"required" example:"Motorola"`
	State string `json:"state" validate:"required,oneof=available in-use inactive" example:"available"`
}

func (r CreateDeviceRequest) Validate() error {
	if !validDeviceStatuses[r.State] {
		return fmt.Errorf("invalid device state %s", r.State)
	}
	return nil
}

type DeviceResponse struct {
	ID        string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string     `json:"name" example:"iPhone 13"`
	Brand     string     `json:"brand" example:"Apple"`
	State     string     `json:"state" example:"available"`
	CreatedAt time.Time  `json:"created_at" example:"2025-08-31T21:00:00Z"`
	UpdatedAt *time.Time `json:"updated_at" example:"2025-08-31T21:00:00Z"`
	DeletedAt *time.Time `json:"deleted_at" example:"null"`
}

type UpdateDeviceRequest struct {
	Name  string `json:"name" example:"Galaxy S21"`
	Brand string `json:"brand" example:"Samsung"`
	State string `json:"state" validate:"required,oneof=available in-use inactive" example:"in-use"`
}

func (r UpdateDeviceRequest) Validate() error {
	if !validDeviceStatuses[r.State] {
		return fmt.Errorf("invalid device state %s", r.State)
	}
	return nil
}
