package di

import (
	"fmt"
	"os"

	"github.com/TimDebug/FitByte/database"
	"github.com/TimDebug/FitByte/domain"
	"github.com/TimDebug/FitByte/handler"
	"github.com/TimDebug/FitByte/infrastructure/storage"
	"github.com/TimDebug/FitByte/logger"
	"github.com/TimDebug/FitByte/repository"
	"github.com/TimDebug/FitByte/service"
	"github.com/jackc/pgx/v5/pgxpool"

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
	do.Provide[repository.UserRepository](Injector, repository.NewUserRepositoryInject)

	// Setup Services
	do.Provide[service.UserService](Injector, service.NewUserServiceInject)

	// Setup Handlers
	do.Provide[handler.AuthorizationHandler](Injector, handler.NewHandlerInject)
	do.Provide[handler.UserHandler](Injector, handler.NewUserHandlerInject)
}
