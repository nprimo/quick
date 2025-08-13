package db

import (
	"database/sql"
	"embed"
	"fmt"
)

//go:embed migrations/*.up.sql
var upFs embed.FS

//go:embed seed/*.sql
var seedFs embed.FS

const (
	MIGRATIONS_DIR = "migrations"
	SEED_DIR       = "seed"
)

func Init(dbConn *sql.DB) error {
	upMigrationsEntries, err := upFs.ReadDir(MIGRATIONS_DIR)
	if err != nil {
		return err
	}
	for _, entry := range upMigrationsEntries {
		cont, err := upFs.ReadFile(fmt.Sprintf("%s/%s",
			MIGRATIONS_DIR, entry.Name()))
		if err != nil {
			return err
		}
		if _, err := dbConn.Exec(string(cont)); err != nil {
			return err
		}
	}
	return nil
}

func Seed(dbConn *sql.DB) error {
	entries, err := seedFs.ReadDir(SEED_DIR)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		cont, err := seedFs.ReadFile(fmt.Sprintf("%s/%s",
			SEED_DIR, entry.Name()))
		if err != nil {
			return err
		}
		if _, err := dbConn.Exec(string(cont)); err != nil {
			return err
		}
	}
	return nil
}
