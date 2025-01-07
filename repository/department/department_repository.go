package departmentRepository

import (
	"context"
	"fmt"

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
) (*entity.Department, error) {
	query := `
		INSERT INTO department (departmentid, departmentname)
		VALUES (DEFAULT, $1)
		RETURNING departmentid, departmentname
	`
	row := r.db.QueryRow(ctx, query, name)
	var departmentID int
	var departmentName string
	err := row.Scan(&departmentID, &departmentName)
	if err != nil {
		return nil, err
	}
	result := entity.Department{
		Id:   fmt.Sprintf("%d", departmentID),
		Name: departmentName,
	}
	return &result, nil
}

func (r *DepartmentRepository) GetAll(
	ctx context.Context,
	limit int,
	offset int,
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
