package bot

import (
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"

	"github.com/sysnote8main/readmyvc/internal/diswrap"
)

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	author := m.Author
	if author.Bot {
		return
	}

	slog.Debug("Message fired!", slog.String("userId", author.ID))

	args := strings.Split(m.Content, " ")
	if args[0] == prefix {
		switch args[1] {
		case "s":
			slog.Debug("Starting tts", slog.String("guildId", m.GuildID), slog.String("userId", author.ID))
			userstate, _ := s.State.VoiceState(m.GuildID, author.ID)
			if userstate == nil {
				diswrap.SendErrorEmbed(
					s,
					m.ChannelID,
					"Failed to connect vc",
					"接続に失敗しました。\n呼び出す前に、VCに接続しているか確認してください。",
				)
				return
			}
			err := vcManager.Connect(s, m.GuildID, userstate.ChannelID, m.ChannelID)
			if err != nil {
				slog.Error("Failed to connect vc. check log")
			}
			return
		case "e":
			err := vcManager.Disconnect(s, m.GuildID, m.ChannelID)
			if err != nil {
				slog.Error("Failed to disconnect from vc. check log")
			}
			return
		}
	} else {
		vcData := vcManager.GetVCData(m.GuildID)
		if vcData == nil {
			slog.Debug("VCData is nil", slog.String("userId", author.ID))
			return
		}
		if vcData.TextChId != m.ChannelID {
			slog.Debug("Text Channel's Id isn't matched", slog.String("userId", author.ID))
			return
		}
		// TODO set speakerId
		slog.Debug("TTS fired!", slog.String("userId", author.ID))
		msgContent := m.Content
		// TODO support change truncate size
		if utf8.RuneCountInString(msgContent) > 30 {
			slice := []rune(msgContent)
			strarr := []string{string(slice[:30]), "いかりゃく"}
			msgContent = strings.Join(strarr, " ")
		}
		voiceVox.TTS(msgContent, 8, vcData)
	}
}
