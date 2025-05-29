package storage

import (
	"context"
	"errors"
)

var (
	ErrConnectStorage = errors.New("failed to connect storage")
)

type Storage interface {
	Save(ctx context.Context, quote string) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (string, error)
	List(ctx context.Context) ([]string, error)
	Close()
}
