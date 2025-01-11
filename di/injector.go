package di

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levensspel/go-gin-template/database"
	"github.com/levensspel/go-gin-template/domain"
	authHandler "github.com/levensspel/go-gin-template/handler/auth"
	departmentHandler "github.com/levensspel/go-gin-template/handler/department"
	employeeHandler "github.com/levensspel/go-gin-template/handler/employee"
	fileHandler "github.com/levensspel/go-gin-template/handler/file"
	userHandler "github.com/levensspel/go-gin-template/handler/user"
	"github.com/levensspel/go-gin-template/infrastructure/storage"
	"github.com/levensspel/go-gin-template/logger"
	departmentService "github.com/levensspel/go-gin-template/service/department"
	user_service "github.com/levensspel/go-gin-template/service/employee"
	fileService "github.com/levensspel/go-gin-template/service/file"
	userService "github.com/levensspel/go-gin-template/service/user"

	departmentRepository "github.com/levensspel/go-gin-template/repository/department"
	repositories "github.com/levensspel/go-gin-template/repository/employee"
	fileRepository "github.com/levensspel/go-gin-template/repository/file"
	userRepository "github.com/levensspel/go-gin-template/repository/user"

	"github.com/samber/do/v2"
)

var Injector *do.RootScope

func init() {
	Injector = do.New()

	// Jika ada dependensi, tolong tambahkan sesuai dengan hirarki
	// Setup client
	envMode := os.Getenv("MODE")
	fmt.Print("Mode :%s", envMode)
	if envMode == "DEBUG" {
		do.Provide[domain.StorageClient](Injector, storage.NewMockStorageClientInject)
	} else {
		do.Provide[domain.StorageClient](Injector, storage.NewS3StorageClientInject)
	}

	// Setup database connection
	do.Provide[*pgxpool.Pool](Injector, database.NewUserRepositoryInject)
	// setup logger
	do.Provide[logger.LogHandler](Injector, logger.NewlogHandlerInject)

	// Setup repositories
	// UserRepository
	do.Provide[userRepository.UserRepository](Injector, userRepository.NewUserRepositoryInject)
	do.Provide[departmentRepository.DepartmentRepository](Injector, departmentRepository.NewInject)
	do.Provide[repositories.EmployeeRepository](Injector, repositories.NewEmployeeRepositoryInject)
	do.Provide[fileRepository.FileRepository](Injector, fileRepository.NewInject)

	// Setup Services
	do.Provide[userService.UserService](Injector, userService.NewUserServiceInject)
	do.Provide[departmentService.DepartmentService](Injector, departmentService.NewInject)
	do.Provide[user_service.EmployeeService](Injector, user_service.NewEmployeeServiceInject)
	do.Provide[fileService.FileService](Injector, fileService.NewInject)

	// Setup Handlers
	do.Provide[userHandler.UserHandler](Injector, userHandler.NewUserHandlerInject)
	do.Provide[authHandler.AuthorizationHandler](Injector, authHandler.NewHandlerInject)
	do.Provide[departmentHandler.DepartmentHandler](Injector, departmentHandler.NewInject)
	do.Provide[employeeHandler.EmployeeHandler](Injector, employeeHandler.NewEmployeeHandlerInject)
	do.Provide[fileHandler.FileHandler](Injector, fileHandler.NewInject)

}
