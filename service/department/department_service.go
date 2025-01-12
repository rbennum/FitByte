package departmentService

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/levensspel/go-gin-template/cache"

	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	repositories "github.com/levensspel/go-gin-template/repository/department"
	"github.com/samber/do/v2"
)

const (
	defaultTtl = 5 * time.Minute // Department data won't more stale than 5 mins
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

func NewInject(i do.Injector) (DepartmentService, error) {
	_repo := do.MustInvoke[repositories.DepartmentRepository](i)
	_logger := do.MustInvoke[logger.LogHandler](i)
	return New(_repo, &_logger), nil
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

	invalidateCache()

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
	cacheKey := s.generateCacheKey(managerID, input)

	// Check cache
	cachedDepartments, found := cache.GetAsMapArray(cacheKey)
	if found {
		result := make([]dto.ResponseSingleDepartment, len(cachedDepartments))
		for i, department := range cachedDepartments {
			result[i] = dto.ResponseSingleDepartment{
				DepartmentID:   department["id"],
				DepartmentName: department["name"],
			}
		}
		return result, nil
	}

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

	// Put to cache
	costMultiplier := s.calculateCostMultiplier(input)
	departmentsToCache := make([]map[string]string, len(rows))
	for i, row := range rows {
		departmentsToCache[i] = map[string]string{
			"id":   row.Id,
			"name": row.Name,
		}
	}
	cache.SetAsMapArrayWithTtlAndCostMultiplier(cacheKey, departmentsToCache, costMultiplier, defaultTtl)

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
	if len(name) < 4 || len(name) > 33 {
		return dto.ResponseSingleDepartment{}, helper.ErrBadRequest
	}
	row, err := s.repo.Update(context.Background(), name, id, managerID)
	if err != nil {
		s.logger.Error(
			fmt.Sprintf("Error fetching rows: %v", err),
			helper.DepartmentServicePatch,
			err,
		)
		if err == helper.ErrNotFound {
			return dto.ResponseSingleDepartment{}, helper.NewErrorResponse(
				http.StatusNotFound,
				err.Error(),
			)
		} else if err == helper.ErrConflict {
			return dto.ResponseSingleDepartment{}, helper.NewErrorResponse(
				http.StatusConflict,
				err.Error(),
			)
		} else {
			return dto.ResponseSingleDepartment{}, helper.NewErrorResponse(
				http.StatusInternalServerError,
				err.Error(),
			)
		}
	}

	invalidateCache()

	result := dto.ResponseSingleDepartment{
		DepartmentID:   row.Id,
		DepartmentName: row.Name,
	}
	return result, nil
}

func (s *service) Delete(id string, managerID string) error {
	err := s.repo.Delete(context.Background(), id, managerID)
	if err != nil {
		s.logger.Error(
			err.Error(),
			helper.DepartmentServiceDelete,
			err,
		)
	}

	invalidateCache()

	return err
}

func (s *service) generateCacheKey(managerId string, input dto.RequestDepartment) string {
	namespaceVersion := cache.DepartmentNamespaceVersion.Load()

	// Serialize params into a string (e.g., "name=Jono&gender=male")
	var filterParts []string

	filterParts = append(filterParts, fmt.Sprintf("limit=%d", input.Limit))
	filterParts = append(filterParts, fmt.Sprintf("offset=%d", input.Offset))
	filterParts = append(filterParts, fmt.Sprintf("managerId=%s", managerId))
	filterParts = append(filterParts, fmt.Sprintf("name=%s", input.DepartmentName))

	filters := strings.Join(filterParts, "&")
	return fmt.Sprintf(cache.CacheDepartmentsWithParams, namespaceVersion, filters)
}

func (s *service) calculateCostMultiplier(input dto.RequestDepartment) int {
	// The more likely it is to be searched, the more beneficial it is to cache
	// By setting higher cost, it will be more likely to be cached and less likely to be evicted

	// No filter
	noFilter := input.DepartmentName == ""
	if noFilter {
		if input.Offset == 1 {
			// First page
			return 4
		} else if input.Offset == 2 {
			// Second page
			return 3
		} else {
			// Subsequent pages
			return 2
		}
	} else {
		return 1
	}
}

func invalidateCache() {
	cache.DepartmentNamespaceVersion.Add(1)
}
