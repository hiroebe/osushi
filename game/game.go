package game

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

const (
	fontSize     = 16
	playerOffset = 64
	groundY      = 16
)

var (
	screenWidth  int
	screenHeight int
)

var arcadeFont font.Face

func init() {
	rand.Seed(time.Now().UnixNano())

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
	jumpHeightRecord int
	jumpLendthRecord int

	newRecordSound *NewRecordSound
}

func NewGame() (*Game, error) {
	return &Game{
		player:         &Player{},
		ground:         &Ground{},
		scale:          1,
		newRecordSound: NewNewRecordSound(),
	}, nil
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.player.Update(g.ground.At(g.player.x))
	g.scale = float64(screenHeight) / (g.player.y + playerOffset*4)
	if g.scale > 1 {
		g.scale = 1
	}
	g.ground.Update(g.player.x-playerOffset, g.scale)
	g.updateRecord()

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
		y := fontSize * (i + 2)
		text.Draw(screen, t, arcadeFont, x, y, color.Black)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	screenWidth = outsideWidth
	screenHeight = outsideHeight
	return screenWidth, screenHeight
}
