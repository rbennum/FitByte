package departmentRepository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levensspel/go-gin-template/entity"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/samber/do/v2"
)

type DepartmentRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) DepartmentRepository {
	return DepartmentRepository{db: db}
}

func NewInject(i do.Injector) (DepartmentRepository, error) {
	db := do.MustInvoke[*pgxpool.Pool](i)
	return New(db), nil
}

func (r *DepartmentRepository) Create(
	ctx context.Context,
	name string,
	managerID string,
) (*entity.Department, error) {
	// check if the department exists
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM department 
			WHERE
				departmentname = $1
				AND managerid = $2
				AND isdeleted = FALSE
		);
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, name, managerID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		// the said department already exists
		return nil, helper.ErrConflict
	}
	// insert the newly created name into DB
	query = `
		INSERT INTO department (departmentid, departmentname, managerid)
		VALUES (DEFAULT, $1, $2)
		RETURNING departmentid, departmentname
	`
	row := r.db.QueryRow(ctx, query, name, managerID)
	var departmentID string
	var departmentName string
	err = row.Scan(&departmentID, &departmentName)
	if err != nil {
		return nil, err
	}
	result := entity.Department{
		Id:   departmentID,
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
	rows, err := r.db.Query(ctx, query, managerID, fmt.Sprintf("%s%s%s", "%", name, "%"), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var departments []entity.Department
	for rows.Next() {
		var departmentID string
		var departmentName string
		if err := rows.Scan(&departmentID, &departmentName); err != nil {
			return departments, err
		}
		var dept entity.Department
		dept.Id = departmentID
		dept.Name = departmentName
		departments = append(departments, dept)
	}
	return departments, nil
}

func (r *DepartmentRepository) Update(
	ctx context.Context,
	name string,
	deptID string,
	managerID string,
) (*entity.Department, error) {
	query := `
		UPDATE department
		SET departmentname = $1,
			updatedon = CURRENT_TIMESTAMP
		WHERE 
			departmentid = $2
			AND managerid = $3
			AND isdeleted = FALSE
		RETURNING departmentid, departmentname;
	`
	var returnedID string
	var returnedName string
	err := r.db.QueryRow(ctx, query, name, deptID, managerID).Scan(&returnedID, &returnedName)
	if returnedID == "" || returnedName == "" {
		return nil, helper.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	result := entity.Department{}
	result.Id = returnedID
	result.Name = returnedName
	return &result, nil
}

func (r *DepartmentRepository) Delete(
	ctx context.Context,
	deptID string,
	managerID string,
) error {
	// check if the department exists
	var deptName string
	query := `
		SELECT 
			departmentname
		FROM department 
		WHERE 
			departmentid = $1
			AND managerid = $2
			AND isdeleted = FALSE;
	`
	err := r.db.QueryRow(ctx, query, deptID, managerID).Scan(&deptName)
	if deptName == "" {
		return helper.ErrNotFound
	}
	if err != nil {
		return err
	}
	// check if the department has employees assigned
	var employeeCount int64
	query = `
		SELECT COUNT(*)
		FROM employees
		WHERE departmentid = $1;
	`
	err = r.db.QueryRow(ctx, query, deptID).Scan(&employeeCount)
	if err != nil {
		return err
	}
	if employeeCount > 0 {
		return helper.ErrConflict
	}
	// update the isdeleted flag
	query = `
		UPDATE department
		SET isdeleted = TRUE
		WHERE 
			departmentid = $1
			AND managerid = $2
			AND isdeleted = FALSE;
	`
	_, err = r.db.Exec(ctx, query, deptID, managerID)
	return err
}
