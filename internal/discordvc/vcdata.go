package discordvc

import (
	"os"
	"sync/atomic"

	"github.com/adrianbrad/queue"
	"github.com/bwmarrin/discordgo"
	"github.com/daystram/dgvoice"
)

type VCData struct {
	conn      *discordgo.VoiceConnection
	TextChId  string
	queue     *queue.Blocking[string]
	isPlaying atomic.Bool
}

func NewVCData(vcConn *discordgo.VoiceConnection, textChId string) *VCData {
	return &VCData{
		conn:      vcConn,
		TextChId:  textChId,
		queue:     queue.NewBlocking([]string{}),
		isPlaying: atomic.Bool{},
	}
}

func (d *VCData) AddQueueAndPlay(filePath string) {
	// add to queue
	d.queue.OfferWait(filePath)
	d.playQueue()
}

func (d *VCData) playQueue() {
	// if on playing, returns
	if d.isPlaying.Load() {
		return
	}

	// start playing all elements of queue
	vc := d.conn

	// start speaking
	vc.Speaking(true)

	// on finish, stop to notify speaking
	defer vc.Speaking(false)

	defer func() {
		if d.queue.Size() > 0 {
			d.playQueue()
		}
	}()

	path := d.queue.GetWait()

	// remove temp wav file
	defer os.Remove(path)

	dgvoice.PlayAudioFile(d.conn, path, make(chan bool))
}

func (d *VCData) Disconnect() error {
	return d.conn.Disconnect()
}
