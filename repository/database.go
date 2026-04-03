package repository

import (
	"database/sql"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

type databaseErrorString struct {
	s string
}

func (e databaseErrorString) Error() string {
	return e.s
}

func InitPostgres() (*sql.DB, error) {
	if err := godotenv.Load(".env.postgres"); err != nil {
		return nil, databaseErrorString{"Error loading .env.postgres: " + err.Error()}
	}

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

	return db, nil
}
