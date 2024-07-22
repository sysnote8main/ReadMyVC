package script

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

var prefix = "!tts"

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	author := m.Author
	if !author.Bot {
		if strings.HasPrefix(m.Content, prefix) {
			command := strings.Split(m.Content, " ")
			switch command[1] {
			case "s":
				Connect(s, m)
			case "e":
				Disconnect(s, m)
			}
		} else {
			TTS(m)
		}
	}
}
