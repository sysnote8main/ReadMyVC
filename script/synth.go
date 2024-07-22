package script

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

type vcData struct {
	connection *discordgo.VoiceConnection
	channelID  string
	queue      *[]string
}

func getAudioData(s string) ([]byte, error) {
	// Truncate string for speak fluently
	str := s
	if utf8.RuneCountInString(str) > 30 {
		slice := []rune(str)
		strarr := []string{string(slice[:30]), "いかりゃく"}
		str = strings.Join(strarr, " ")
	}

	// === Audio Query ===
	// Generate url
	urlArr := []string{
		"http://localhost:50021/", // TODO able to change voicevox server
		"audio_query?text=",
		url.QueryEscape(str),
		"&speaker=8", // TODO Add Speaker Id settings
	}
	url := strings.Join(urlArr, "")

	// Build request
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("accept", "application/json")

	// Send request
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil { // On request error
		fmt.Println("Failed to request audio query", err)
		return nil, err
	}

	// === Synthesis Audio ===

	// Generate url
	// TODO can tune up voice parameter
	synthUrl := "http://localhost:50021/synthesis?speaker=8&enable_interrogative_upspeak=true"

	// Build request
	synthReq, _ := http.NewRequest("POST", synthUrl, res.Body)
	synthReq.Header.Set("accept", "audio/wav")
	synthReq.Header.Set("Content-type", "application/json")

	// Send request
	synthRes, err := client.Do(synthReq)
	if err != nil {
		fmt.Println("Failed to request synthesis", err)
		return nil, err
	}

	// === Return Bytes ===
	defer synthRes.Body.Close()
	buff := bytes.NewBuffer(nil)
	if _, err := io.Copy(buff, synthRes.Body); err != nil {
		fmt.Println("Failed to move data to buffer", err)
		return nil, err
	}
	return buff.Bytes(), nil // On success
}

func makeWaveFile(b []byte, guildId string) (string, error) {
	max := new(big.Int)
	max.SetInt64(int64(1000000))
	r, err := rand.Int(rand.Reader, max)
	if err != nil {
		fmt.Println("Failed to generate random number", err)
		return "", err
	}
	path := fmt.Sprintf("%s_%d.wav", guildId, r) // TODO move filename to hash
	file, _ := os.Create(path)
	defer func() {
		file.Close()
	}()
	file.Write(b)
	return path, nil
}

var vcMap map[string]vcData

func Connect(s *discordgo.Session, m *discordgo.MessageCreate) {
	if vcMap == nil {
		vcMap = make(map[string]vcData)
	}

	// check user in vc
	userstate, _ := s.State.VoiceState(m.GuildID, m.Author.ID)
	if userstate == nil {
		SendEmbed(s, m.ChannelID, "接続に失敗しました。", "呼び出す前にVCに接続しているか確認してください。")
		return
	}

	_, ok := vcMap[m.GuildID]
	if ok {
		SendEmbed(s, m.ChannelID, "接続は不要です。", "すでに、VCに接続しています。")
		return
	}

	vcsession, err := s.ChannelVoiceJoin(m.GuildID, userstate.ChannelID, false, false)
	if err != nil {
		SendEmbed(s, m.ChannelID, "エラーが発生しました。", "問題が発生したため、接続できませんでした。")
		fmt.Println("[Discord] Failed to connect vc", err)
		return
	} else {
		// get channel info
		textCh, _ := s.Channel(m.ChannelID)
		vcCh, _ := s.Channel(userstate.ChannelID)

		// generate embed field
		field := []*discordgo.MessageEmbedField{
			{
				Name:  "From",
				Value: textCh.Mention(),
			},
			{
				Name:  "To",
				Value: vcCh.Mention(),
			},
		}

		// send embed
		SendEmbedWithField(
			s,
			m.ChannelID,
			"読み上げ開始",
			"対象チャンネルでの読み上げを開始しました。",
			field,
		)
		slice := make([]string, 0, 10)
		newData := vcData{vcsession, m.ChannelID, &slice}
		vcMap[m.GuildID] = newData
	}
}

func Disconnect(s *discordgo.Session, m *discordgo.MessageCreate) {
	v, ok := vcMap[m.GuildID]
	if !ok {
		SendEmbed(s, m.ChannelID, "切断に失敗しました。", "BotはどこのVCにも参加していません！")
		return
	}

	err := v.connection.Disconnect()
	if err != nil {
		SendEmbed(s, m.ChannelID, "切断に失敗しました。", "問題が発生したため、切断に失敗しました。")
		fmt.Println("[Discord] Failed to disconnect vc", err)
		return
	} else {
		delete(vcMap, m.GuildID)
		SendEmbed(s, m.ChannelID, "退出完了", "VCから切断しました！")
	}
}

func play(v *vcData, path string, force bool) {
	// Add to queue
	if len(*v.queue) > 0 && !force {
		*v.queue = append(*v.queue, path)
		return
	}
	*v.queue = append(*v.queue, path)

	// vc setup
	vc := v.connection
	vc.Speaking(true)
	defer vc.Speaking(false)

	// Remove file
	defer os.Remove(path)

	// Go to next on end
	defer func(v *vcData) {
		*v.queue = (*v.queue)[1:]
		if len(*v.queue) > 0 {
			pt := (*v.queue)[0]
			*v.queue = (*v.queue)[1:]
			play(v, pt, true)
		}
	}(v)
	dgvoice.PlayAudioFile(vc, path, make(chan bool))
}

func TTS(m *discordgo.MessageCreate) {
	go func() {
		v, ok := vcMap[m.GuildID]
		if ok {
			if v.channelID == m.ChannelID {
				b, err := getAudioData(m.Content)
				if err != nil {
					fmt.Println("Failed to get audio data", err)
					return
				}
				var path string
				path, err = makeWaveFile(b, m.GuildID)
				if err != nil {
					fmt.Println("Failed to create wave file", err)
					return
				}
				play(&v, path, false)
			}
		}
	}()
}
