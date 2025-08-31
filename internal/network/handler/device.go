package handler

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tiagos4ntos/device-manager/internal/domain/device"
	"github.com/tiagos4ntos/device-manager/internal/domain/device/entity"
	"github.com/tiagos4ntos/device-manager/internal/network/dto"
	errorhandler "github.com/tiagos4ntos/device-manager/internal/network/errors"
)

type DeviceHandler interface {
	List() echo.HandlerFunc
	GetByID() echo.HandlerFunc
	Create() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type deviceHandler struct {
	deviceService device.DeviceService
}

func NewDeviceHandler(service device.DeviceService) DeviceHandler {
	return &deviceHandler{
		deviceService: service,
	}

}

func (h *deviceHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		devices, err := h.deviceService.List(context.Background())

		if err != nil {
			return errorhandler.Handle(c, err)
		}

		result := make([]dto.DeviceResponse, 0)
		for _, device := range devices {
			d := dto.DeviceResponse{
				ID:        device.ID.String(),
				Name:      device.Name,
				Brand:     device.Brand,
				State:     string(device.State),
				CreatedAt: device.CreatedAt,
				UpdatedAt: device.UpdatedAt,
				DeletedAt: device.DeletedAt,
			}
			result = append(result, d)
		}

		return c.JSON(http.StatusOK, result)
	}
}

func (h *deviceHandler) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		//validate id is not empty and is a valid uuid
		deviceID, err := parseAndValidateDeviceId(id)
		if err != nil {
			return errorhandler.Handle(c, err)
		}

		// fetch device by id
		device, err := h.deviceService.GetByID(context.Background(), deviceID)

		if err != nil {
			return errorhandler.Handle(c, err)
		}

		result := dto.DeviceResponse{
			ID:        device.ID.String(),
			Name:      device.Name,
			Brand:     device.Brand,
			State:     string(device.State),
			CreatedAt: device.CreatedAt,
			UpdatedAt: device.UpdatedAt,
			DeletedAt: device.DeletedAt,
		}

		return c.JSON(http.StatusOK, result)
	}
}

func (h *deviceHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req dto.CreateDeviceRequest
		if err := c.Bind(&req); err != nil {
			return errorhandler.Handle(c, errorhandler.NewApiError(errorhandler.ErrInvalid, "you must inform all required parameters", nil))
		}

		if err := req.Validate(); err != nil {
			return errorhandler.Handle(c, errorhandler.NewApiError(errorhandler.ErrInvalid, "validation error", err))
		}

		device := entity.Device{
			Name:  req.Name,
			Brand: req.Brand,
			State: entity.DeviceState(req.State),
		}

		deviceCreated, err := h.deviceService.Create(context.Background(), device)

		if err != nil {
			return errorhandler.Handle(c, err)
		}

		result := dto.DeviceResponse{
			ID:        deviceCreated.ID.String(),
			Name:      deviceCreated.Name,
			Brand:     deviceCreated.Brand,
			State:     deviceCreated.State.String(),
			CreatedAt: deviceCreated.CreatedAt,
		}

		return c.JSON(http.StatusCreated, result)
	}
}

func (h *deviceHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		deviceID, err := parseAndValidateDeviceId(id)
		if err != nil {
			return errorhandler.Handle(c, err)
		}

		var req dto.UpdateDeviceRequest
		if err := c.Bind(&req); err != nil {
			return errorhandler.Handle(c, errorhandler.NewApiError(errorhandler.ErrInvalid, "you must inform all required parameters", err))
		}

		if err := req.Validate(); err != nil {
			return errorhandler.Handle(c, errorhandler.NewApiError(errorhandler.ErrInvalid, "validation error", err))
		}

		device := entity.Device{
			ID:    deviceID,
			Name:  req.Name,
			Brand: req.Brand,
			State: entity.DeviceState(req.State),
		}

		updatedDevice, err := h.deviceService.Update(context.Background(), device)

		if err != nil {
			return errorhandler.Handle(c, err)
		}

		result := dto.DeviceResponse{
			ID:        updatedDevice.ID.String(),
			Name:      updatedDevice.Name,
			Brand:     updatedDevice.Brand,
			State:     updatedDevice.State.String(),
			CreatedAt: updatedDevice.CreatedAt,
			UpdatedAt: updatedDevice.UpdatedAt,
			DeletedAt: updatedDevice.DeletedAt,
		}
		return c.JSON(http.StatusOK, result)
	}
}

func (h *deviceHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		//validate id is not empty and is a valid uuid
		deviceID, err := parseAndValidateDeviceId(id)
		if err != nil {
			return errorhandler.Handle(c, err)
		}

		// fetch device by id
		err = h.deviceService.Delete(context.Background(), deviceID)

		if err != nil {
			return errorhandler.Handle(c, err)
		}

		return c.JSON(http.StatusNoContent, nil)
	}
}

func parseAndValidateDeviceId(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.Nil, errorhandler.NewApiError(errorhandler.ErrInvalid, "you must inform the device id", nil)
	}
	deviceID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, errorhandler.NewApiError(errorhandler.ErrInvalid, "invalid device id format, must be an uuid", nil)
	}
	return deviceID, nil
}
