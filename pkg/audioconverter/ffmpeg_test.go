package audioconverter_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/marcusadriano/tgbot-stt/pkg/audioconverter"
	"github.com/marcusadriano/tgbot-stt/pkg/fileserver"
	"github.com/marcusadriano/tgbot-stt/pkg/mocks"
)

type cmdRunnerMock struct {
}

func (c *cmdRunnerMock) Run(name string, args ...string) error {
	return nil
}

func TestFfmpegToMp3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const mp3FileName = "fileName.mp3"

	fileServer := mocks.NewMockFileserver(ctrl)
	fileServer.
		EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(&fileserver.FilePath{Path: "fileName"}, nil).
		Times(1)

	fileServer.
		EXPECT().
		Read(gomock.Any(), mp3FileName).
		Return(&fileserver.File{Data: []byte("fileData")}, nil).
		Times(1)

	fileServer.
		EXPECT().
		Delete(gomock.Any(), gomock.Eq(mp3FileName)).
		Return(nil).
		Times(1)

	converter := audioconverter.NewFfmpegWithCmdRunner(fileServer, &cmdRunnerMock{})
	result, err := converter.ToMp3(context.TODO(), []byte("fileData"), "fileName")

	if err != nil {
		t.Fatalf("Error to convert file: %s", err)
	}

	if string(result.Data) != "fileData" {
		t.Fatalf("Expected fileData, got %s", result.Data)
	}

	if result.Filename != mp3FileName {
		t.Fatalf("Expected fileName.mp3, got %s", result.Filename)
	}
}
