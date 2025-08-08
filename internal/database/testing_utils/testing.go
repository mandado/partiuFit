package testing_utils

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"partiuFit/internal/database"
	"partiuFit/migrations"
)

func SetupTestDB() (*sql.DB, error) {
	godotenv.Load("../../.env.testing")
	db, err := database.Open(os.Getenv("DATABASE_URL"))

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	//err = Migrate(db, "../../migrations")
	err = database.MigrateFS(db, migrations.FS, migrations.FSPath)

	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	err = truncateTables(db)

	if err != nil {
		return nil, fmt.Errorf("failed to truncate tables: %w", err)
	}

	return db, nil
}

func TeardownTestDB(db *sql.DB) error {
	err := truncateTables(db)

	if err != nil {
		return fmt.Errorf("failed to run truncate tables: %w", err)
	}

	return db.Close()
}
func truncateTables(db *sql.DB) error {
	_, err := db.Exec("truncate workouts, workout_entries, users cascade")

	if err != nil {
		return fmt.Errorf("failed to truncate tables: %w", err)
	}

	return nil
}
