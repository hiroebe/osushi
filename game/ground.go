package game

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

const (
	groundY           = 16
	minMountainWidth  = 200
	minMountainHeight = 100
	maxMountainWidth  = 500
	maxMountainHeight = 300
)

var (
	mountainBaseImg     *ebiten.Image
	undergroundBaseImg  *ebiten.Image
	surfaceColorBaseImg *ebiten.Image

	groundSurfaceColor = color.NRGBA{0x00, 0x99, 0x00, 0xff}
	undergroundColor1  = color.NRGBA{0xcc, 0x99, 0x00, 0xff}
	undergroundColor2  = color.NRGBA{0x99, 0x66, 0x00, 0xff}
)

func init() {
	initMountainBaseImg()
	initGroundBaseImg()
	surfaceColorBaseImg, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	surfaceColorBaseImg.Fill(groundSurfaceColor)
}

func initMountainBaseImg() {
	const size = 512
	mountainBaseImg, _ = ebiten.NewImage(size, size, ebiten.FilterDefault)
	for x := 0; x < size; x++ {
		y := int(float64(size) / 2 * (1 + math.Cos(2*math.Pi/float64(size)*float64(x))))
		for y := y; y < size; y++ {
			mountainBaseImg.Set(x, y, groundSurfaceColor)
		}
	}
}

func initGroundBaseImg() {
	const (
		size     = 128
		cellSize = 32
		dotSize  = 8
	)
	undergroundBaseImg, _ = ebiten.NewImage(size, size, ebiten.FilterDefault)
	undergroundBaseImg.Fill(undergroundColor1)
	groundDotImg, _ := ebiten.NewImage(dotSize, dotSize, ebiten.FilterDefault)
	groundDotImg.Fill(undergroundColor2)
	for i := 0; i < size/cellSize; i++ {
		for j := 0; j < size/cellSize; j++ {
			centerX := float64(cellSize*i + cellSize/2)
			centerY := float64(cellSize*j + cellSize/2)
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(centerX-dotSize, centerY-dotSize)
			undergroundBaseImg.DrawImage(groundDotImg, opts)
		}
	}
}

type Mountain struct {
	startX        float64
	width, height float64
}

func NewRandomMountain(startX float64) *Mountain {
	width := minMountainWidth + rand.Float64()*(maxMountainWidth-minMountainWidth)
	height := minMountainHeight + rand.Float64()*(maxMountainHeight-minMountainHeight)
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
	mountains []*Mountain
	screenX   float64
	img       *ebiten.Image
}

func (g *Ground) At(x float64) (y, grad float64) {
	for _, m := range g.mountains {
		if x >= m.StartX() && x <= m.EndX() {
			return m.At(x)
		}
	}
	return 0, 0
}

func (g *Ground) Update(screenX, scale float64) {
	g.screenX = screenX
	if g.mountains == nil {
		g.mountains = make([]*Mountain, 0, 64)
	}
	if len(g.mountains) == 0 {
		m := &Mountain{startX: -maxMountainWidth / 2, width: maxMountainWidth, height: maxMountainHeight}
		g.mountains = append(g.mountains, m)
	}
	if g.mountains[0].EndX() < screenX {
		copy(g.mountains, g.mountains[1:])
		g.mountains = g.mountains[:len(g.mountains)-1]
	}
	for {
		lastX := g.mountains[len(g.mountains)-1].EndX()
		if lastX >= float64(screenWidth)/scale+screenX {
			break
		}
		g.mountains = append(g.mountains, NewRandomMountain(lastX))
	}
}

func (g *Ground) Draw(screen *ebiten.Image, scale float64) {
	w, h := screen.Size()
	if g.img == nil || !g.checkImgSize(g.img, w, h) {
		g.img, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
	} else {
		g.img.Clear()
	}

	g.drawUnderground(g.img, scale)
	g.drawGroundPattern(g.img, scale)
	g.drawGroundSurface(g.img, scale)

	screen.DrawImage(g.img, &ebiten.DrawImageOptions{})
}

func (g *Ground) drawUnderground(dstImg *ebiten.Image, scale float64) {
	for _, m := range g.mountains {
		g.drawMountain(dstImg, m, scale, 0.8, 0)
	}
}

func (g *Ground) drawGroundPattern(dstImg *ebiten.Image, scale float64) {
	dw, dh := dstImg.Size()
	srcW, srcH := undergroundBaseImg.Size()
	offset := float64((dw+int(g.screenX))%dw%srcW) * scale
	sw := float64(srcW) * scale
	sh := float64(srcH) * scale

	lastX := (float64(dw) + offset) / sw
	lastY := float64(dh) / sh

	opts := &ebiten.DrawImageOptions{}
	opts.CompositeMode = ebiten.CompositeModeSourceIn
	for x := 0.0; x <= lastX; x++ {
		for y := 0.0; y <= lastY; y++ {
			dx := x*sw - offset
			dy := float64(dh) - (y+1)*sh
			opts.GeoM.Reset()
			opts.GeoM.Scale(scale, scale)
			opts.GeoM.Translate(dx, dy)
			dstImg.DrawImage(undergroundBaseImg, opts)
		}
	}
}

func (g *Ground) drawGroundSurface(dstImg *ebiten.Image, scale float64) {
	y := float64(screenHeight) - groundY*scale
	opts := &ebiten.DrawImageOptions{}
	opts.CompositeMode = ebiten.CompositeModeDestinationOver
	opts.GeoM.Scale(float64(screenWidth), groundY*scale)
	opts.GeoM.Translate(0, y)
	dstImg.DrawImage(surfaceColorBaseImg, opts)

	for _, m := range g.mountains {
		g.drawMountain(dstImg, m, scale, 1, groundY)
	}
}

func (g *Ground) drawMountain(dstImg *ebiten.Image, m *Mountain, scale, mtScale, offsetY float64) {
	w, h := mountainBaseImg.Size()
	x := (m.StartX()-g.screenX)*scale + m.Width()*(1-mtScale)*scale/2
	y := float64(screenHeight) - offsetY*scale - m.Height()*scale*mtScale

	opts := &ebiten.DrawImageOptions{}
	opts.CompositeMode = ebiten.CompositeModeDestinationOver
	opts.GeoM.Scale(m.Width()/float64(w)*scale*mtScale, m.Height()/float64(h)*scale*mtScale)
	opts.GeoM.Translate(x, y)
	dstImg.DrawImage(mountainBaseImg, opts)
}

func (g *Ground) checkImgSize(img *ebiten.Image, w, h int) bool {
	w0, h0 := img.Size()
	return w == w0 && h == h0
}
