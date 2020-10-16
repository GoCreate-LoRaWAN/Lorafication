// Package db handles all database related integrations.
package db

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Config represents all of the information that needs to be passed to a postgres connection
// string in order to resolve and connect to a database.
type Config struct {
	User string
	Pass string
	Name string
	Host string
	Port int
}

// NewConnection returns a new database connection with the schema applied, if not already
// applied.
func NewConnection(logger *zap.Logger, cfg Config) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	conn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		cfg.User, cfg.Pass, cfg.Name, cfg.Host, cfg.Port)

	logger.Info("connecting to postgres database...")
	if db, err = sqlx.Connect("postgres", conn); err != nil {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()

		for range ticker.C {
			if db, err = sqlx.Connect("postgres", conn); err == nil {
				break
			}
		}
	}
	logger.Info("connected to postgres database")

	logger.Info("verifying postgres connection...")
	if err := db.Ping(); err != nil {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()

		for range ticker.C {
			if err := db.Ping(); err == nil {
				break
			}
		}
	}
	logger.Info("verified postgres connection")

	if _, err = db.Exec(schema); err != nil {
		return nil, errors.Wrap(err, "apply database schema")
	}

	return db, nil
}
