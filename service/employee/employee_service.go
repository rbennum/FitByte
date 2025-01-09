package user_service

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	repositories "github.com/levensspel/go-gin-template/repository/employee"
	"github.com/samber/do/v2"
)

type EmployeeService interface {
	Create(ctx context.Context, input dto.EmployeePayload, managerId string) error
	GetAll(ctx context.Context, input dto.GetEmployeesRequest) ([]dto.EmployeePayload, error)
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

	err = s.employeeRepo.IsDepartmentOwnedByManager(ctx, txPool, input.DepartmentID, managerId)
	if err != nil {
		s.logger.Error(err.Error(), helper.EmployeeServiceGet, err)
		return err
	}

	err = s.employeeRepo.IsIdentityNumberAvailable(ctx, txPool, input.IdentityNumber, managerId)
	if err != nil {
		s.logger.Error(err.Error(), helper.EmployeeServiceGet, err)
		return err
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
