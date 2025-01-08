package user_service

import (
	"context"
	"strings"

	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	repositories "github.com/levensspel/go-gin-template/repository/employee"
	"github.com/samber/do/v2"
)

type EmployeeService interface {
	Create(input dto.EmployeePayload, managerId string) error
	GetAll(input dto.GetEmployeesRequest) ([]dto.EmployeePayload, error)
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

func NewEmployeeServiceInject(i do.Injector) (EmployeeService, error) {
	_repo := do.MustInvoke[repositories.EmployeeRepository](i)
	_logger := do.MustInvoke[logger.LogHandler](i)
	return NewEmployeeService(_repo, &_logger), nil
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

func (s *service) GetAll(input dto.GetEmployeesRequest) ([]dto.EmployeePayload, error) {
	employees, err := s.employeeRepo.GetAll(context.Background(), &input)
	if err != nil {
		s.logger.Error(err.Error(), helper.EmployeeServiceGet, input)
		return []dto.EmployeePayload{}, err
	}

	return employees, nil
}
