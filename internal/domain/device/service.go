package device

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/tiagos4ntos/device-manager/internal/domain/device/entity"
	"github.com/tiagos4ntos/device-manager/internal/domain/device/repository"
)

type DeviceService interface {
	List(ctx context.Context) ([]entity.Device, error)
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

func (s *deviceService) List(ctx context.Context) ([]entity.Device, error) {
	devices, err := s.repo.ListDevices(ctx)
	if err != nil {
		log.Printf("Error listing devices: (%v)", err)
		return nil, errors.New("something went wrong while listing devices")
	}
	return devices, nil
}

func (s *deviceService) GetByID(ctx context.Context, id uuid.UUID) (entity.Device, error) {
	device, err := s.repo.GetDeviceByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("device not found: (%v)", id)
			return entity.Device{}, errors.New("device not found")
		}

		log.Printf("error when search device by id: (%v)", id)
		return entity.Device{}, errors.New("error on searching device by id")
	}
	return device, nil
}

func (s *deviceService) Create(ctx context.Context, device entity.Device) (entity.Device, error) {
	device.ID = uuid.New()
	err := s.repo.CreateDevice(ctx, &device)
	if err != nil {
		log.Printf("error creating device: (%v)", err)
		return device, errors.New("something went wrong while create device")
	}
	return device, nil
}

func (s *deviceService) Update(ctx context.Context, device entity.Device) (entity.Device, error) {
	//TODO: Refactor this method to avoid possible race conditions when 2 requests to update a same device are made simultaneously

	baseDevice, err := s.repo.GetDeviceByID(ctx, device.ID)
	if err != nil {
		log.Printf("error getting device by id: (%v)", err)
		return entity.Device{}, errors.New("something went wrong while retrieving device")
	}

	//If device is in use, only status can be updated
	if baseDevice.State == entity.InUse {
		if baseDevice.State == device.State {
			return entity.Device{}, errors.New("device is in use and cannot be updated")
		}

		//Update only status
		device, err = s.repo.UpdateDeviceState(ctx, device.ID, device.State)
		if err != nil {
			log.Printf("error updating device state: (%v)", err)
			return entity.Device{}, errors.New("something went wrong while update device state")
		}
		return device, nil
	}

	//Fully update
	err = s.repo.FullyUpdateDevice(ctx, &device)
	if err != nil {
		log.Printf("error fully updating device: (%v)", err)
		return entity.Device{}, errors.New("something went wrong while fully update device")
	}
	return device, nil
}

func (s *deviceService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteDevice(ctx, id)
	if err != nil {
		log.Printf("error deleting device: (%v)", err)
		return errors.New("something went wrong while delete device")
	}
	return nil
}
