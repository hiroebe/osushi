package mobile

import (
	"log"

	"github.com/hajimehoshi/ebiten/mobile"
	"github.com/hiroebe/osushi/game"
)

//go:generate env GO111MODULE=off ebitenmobile bind -target android -javapkg com.hiroebe.osushi -o ./android/osushi/osushi.aar .

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
