package transcript

type Transcriptor interface {
	Transcript(fileData []byte, fileName string) (*Transcription, error)
}
