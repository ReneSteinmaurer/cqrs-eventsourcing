package db

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func StartMigration(connectionString string) error {
	fmt.Println("Starte mit der Migration der Datenbank..")

	m, err := migrate.New("file://./db/migrations", connectionString)
	if err != nil {
		return fmt.Errorf("Fehler beim Erstellen des Migrate-Objekts: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("Fehler beim Ausführen der Migrationen: %w", err)
	}

	return nil
}
