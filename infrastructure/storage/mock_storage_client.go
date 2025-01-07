package storage

import (
	"github.com/levensspel/go-gin-template/domain"
	"github.com/samber/do/v2"
)

type MockStorageClient struct {
	// TODO: Complete
}

func NewMockStorageClient() domain.StorageClient {
	return MockStorageClient{}
}

func NewMockStorageClientInject(i do.Injector) (domain.StorageClient, error) {
	return NewMockStorageClient(), nil
}
