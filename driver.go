package storage

import (
	"context"
	"io"
)

type Storage interface {
	Type() string
	Upload(ctx context.Context, reader io.Reader, key ...string) (string, error)
	SignAccessURL(ctx context.Context, objectKey string, ttl ...int) (string, error)
}
