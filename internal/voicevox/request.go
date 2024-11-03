package voicevox

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/sysnote8main/readmyvc/internal/easyhttp"
)

func (i VoiceVoxInstance) requestAudioQuery(text string, speakerId int) ([]byte, error) {
	slog.Debug("Requesting audio query...", slog.Int("speakerId", speakerId))
	req, err := easyhttp.RequestPost(i.getUrl(fmt.Sprintf("audio_query?text=%s&speaker=%d", url.QueryEscape(text), speakerId)), nil)
	if err != nil {
		slog.Error("Failed to create post request", slog.Any("error", err))
		return nil, err
	}
	req.Header.Set("accept", "application/json")

	res, err := easyhttp.Do(req)
	if err != nil {
		slog.Error("Failed to request audio query", slog.Any("error", err))
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Failed to read request body", slog.Any("error", err))
		return nil, err
	}

	slog.Debug("Audio query request completed!", slog.Int("speakerId", speakerId))

	return b, nil
}

func (i VoiceVoxInstance) requestSynth(audioQueryData *[]byte, speakerId int) ([]byte, error) {
	slog.Debug("Requesting synth", slog.Int("speakerId", speakerId))
	req, err := easyhttp.RequestPost(i.getUrl(fmt.Sprintf("synthesis?speaker=%d&enable_interrogative_upspeak=true", speakerId)), bytes.NewReader(*audioQueryData))
	if err != nil {
		slog.Error("Failed to create post request", slog.Any("error", err))
		return nil, err
	}

	req.Header.Set("accept", "audio/wav")
	req.Header.Set("Content-type", "application/json")

	res, err := easyhttp.Do(req)
	if err != nil {
		slog.Error("Failed to request synth", slog.Any("error", err))
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Failed to read response body", slog.Any("error", err))
		return nil, err
	}

	slog.Debug("Synth completed", slog.Int("speakerId", speakerId))

	return b, nil
}

func (i VoiceVoxInstance) DoSynthAndSave(text string, speakerId int) (*string, error) {
	_audioQuery, err := i.requestAudioQuery(text, speakerId)
	if err != nil {
		return nil, err
	}

	_audioData, err := i.requestSynth(&_audioQuery, speakerId)
	if err != nil {
		return nil, err
	}

	fileName := string(time.Now().UnixMicro()) + ".wav"
	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	f.Write(_audioData)

	return &fileName, nil
}
