package osushi

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type Player struct {
	x, y      float64
	vx, vy    float64
	isJumping bool
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

	if isKeyJustReleased() && grad > 0 {
		p.vy += gravity
		p.isJumping = true
		return
	}

	g *= grad / obl
	p.vx += g / obl
	p.vy += g * grad / obl
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
	const size = 32
	x := (playerOffset - size/2) / scale
	y := screenHeight - (size+p.y)/scale
	ebitenutil.DrawRect(screen, x, y, size/scale, size/scale, color.NRGBA{0xff, 0x00, 0x00, 0xff})
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
