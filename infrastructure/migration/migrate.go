package migration

import (
	"fmt"

	"github.com/TimDebug/FitByte/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const MIGRATION_FILE_PATH = "file://database/migrations"

func AutoMigrate() {
	cfg := config.DatabaseMigrateUrl()
	migrate, err := migrate.New(MIGRATION_FILE_PATH, cfg)

	if err != nil {
		message := fmt.Sprintf("Error creating migrate instance: %v", err)
		panic(message)
	}

	if err := migrate.Up(); err != nil {
		if err.Error() == "no change" {
			fmt.Println("No new migrations to apply.")
		} else {
			message := fmt.Sprintf("Migration failed: %v", err)
			panic(message)
		}
	}
	fmt.Println("Migration completed successfully!")
}
