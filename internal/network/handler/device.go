package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tiagos4ntos/device-manager/internal/domain/device"
	"github.com/tiagos4ntos/device-manager/internal/domain/device/entity"
	"github.com/tiagos4ntos/device-manager/internal/network/dto"
)

type DeviceHandler interface {
	List() echo.HandlerFunc
	Create() echo.HandlerFunc
}

type deviceHandler struct {
	deviceService device.DeviceService
}

func NewDeviceHandler(service device.DeviceService) DeviceHandler {
	return &deviceHandler{
		deviceService: service,
	}

}

func (h *deviceHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req dto.CreateDeviceRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "you must inform all required parameters")
		}

		device := entity.Device{
			Name:  req.Name,
			Brand: req.Brand,
			State: entity.DeviceStatus(req.State),
		}

		deviceCreated, err := h.deviceService.Create(context.Background(), device)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		result := dto.DeviceResponse{
			ID:        deviceCreated.ID.String(),
			Name:      deviceCreated.Name,
			Brand:     deviceCreated.Brand,
			State:     string(deviceCreated.State),
			CreatedAt: device.CreatedAt,
		}

		return c.JSON(http.StatusCreated, result)
	}
}

func (h *deviceHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		devices, err := h.deviceService.List(context.Background())

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
