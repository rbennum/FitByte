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
	managerID string,
) (*entity.Department, error) {
	query := `
		INSERT INTO department (departmentid, departmentname, managerid)
		VALUES (DEFAULT, $1, $2)
		RETURNING departmentid, departmentname
	`
	row := r.db.QueryRow(ctx, query, name, managerID)
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
	name string,
	limit int,
	offset int,
	managerID string,
) ([]entity.Department, error) {
	query := `
		SELECT departmentid, departmentname
		FROM department
		WHERE 
			managerid = $1
			AND departmentname ILIKE $2
			AND isdeleted = FALSE
		LIMIT $3 OFFSET $4;
	`
	rows, err := r.db.Query(ctx, query, managerID, name, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var departments []entity.Department
	for rows.Next() {
		var departmentID int
		var departmentName string
		if err := rows.Scan(&departmentID, &departmentName); err != nil {
			return departments, err
		}
		var dept entity.Department
		dept.Id = fmt.Sprintf("%d", departmentID)
		dept.Name = departmentName
		departments = append(departments, dept)
	}
	return departments, nil
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
