package repository

import (
	"database/sql"
	"embed"
	"log"

	"github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type databaseErrorString struct {
	s string
}

func (e databaseErrorString) Error() string {
	return e.s
}

//go:embed migrations/*.sql
var migrationFiles embed.FS

func migrateDB(db *sql.DB) error {
	source, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	} else if err == nil {
		log.Println("Database migrated successfully.")
	}
	return nil
}

func InitPostgres() (*Queries, error) {
	cfg, err := pq.NewConfig("")
	if err != nil {
		return nil, databaseErrorString{"Error creating postgres config: " + err.Error()}
	}

	c, err := pq.NewConnectorConfig(cfg)
	if err != nil {
		return nil, databaseErrorString{"Error creating postgres connection: " + err.Error()}
	}

	db := sql.OpenDB(c)

	err = db.Ping()
	if err != nil {
		return nil, databaseErrorString{"Error connecting to the database: " + err.Error()}
	}

	err = migrateDB(db)
	if err != nil {
		return nil, databaseErrorString{"Error migrating the database: " + err.Error()}
	}

	return New(db), nil
}
