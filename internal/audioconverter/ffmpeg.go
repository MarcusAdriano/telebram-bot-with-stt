package audioconverter

import (
	"context"
	"os/exec"

	"github.com/marcusadriano/sound-stt-tgbot/internal/fileserver"
)

type ffmpeg struct {
	fileServer fileserver.Fileserver
	CmdRunner  FfmpegCmdRunner
}

func NewFfmpeg(fs fileserver.Fileserver) AudioConverter {
	return &ffmpeg{
		fileServer: fs,
		CmdRunner:  &defaultCmdRunner{},
	}
}

func NewFfmpegWithCmdRunner(fs fileserver.Fileserver, cmdRunner FfmpegCmdRunner) AudioConverter {
	return &ffmpeg{
		fileServer: fs,
		CmdRunner:  cmdRunner,
	}
}

func (f *ffmpeg) ToMp3(ctx context.Context, fileData []byte, fileName string) (*Result, error) {

	fpath, err := f.fileServer.Save(ctx, fileserver.File{Name: fileName, Data: fileData})
	if err != nil {
		return nil, err
	}

	outputFilePath := fpath.Path + ".mp3"
	f.CmdRunner.Run("ffmpeg", "-i", fpath.Path, "-f", "mp3", "-ab", "192000", "-vn", outputFilePath)

	file, err := f.fileServer.Read(ctx, outputFilePath)
	if err != nil {
		return nil, err
	}

	_ = f.fileServer.Delete(ctx, outputFilePath)

	return &Result{Data: file.Data}, nil
}

type FfmpegCmdRunner interface {
	Run(name string, args ...string) error
}

type defaultCmdRunner struct {
}

func (d *defaultCmdRunner) Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}
