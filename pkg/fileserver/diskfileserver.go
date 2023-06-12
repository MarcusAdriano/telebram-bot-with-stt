package fileserver

import (
	"context"
	"os"
	"path/filepath"

	"github.com/marcusadriano/tgbot-stt/internal/logger"
)

type diskFileserver struct {
	path string
}

func NewDiskFileserver(path string) Fileserver {
	return &diskFileserver{
		path: path,
	}
}

func (d *diskFileserver) Read(ctx context.Context, name string) (*File, error) {
	logger.Log(ctx).Info().Msgf("Reading file %s", name)

	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	return &File{
		Name: name,
		Data: data,
	}, nil
}

func (d *diskFileserver) Save(ctx context.Context, file File) (*FilePath, error) {

	logger.Log(ctx).Info().Msgf("Saving file %s at %s", file.Name, d.path)

	fileFullPath := d.path + "/" + file.Name
	if err := os.MkdirAll(filepath.Dir(fileFullPath), 0770); err != nil {
		return nil, err
	}
	fd, err := os.Create(fileFullPath)
	if err != nil {
		logger.Log(ctx).Error().Msgf("Error to create file %s", err.Error())
		return nil, err
	}
	defer func() {
		if err := fd.Close(); err != nil {
			logger.Log(ctx).Error().Msgf("Error to close file %s", err.Error())
		}
	}()

	_, _ = fd.Write(file.Data)

	return &FilePath{
		Path: fileFullPath,
	}, nil
}

func (d *diskFileserver) Delete(ctx context.Context, fileName string) error {
	logger.Log(ctx).Info().Msgf("Deleting file %s", fileName)
	return os.Remove(fileName)
}
