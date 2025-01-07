package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levensspel/go-gin-template/config"
	"github.com/samber/do/v2"
)

func NewUserRepositoryInject(i do.Injector) (*pgxpool.Pool, error) {
	config := config.LoadConfig()
	databaseURL := config.DatabaseURL

	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return nil, err
	}
	log.Println("Connected to the database successfully")
	return db, nil
}

func Connect(databaseURL string) *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	log.Println("Connected to the database successfully")
	return db
}
