package di

import (
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levensspel/go-gin-template/database"
	"github.com/levensspel/go-gin-template/domain"
	authHandler "github.com/levensspel/go-gin-template/handler/auth"
	departmentHandler "github.com/levensspel/go-gin-template/handler/department"
	employeeHandler "github.com/levensspel/go-gin-template/handler/employee"
	userHandler "github.com/levensspel/go-gin-template/handler/user"
	"github.com/levensspel/go-gin-template/infrastructure/storage"
	"github.com/levensspel/go-gin-template/logger"
	departmentService "github.com/levensspel/go-gin-template/service/department"
	user_service "github.com/levensspel/go-gin-template/service/employee"
	userService "github.com/levensspel/go-gin-template/service/user"

	departmentRepository "github.com/levensspel/go-gin-template/repository/department"
	repositories "github.com/levensspel/go-gin-template/repository/employee"
	userRepository "github.com/levensspel/go-gin-template/repository/user"

	"github.com/samber/do/v2"
)

var Injector *do.RootScope

func init() {
	Injector = do.New()

	// Jika ada dependensi, tolong tambahkan sesuai dengan hirarki

	// Setup dependensi-depenensi dasar sebuah service

	// Setup database connection
	do.Provide[*pgxpool.Pool](Injector, database.NewUserRepositoryInject)
	// setup logger
	do.Provide[logger.LogHandler](Injector, logger.NewlogHandlerInject)

	// Setup repositories
	// UserRepository
	do.Provide[userRepository.UserRepository](Injector, userRepository.NewUserRepositoryInject)
	do.Provide[departmentRepository.DepartmentRepository](Injector, departmentRepository.NewInject)
	do.Provide[repositories.EmployeeRepository](Injector, repositories.NewEmployeeRepositoryInject)

	// Setup Services
	do.Provide[userService.UserService](Injector, userService.NewUserServiceInject)
	do.Provide[departmentService.DepartmentService](Injector, departmentService.NewInject)
	do.Provide[user_service.EmployeeService](Injector, user_service.NewEmployeeServiceInject)

	// Setup Handlers
	do.Provide[userHandler.UserHandler](Injector, userHandler.NewUserHandlerInject)
	do.Provide[authHandler.AuthorizationHandler](Injector, authHandler.NewHandlerInject)
	do.Provide[departmentHandler.DepartmentHandler](Injector, departmentHandler.NewInject)
	do.Provide[employeeHandler.EmployeeHandler](Injector, employeeHandler.NewEmployeeHandlerInject)

	// Setup client
	envMode := os.Getenv("MODE")
	if envMode == "DEBUG" {
		do.Provide[domain.StorageClient](Injector, storage.NewMockStorageClientInject)
	} else {
		do.Provide[domain.StorageClient](Injector, storage.NewS3StorageClientInject)
	}
}
