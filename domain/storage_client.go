package domain

import "context"

type StorageClient interface {
	// PutFile puts a file to the storage.
	// The key is the filename or path in the storage.
	// The mimeType is the content type of the file.
	// The fileContent is the actual content of the file.
	// The isPublic flag determines whether the file is publicly accessible.
	// It returns the URL of the uploaded file on success, or an error on failure.
	PutFile(
		ctx context.Context,
		key string,
		mimeType string,
		fileContent []byte,
		isPublic bool,
	) (string, error)
	// GetFileContent retrieves the content of a file from the storage.
	// The key is the filename or path in the storage.
	// It returns the content of the file on success, or an error on failure.
	GetFileContent(ctx context.Context, key string) ([]byte, error)
	// GetUrl generates the complete URL for a file in the storage.
	// The key is the filename or path in the storage.
	// It returns the file's URL as a string.
	GetUrl(key string) string
}
