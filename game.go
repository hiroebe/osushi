package osushi

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth      = 640
	screenHeight     = 480
	playerOffset     = 32
	groundY          = 16
	minVx            = 1
	gravity          = 0.1
	minMoutainWidth  = 128
	minMoutainHeight = 64
	maxMoutainWidth  = 256
	maxMoutainHeight = 256
)

type Game struct {
	player *Player
	ground *Ground
	scale  float64
}

func NewGame() (*Game, error) {
	return &Game{
		player: &Player{},
		ground: &Ground{},
		scale:  1,
	}, nil
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.player.Update(g.ground.At(g.player.x))
	g.scale = (g.player.y + playerOffset*4) / screenHeight
	if g.scale < 1 {
		g.scale = 1
	}
	g.ground.Update(g.player.x-playerOffset, g.scale)

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.Fill(color.White)
	g.player.Draw(screen, g.scale)
	g.ground.Draw(screen, g.scale)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return screenWidth, screenHeight
}
