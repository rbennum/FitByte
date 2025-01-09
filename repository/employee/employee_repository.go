package repositories

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/entity"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/samber/do/v2"
)

type EmployeeRepository struct {
	db *pgxpool.Pool
}

func NewEmployeeRepository(db *pgxpool.Pool) EmployeeRepository {
	return EmployeeRepository{db: db}
}

func NewEmployeeRepositoryInject(i do.Injector) (EmployeeRepository, error) {
	db := do.MustInvoke[*pgxpool.Pool](i)
	return NewEmployeeRepository(db), nil
}

func (r *EmployeeRepository) IsDepartmentOwnedByManager(ctx context.Context, pool *pgxpool.Tx, departmentId, managerId string) (bool, error) {
	query := "SELECT 1 FROM department WHERE departmentId = $1 AND managerId = $2;"

	rows, err := pool.Exec(ctx, query, departmentId, managerId)
	if err != nil {
		return false, err
	}

	return rows.RowsAffected() > 0, nil
}

func (r *EmployeeRepository) IsIdentityNumberExist(ctx context.Context, pool *pgxpool.Tx, identityNumber, managerId string) (bool, error) {
	query := `
		SELECT 1 
		FROM employees e
		JOIN department d
		ON e.departmentId = d.departmentId
		WHERE
			e.identityNumber = $1
			AND d.managerId = $2
			AND d.isDeleted = FALSE;
	`
	rows, err := pool.Exec(ctx, query, identityNumber, managerId)
	if err != nil {
		return false, err
	}

	return rows.RowsAffected() > 0, nil
}

func (r *EmployeeRepository) GetEmployeeIdIfExist(ctx context.Context, pool *pgxpool.Tx, identityNumber, managerId string) (string, error) {
	query := `
		SELECT e.id 
		FROM employees e
		JOIN department d
		ON e.departmentId = d.departmentId
		WHERE
			e.identityNumber = $1
			AND d.managerId = $2
			AND d.isDeleted = FALSE;
	`
	rows, err := pool.Query(ctx, query, identityNumber, managerId)
	if err != nil {
		return "", err
	}

	var id string
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return "", err
		}
	}

	return id, nil
}

func (r *EmployeeRepository) Insert(ctx context.Context, pool *pgxpool.Tx, input *dto.EmployeePayload, managerId string) error {
	// Check if department ID is owned by the valid manager
	// altogether with the insertion only if its valid within single query.
	query := `
		INSERT INTO employees (
			identityNumber,
			name,
			employeeImageUri,
			gender,
			departmentId
		)
		VALUES ($1, $2, $3, $4, $5);
	`
	rows, err := pool.Exec(
		ctx,
		query,
		input.IdentityNumber,
		input.Name,
		input.EmployeeImageUri,
		input.Gender,
		input.DepartmentID,
	)

	if err != nil {
		return err
	}

	if rows.RowsAffected() < 1 {
		return helper.ErrInvalidDepartmentId
	}

	return nil
}

func (r *EmployeeRepository) Update(ctx context.Context, pool *pgxpool.Tx, input entity.Employee) (dto.EmployeePayload, error) {
	query := "UPDATE employees SET "
	argIndex := 1
	var args []interface{}

	// `UPDATE employees
	// SET
	// 	identityNumber = $1,
	// 	name = $2,
	// 	employeeImageUri = $3,
	// 	gender = $4,
	// 	departmentId = $5
	// WHERE id = $6
	// RETURNING identityNumber, name, employeeImageUri, gender, departmentId`

	if input.Id == "" {
		return dto.EmployeePayload{}, helper.ErrBadRequest
	}

	if input.IdentityNumber != "" {
		query += fmt.Sprintf(" identityNumber = $%d,", argIndex)
		args = append(args, input.IdentityNumber)
		argIndex++
	}
	if input.Name != "" {
		query += fmt.Sprintf(" name = $%d,", argIndex)
		args = append(args, input.Name)
		argIndex++
	}
	if input.EmployeeImageUri != "" {
		query += fmt.Sprintf(" employeeImageUri = $%d,", argIndex)
		args = append(args, input.EmployeeImageUri)
		argIndex++
	}
	if input.Gender != "" {
		query += fmt.Sprintf(" gender = $%d,", argIndex)
		args = append(args, input.Gender)
		argIndex++
	}
	if input.DepartmentId != "" {
		query += fmt.Sprintf(" departmentId = $%d", argIndex)
		args = append(args, input.DepartmentId)
		argIndex++
	}

	query = strings.TrimRight(query, ",") + fmt.Sprintf(" WHERE employees.id = $%d ", argIndex)
	args = append(args, input.Id)

	query += "RETURNING identityNumber, name, employeeImageUri, gender, departmentId"

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return dto.EmployeePayload{}, err
	}

	var employee dto.EmployeePayload

	for rows.Next() {
		err = rows.Scan(
			&employee.IdentityNumber,
			&employee.Name,
			&employee.EmployeeImageUri,
			&employee.Gender,
			&employee.DepartmentID,
		)
		if err != nil {
			return dto.EmployeePayload{}, err
		}
	}
	fmt.Println(
		employee.IdentityNumber,
		employee.Name,
		employee.EmployeeImageUri,
		employee.Gender,
		employee.DepartmentID,
	)

	return employee, nil
}

func (r *EmployeeRepository) GetAll(ctx context.Context, input *dto.GetEmployeesRequest) ([]dto.EmployeePayload, error) {
	query := "SELECT e.identityNumber, e.name, e.employeeImageUri, e.gender, e.departmentId" // 'e' refer to 'employee e' which will be appended later
	conditions := "WHERE m.managerId = $1"                                                   // 'u' refer to 'manager u' which will be appended later
	argIndex := 2
	var args []interface{}
	args = append(args, input.ManagerID)

	// `SELECT
	// 	e.identity_number,
	// 	e.name,
	// 	e.employeeImageUri,
	// 	e.gender,
	// 	e.department_id
	// FROM employees
	// WHERE
	//  manager_id = $1
	// 	identity_number ILIKE $2%
	// 	AND name ILIKE %$3%
	// 	AND gender = $4
	// 	AND department_id = $5
	// LIMIT $5
	// OFFSET $6`

	if input.IdentityNumber != "" {
		args = append(args, input.IdentityNumber)
		conditions += fmt.Sprintf(" AND LOWER(e.identityNumber) ILIKE $%d || '%s'", argIndex, "%") // eg. AND identity_number ILIKE $2 || '%'
		argIndex++
	}
	if input.Name != "" {
		args = append(args, input.Name)
		conditions += fmt.Sprintf(" AND e.name ILIKE '%s' || $%d || '%s'", "%", argIndex, "%") // eg. AND name ILIKE %$2%
		argIndex++
	}
	if input.Gender != "" {
		args = append(args, input.Gender)
		conditions += fmt.Sprintf(" AND e.gender = $%d", argIndex)
		argIndex++
	}
	if input.DepartmentID != "" {
		args = append(args, input.DepartmentID)
		conditions += fmt.Sprintf(" AND e.departmentId = $%d", argIndex)
		argIndex++
	}
	query = strings.TrimRight(query, ",") + " FROM employees AS e LEFT JOIN department d ON e.departmentId = d.departmentId LEFT JOIN manager m ON d.managerId = m.managerId "

	args = append(args, input.Limit)
	conditions += fmt.Sprintf(" LIMIT $%d", argIndex)
	argIndex++

	args = append(args, input.Offset)
	conditions += fmt.Sprintf(" OFFSET $%d;", argIndex)

	query += conditions

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var employees []dto.EmployeePayload
	for rows.Next() {
		var employee dto.EmployeePayload
		err := rows.Scan(
			&employee.IdentityNumber,
			&employee.Name,
			&employee.EmployeeImageUri,
			&employee.Gender,
			&employee.DepartmentID,
		)
		if err != nil {
			log.Printf("Failed to scan row: %v\n", err)
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

func (r *EmployeeRepository) Delete(ctx context.Context, pool *pgxpool.Tx, employeeId string) error {
	query := `DELETE FROM employees WHERE id = $1`

	_, err := pool.Exec(ctx, query, employeeId)
	if err != nil {
		return err
	}

	return nil
}
