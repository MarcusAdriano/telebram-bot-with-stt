package audioconverter

import (
	"context"
	"os/exec"

	"github.com/marcusadriano/sound-stt-tgbot/internal/fileserver"
	"github.com/rs/zerolog"
)

type ffmpeg struct {
	tmpDir     string
	fileServer fileserver.Fileserver
}

func NewFfmpeg(logger *zerolog.Logger) AudioConverter {
	tmpDir := "/tmp"
	return &ffmpeg{
		tmpDir:     tmpDir,
		fileServer: fileserver.NewDiskFileserver(logger, tmpDir),
	}
}

func (f *ffmpeg) ToMp3(ctx context.Context, fileData []byte, fileName string) (*Result, error) {

	fpath, err := f.fileServer.Save(ctx, fileserver.File{Name: fileName, Data: fileData})
	if err != nil {
		return nil, err
	}

	outputFilePath := f.tmpDir + "/" + fileName + ".mp3"
	cmd := exec.Command("ffmpeg", "-i", fpath.Path, "-f", "mp3", "-ab", "192000", "-vn", outputFilePath)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	file, err := f.fileServer.Read(ctx, outputFilePath)
	if err != nil {
		return nil, err
	}

	return &Result{Data: file.Data}, nil
}
