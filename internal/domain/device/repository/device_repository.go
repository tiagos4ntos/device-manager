package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tiagos4ntos/device-manager/internal/domain/device/entity"
)

//go:generate mockgen -source=device_repository.go -destination=../../mocks/device_repository_mock.go -package=mocks

type DeviceRepository interface {
	CreateDevice(ctx context.Context, device *entity.Device) error
	GetDeviceByID(ctx context.Context, id uuid.UUID) (entity.Device, error)
	FullyUpdateDevice(ctx context.Context, device *entity.Device) error
	UpdateDeviceState(ctx context.Context, deviceID uuid.UUID, newState entity.DeviceState) (entity.Device, error)
	DeleteDevice(ctx context.Context, id uuid.UUID) error
	ListDevices(ctx context.Context) ([]entity.Device, error)
}

type postegresDeviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *postegresDeviceRepository {
	return &postegresDeviceRepository{db: db}
}

func (r *postegresDeviceRepository) CreateDevice(ctx context.Context, device *entity.Device) error {
	const query = `
	INSERT INTO devices (id, name, brand, state)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at, deleted_at;`

	statment, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer statment.Close()

	err = statment.
		QueryRowContext(ctx,
			device.ID,
			device.Name,
			device.Brand,
			device.State.String(),
		).
		Scan(&device.ID, &device.CreatedAt, &device.UpdatedAt, &device.DeletedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *postegresDeviceRepository) GetDeviceByID(ctx context.Context, id uuid.UUID) (entity.Device, error) {
	var device entity.Device

	query := `
	SELECT id, name, brand, state, created_at, updated_at, deleted_at
	FROM devices
	WHERE id = $1 AND deleted_at IS NULL;`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return device, err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(
		ctx,
		id.String(),
	).Scan(
		&device.ID,
		&device.Name,
		&device.Brand,
		&device.State,
		&device.CreatedAt,
		&device.UpdatedAt,
		&device.DeletedAt)
	if err != nil {
		return device, err
	}

	return device, nil
}

func (r *postegresDeviceRepository) FullyUpdateDevice(ctx context.Context, device *entity.Device) error {
	query := `
	UPDATE devices SET 
		name = $2,
		brand = $3,
		state = $4,
		updated_at = now()
	WHERE id = $1 AND deleted_at IS NULL
	RETURNING id, created_at, updated_at, deleted_at;`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(
		ctx,
		device.ID,
		device.Name,
		device.Brand,
		device.State.String(),
	).Scan(
		&device.ID,
		&device.CreatedAt,
		&device.UpdatedAt,
		&device.DeletedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *postegresDeviceRepository) UpdateDeviceState(ctx context.Context, deviceID uuid.UUID, newStatus entity.DeviceState) (entity.Device, error) {
	var device entity.Device
	query := `
	UPDATE devices SET 
		state = $2,
		updated_at = now()
	WHERE id = $1 AND deleted_at IS NULL
	RETURNING id, name, brand, state, created_at, updated_at, deleted_at;`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return device, err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(
		ctx,
		deviceID.String(),
		newStatus.String(),
	).Scan(
		&device.ID,
		&device.Name,
		&device.Brand,
		&device.State,
		&device.CreatedAt,
		&device.UpdatedAt,
		&device.DeletedAt,
	)

	if err != nil {
		return device, err
	}

	return device, nil
}

func (r *postegresDeviceRepository) DeleteDevice(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE devices SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL AND state <> 'in-use';`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(
		ctx,
		id)
	if err != nil {
		return err
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowCount <= 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *postegresDeviceRepository) ListDevices(ctx context.Context) ([]entity.Device, error) {
	var devices []entity.Device

	query := `SELECT id, name, brand, state, created_at, updated_at, deleted_at
	FROM devices
	WHERE deleted_at IS NULL
	ORDER BY name;`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d entity.Device
		err = rows.Scan(
			&d.ID,
			&d.Name,
			&d.Brand,
			&d.State,
			&d.CreatedAt,
			&d.UpdatedAt,
			&d.DeletedAt)

		if err != nil {
			return nil, err
		}

		devices = append(devices, d)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}
