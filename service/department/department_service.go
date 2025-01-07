package departmentService

import (
	"context"
	"fmt"

	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	repositories "github.com/levensspel/go-gin-template/repository/department"
)

type DepartmentService interface {
	Create(input dto.RequestDepartment) (dto.ResponseSingleDepartment, error)
	GetAll(name string, limit int, offset int) ([]dto.ResponseSingleDepartment, error)
	Update(name string, id string) (dto.ResponseSingleDepartment, error)
	Delete(id string) error
}

type service struct {
	repo   repositories.DepartmentRepository
	logger logger.Logger
}

func New(repo repositories.DepartmentRepository, logger logger.Logger) DepartmentService {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) Create(
	input dto.RequestDepartment,
) (dto.ResponseSingleDepartment, error) {
	if len(input.DepartmentName) < 4 || len(input.DepartmentName) > 33 {
		return dto.ResponseSingleDepartment{}, helper.ErrBadRequest
	}
	row, err := s.repo.Create(context.Background(), input.DepartmentName)
	if err != nil {
		s.logger.Error(
			fmt.Sprintf("Error fetching rows: %v", err),
			helper.DepartmentServiceCreate,
			err,
		)
		return dto.ResponseSingleDepartment{}, err
	}
	result := dto.ResponseSingleDepartment{
		DepartmentID:   row.Id,
		DepartmentName: row.Name,
	}
	return result, nil
}

func (s *service) GetAll(
	name string,
	limit int,
	offset int,
) ([]dto.ResponseSingleDepartment, error) {
	return nil, nil
}

func (s *service) Update(
	name string,
	id string,
) (dto.ResponseSingleDepartment, error) {
	return dto.ResponseSingleDepartment{}, nil
}

func (s *service) Delete(id string) error {
	return nil
}
