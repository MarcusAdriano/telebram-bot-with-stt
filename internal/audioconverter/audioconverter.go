package audioconverter

import "context"

type AudioConverter interface {
	ToMp3(ctx context.Context, fileData []byte, fileName string) (*Result, error)
}

type Result struct {
	Data     []byte
	Filename string
}
