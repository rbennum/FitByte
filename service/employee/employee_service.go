package user_service

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/entity"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	repositories "github.com/levensspel/go-gin-template/repository/employee"
	"github.com/samber/do/v2"
)

type EmployeeService interface {
	Create(ctx context.Context, input dto.EmployeePayload, managerId string) error
	GetAll(ctx context.Context, input dto.GetEmployeesRequest) ([]dto.EmployeePayload, error)
	Update(ctx context.Context, input dto.UpdateEmployeeRequest, identityNumber, managerId string) (dto.EmployeePayload, error)
	Delete(ctx context.Context, identityNumber, managerId string) error
}

type service struct {
	dbPool       *pgxpool.Pool
	employeeRepo repositories.EmployeeRepository
	logger       logger.Logger
}

func NewEmployeeService(
	dbPool *pgxpool.Pool,
	employeeRepo repositories.EmployeeRepository,
	logger logger.Logger,
) EmployeeService {
	return &service{
		dbPool:       dbPool,
		employeeRepo: employeeRepo,
		logger:       logger,
	}
}

func NewEmployeeServiceInject(i do.Injector) (EmployeeService, error) {
	_dbPool := do.MustInvoke[*pgxpool.Pool](i)
	_repo := do.MustInvoke[repositories.EmployeeRepository](i)
	_logger := do.MustInvoke[logger.LogHandler](i)
	return NewEmployeeService(_dbPool, _repo, &_logger), nil
}

func (s *service) Create(ctx context.Context, input dto.EmployeePayload, managerId string) error {
	pool, err := s.dbPool.Begin(ctx)
	if err != nil {
		return helper.ErrInternalServer
	}
	txPool := pool.(*pgxpool.Tx)
	defer helper.RollbackOrCommit(ctx, txPool)

	isOwnedByManager, err := s.employeeRepo.IsDepartmentOwnedByManager(ctx, txPool, input.DepartmentID, managerId)
	if err != nil {
		s.logger.Error(err.Error(), helper.EmployeeServiceGet, err)
		return err
	}

	if !isOwnedByManager {
		return helper.ErrInvalidDepartmentId
	}

	exist, err := s.employeeRepo.IsIdentityNumberExist(ctx, txPool, input.IdentityNumber, managerId)
	if err != nil {
		s.logger.Error(err.Error(), helper.EmployeeServiceGet, err)
		return err
	}
	if exist {
		return helper.ErrConflictIdentityNumber
	}

	err = s.employeeRepo.Insert(ctx, txPool, &input, managerId)
	if err != nil {
		s.logger.Error(err.Error(), helper.EmployeeServiceGet, err)
		if strings.Contains(err.Error(), "23505") {
			return helper.ErrConflict
		}

		return err
	}

	return nil
}

func (s *service) GetAll(ctx context.Context, input dto.GetEmployeesRequest) ([]dto.EmployeePayload, error) {
	employees, err := s.employeeRepo.GetAll(ctx, &input)
	if err != nil {
		s.logger.Error(err.Error(), helper.EmployeeServiceGet, input)
		return []dto.EmployeePayload{}, err
	}

	return employees, nil
}

func (s *service) Update(ctx context.Context, input dto.UpdateEmployeeRequest, identityNumber, managerId string) (dto.EmployeePayload, error) {
	pool, err := s.dbPool.Begin(ctx)
	if err != nil {
		return dto.EmployeePayload{}, helper.ErrInternalServer
	}
	txPool := pool.(*pgxpool.Tx)
	defer helper.RollbackOrCommit(ctx, txPool)

	if input.DepartmentID != "" {
		isOwnedByManager, err := s.employeeRepo.IsDepartmentOwnedByManager(ctx, txPool, input.DepartmentID, managerId)
		if err != nil {
			s.logger.Error(err.Error(), helper.EmployeeServiceUpdate, err)
			return dto.EmployeePayload{}, err
		}

		if !isOwnedByManager {
			return dto.EmployeePayload{}, helper.ErrInvalidDepartmentId
		}
	}

	id, err := s.employeeRepo.GetEmployeeIdIfExist(ctx, txPool, identityNumber, managerId)
	if err != nil {
		s.logger.Error(err.Error(), helper.EmployeeServiceUpdate, err)
		return dto.EmployeePayload{}, err
	}
	if id != "" {
		return dto.EmployeePayload{}, helper.ErrConflict
	}

	employee, err := s.employeeRepo.Update(ctx, txPool, entity.Employee{
		Id:               id,
		Name:             input.Name,
		IdentityNumber:   input.IdentityNumber,
		EmployeeImageUri: input.EmployeeImageUri,
		Gender:           input.Gender,
		DepartmentId:     input.DepartmentID,
	})
	if err != nil {
		s.logger.Error(err.Error(), helper.EmployeeServiceUpdate, err)
		if strings.Contains(err.Error(), "23505") {
			return dto.EmployeePayload{}, helper.ErrConflict
		}

		return dto.EmployeePayload{}, err
	}
	return dto.EmployeePayload{
		Name:             employee.Name,
		IdentityNumber:   employee.IdentityNumber,
		EmployeeImageUri: employee.EmployeeImageUri,
		Gender:           employee.Gender,
		DepartmentID:     employee.DepartmentID,
	}, nil
}

func (s *service) Delete(ctx context.Context, identityNumber, managerId string) error {
	pool, err := s.dbPool.Begin(ctx)
	if err != nil {
		return helper.ErrInternalServer
	}
	txPool := pool.(*pgxpool.Tx)
	defer helper.RollbackOrCommit(ctx, txPool)

	id, err := s.employeeRepo.GetEmployeeIdIfExist(ctx, txPool, identityNumber, managerId)
	if err != nil {
		s.logger.Error(err.Error(), helper.EmployeeServiceDelete, err)
		return err
	}
	if id == "" {
		return helper.ErrNotFound
	}

	err = s.employeeRepo.Delete(ctx, txPool, id)
	if err != nil {
		s.logger.Error(err.Error(), helper.EmployeeServiceDelete, err)

		return err
	}

	return nil
}
