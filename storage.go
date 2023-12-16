package storage

import (
	"errors"
)

import "github.com/zoueature/config"

type storageDriver func(cfg config.StorageConfig) Storage

var driver storageDriver

func RegisterStorageDriver(c storageDriver) {
	driver = c
}

// New 实例化存储客户端
func New(cfg config.StorageConfig) Storage {
	if driver == nil {
		panic(errors.New("storage driver not register"))
	}
	s := driver(cfg)
	return s
}
