package discordvc

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"

	"github.com/sysnote8main/readmyvc/internal/diswrap"
)

type vcManager struct {
	vcMap map[string]*VCData // guildId -> vcData
}

func NewVCManager() vcManager {
	return vcManager{
		vcMap: make(map[string]*VCData),
	}
}

func (m vcManager) IsVCConnected(guildId string) bool {
	_, ok := m.vcMap[guildId]
	return ok
}

func (m vcManager) Connect(session *discordgo.Session, guildId, vcChId, textChId string) error {
	textCh, err := session.Channel(textChId)
	if err != nil {
		diswrap.SendErrorEmbed(
			session,
			textChId,
			"An error occurred.",
			"予想外のエラーが発生しました。\n開発者にお問い合わせください。",
		)
		return err
	}
	vcCh, err := session.Channel(vcChId)
	if err != nil {
		diswrap.SendErrorEmbed(
			session,
			textChId,
			"An error occurred.",
			"予想外のエラーが発生しました。\n開発者にお問い合わせください。",
		)
		return err
	}
	if m.IsVCConnected(guildId) {
		err = m.vcMap[guildId].Disconnect()
		if err != nil {
			slog.Error("Failed to disconnect from vc", slog.Any("error", err))
			return err
		}
		diswrap.SendWarnEmbed(
			session,
			m.vcMap[guildId].TextChId,
			"Bot was moved",
			"Botは別のVCに移動しました。",
			// 移動先は、権限等の関係もあるのを考慮してとりあえず出さない実装
		)
		delete(m.vcMap, guildId)
	}

	vcConn, err := session.ChannelVoiceJoin(guildId, vcChId, false, false)
	if err != nil {
		diswrap.SendErrorEmbed(
			session,
			textChId,
			"Failed to join VC",
			"VCの接続に失敗しました。\n開発者にお問い合わせください。",
		)
		return fmt.Errorf("failed to connect vc: %v", err)
	}

	m.vcMap[guildId] = NewVCData(vcConn, textChId)

	diswrap.SendSuccessEmbed(
		session,
		textChId,
		"Start reading",
		"読み上げを開始しました。",
		&discordgo.MessageEmbedField{
			Name:  "From",
			Value: textCh.Mention(),
		},
		&discordgo.MessageEmbedField{
			Name:  "To",
			Value: vcCh.Mention(),
		},
	)

	return nil
}

func (m vcManager) Disconnect(session *discordgo.Session, guildId, cmdFiredChId string) error {
	if m.IsVCConnected(guildId) {
		vcdata := m.vcMap[guildId]
		err := vcdata.Disconnect()
		if err != nil {
			slog.Error("Failed to disconnect from vc", slog.Any("error", err))
			diswrap.SendErrorEmbed(
				session,
				cmdFiredChId,
				"An error occurred.",
				"予想外のエラーが発生しました。\n開発者にお問い合わせください。",
			)
			return err
		}
		delete(m.vcMap, guildId)
		diswrap.SendSuccessEmbed(
			session,
			cmdFiredChId,
			"Disconnected!",
			"正常に切断しました。\nご利用ありがとうございました！",
		)
		return nil
	}

	diswrap.SendWarnEmbed(
		session,
		cmdFiredChId,
		"No vc to disconnect",
		"このサーバー内のどのVCにも参加していないようです。",
	)
	return nil
}

func (m vcManager) GetVCData(guildId string) *VCData {
	return m.vcMap[guildId]
}

func (m vcManager) ChangeTextCh(session *discordgo.Session, guildId string, targetTextChId string) error {
	diswrap.SendWarnEmbed(
		session,
		targetTextChId,
		"Not implemented yet",
		"まだ実装されてないです。ごめんね",
	)
	return nil
}
