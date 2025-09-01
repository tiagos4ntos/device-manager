package device

import (
	"context"
	"database/sql"
	goerrors "errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tiagos4ntos/device-manager/internal/domain/device/entity"
	"github.com/tiagos4ntos/device-manager/internal/domain/device/errors"
	"github.com/tiagos4ntos/device-manager/internal/domain/device/repository"
)

type DeviceService interface {
	List(ctx context.Context, params map[string]any) ([]entity.Device, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.Device, error)
	Create(ctx context.Context, device entity.Device) (entity.Device, error)
	Update(ctx context.Context, device entity.Device) (entity.Device, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type deviceService struct {
	repo repository.DeviceRepository
}

func NewDeviceService(repo repository.DeviceRepository) *deviceService {
	return &deviceService{repo: repo}
}

func (s *deviceService) List(ctx context.Context, params map[string]any) ([]entity.Device, error) {
	devices, err := s.repo.ListDevices(ctx, params)
	if err != nil {
		return nil, errors.NewDeviceError(errors.ErrInternal, "something went wrong while listing devices", err)
	}
	return devices, nil
}

func (s *deviceService) GetByID(ctx context.Context, id uuid.UUID) (entity.Device, error) {
	device, err := s.repo.GetDeviceByID(ctx, id)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return entity.Device{}, errors.NewDeviceError(errors.ErrNotFound, "device not found", err)
		}

		return entity.Device{}, errors.NewDeviceError(errors.ErrInternal, "something went wrong while retrieving device", err)
	}
	return device, nil
}

func (s *deviceService) Create(ctx context.Context, device entity.Device) (entity.Device, error) {
	device.ID = uuid.New()
	err := s.repo.CreateDevice(ctx, &device)
	if err != nil {
		return device, errors.NewDeviceError(errors.ErrInternal, "something went wrong while creating device", err)
	}
	return device, nil
}

func (s *deviceService) Update(ctx context.Context, device entity.Device) (entity.Device, error) {
	//TODO: Refactor this method to avoid possible race conditions when 2 requests to update a same device are made simultaneously

	baseDevice, err := s.repo.GetDeviceByID(ctx, device.ID)
	if err != nil {
		return entity.Device{}, errors.NewDeviceError(errors.ErrNotFound, "something went wrong while retrieving device", err)
	}

	//If device is in use, only status can be updated
	if baseDevice.State == entity.InUse {
		if baseDevice.State == device.State {
			return entity.Device{}, errors.NewDeviceError(errors.ErrInvalid, "device is in use and cannot be updated", fmt.Errorf("device state is the same: %s", device.State))
		}

		//Update only status
		device, err = s.repo.UpdateDeviceState(ctx, device.ID, device.State)
		if err != nil {
			return entity.Device{}, errors.NewDeviceError(errors.ErrInternal, "something went wrong while update device state", fmt.Errorf("error updating device state: %v", err))
		}
		return device, nil
	}

	//Fully update
	err = s.repo.FullyUpdateDevice(ctx, &device)
	if err != nil {
		return entity.Device{}, errors.NewDeviceError(errors.ErrInternal, "something went wrong while fully update device", err)
	}
	return device, nil
}

func (s *deviceService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteDevice(ctx, id)
	if err != nil {
		return errors.NewDeviceError(errors.ErrInternal, "something went wrong while delete device", err)
	}
	return nil
}
