package infrastructure

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

var (
	AWS_ACCESS_KEY_ID     = os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY = os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWS_REGION            = os.Getenv("AWS_REGION")
)

func NewAws() aws.Config {
	sdkConfig, err := awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithRegion(AWS_REGION),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			AWS_ACCESS_KEY_ID,
			AWS_SECRET_ACCESS_KEY,
			"",
		)),
	)
	if err != nil {
		panic(err)
	}
	return sdkConfig
}
