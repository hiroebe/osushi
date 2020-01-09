package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hiroebe/osushi/game"
)

func main() {
	game, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Osushi")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
