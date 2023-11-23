package storage

import (
	"context"
	"io"
)

type Storage interface {
	Type() string
	Upload(ctx context.Context, key string, reader io.Reader, bucket ...string) (string, error)
}
