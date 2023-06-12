package audioconverter

import (
	"context"
	"os/exec"

	fileserver2 "github.com/marcusadriano/tgbot-stt/pkg/fileserver"
)

type ffmpeg struct {
	fileServer fileserver2.Fileserver
	CmdRunner  FfmpegCmdRunner
}

func NewFfmpeg(fs fileserver2.Fileserver) AudioConverter {
	return &ffmpeg{
		fileServer: fs,
		CmdRunner:  &defaultCmdRunner{},
	}
}

func NewFfmpegWithCmdRunner(fs fileserver2.Fileserver, cmdRunner FfmpegCmdRunner) AudioConverter {
	return &ffmpeg{
		fileServer: fs,
		CmdRunner:  cmdRunner,
	}
}

func (f *ffmpeg) ToMp3(ctx context.Context, fileData []byte, fileName string) (*Result, error) {

	fpath, err := f.fileServer.Save(ctx, fileserver2.File{Name: fileName, Data: fileData})
	if err != nil {
		return nil, err
	}

	outputFilePath := fpath.Path + ".mp3"
	err = f.CmdRunner.Run("ffmpeg", "-i", fpath.Path, "-f", "mp3", "-ab", "192000", "-vn", outputFilePath)
	if err != nil {
		return nil, err
	}

	file, err := f.fileServer.Read(ctx, outputFilePath)
	if err != nil {
		return nil, err
	}

	_ = f.fileServer.Delete(ctx, outputFilePath)

	return &Result{Data: file.Data, Filename: outputFilePath}, nil
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
