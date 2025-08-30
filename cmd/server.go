package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tiagos4ntos/device-manager/internal/config"
	"github.com/tiagos4ntos/device-manager/internal/database"
)

func main() {

	cfg := config.LoadConfig()

	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	psqlConn, err := database.NewPostgresDB(cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseUser, cfg.DatabasePass, cfg.DatabaseName)
	if err != nil {
		log.Fatalf("failed to connect to database: (%v) ", err.Error())
	}
	defer psqlConn.Close()

	// Run DB migrations
	database.MigrateUp(psqlConn)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "I'm fine")
	})
	e.Logger.Fatal(e.Start(":" + cfg.ServerPort))
}
