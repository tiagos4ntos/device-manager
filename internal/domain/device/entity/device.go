package entity

import (
	"time"

	"github.com/google/uuid"
)

type DeviceStatus string

const (
	Available DeviceStatus = "available"
	InUse     DeviceStatus = "in-use"
	Inactive  DeviceStatus = "inactive"
)

func (ds DeviceStatus) String() string {
	return string(ds)
}

type Device struct {
	ID        uuid.UUID    `json:"id"`
	Name      string       `json:"name"`
	Brand     string       `json:"brand"`
	State     DeviceStatus `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt *time.Time   `json:"updated_at"`
	DeletedAt *time.Time   `json:"deleted_at"`
}
