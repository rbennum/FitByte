package di

import (
	"github.com/levensspel/go-gin-template/domain"
	"github.com/levensspel/go-gin-template/infrastructure/storage"
	"github.com/samber/do/v2"
	"os"
)

var Injector *do.RootScope

func init() {
	Injector = do.New()

	// Setup repositories

	// Setup client
	envMode := os.Getenv("MODE")
	if envMode == "LOCAL" {
		do.Provide[domain.StorageClient](Injector, storage.NewMockStorageClientInject)
	} else {
		do.Provide[domain.StorageClient](Injector, storage.NewS3ClientInject)
	}
}
