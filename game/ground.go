package game

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	minMoutainWidth  = 200
	minMoutainHeight = 100
	maxMoutainWidth  = 500
	maxMoutainHeight = 300
)

var (
	moutainBaseImg *ebiten.Image
	groundBaseImg  *ebiten.Image

	groundSurfaceColor = color.NRGBA{0x00, 0x99, 0x00, 0xff}
	undergroundColor1  = color.NRGBA{0xcc, 0x99, 0x00, 0xff}
	undergroundColor2  = color.NRGBA{0x99, 0x66, 0x00, 0xff}
)

func init() {
	initMountainBaseImg()
	initGroundBaseImg()
}

func initMountainBaseImg() {
	const size = 512
	moutainBaseImg, _ = ebiten.NewImage(size, size, ebiten.FilterDefault)
	for x := 0; x < size; x++ {
		y := int(float64(size) / 2 * (1 + math.Cos(2*math.Pi/float64(size)*float64(x))))
		for y := y; y < size; y++ {
			moutainBaseImg.Set(x, y, groundSurfaceColor)
		}
	}
}

func initGroundBaseImg() {
	const size = 32
	const r = 4
	groundBaseImg, _ = ebiten.NewImage(size, size, ebiten.FilterDefault)
	groundBaseImg.Fill(undergroundColor1)
	for dx := -r; dx <= r; dx++ {
		for dy := -r; dy <= r; dy++ {
			groundBaseImg.Set(size/2+dx, size/2+dy, undergroundColor2)
		}
	}
}

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

	groundPatternImg *ebiten.Image
	groundSurfaceImg *ebiten.Image
	undergroundImg   *ebiten.Image
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
		if lastX >= float64(screenWidth)*scale+screenX {
			break
		}
		g.moutains = append(g.moutains, NewRandomMountain(lastX))
	}
}

func (g *Ground) Draw(screen *ebiten.Image, scale float64) {
	w, h := screen.Size()
	if g.groundPatternImg == nil || !g.checkImgSize(g.groundPatternImg, w, h) {
		g.groundPatternImg, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
	}
	if g.groundSurfaceImg == nil || !g.checkImgSize(g.groundSurfaceImg, w, h) {
		g.groundSurfaceImg, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
	}
	if g.undergroundImg == nil || !g.checkImgSize(g.undergroundImg, w, h) {
		g.undergroundImg, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
	}
	g.groundPatternImg.Clear()
	g.groundSurfaceImg.Clear()
	g.undergroundImg.Clear()

	g.drawGroundPattern(g.groundPatternImg, scale)
	g.drawGroundSurface(g.groundSurfaceImg, scale)
	g.drawUnderground(g.undergroundImg, scale)

	opts := &ebiten.DrawImageOptions{}

	opts.CompositeMode = ebiten.CompositeModeSourceAtop
	g.undergroundImg.DrawImage(g.groundPatternImg, opts)

	opts.CompositeMode = ebiten.CompositeModeSourceOver
	screen.DrawImage(g.groundSurfaceImg, &ebiten.DrawImageOptions{})
	screen.DrawImage(g.undergroundImg, &ebiten.DrawImageOptions{})
}

func (g *Ground) checkImgSize(img *ebiten.Image, w, h int) bool {
	w0, h0 := img.Size()
	return w == w0 && h == h0
}

func (g *Ground) drawGroundPattern(dstImg *ebiten.Image, scale float64) {
	w, h := dstImg.Size()
	offset := float64((w+int(g.screenX))%w) / scale
	srcW, srcH := groundBaseImg.Size()
	srcWs := float64(srcW) / scale
	srcHs := float64(srcH) / scale
	for x := 0.0; x <= float64(w)/srcWs+offset; x++ {
		for y := 0.0; y <= float64(h)/srcHs; y++ {
			dx := x*srcWs - offset
			dy := float64(h) - (y+1)*srcHs
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Scale(1/scale, 1/scale)
			opts.GeoM.Translate(dx, dy)
			dstImg.DrawImage(groundBaseImg, opts)
		}
	}
}

func (g *Ground) drawGroundSurface(dstImg *ebiten.Image, scale float64) {
	y := float64(screenHeight) - groundY/scale
	ebitenutil.DrawRect(dstImg, 0, y, float64(screenWidth), groundY/scale, groundSurfaceColor)
	for _, m := range g.moutains {
		g.drawMoutain(dstImg, m, scale, 1, groundY)
	}
}

func (g *Ground) drawUnderground(dstImg *ebiten.Image, scale float64) {
	for _, m := range g.moutains {
		g.drawMoutain(dstImg, m, scale, 1.2, 0)
	}
}

func (g *Ground) drawMoutain(screen *ebiten.Image, m *Mountain, scale, mountainScale, offsetY float64) {
	w, h := moutainBaseImg.Size()
	x := (m.StartX()-g.screenX)/scale + m.Width()*(1-1/mountainScale)/scale/2
	y := float64(screenHeight) - offsetY/scale - m.Height()/scale/mountainScale

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(m.Width()/float64(w)/scale/mountainScale, m.Height()/float64(h)/scale/mountainScale)
	opts.GeoM.Translate(x, y)
	screen.DrawImage(moutainBaseImg, opts)
}
