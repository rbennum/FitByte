package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/levensspel/go-gin-template/domain"
	"github.com/levensspel/go-gin-template/infrastructure"
	"github.com/samber/do/v2"
	"io"
	"log"
	"os"
	"sync"
)

var (
	AWS_REGION = os.Getenv("AWS_REGION")
	AWS_BUCKET = os.Getenv("AWS_BUCKET")

	s3StorageClientOnce     sync.Once
	s3StorageClientInstance *S3StorageClient
)

type S3StorageClient struct {
	s3Downloader *manager.Downloader
	s3Uploader   *manager.Uploader
	s3           *s3.Client
	sts          *sts.Client
}

func (s S3StorageClient) PutFile(
	ctx context.Context,
	key string,
	mimeType string,
	fileContent []byte,
	isPublic bool,
) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:        aws.String(AWS_BUCKET),
		Key:           aws.String(key),
		Body:          bytes.NewReader(fileContent),
		ContentLength: aws.Int64(int64(len(fileContent))),
		ContentType:   aws.String(mimeType),
		ACL: func() types.ObjectCannedACL {
			if isPublic {
				return types.ObjectCannedACLPublicRead
			}
			return types.ObjectCannedACLPrivate
		}(),
	}
	_, err := s.s3.PutObject(ctx, input)
	if err != nil {
		return "", err
	}

	return s.GetUrl(key), nil
}

func (s S3StorageClient) GetFileContent(ctx context.Context, key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(AWS_BUCKET),
		Key:    aws.String(key),
	}
	output, err := s.s3.GetObject(ctx, input)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := output.Body.Close(); closeErr != nil {
			log.Printf("failed to close body: %v", closeErr)
		}
	}()

	bytes, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (s S3StorageClient) GetUrl(key string) string {
	return fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s",
		AWS_BUCKET,
		AWS_REGION,
		key,
	)
}

func NewS3StorageClient() domain.StorageClient {
	s3StorageClientOnce.Do(func() {
		sdkConfig := infrastructure.NewAws()
		_s3 := s3.NewFromConfig(sdkConfig)
		downloader := manager.NewDownloader(_s3)
		uploader := manager.NewUploader(_s3)
		_sts := sts.NewFromConfig(sdkConfig)
		s3StorageClientInstance = &S3StorageClient{
			s3Downloader: downloader,
			s3Uploader:   uploader,
			s3:           _s3,
			sts:          _sts,
		}
	})
	return s3StorageClientInstance
}

func NewS3StorageClientInject(i do.Injector) (domain.StorageClient, error) {
	return NewS3StorageClient(), nil
}
