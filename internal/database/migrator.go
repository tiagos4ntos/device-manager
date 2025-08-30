package database

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateUp(dbConnection *sql.DB) {

	driver, err := postgres.WithInstance(dbConnection, &postgres.Config{})
	if err != nil {
		log.Fatalf("Migration connection failed: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/migrations",
		"postgres", driver)

	if err != nil {
		log.Fatalf("Migration init failed: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations applied successfully.")
}
