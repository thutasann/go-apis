package main

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

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

	g := &Game{
		gameOver: true,
		snake: []Point{
			{
				x: screenWidth / gridSize / 2,
				y: screenHeight / gridSize / 2,
			},
			{
				x: screenWidth/gridSize/2 - 1,
				y: screenHeight / gridSize / 2,
			},
		},
		direction: Point{x: 1, y: 0},
	}

	g.spwanFood()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Go Snake")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
