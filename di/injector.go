package di

import (
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/levensspel/go-gin-template/database"
	"github.com/levensspel/go-gin-template/domain"
	authHandler "github.com/levensspel/go-gin-template/handler/auth"
	userHandler "github.com/levensspel/go-gin-template/handler/user"
	"github.com/levensspel/go-gin-template/infrastructure/storage"
	"github.com/levensspel/go-gin-template/logger"
	userService "github.com/levensspel/go-gin-template/service/user"

	userRepository "github.com/levensspel/go-gin-template/repository/user"

	"github.com/samber/do/v2"
)

var Injector *do.RootScope

func init() {
	Injector = do.New()

	// Setup database connection
	do.Provide[*pgxpool.Pool](Injector, database.NewUserRepositoryInject)
	// setup logger
	do.Provide[logger.LogHandler](Injector, logger.NewlogHandlerInject)

	// Setup repositories
	// UserRepository
	do.Provide[userRepository.UserRepository](Injector, userRepository.NewUserRepositoryInject)

	// Services
	do.Provide[userService.UserService](Injector, userService.NewUserServiceInject)

	// Handlers
	do.Provide[userHandler.UserHandler](Injector, userHandler.NewUserHandlerInject)
	do.Provide[authHandler.AuthorizationHandler](Injector, authHandler.NewHandlerInject)

	// Setup client
	envMode := os.Getenv("MODE")
	if envMode == "DEBUG" {
		do.Provide[domain.StorageClient](Injector, storage.NewMockStorageClientInject)
	} else {
		do.Provide[domain.StorageClient](Injector, storage.NewS3StorageClientInject)
	}
}
