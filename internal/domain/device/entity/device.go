package entity

import (
	"time"

	"github.com/google/uuid"
)

type DeviceState string

const (
	Available DeviceState = "available"
	InUse     DeviceState = "in-use"
	Inactive  DeviceState = "inactive"
)

func (ds DeviceState) String() string {
	return string(ds)
}

type Device struct {
	ID        uuid.UUID   `json:"id"`
	Name      string      `json:"name"`
	Brand     string      `json:"brand"`
	State     DeviceState `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt *time.Time  `json:"updated_at"`
	DeletedAt *time.Time  `json:"deleted_at"`
}
