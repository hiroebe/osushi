package osushi

import (
	"image"
	"math"
	"net/http"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	_ "github.com/hiroebe/osushi/statik"
	"github.com/rakyll/statik/fs"
)

const (
	minVx   = 1
	gravity = 0.1
)

var (
	gopherImageNormal     *ebiten.Image
	gopherImageAcceralate *ebiten.Image
	gopherImageFly1       *ebiten.Image
	gopherImageFly2       *ebiten.Image
)

func init() {
	statikFs, err := fs.New()
	if err != nil {
		panic(err)
	}
	gopherImageNormal = mustLoadImage(statikFs, "/gopher-normal.png")
	gopherImageAcceralate = mustLoadImage(statikFs, "/gopher-acceralate.png")
	gopherImageFly1 = mustLoadImage(statikFs, "/gopher-fly-1.png")
	gopherImageFly2 = mustLoadImage(statikFs, "/gopher-fly-2.png")
}

func mustLoadImage(fs http.FileSystem, name string) *ebiten.Image {
	f, err := fs.Open(name)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	ebitenImg, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	return ebitenImg
}

type Player struct {
	x, y      float64
	vx, vy    float64
	isJumping bool
	img       *ebiten.Image
	imgFrames int
}

func (p *Player) Update(gy, grad float64) {
	p.updateV(grad)

	if p.vx < minVx {
		p.vx = minVx
	}
	p.x += p.vx
	p.y += p.vy

	if !p.isJumping || p.y < gy {
		p.y = gy
		if p.isJumping {
			p.land(grad)
		}
	}

	p.updateImg()
}

func (p *Player) updateV(grad float64) {
	g := -gravity
	if isKeyPressed() {
		g *= 2
	}
	if p.isJumping {
		p.vy += g
		return
	}

	obl := math.Sqrt(1 + grad*grad)
	v := math.Sqrt(p.vx*p.vx + p.vy*p.vy)
	p.vx = v / obl
	p.vy = v * grad / obl

	if isKeyJustReleased() && grad >= -0.1 {
		p.vy += gravity
		p.isJumping = true
		return
	}

	g *= grad / obl
	p.vx += g / obl
	p.vy += g * grad / obl
}

func (p *Player) updateImg() {
	if isKeyPressed() {
		p.img = gopherImageAcceralate
		return
	}
	if !p.isJumping {
		p.img = gopherImageNormal
		return
	}
	p.imgFrames++
	if p.imgFrames < 10 {
		return
	}
	p.imgFrames = 0
	if p.img == gopherImageFly1 {
		p.img = gopherImageFly2
	} else {
		p.img = gopherImageFly1
	}

}

func (p *Player) land(grad float64) {
	p.isJumping = false
	obl := math.Sqrt(1 + grad*grad)
	dv := (p.vx + p.vy*grad) / obl
	if dv < 0 {
		p.vx = 0
		p.vy = 0
		return
	}
	p.vx = dv / obl
	p.vy = dv * grad / obl
}

func (p *Player) Draw(screen *ebiten.Image, scale float64) {
	w, h := p.img.Size()
	x := playerOffset / scale
	y := screenHeight - p.y/scale + float64(h)/10
	grad := -p.vy / p.vx

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-float64(w)/2, -float64(h))
	opts.GeoM.Rotate(math.Atan(grad))
	opts.GeoM.Translate(x, y)

	screen.DrawImage(p.img, opts)
}

func isKeyPressed() bool {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		return true
	}
	return false
}

func isKeyJustReleased() bool {
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		return true
	}
	return false
}
