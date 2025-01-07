package fileService

import (
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/logger"
	repositories "github.com/levensspel/go-gin-template/repository/file"
)

type FileService interface {
	Upload(input dto.FileUploadRequestPayload) (dto.FileUploadRespondPayload, error)
	DeleteByID(fileid string) error
}

type fileService struct {
	repo   repositories.FileRepository
	logger logger.Logger
}

func NewFileService(
	repo repositories.FileRepository,
	logger logger.Logger,
) FileService {
	return &fileService{
		repo:   repo,
		logger: logger,
	}
}

func (s *fileService) Upload(input dto.FileUploadRequestPayload) (dto.FileUploadRespondPayload, error) {
	return dto.FileUploadRespondPayload{}, nil
}

func (s *fileService) DeleteByID(fileid string) error {
	return nil
}
