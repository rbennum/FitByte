package di

import (
	"fmt"
	"os"

	"github.com/TimDebug/FitByte/database"
	"github.com/TimDebug/FitByte/domain"
	authHandler "github.com/TimDebug/FitByte/handler/auth"
	departmentHandler "github.com/TimDebug/FitByte/handler/department"
	employeeHandler "github.com/TimDebug/FitByte/handler/employee"
	fileHandler "github.com/TimDebug/FitByte/handler/file"
	userHandler "github.com/TimDebug/FitByte/handler/user"
	"github.com/TimDebug/FitByte/infrastructure/storage"
	"github.com/TimDebug/FitByte/logger"
	departmentService "github.com/TimDebug/FitByte/service/department"
	user_service "github.com/TimDebug/FitByte/service/employee"
	fileService "github.com/TimDebug/FitByte/service/file"
	userService "github.com/TimDebug/FitByte/service/user"
	"github.com/jackc/pgx/v5/pgxpool"

	departmentRepository "github.com/TimDebug/FitByte/repository/department"
	repositories "github.com/TimDebug/FitByte/repository/employee"
	fileRepository "github.com/TimDebug/FitByte/repository/file"
	userRepository "github.com/TimDebug/FitByte/repository/user"

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
