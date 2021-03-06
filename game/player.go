package game

import (
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	minV     = 2
	gravity  = 0.05
	friction = 0.02
)

type Player struct {
	jumpSound *JumpSound

	x, y      float64
	vx, vy    float64
	isJumping bool

	jumpHeight float64
	jumpLength float64
	jumpStartX float64

	img       *ebiten.Image
	imgFrames int
}

func (p *Player) Update(gy, grad float64) {
	obl := math.Sqrt(1 + grad*grad)

	p.updateV(grad, obl)

	p.x += p.vx
	p.y += p.vy

	if !p.isJumping || p.y < gy {
		p.y = gy
		if p.isJumping {
			p.land(grad, obl)
		}
	}

	p.updateJumpScore()
	p.updateImg()
}

func (p *Player) updateV(grad, obl float64) {
	g := -gravity
	if isKeyPressed() {
		g *= 3
	}
	if p.isJumping {
		p.vy += g
		return
	}

	v := math.Sqrt(p.vx*p.vx+p.vy*p.vy) + g*grad/obl - friction/obl
	if v < minV {
		v = minV
	}
	p.vx = v / obl
	p.vy = v * grad / obl

	if isKeyJustReleased() {
		p.jump(grad, obl)
		return
	}
}

func (p *Player) jump(grad, obl float64) {
	p.jumpSound.Start()

	p.isJumping = true
	p.jumpStartX = p.x

	p.vy += gravity / obl
}

func (p *Player) land(grad, obl float64) {
	p.jumpSound.Stop()

	p.isJumping = false

	dv := (p.vx + p.vy*grad) / obl
	if dv < 0 {
		p.vx = 0
		p.vy = 0
		return
	}
	if p.jumpLength > minMountainWidth {
		dv *= 1.1
	}
	p.vx = dv / obl
	p.vy = dv * grad / obl
}

func (p *Player) updateJumpScore() {
	if !p.isJumping {
		p.jumpHeight = 0
		p.jumpLength = 0
		return
	}
	p.jumpLength = p.x - p.jumpStartX
	if p.y > p.jumpHeight {
		p.jumpHeight = p.y
	}
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

func (p *Player) Draw(screen *ebiten.Image, scale float64) {
	w, h := p.img.Size()
	x := playerOffset * scale
	y := float64(screenHeight) - p.y*scale + float64(h)/10
	grad := -p.vy / p.vx

	opts := &ebiten.DrawImageOptions{}
	opts.Filter = ebiten.FilterLinear
	opts.GeoM.Translate(-float64(w)/2, -float64(h))
	opts.GeoM.Rotate(math.Atan(grad))
	opts.GeoM.Scale(scale, scale)
	opts.GeoM.Translate(x, y)

	screen.DrawImage(p.img, opts)
}

var touching bool

func isKeyPressed() bool {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		return true
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return true
	}
	if len(ebiten.TouchIDs()) > 0 {
		touching = true
		return true
	}
	return false
}

func isKeyJustReleased() bool {
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		return true
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return true
	}
	if touching && len(ebiten.TouchIDs()) == 0 {
		touching = false
		return true
	}
	return false
}
