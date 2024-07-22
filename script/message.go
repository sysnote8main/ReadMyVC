package script

import "github.com/bwmarrin/discordgo"

func SendMessage(s *discordgo.Session, channelID, msg string) (e error) {
	_, err := s.ChannelMessageSend(channelID, msg)
	return err
}

var embedColor = 0x880088

func SendEmbedWithField(s *discordgo.Session, channelID, title, description string, field []*discordgo.MessageEmbedField) (e error) {
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       embedColor,
		Title:       title,
		Description: description,
		Fields:      field,
	}
	_, err := s.ChannelMessageSendEmbed(channelID, embed)
	return err
}

func SendEmbed(s *discordgo.Session, channelID, title, description string) (e error) {
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       embedColor,
		Title:       title,
		Description: description,
	}
	_, err := s.ChannelMessageSendEmbed(channelID, embed)
	return err
}
