package repositories

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levensspel/go-gin-template/dto"
)

type EmployeeRepository struct {
	db *pgxpool.Pool
}

func NewEmployeeRepository(db *pgxpool.Pool) EmployeeRepository {
	return EmployeeRepository{db: db}
}

func (r *EmployeeRepository) GetEmployees(ctx context.Context, input *dto.GetEmployeesRequest) ([]dto.GetEmployeeResponseItem, error) {
	//	The constructed query if all the parameters are sent:
	// `SELECT
	// 		identity_number,
	// 		name,
	// 		gender,
	// 		image_uri,
	// 		department_id
	// FROM employees
	// WHERE
	//  	manager_id = $1
	// 		AND LOWER(identity_number) ILIKE $2 || '%'
	// 		AND name ILIKE '%' || $2 || '%'"
	// 		AND gender = $4
	// 		AND department_id = $5
	// LIMIT $5
	// OFFSET $6`

	query := "SELECT identitynumber, name, gender, employeeimageuri, departmentid"
	conditions := "WHERE managerid = $1"
	argIndex := 2
	var args []interface{}
	args = append(args, input.ManagerID)

	if input.IdentityNumber != "" {
		args = append(args, input.IdentityNumber)
		conditions += fmt.Sprintf(" AND LOWER(identitynumber) ILIKE $%d || '%s'", argIndex, "%") // eg. "AND LOWER(identity_number) ILIKE $2 || '%'"
		argIndex++
	}
	if input.Name != "" {
		args = append(args, input.Name)
		conditions += fmt.Sprintf(" AND name ILIKE '%s' || $%d || '%s'", "%", argIndex, "%") // eg. "AND name ILIKE '%' || $2 || '%'"
		argIndex++
	}
	if input.Gender != "" {
		args = append(args, input.Gender)
		conditions += fmt.Sprintf(" AND gender = $%d", argIndex)
		argIndex++
	}
	if input.DepartmentID != "" {
		args = append(args, input.DepartmentID)
		conditions += fmt.Sprintf(" AND departmentid = $%d", argIndex)
		argIndex++
	}
	query = strings.TrimRight(query, ",") + " from employees "

	args = append(args, input.Limit)
	conditions += fmt.Sprintf(" LIMIT $%d", argIndex)
	argIndex++

	args = append(args, input.Offset)
	conditions += fmt.Sprintf(" OFFSET $%d;", argIndex)

	query += conditions
	fmt.Println(query)
	fmt.Println(args...)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var employees []dto.GetEmployeeResponseItem
	for rows.Next() {
		var employee dto.GetEmployeeResponseItem

		// The employee's attributes depends on the order of written columns in the SQL query
		err := rows.Scan(
			&employee.IdentityNumber,
			&employee.Name,
			&employee.Gender,
			&employee.EmployeeImageUri,
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
