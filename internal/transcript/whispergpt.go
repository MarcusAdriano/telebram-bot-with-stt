package transcript

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/rs/zerolog"
)

type whisperGptTranscriptor struct {
	token  string
	apiUrl string
	logger *zerolog.Logger
}

func NewWhisperGptTranscriptor(logger *zerolog.Logger, token string) Transcriptor {
	return &whisperGptTranscriptor{
		token:  token,
		apiUrl: "https://api.openai.com/v1/audio/transcriptions",
		logger: logger,
	}
}

type WhisperResponse struct {
	Text  string `json:"text,omitempty"`
	Error struct {
		Message string `json:"message,omitempty"`
		Type    string `json:"type,omitempty"`
	} `json:"error"`
}

func (w *whisperGptTranscriptor) Transcript(fileData []byte, fileName string) (*Transcription, error) {

	buf := new(bytes.Buffer)
	formWriter := multipart.NewWriter(buf)

	file, _ := formWriter.CreateFormFile("file", fileName)
	file.Write(fileData)

	model, _ := formWriter.CreateFormField("model")
	model.Write([]byte("whisper-1"))

	formWriter.Close()
	req, err := http.NewRequest("POST", w.apiUrl, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", formWriter.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+w.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	w.logger.Info().Msgf("chatGPT response status code is %s", resp.Status)

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var transcript WhisperResponse
	json.Unmarshal(respBody, &transcript)
	if transcript.Error.Message != "" {
		return nil, fmt.Errorf("error to transcript file %s - error typer: %s", transcript.Error.Message, transcript.Error.Type)
	}

	return &Transcription{Text: transcript.Text}, nil

}
