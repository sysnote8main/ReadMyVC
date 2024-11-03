package diswrap

import "github.com/bwmarrin/discordgo"

var (
	COLOR_SUCCESS = 0x57F287
	COLOR_WARN    = 0xFEE75C
	COLOR_ERROR   = 0xED4245
)

func SendEmbed(session *discordgo.Session, chId, title, description string, embedColor int, fields ...*discordgo.MessageEmbedField) (*discordgo.Message, error) {
	return session.ChannelMessageSendEmbed(chId, &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       embedColor,
		Title:       title,
		Description: description,
		Fields:      fields,
	})
}

func SendSuccessEmbed(session *discordgo.Session, chId, title, description string, fields ...*discordgo.MessageEmbedField) (*discordgo.Message, error) {
	return SendEmbed(session, chId, title, description, COLOR_SUCCESS, fields...)
}

func SendWarnEmbed(session *discordgo.Session, chId, title, description string, fields ...*discordgo.MessageEmbedField) (*discordgo.Message, error) {
	return SendEmbed(session, chId, title, description, COLOR_WARN, fields...)
}

func SendErrorEmbed(session *discordgo.Session, chId, title, description string, fields ...*discordgo.MessageEmbedField) (*discordgo.Message, error) {
	return SendEmbed(session, chId, title, description, COLOR_ERROR, fields...)
}
