package department_service

import (
	"github.com/levensspel/go-gin-template/dto"
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

func New(logger logger.Logger) DepartmentService {
	return &service{
		logger: logger,
	}
}

func (s *service) Create(
	input dto.RequestDepartment,
) (dto.ResponseSingleDepartment, error) {
	return dto.ResponseSingleDepartment{}, nil
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
