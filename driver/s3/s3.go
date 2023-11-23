package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jiebutech/config"
	"github.com/jiebutech/storage"
	"io"
)

func init() {
	storage.RegisterStorageDriver(NewS3Storage)
}

type client struct {
	cli           *s3.Client
	defaultBucket string
	domain        string
}

func NewS3Storage(cfg config.StorageConfig) storage.Storage {
	s3Client := s3.New(s3.Options{Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.AccessSecret, ""))})
	return &client{
		cli:           s3Client,
		defaultBucket: cfg.DefaultBucket,
		domain:        cfg.Domain,
	}
}

func (c *client) Type() string {
	return "S3"
}

func (c *client) bucket(bucket ...string) *string {
	b := c.defaultBucket
	if len(bucket) > 0 && bucket[0] != "" {
		b = bucket[0]
	}
	return aws.String(b)
}

func (c *client) Upload(ctx context.Context, key string, reader io.Reader, bucket ...string) (string, error) {

	_, err := c.cli.PutObject(ctx, &s3.PutObjectInput{
		Key:    aws.String(key),
		Bucket: c.bucket(bucket...),
		Body:   reader,
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", c.domain, key), nil
}
