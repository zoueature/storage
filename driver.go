package storage

import (
	"context"
	"io"
)

type KeyOperate func(string2 *string)

type Storage interface {
	Type() string
	Upload(ctx context.Context, reader io.Reader, key ...KeyOperate) (string, error)
	SignAccessURL(ctx context.Context, objectKey string, ttl ...int) (string, error)
}
