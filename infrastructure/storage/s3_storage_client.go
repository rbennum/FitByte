package storage

import (
	"github.com/levensspel/go-gin-template/domain"
	"github.com/samber/do/v2"
)

type S3Client struct {
	// TODO: Complete
}

func NewS3Client() domain.StorageClient {
	return S3Client{}
}

func NewS3ClientInject(i do.Injector) (domain.StorageClient, error) {
	return NewS3Client(), nil
}
