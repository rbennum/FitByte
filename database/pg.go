package database

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levensspel/go-gin-template/config"
	"github.com/samber/do/v2"
)

func NewUserRepositoryInject(i do.Injector) (*pgxpool.Pool, error) {
	config := config.LoadConfig()
	databaseURL := config.DatabaseURL

	tempPool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer tempPool.Close()

	var maxConnectionText string
	row := tempPool.QueryRow(context.Background(), "SHOW max_connections")
	err = row.Scan(&maxConnectionText)
	if err != nil {
		panic(err)
	}

	maxConnection, err := strconv.Atoi(maxConnectionText)
	if err != nil {
		panic(err)
	}

	maxOpenConnects := int(float64(maxConnection) * 0.8)

	fmt.Println(maxOpenConnects)
	fmt.Println(maxConnection)

	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}

	poolConfig.MaxConns = int32(maxOpenConnects)

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
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
