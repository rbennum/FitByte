package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/levensspel/go-gin-template/domain"
	"github.com/samber/do/v2"
)

type MockStorageClient struct{}

func (m MockStorageClient) PutFile(
	ctx context.Context,
	key string,
	mimeType string,
	fileContent []byte,
	isPublic bool,
) (string, error) {
	// Simulasi kegagalan jika key mengandung "mock_failed"
	if strings.Contains(key, "mock_failed") {
		return "", errors.New("Failed to put file")
	}

	// Tentukan lokasi penyimpanan file
	savePath := fmt.Sprintf("./.uploads/%s", key)

	// Buat direktori jika belum ada
	if err := os.MkdirAll("./.uploads", os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Simpan file ke direktori
	if err := os.WriteFile(savePath, fileContent, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Kembalikan URL file yang disimpan
	return m.GetUrl(key), nil
}

func (m MockStorageClient) GetFileContent(ctx context.Context, key string) ([]byte, error) {
	if strings.Contains(key, "mock_failed") {
		return nil, errors.New("Failed to get file content")
	}

	return []byte{0x01, 0x02, 0x03}, nil
}

func (m MockStorageClient) GetUrl(key string) string {
	// Buat URI untuk file
	baseURL := "/uploads"
	return fmt.Sprintf("%s/%s", baseURL, key)
}

func NewMockStorageClient() domain.StorageClient {
	return MockStorageClient{}
}

func NewMockStorageClientInject(i do.Injector) (domain.StorageClient, error) {
	return NewMockStorageClient(), nil
}
