package migration

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/levensspel/go-gin-template/config"
)

const MIGRATION_FILE_PATH = "file://database/migrations"

func AutoMigrate() {
	cfg := config.DatabaseMigrateUrl()
	migrate, err := migrate.New(MIGRATION_FILE_PATH, cfg)

	if err != nil {
		log.Fatalf("Error creating migrate instance: %v", err)
	}

	if err := migrate.Up(); err != nil {
		if err.Error() == "no change" {
			log.Println("No new migrations to apply.")
		} else {
			log.Fatalf("Migration failed: %v", err)
		}
	}
	log.Println("Migration completed successfully!")
}
