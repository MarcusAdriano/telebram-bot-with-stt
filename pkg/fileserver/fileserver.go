package fileserver

import (
	"context"
)

type Fileserver interface {
	Save(ctx context.Context, file File) (*FilePath, error)
	Read(ctx context.Context, fileName string) (*File, error)
	Delete(ctx context.Context, fileName string) error
}
