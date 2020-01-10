package game

import (
	"errors"
	"io"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/audio"
)

const (
	sampleRate = 44100
)

var audioContext *audio.Context

func init() {
	var err error
	audioContext, err = audio.NewContext(sampleRate)
	if err != nil {
		log.Fatal(err)
	}
}

type JumpSound struct {
	player *audio.Player
	timer  *time.Timer
}

func NewJumpSound() *JumpSound {
	s := &JumpSound{}

	var err error
	s.player, err = audio.NewPlayer(audio.CurrentContext(), s.wave(440))
	if err != nil {
		log.Println(err)
		return nil
	}
	return s
}

func (s *JumpSound) Start() {
	if s.player.Current() != 0 {
		s.player.Rewind()
	}

	// Do not play sound for a short (< 100ms) jump
	s.timer = time.AfterFunc(100*time.Millisecond, func() {
		s.player.Play()
	})
}

func (s *JumpSound) Stop() {
	if !s.timer.Stop() {
		s.player.Pause()
	}
}

func (s *JumpSound) wave(freq float64) *Wave {
	p0 := 0.15
	p1 := 0.3
	x0 := (1 - p1) / (p1 - p0) * sampleRate
	x1 := p0 * x0

	return NewWave(freq, func(x, lambda float64) float64 {
		r := (x + x1) / (x + x0)
		return math.Sin(2 * math.Pi / lambda * x * r)
	})
}

var baseFreq = []float64{
	261.626,
	293.665,
	329.628,
	349.228,
	391.995,
	440.000,
	493.883,
}

type NewRecordSound struct {
	freqIdx int
	octave  int
	players map[float64]*audio.Player
}

func NewNewRecordSound() *NewRecordSound {
	s := &NewRecordSound{
		freqIdx: 0,
		octave:  1,
		players: make(map[float64]*audio.Player, 32),
	}
	for _, freq := range baseFreq {
		for i := 1; i <= 3; i++ {
			s.createPlayer(freq * float64(i))
		}
	}

	return s
}

func (s *NewRecordSound) Update() {
	freq := baseFreq[s.freqIdx] * float64(s.octave)
	p, ok := s.players[freq]
	if !ok {
		p = s.createPlayer(freq)
	}
	if p.Current() != 0 {
		p.Rewind()
	}
	p.Play()
	time.AfterFunc(100*time.Millisecond, func() {
		p.Pause()
	})

	s.freqIdx++
	if s.freqIdx >= len(baseFreq) {
		s.freqIdx = 0
		s.octave *= 2
	}
}

func (s *NewRecordSound) Reset() {
	s.freqIdx = 0
	s.octave = 1
}

func (s *NewRecordSound) createPlayer(freq float64) *audio.Player {
	p, err := audio.NewPlayer(audio.CurrentContext(), s.wave(freq))
	if err != nil {
		log.Println(err)
		return nil
	}

	s.players[freq] = p
	return p
}

func (s *NewRecordSound) wave(freq float64) *Wave {
	return NewWave(freq, func(x, lambda float64) float64 {
		return math.Sin(2 * math.Pi / lambda * x)
	})
}

type Wave struct {
	freq float64
	f    func(x, lambda float64) float64
	pos  int64
}

func NewWave(freq float64, f func(x, lambda float64) float64) *Wave {
	return &Wave{
		freq: freq,
		f:    f,
	}
}

func (w *Wave) Read(buf []byte) (int, error) {
	n := len(buf) / 4
	for i := 0; i < n; i++ {
		val := w.f(float64(w.pos), float64(sampleRate)/w.freq)
		b := int16(val * 0.3 * math.MaxInt16)
		idx := i * 4
		buf[idx] = byte(b)
		buf[idx+1] = byte(b >> 8)
		buf[idx+2] = byte(b)
		buf[idx+3] = byte(b >> 8)
		w.pos++
	}

	return n * 4, nil
}

func (w *Wave) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		w.pos = offset
	case io.SeekCurrent:
		w.pos += offset
	case io.SeekEnd:
		return 0, errors.New("SeekEnd: End of wave is not defined")
	}

	return w.pos, nil
}

func (w *Wave) Close() error {
	return nil
}
