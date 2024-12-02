package postgres

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func Migrate(conn *sqlx.DB) error {
	driver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "failed to create postgres driver")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", 
		"postgres",
		driver,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create migration instance")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed to apply migrations")
	}

	return nil
}
