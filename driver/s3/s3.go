package s3

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/zoueature/config"
	"github.com/zoueature/storage"
	"io"
	"net/http"
	"time"
)

const defaultAccessTTL = 300

func init() {
	storage.RegisterStorageDriver(NewS3Storage)
}

type client struct {
	cli       *s3.Client
	preSigner *s3.PresignClient
	bucket    string
	domain    string
}

func NewS3Storage(cfg config.StorageConfig) storage.Storage {
	s3Client := s3.New(s3.Options{
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.AccessSecret, "")),
		Region:      cfg.Region,
	})
	return &client{
		cli:       s3Client,
		preSigner: s3.NewPresignClient(s3Client),
		bucket:    cfg.Bucket,
		domain:    cfg.Domain,
	}
}

func (c *client) Type() string {
	return "S3"
}

type keyOption func(string2 *string)

func (c *client) Upload(ctx context.Context, reader io.Reader, keyOps ...storage.KeyOperate) (string, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	h := md5.New()
	h.Write(content)
	hashBytes := h.Sum(nil)
	objectName := hex.EncodeToString(hashBytes)
	for _, opKey := range keyOps {
		opKey(&objectName)
	}
	mimeType := http.DetectContentType(content)
	_, err = c.cli.PutObject(ctx, &s3.PutObjectInput{
		Key:         aws.String(objectName),
		Bucket:      aws.String(c.bucket),
		Body:        bytes.NewReader(content),
		ContentType: aws.String(mimeType),
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", c.domain, objectName), nil
}

func (c *client) SignAccessURL(ctx context.Context, objectKey string, ttl ...int) (string, error) {
	expireTime := defaultAccessTTL
	if len(ttl) > 0 && ttl[0] != 0 {
		expireTime = ttl[0]
	}
	request, err := c.preSigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expireTime) * time.Second
	})
	if err != nil {
		return "", err
	}
	return request.URL, nil
}
