package main

import (
	"log"
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
		log.Fatal(err)
	}
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Osushi")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
