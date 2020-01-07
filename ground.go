package osushi

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Mountain struct {
	startX        float64
	width, height float64
}

func NewRandomMountain(startX float64) *Mountain {
	width := minMoutainWidth + rand.Float64()*(maxMoutainWidth-minMoutainWidth)
	height := minMoutainHeight + rand.Float64()*(maxMoutainHeight-minMoutainHeight)
	return &Mountain{startX: startX, width: width, height: height}
}

func (m *Mountain) StartX() float64 {
	return m.startX
}

func (m *Mountain) EndX() float64 {
	return m.startX + m.width
}

func (m *Mountain) TopX() float64 {
	return m.startX + m.width/2
}

func (m *Mountain) Width() float64 {
	return m.width
}

func (m *Mountain) Height() float64 {
	return m.height
}

func (m *Mountain) At(x float64) (y, grad float64) {
	x -= m.StartX()
	y = groundY + m.Height()/2*(1-math.Cos(2*math.Pi/m.Width()*x))
	grad = m.Height() / m.Width() * math.Pi * math.Sin(2*math.Pi/m.Width()*x)
	return y, grad
}

type Ground struct {
	moutains []*Mountain
	screenX  float64
	baseImg  *ebiten.Image
}

func (g *Ground) At(x float64) (y, grad float64) {
	for _, m := range g.moutains {
		if x >= m.StartX() && x <= m.EndX() {
			return m.At(x)
		}
	}
	return 0, 0
}

func (g *Ground) Update(screenX, scale float64) {
	g.screenX = screenX
	if g.moutains == nil {
		g.moutains = make([]*Mountain, 0, 64)
	}
	if len(g.moutains) == 0 {
		m := &Mountain{startX: -maxMoutainWidth / 2, width: maxMoutainWidth, height: maxMoutainHeight}
		g.moutains = append(g.moutains, m)
	}
	if g.moutains[0].EndX() < screenX {
		copy(g.moutains, g.moutains[1:])
		g.moutains = g.moutains[:len(g.moutains)-1]
	}
	for {
		lastX := g.moutains[len(g.moutains)-1].EndX()
		if lastX >= screenWidth*scale+screenX {
			break
		}
		g.moutains = append(g.moutains, NewRandomMountain(lastX))
	}
}

func (g *Ground) Draw(screen *ebiten.Image, scale float64) {
	clr := color.NRGBA{0x00, 0xff, 0x00, 0xff}
	if g.baseImg == nil {
		w, h := screen.Size()
		g.baseImg, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
		for x := 0; x < w; x++ {
			y := int(float64(h) / 2 * (1 + math.Cos(2*math.Pi/float64(w)*float64(x))))
			for y := y; y < h; y++ {
				g.baseImg.Set(x, y, clr)
			}
		}
	}
	ebitenutil.DrawRect(screen, 0, screenHeight-groundY/scale, screenWidth, groundY/scale, clr)
	for _, m := range g.moutains {
		g.drawMoutain(screen, m, scale)
	}
}

func (g *Ground) drawMoutain(screen *ebiten.Image, m *Mountain, scale float64) {
	w, h := g.baseImg.Size()
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(m.Width()/float64(w)/scale, m.Height()/float64(h)/scale)
	x := (m.StartX() - g.screenX) / scale
	y := float64(screenHeight) - groundY/scale - m.Height()/scale
	opts.GeoM.Translate(x, y)
	screen.DrawImage(g.baseImg, opts)
}
