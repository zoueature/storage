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
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

func (c *client) Upload(ctx context.Context, reader io.Reader, keyOps ...storage.KeyOperate) (string, error) {
	content, name, err := objectName(reader, keyOps...)
	if err != nil {
		return "", err
	}
	mimeType := http.DetectContentType(content)
	_, err = c.cli.PutObject(ctx, &s3.PutObjectInput{
		Key:         aws.String(name),
		Bucket:      aws.String(c.bucket),
		Body:        bytes.NewReader(content),
		ContentType: aws.String(mimeType),
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", c.domain, name), nil
}

// UploadByCustomKey 自定义key上传
func (c *client) UploadByCustomKey(ctx context.Context, reader io.Reader, objectKey string) (string, error) {
	_, err := c.cli.PutObject(ctx, &s3.PutObjectInput{
		Key:    aws.String(objectKey),
		Bucket: aws.String(c.bucket),
		Body:   reader,
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", c.domain, objectKey), nil
}

func objectName(reader io.Reader, keyOps ...storage.KeyOperate) ([]byte, string, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, "", err
	}
	h := md5.New()
	h.Write(content)
	hashBytes := h.Sum(nil)
	name := hex.EncodeToString(hashBytes)
	for _, opKey := range keyOps {
		opKey(&name)
	}
	return content, name, nil
}

func (c *client) UploadToPublic(ctx context.Context, reader io.Reader, keyOps ...storage.KeyOperate) (string, error) {
	content, name, err := objectName(reader, keyOps...)

	mimeType := http.DetectContentType(content)
	_, err = c.cli.PutObject(ctx, &s3.PutObjectInput{
		Key:         aws.String(name),
		Bucket:      aws.String(c.bucket),
		Body:        bytes.NewReader(content),
		ContentType: aws.String(mimeType),
		ACL:         types.ObjectCannedACLPublicRead,
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

// GetContent 获取存储内容
func (c *client) GetContent(ctx context.Context, objectKey string) ([]byte, error) {
	resp, err := c.cli.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}
