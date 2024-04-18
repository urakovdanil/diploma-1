package postgres

import (
	"database/sql"
	"diploma-1/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

func migrateUp(su *config.StartUp) error {
	db, err := sql.Open("pgx", su.DatabaseURI)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	defer db.Close()
	if err := goose.Up(db, su.MigrationsFolder); err != nil {
		return err
	}
	return nil
}
