package pg

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"

	"github.com/Alexey-step/rocket-factory/platform/pkg/migrator"
)

type pgMigrator struct {
	db            *sql.DB
	migrationsDir string
}

func NewMigrator(db *sql.DB, migrationDir string) migrator.Migrator {
	return &pgMigrator{
		db:            db,
		migrationsDir: migrationDir,
	}
}

func (m *pgMigrator) Up(_ context.Context) error {
	err := goose.Up(m.db, m.migrationsDir)
	if err != nil {
		return err
	}
	return nil
}
