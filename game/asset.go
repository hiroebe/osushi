package game

import (
	"image"
	_ "image/png"
	"log"
	"net/http"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	_ "github.com/hiroebe/osushi/game/statik"
	"github.com/rakyll/statik/fs"
	"golang.org/x/image/font"
)

//go:generate statik -m -src images

var (
	gopherImageNormal     *ebiten.Image
	gopherImageAcceralate *ebiten.Image
	gopherImageFly1       *ebiten.Image
	gopherImageFly2       *ebiten.Image
	soundIconOn           *ebiten.Image
	soundIconOff          *ebiten.Image

	arcadeFont font.Face
)

func init() {
	statikFs, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	gopherImageNormal = mustLoadImage(statikFs, "/gopher-normal.png")
	gopherImageAcceralate = mustLoadImage(statikFs, "/gopher-acceralate.png")
	gopherImageFly1 = mustLoadImage(statikFs, "/gopher-fly-1.png")
	gopherImageFly2 = mustLoadImage(statikFs, "/gopher-fly-2.png")
	soundIconOn = mustLoadImage(statikFs, "/volume-on.png")
	soundIconOff = mustLoadImage(statikFs, "/volume-off.png")

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

func mustLoadImage(fs http.FileSystem, name string) *ebiten.Image {
	f, err := fs.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	ebitenImg, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	return ebitenImg
}
