package router

import (
	"github.com/labstack/echo/v4"
	"github.com/tiagos4ntos/device-manager/internal/network/handler"
)

func RegisterRoutes(e *echo.Echo, dh handler.DeviceHandler) {
	e.POST("/devices", dh.Create())
	e.GET("/devices", dh.List())
	e.GET("/devices/:id", dh.GetByID())
	e.PUT("/devices/:id", dh.Update())
	e.DELETE("/devices/:id", dh.Delete())
}
