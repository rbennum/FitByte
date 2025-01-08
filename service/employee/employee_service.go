package user_service

import (
	"context"
	"strings"

	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	repositories "github.com/levensspel/go-gin-template/repository/employee"
)

type EmployeeService interface {
	Create(input dto.EmployeePayload, managerId string) error
	GetEmployees(input dto.GetEmployeesRequest) ([]dto.EmployeePayload, error)
}

type service struct {
	employeeRepo repositories.EmployeeRepository
	logger       logger.Logger
}

func NewEmployeeService(
	employeeRepo repositories.EmployeeRepository,
	logger logger.Logger,
) EmployeeService {
	return &service{
		employeeRepo: employeeRepo,
		logger:       logger,
	}
}

func (s *service) Create(input dto.EmployeePayload, managerId string) error {
	err := s.employeeRepo.Create(context.Background(), &input, managerId)
	if err != nil {
		s.logger.Error(err.Error(), helper.EmployeeServiceGet, err)
		if strings.Contains(err.Error(), "23505") {
			return helper.ErrConflict
		}

		return err
	}

	return nil
}

func (s *service) GetEmployees(input dto.GetEmployeesRequest) ([]dto.EmployeePayload, error) {
	param := &dto.GetEmployeesRequest{
		Limit:          input.Limit,
		Offset:         input.Offset,
		IdentityNumber: input.IdentityNumber,
		Name:           input.Name,
		Gender:         input.Gender,
		DepartmentID:   input.DepartmentID,
		ManagerID:      input.ManagerID,
	}

	employees, err := s.employeeRepo.GetEmployees(context.Background(), param)
	if err != nil {
		s.logger.Error(err.Error(), helper.EmployeeServiceGet, param)
		return []dto.EmployeePayload{}, err
	}

	return employees, nil
}
