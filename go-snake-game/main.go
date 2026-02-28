package main

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/thuta/go-snake/common"
	"github.com/thuta/go-snake/entity"
	"github.com/thuta/go-snake/game"
	"github.com/thuta/go-snake/math"
)

var mplusFaceSource *text.GoTextFaceSource

// Building a Snake Game with Go and Ebiten
func main() {
	s, err := text.NewGoTextFaceSource(
		bytes.NewReader(
			fonts.MPlus1pRegular_ttf,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s

	world := game.NewWorld()
	world.AddEntity(
		entity.NewPlayer(
			math.Point{
				X: common.ScreenWidth / common.GridSize / 2,
				Y: common.ScreenHeight / common.GridSize / 2,
			},
			math.DirRight,
		),
	)
	world.AddEntity(entity.NewFood())

	g := &Game{
		world: world,
	}

	ebiten.SetWindowSize(common.ScreenWidth, common.ScreenHeight)
	ebiten.SetWindowTitle("Go Snake")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
