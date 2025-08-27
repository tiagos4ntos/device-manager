package entity

import (
	"time"

	"github.com/google/uuid"
)

type DeviceStatus string

const (
	Avaiable DeviceStatus = "available"
	InUse    DeviceStatus = "in-use"
	Inactive DeviceStatus = "inactive"
)

func (ds DeviceStatus) String() string {
	return string(ds)
}

type Device struct {
	ID           uuid.UUID    `json:"id"`
	Name         string       `json:"name"`
	Brand        string       `json:"brand"`
	Status       DeviceStatus `json:"status"`
	CreationTime time.Time    `json:"creation_time"`
}
