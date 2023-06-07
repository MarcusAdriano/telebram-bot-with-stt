package fileserver

import (
	"context"
)

type Fileserver interface {
	Save(context context.Context, file File) (*FilePath, error)
	Read(context context.Context, fileName string) (*File, error)
}
