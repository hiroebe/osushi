package game

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type Element interface {
	Update()
	Draw(screen *ebiten.Image)
	SetPosition(x, y int)
	Position() (x, y int)
	SetSize(w, h int)
	Size() (w, h int)
}

type ElementImpl interface {
	Draw(screen *ebiten.Image, x, y, w, h int)
	Size() (w, h int)
	OnClick()
}

func NewElement(impl ElementImpl) Element {
	return &ElementBase{
		impl: impl,
	}
}

type ElementBase struct {
	impl ElementImpl

	x, y    int
	w, h    int
	touchID int
}

func (e *ElementBase) Update() {
	if ids := inpututil.JustPressedTouchIDs(); len(ids) == 1 {
		id := ids[0]
		x, y := ebiten.TouchPosition(id)
		if e.isInside(x, y) {
			e.touchID = id
		}
	}
	if e.touchID > 0 {
		if !e.isInside(ebiten.TouchPosition(e.touchID)) {
			e.touchID = 0
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		cursorX, cursorY := ebiten.CursorPosition()
		if e.isInside(cursorX, cursorY) {
			e.impl.OnClick()
		}
	}
	if inpututil.IsTouchJustReleased(e.touchID) {
		e.impl.OnClick()
		e.touchID = 0
	}
}

func (e *ElementBase) Draw(screen *ebiten.Image) {
	x, y := e.Position()
	w, h := e.Size()
	e.impl.Draw(screen, x, y, w, h)
}

func (e *ElementBase) SetPosition(x, y int) {
	e.x = x
	e.y = y
}

func (e *ElementBase) Position() (x, y int) {
	return e.x, e.y
}

func (e *ElementBase) SetSize(w, h int) {
	e.w = w
	e.h = h
}

func (e *ElementBase) Size() (w, h int) {
	if e.w == 0 && e.h == 0 {
		return e.impl.Size()
	}
	return e.w, e.h
}

func (e *ElementBase) isInside(cursorX, cursorY int) bool {
	x, y := e.Position()
	w, h := e.Size()
	return x <= cursorX && cursorX < x+w && y <= cursorY && cursorY < y+h
}
