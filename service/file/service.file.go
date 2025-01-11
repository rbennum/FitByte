package fileService

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/levensspel/go-gin-template/domain"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	repositories "github.com/levensspel/go-gin-template/repository/file"
	"github.com/samber/do/v2"
)

type FileService interface {
	Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader) (dto.FileUploadRespondPayload, error)
	DeleteByID(fileid string) error
}

type fileService struct {
	repo          repositories.FileRepository
	logger        logger.Logger
	storageClient domain.StorageClient
}

func NewFileService(
	repo repositories.FileRepository,
	logger logger.Logger,
	storageClient domain.StorageClient,
) FileService {
	return &fileService{
		repo:          repo,
		logger:        logger,
		storageClient: storageClient,
	}
}

func NewInject(i do.Injector) (FileService, error) {
	_db := do.MustInvoke[repositories.FileRepository](i)
	_logger := do.MustInvoke[logger.LogHandler](i)
	_storageClient := do.MustInvoke[domain.StorageClient](i)
	return NewFileService(_db, &_logger, _storageClient), nil
}

func (s *fileService) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader) (dto.FileUploadRespondPayload, error) {
	// Simpan file ke server lokal
	file, err := header.Open()
	if err != nil {
		s.logger.Warn(err.Error(), helper.FunctionCaller("FileHandler.Upload"), header)
		return dto.FileUploadRespondPayload{}, helper.ErrInternalServer
	}
	defer file.Close()

	// Tentukan lokasi penyimpanan file
	savePath := fmt.Sprintf("./.uploads/%s", header.Filename)

	// Buat direktori jika belum ada
	if err := os.MkdirAll("./.uploads", os.ModePerm); err != nil {
		s.logger.Error(err.Error(), helper.FunctionCaller("FileHandler.Upload"))
		return dto.FileUploadRespondPayload{}, helper.ErrInternalServer
	}

	// Simpan file
	out, err := os.Create(savePath)
	if err != nil {
		s.logger.Error(err.Error(), helper.FunctionCaller("FileHandler.Upload"))
		s.logger.Error("Failed to save file", helper.FunctionCaller("FileHandler.Upload"))
		return dto.FileUploadRespondPayload{}, helper.ErrInternalServer
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		s.logger.Error(err.Error(), helper.FunctionCaller("FileHandler.Upload"))
		s.logger.Error("Failed to write file", helper.FunctionCaller("FileHandler.Upload"))
		return dto.FileUploadRespondPayload{}, helper.ErrInternalServer
	}

	return dto.FileUploadRespondPayload{}, nil
}

func (s *fileService) DeleteByID(fileid string) error {
	return nil
}
