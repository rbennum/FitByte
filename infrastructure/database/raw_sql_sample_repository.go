package database

import "github.com/jackc/pgx/v5/pgxpool"

type RawSqlSampleRepository struct {
	db *pgxpool.Pool
}

// Continue according to the need
