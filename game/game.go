package game

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
)

const (
	fontSize     = 16
	iconSize     = 32
	playerOffset = 64
)

var (
	screenWidth  int
	screenHeight int
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type volumeSetter interface {
	SetVolume(volume float64)
}

type soundIcon struct {
	setters []volumeSetter
	isMuted bool
}

func (i *soundIcon) Draw(screen *ebiten.Image, x, y, w, h int) {
	imgW, imgH := i.Size()
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(float64(w)/float64(imgW), float64(h)/float64(imgH))
	opts.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(i.img(), opts)
}

func (i *soundIcon) Size() (w, h int) {
	return i.img().Size()
}

func (i *soundIcon) OnClick() {
	i.setMuted(!i.isMuted)
}

func (i *soundIcon) setMuted(muted bool) {
	i.isMuted = muted
	var volume float64
	if !i.isMuted {
		volume = 1
	}
	for _, setter := range i.setters {
		setter.SetVolume(volume)
	}
}

func (i *soundIcon) img() *ebiten.Image {
	if i.isMuted {
		return soundIconOff
	}
	return soundIconOn
}

type Game struct {
	player           *Player
	ground           *Ground
	soundIcon        Element
	scale            float64
	jumpHeightRecord int
	jumpLendthRecord int

	newRecordSound *NewRecordSound
}

func NewGame() (*Game, error) {
	jumpSound := NewJumpSound()
	newRecordSound := NewNewRecordSound()

	soundIcon := &soundIcon{
		setters: []volumeSetter{jumpSound, newRecordSound},
	}
	soundIcon.setMuted(true)
	soundIconElem := NewElement(soundIcon)
	soundIconElem.SetSize(iconSize, iconSize)

	return &Game{
		player: &Player{
			jumpSound: jumpSound,
		},
		ground:         &Ground{},
		soundIcon:      soundIconElem,
		scale:          1,
		newRecordSound: newRecordSound,
	}, nil
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.player.Update(g.ground.At(g.player.x))
	g.scale = float64(screenHeight) / (g.player.y + playerOffset*4)
	if g.scale > 1 {
		g.scale = 1
	}
	g.ground.Update(g.player.x-playerOffset, g.scale)
	g.soundIcon.Update()
	g.updateRecord()

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.Fill(color.White)
	g.ground.Draw(screen, g.scale)
	g.player.Draw(screen, g.scale)
	g.soundIcon.Draw(screen)
	g.drawScore(screen)

	// ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))

	return nil
}

func (g *Game) updateRecord() {
	if h := int(g.player.jumpHeight); h > g.jumpHeightRecord {
		if h/100 > g.jumpHeightRecord/100 {
			g.newRecordSound.Update()
		}
		g.jumpHeightRecord = h
	} else if h == 0 {
		g.newRecordSound.Reset()
	}
	if l := int(g.player.jumpLength); l > g.jumpLendthRecord {
		g.jumpLendthRecord = l
	}
}

func (g *Game) drawScore(screen *ebiten.Image) {
	texts := []string{
		fmt.Sprintf("Height: %6d (%6d)", int(g.player.jumpHeight), g.jumpHeightRecord),
		fmt.Sprintf("Length: %6d (%6d)", int(g.player.jumpLength), g.jumpLendthRecord),
	}
	for i, t := range texts {
		x := screenWidth - fontSize*len(t)
		y := fontSize*(i+1) + iconSize
		text.Draw(screen, t, arcadeFont, x, y, color.Black)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	screenWidth = outsideWidth
	screenHeight = outsideHeight
	g.soundIcon.SetPosition(screenWidth-iconSize, 0)
	return screenWidth, screenHeight
}
