package storage

import (
	"context"
	"io"
	"time"
)

type IStorageProviderRepository interface {
	Exists(ctx context.Context, path string) bool
	Get(ctx context.Context, path string) (string, error)
	ReadStream(ctx context.Context, path string) (io.ReadCloser, error)
	Put(ctx context.Context, path string, contents interface{}, options ...interface{}) error
	WriteStream(ctx context.Context, path string, reader io.Reader, options ...interface{}) error
	GetVisibility(ctx context.Context, path string) (string, error)
	SetVisibility(ctx context.Context, path, visibility string) error
	Prepend(ctx context.Context, path, data string) error
	Append(ctx context.Context, path, data string) error
	Delete(ctx context.Context, paths ...string) error
	Copy(ctx context.Context, from, to string) error
	Move(ctx context.Context, from, to string) error
	Size(ctx context.Context, path string) (int64, error)
	LastModified(ctx context.Context, path string) (int64, error)
	Files(ctx context.Context, directory string, recursive bool) ([]string, error)
	AllFiles(ctx context.Context, directory string) ([]string, error)
	Directories(ctx context.Context, directory string, recursive bool) ([]string, error)
	AllDirectories(ctx context.Context, directory string) ([]string, error)
	MakeDirectory(ctx context.Context, path string) error
	DeleteDirectory(ctx context.Context, directory string) error
	GenerateURL(ctx context.Context, path string, expires time.Duration) (*string, error)
}
