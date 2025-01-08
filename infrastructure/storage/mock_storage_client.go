package storage

import (
	"context"
	"errors"
	"github.com/levensspel/go-gin-template/domain"
	"github.com/samber/do/v2"
	"strings"
)

type MockStorageClient struct{}

func (m MockStorageClient) PutFile(
	ctx context.Context,
	key string,
	mimeType string,
	fileContent []byte,
	isPublic bool,
) (string, error) {
	if strings.Contains(key, "mock_failed") {
		return "", errors.New("Failed to put file")
	}

	return m.GetUrl(key), nil
}

func (m MockStorageClient) GetFileContent(ctx context.Context, key string) ([]byte, error) {
	if strings.Contains(key, "mock_failed") {
		return nil, errors.New("Failed to get file content")
	}

	return []byte{0x01, 0x02, 0x03}, nil
}

func (m MockStorageClient) GetUrl(key string) string {
	return "https://tmssl.akamaized.net/images/foto/galerie/cristiano-ronaldo-im-trikot-von-portugal-1718197560-139337.jpg"
}

func NewMockStorageClient() domain.StorageClient {
	return MockStorageClient{}
}

func NewMockStorageClientInject(i do.Injector) (domain.StorageClient, error) {
	return NewMockStorageClient(), nil
}
