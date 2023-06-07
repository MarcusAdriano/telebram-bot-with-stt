package fileserver

import (
	"context"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

type diskFileserver struct {
	logger *zerolog.Logger
	path   string
}

func NewDiskFileserver(logger *zerolog.Logger, path string) Fileserver {
	return &diskFileserver{
		logger: logger,
		path:   path,
	}
}

func (d *diskFileserver) Read(_ context.Context, name string) (*File, error) {
	d.logger.Info().Msgf("Reading file %s", name)

	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	return &File{
		Name: name,
		Data: data,
	}, nil
}

func (d *diskFileserver) Save(_ context.Context, file File) (*FilePath, error) {

	d.logger.Info().Msgf("Saving file %s at %s", file.Name, d.path)

	fileFullPath := d.path + "/" + file.Name
	if err := os.MkdirAll(filepath.Dir(fileFullPath), 0770); err != nil {
		return nil, err
	}
	fd, err := os.Create(fileFullPath)
	if err != nil {
		d.logger.Error().Msgf("Error to create file %s", err.Error())
		return nil, err
	}
	defer func() {
		if err := fd.Close(); err != nil {
			d.logger.Error().Msgf("Error to close file %s", err.Error())
		}
	}()

	_, _ = fd.Write(file.Data)

	return &FilePath{
		Path: fileFullPath,
	}, nil
}

func (d *diskFileserver) Delete(_ context.Context, fileName string) error {
	d.logger.Info().Msgf("Deleting file %s", fileName)
	return os.Remove(fileName)
}
