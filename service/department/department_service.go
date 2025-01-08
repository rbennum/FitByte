package departmentService

import (
	"context"
	"fmt"
	"strconv"

	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	repositories "github.com/levensspel/go-gin-template/repository/department"
)

type DepartmentService interface {
	Create(managerID string, input dto.RequestDepartment) (dto.ResponseSingleDepartment, error)
	GetAll(managerID string, input dto.RequestDepartment) ([]dto.ResponseSingleDepartment, error)
	Update(name string, id string, managerID string) (dto.ResponseSingleDepartment, error)
	Delete(id string, managerID string) error
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
	managerID string,
	input dto.RequestDepartment,
) (dto.ResponseSingleDepartment, error) {
	if len(input.DepartmentName) < 4 || len(input.DepartmentName) > 33 {
		return dto.ResponseSingleDepartment{}, helper.ErrBadRequest
	}
	row, err := s.repo.Create(context.Background(), input.DepartmentName, managerID)
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
	managerID string,
	input dto.RequestDepartment,
) ([]dto.ResponseSingleDepartment, error) {
	rows, err := s.repo.GetAll(
		context.Background(),
		input.DepartmentName,
		input.Limit,
		input.Offset,
		managerID,
	)
	if err != nil {
		s.logger.Error(
			fmt.Sprintf("Error fetching rows: %v", err),
			helper.DepartmentServiceGetAll,
			err,
		)
		return nil, err
	}
	results := []dto.ResponseSingleDepartment{}
	for _, item := range rows {
		result := dto.ResponseSingleDepartment{}
		result.DepartmentID = item.Id
		result.DepartmentName = item.Name
		s.logger.Info(result.DepartmentID, helper.DepartmentServiceGetAll)
		s.logger.Info(result.DepartmentName, helper.DepartmentServiceGetAll)
		results = append(results, result)
	}
	return results, nil
}

func (s *service) Update(
	name string,
	id string,
	managerID string,
) (dto.ResponseSingleDepartment, error) {
	deptID, err := strconv.Atoi(id)
	if err != nil {
		s.logger.Error(
			fmt.Sprintf("Invalid int conversion: %v", err),
			helper.DepartmentServicePatch,
			err,
		)
		return dto.ResponseSingleDepartment{}, err
	}
	row, err := s.repo.Update(context.Background(), name, deptID, managerID)
	if err != nil {
		s.logger.Error(
			fmt.Sprintf("Error fetching rows: %v", err),
			helper.DepartmentServicePatch,
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

func (s *service) Delete(id string, managerID string) error {
	deptID, err := strconv.Atoi(id)
	if err != nil {
		s.logger.Error(
			fmt.Sprintf("Invalid int conversion: %v", err),
			helper.DepartmentServiceDelete,
			err,
		)
		return err
	}
	err = s.repo.Delete(context.Background(), deptID, managerID)
	if err != nil {
		s.logger.Error(
			err.Error(),
			helper.DepartmentServiceDelete,
			err,
		)
	}
	return err
}
