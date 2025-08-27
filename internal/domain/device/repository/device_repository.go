package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/tiagos4ntos/device-manager/internal/domain/device/entity"
)

//go:generate mockgen -source=device.go -destination=../../mocks/device_repository_mock.go -package=mocks

type DeviceRepository interface {
	CreateDevice(ctx context.Context, device *entity.Device) error
	GetDeviceByID(ctx context.Context, id string) (*entity.Device, error)
	UpdateDevice(ctx context.Context, device *entity.Device) error
	DeleteDevice(ctx context.Context, id string) error
	ListDevices(ctx context.Context) ([]*entity.Device, error)
}

type postegresDeviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *postegresDeviceRepository {
	return &postegresDeviceRepository{db: db}
}

func (r *postegresDeviceRepository) CreateDevice(ctx context.Context, device *entity.Device) error {
	if device == nil {
		return errors.New("device is nil")
	}

	const query = `
	INSERT INTO devices (id, name, brand, status)
	VALUES ($1, $2, $3, $4)
	RETURNING id, creation_time;`

	statment, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer statment.Close()

	err = statment.
		QueryRowContext(ctx,
			device.ID,
			device.Name,
			device.Brand,
			device.Status.String(),
		).
		Scan(&device.ID, &device.CreationTime)

	if err != nil {
		return fmt.Errorf("scan row: %w", err)
	}

	return nil
}

func (r *postegresDeviceRepository) GetDeviceByID(ctx context.Context, id string) (*entity.Device, error) {
	//TODO: Implementation for retrieving a device by ID from Database
	return nil, nil
}

func (r *postegresDeviceRepository) UpdateDevice(ctx context.Context, device *entity.Device) error {
	//TODO: Implementation for updating a device in Database
	return nil
}

func (r *postegresDeviceRepository) DeleteDevice(ctx context.Context, id string) error {
	//TODO: Implementation for deleting a device in Database
	return nil
}

func (r *postegresDeviceRepository) ListDevices(ctx context.Context) ([]*entity.Device, error) {
	//TODO: Implementation for listing all devices from Database
	return nil, nil
}
