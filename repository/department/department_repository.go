package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levensspel/go-gin-template/entity"
)

type DepartmentRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) DepartmentRepository {
	return DepartmentRepository{db: db}
}

func (r *DepartmentRepository) Create(
	ctx context.Context,
	name string,
	limit int,
	offset int,
) error {
	return nil
}

func (r *DepartmentRepository) GetAll(
	ctx context.Context,
	dept entity.Department,
) ([]entity.Department, error) {
	return nil, nil
}

func (r *DepartmentRepository) Update(
	ctx context.Context,
	dept entity.Department,
) error {
	return nil
}

func (r *DepartmentRepository) Delete(
	ctx context.Context,
	dept entity.Department,
) error {
	return nil
}
