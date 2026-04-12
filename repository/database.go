package repository

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/WASDetchan/wasdetchan-online/core"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

func migrateDB() error {
	source, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return err
	}

	cfg, err := pq.NewConfig("")
	if err != nil {
		return databaseErrorString{"Error creating postgres config: " + err.Error()}
	}

	c, err := pq.NewConnectorConfig(cfg)
	if err != nil {
		return databaseErrorString{"Error creating postgres connection: " + err.Error()}
	}

	db := sql.OpenDB(c)
	defer db.Close()

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
	db, err := pgx.Connect(context.Background(), "")

	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	err = migrateDB()
	if err != nil {
		return nil, databaseErrorString{"Error migrating the database: " + err.Error()}
	}

	return New(db), nil
}

type QueriesKey struct{}

func Middleware(q *Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		core.PushContext(c, QueriesKey{}, q)
		c.Set(QueriesKey{}, q)
	}
}

func GetQueries(c *gin.Context) *Queries {
	queries, succ := c.Get(QueriesKey{})
	if !succ {
		log.Panic("queries not available")
	}
	return queries.(*Queries)
}
