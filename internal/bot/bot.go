package bot

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/sysnote8main/readmyvc/internal/discordvc"
	"github.com/sysnote8main/readmyvc/internal/voicevox"
)

var (
	prefix    = "!tts"
	vcManager = discordvc.NewVCManager()
	voiceVox  = voicevox.VoiceVoxInstance{Host: "localhost:50021"}
)

func Start() {
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		slog.Error("Failed to start discord session", slog.Any("error", err))
		os.Exit(1)
	}

	dg.AddHandler(OnMessageCreate)

	err = dg.Open()
	if err != nil {
		slog.Error("Failed to open bot connection", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Successfully to login!", slog.String("botUsername", dg.State.User.Username))

	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stopBot

	slog.Info("Bot is stopping now...")

	err = dg.Close()
	if err != nil {
		slog.Error("Failed to close connection", slog.Any("error", err))
	}

	slog.Info("Bye bye!")
}
