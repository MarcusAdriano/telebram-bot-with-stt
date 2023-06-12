package transcript

import "context"

type Transcriptor interface {
	Transcript(ctx context.Context, fileData []byte, fileName string) (*Transcription, error)
}
