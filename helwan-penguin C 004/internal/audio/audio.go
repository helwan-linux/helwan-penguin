package audio

import (
	"bytes"
	"embed"
	"io/fs"

	ebitenaudio "github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

//go:embed *.wav
var soundFiles embed.FS

// Manager owns small built-in effects; no external sound files or services are needed.
type Manager struct {
	context *ebitenaudio.Context
	players map[string]*ebitenaudio.Player
}

func NewManager() *Manager {
	m := &Manager{
		context: ebitenaudio.NewContext(44100),
		players: make(map[string]*ebitenaudio.Player),
	}
	for _, name := range []string{"jump", "collect", "hit", "stomp", "boss", "win", "level"} {
		data, err := fs.ReadFile(soundFiles, name+".wav")
		if err != nil {
			continue
		}
		stream, err := wav.DecodeWithSampleRate(m.context.SampleRate(), bytes.NewReader(data))
		if err != nil {
			continue
		}
		player, err := ebitenaudio.NewPlayer(m.context, stream)
		if err != nil {
			continue
		}
		player.SetVolume(0.45)
		m.players[name] = player
	}
	return m
}

func (m *Manager) Play(name string) {
	if m == nil {
		return
	}
	p := m.players[name]
	if p == nil {
		return
	}
	if p.IsPlaying() {
		p.Pause()
	}
	if err := p.Rewind(); err != nil {
		return
	}
	p.Play()
}
