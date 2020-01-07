package main

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hiroebe/osushi"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	game, err := osushi.NewGame()
	if err != nil {
		panic(err)
	}
	ebiten.SetWindowTitle("title")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
