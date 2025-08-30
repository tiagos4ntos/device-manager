package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tiagos4ntos/device-manager/internal/config"
	"github.com/tiagos4ntos/device-manager/internal/database"
	"github.com/tiagos4ntos/device-manager/internal/domain/device"
	"github.com/tiagos4ntos/device-manager/internal/domain/device/repository"
	"github.com/tiagos4ntos/device-manager/internal/network/handler"
	"github.com/tiagos4ntos/device-manager/internal/network/router"
)

func main() {
	cfg := config.LoadConfig()

	log.Printf("%v:starting...", cfg.AppName)

	// validate config
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	// initialize database connection
	psqlConn, err := database.NewPostgresDB(cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseUser, cfg.DatabasePass, cfg.DatabaseName)
	if err != nil {
		log.Fatalf("failed to connect to database: (%v) ", err.Error())
	}
	defer psqlConn.Close()

	// initialize device repository
	deviceRepository := repository.NewDeviceRepository(psqlConn)

	// initialize device service
	deviceService := device.NewDeviceService(deviceRepository)

	// initialize echo server
	e := echo.New()

	//echo settings
	e.Debug = false
	e.DisableHTTP2 = true
	e.HideBanner = true
	e.HidePort = true
	e.Server.ReadTimeout = time.Duration(cfg.HttpTimeout) * time.Second
	e.Server.WriteTimeout = time.Duration(cfg.HttpTimeout) * time.Second

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//initialize device handler
	deviceHandler := handler.NewDeviceHandler(deviceService)

	router.RegisterRoutes(e, deviceHandler)

	// Create a context that cancels on SIGINT/SIGTERM/os.Interrupt
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	go func() {
		log.Printf("%v:ready...", cfg.AppName)
		if err := e.Start(":" + cfg.ServerPort); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Server error:", err)
		}
	}()
	<-ctx.Done()

	e.Logger.Infof("Gracefully shuting down %v...", cfg.AppName)
	if err := e.Shutdown(context.Background()); err != nil {
		e.Logger.Fatal("Shutdown failed:", err)
	}
}
