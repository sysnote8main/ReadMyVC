package voicevox

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/sysnote8main/readmyvc/internal/discordvc"
)

type VoiceVoxInstance struct {
	Host string
}

func (i VoiceVoxInstance) getUrl(path ...string) string {
	_url := fmt.Sprintf("http://%s/%s", i.Host, strings.Join(path, "/"))
	slog.Debug("URL generated!", slog.String("url", _url))
	return _url
}

func (i VoiceVoxInstance) TTS(text string, speakerId int, vcData *discordvc.VCData) {
	go func() {
		queryData, err := i.requestAudioQuery(text, speakerId)
		if err != nil {
			slog.Error("Failed to request audio query", slog.Any("error", err))
			return
		}
		audioData, err := i.requestSynth(&queryData, speakerId)
		if err != nil {
			slog.Error("Failed to request synth", slog.Any("error", err))
			return
		}
		fileName := string(time.Now().Nanosecond()) + ".wav"
		f, err := os.Create(fileName)
		if err != nil {
			slog.Error("Failed to create file", slog.Any("error", err))
			return
		}
		f.Write(audioData)
		f.Close()
		vcData.AddQueueAndPlay(fileName)
	}()
}
