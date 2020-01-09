package mobile

import (
	"log"

	"github.com/hajimehoshi/ebiten/mobile"
	"github.com/hiroebe/osushi/game"
)

func init() {
	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	mobile.SetGame(g)
}

// Dummy is a dummy exported function.
//
// gomobile doesn't compile a package that doesn't include any exported function.
// Dummy forces gomobile to compile this package.
func Dummy() {}
