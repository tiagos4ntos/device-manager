package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

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

// List godoc
// @Summary List devices
// @Description Get all devices
// @Tags devices
// @Accept json
// @Produce json
// @Param        brand   path      string  false  "Brand name: eg. Apple"
// @Param        state   path      string  false  "State, must be one of: available, in-use, inactive"
// @Success 200 {array} dto.DeviceResponse
// @Failure 500 {object} errors.DefaultErrorResult
// @Router /devices [get]
func (h *deviceHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		params, err := validateAndParseListParams(c)
		if err != nil {
			return errorhandler.Handle(c, err)
		}

		devices, err := h.deviceService.List(context.Background(), params)

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

// GetDeviceByID godoc
// @Summary      Get device by ID
// @Description  Returns a single device by its ID
// @Tags         devices
// @Produce      json
// @Param        id   path      string  true  "Device ID"
// @Success      200  {object}  dto.DeviceResponse
// @Failure      400  {object}  errors.DefaultErrorResult
// @Failure      404  {object}  errors.DefaultErrorResult
// @Failure      500  {object}  errors.DefaultErrorResult
// @Router       /devices/{id} [get]
func (h *deviceHandler) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		//validate id is not empty and is a valid uuid
		deviceID, err := validateAndParseDeviceId(id)
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

// CreateDevice godoc
// @Summary      Create a new device
// @Description  Registers a new device on the database with the provided information
// @Tags         devices
// @Accept       json
// @Produce      json
// @Param        device  body      dto.CreateDeviceRequest  true  "Device payload"
// @Success      201     {object}  dto.DeviceResponse
// @Failure      400     {object}  errors.DefaultErrorResult
// @Failure      500     {object}  errors.DefaultErrorResult
// @Router       /devices [post]
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

// Update godoc
// @Summary      Updates device data by ID
// @Description  Update an existing device name and brand but only when the state is not "in-use", the fiel state can be updated anytime
// @Tags         devices
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Device ID"
// @Param        device  body      dto.UpdateDeviceRequest  true  "Updated device payload"
// @Success      200     {object}  dto.DeviceResponse
// @Failure      400     {object}  errors.DefaultErrorResult
// @Failure      500     {object}  errors.DefaultErrorResult
// @Router       /devices/{id} [put]
func (h *deviceHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		deviceID, err := validateAndParseDeviceId(id)
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

// DeleteDevice godoc
// @Summary      Delete a device
// @Description  Removes a device by ID, only devices that are not "in-use" can be deleted
// @Tags         devices
// @Produce      json
// @Param        id   path      string  true  "Device ID"
// @Success      204  "No Content"
// @Failure      404  {object}  errors.DefaultErrorResult
// @Failure      500  {object}  errors.DefaultErrorResult
// @Router       /devices/{id} [delete]
func (h *deviceHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		//validate id is not empty and is a valid uuid
		deviceID, err := validateAndParseDeviceId(id)
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

func validateAndParseListParams(c echo.Context) (map[string]any, error) {
	allowedParams := map[string]bool{
		"brand": true,
		"state": true,
	}

	parsedQuery, err := url.ParseQuery(c.Request().URL.RawQuery)
	if err != nil {
		fmt.Println(err)
		return nil, errorhandler.NewApiError(errorhandler.ErrInvalid, fmt.Sprintf("invalid query string: %v", err.Error()), nil)
	}

	for param := range parsedQuery {
		if !allowedParams[param] {
			return nil, errorhandler.NewApiError(errorhandler.ErrInvalid, fmt.Sprintf("invalid parameter: %s", param), nil)
		}
	}

	if len(c.QueryParams()) <= 0 {
		return map[string]any{}, nil
	}

	brandParam := c.QueryParam("brand")
	stateParam := c.QueryParam("state")

	if strings.Contains(c.Request().URL.String(), "brand=") && (brandParam == "" || !hasValidValue(brandParam)) {
		return nil, errorhandler.NewApiError(errorhandler.ErrInvalid, "invalid brand filter", nil)
	}

	if stateParam != "" && !validateListDeviceStateFilter(stateParam) {
		return nil, errorhandler.NewApiError(errorhandler.ErrInvalid, "invalid state filter, must be one of: available, in-use, inactive", nil)
	}

	return map[string]any{
		"brand": getStringOrNull(brandParam),
		"state": getStringOrNull(stateParam),
	}, nil

}

func validateAndParseDeviceId(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.Nil, errorhandler.NewApiError(errorhandler.ErrInvalid, "you must inform the device id", nil)
	}
	deviceID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, errorhandler.NewApiError(errorhandler.ErrInvalid, "invalid device id format, must be an uuid", nil)
	}
	return deviceID, nil
}

func validateListDeviceStateFilter(stateParam string) bool {
	validStates := map[string]bool{
		"available": true,
		"in-use":    true,
		"inactive":  true,
	}
	return validStates[stateParam]
}

func getStringOrNull(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}

func hasValidValue(param string) bool {
	validParam := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return validParam.MatchString(param)
}
