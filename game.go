package osushi

import (
	"fmt"
	"image/color"
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

const (
	screenWidth  = 640
	screenHeight = 480
	fontSize     = 16
	playerOffset = 32
	groundY      = 16
)

var arcadeFont font.Face

func init() {
	tt, err := truetype.Parse(fonts.ArcadeN_ttf)
	if err != nil {
		log.Fatal(err)
	}
	arcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

type Game struct {
	player           *Player
	ground           *Ground
	scale            float64
	maxHeight        int
	maxHeightRecord  int
	jumpLendth       int
	jumpLendthRecord int
	jumpStartX       int
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
	g.updateHeight()
	g.updateLength()

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.Fill(color.White)
	g.ground.Draw(screen, g.scale)
	g.player.Draw(screen, g.scale)
	g.drawScore(screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))

	return nil
}

func (g *Game) updateHeight() {
	if !g.player.isJumping {
		g.maxHeight = 0
		return
	}
	h := int(g.player.y) - groundY
	if h > g.maxHeight {
		g.maxHeight = h
		if h > g.maxHeightRecord {
			g.maxHeightRecord = h
		}
	}
}

func (g *Game) updateLength() {
	if !g.player.isJumping {
		g.jumpLendth = 0
		g.jumpStartX = 0
		return
	}
	if g.jumpStartX <= 0 {
		g.jumpStartX = int(g.player.x)
	}
	g.jumpLendth = int(g.player.x) - g.jumpStartX
	if g.jumpLendth > g.jumpLendthRecord {
		g.jumpLendthRecord = g.jumpLendth
	}
}

func (g *Game) drawScore(screen *ebiten.Image) {
	texts := []string{
		fmt.Sprintf("Height: %6d (%6d)", g.maxHeight, g.maxHeightRecord),
		fmt.Sprintf("Length: %6d (%6d)", g.jumpLendth, g.jumpLendthRecord),
	}
	for i, t := range texts {
		x := screenWidth - fontSize*len(t)
		y := fontSize * (i + 2)
		text.Draw(screen, t, arcadeFont, x, y, color.Black)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return screenWidth, screenHeight
}
