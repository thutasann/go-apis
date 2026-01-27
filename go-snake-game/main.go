package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
	gridSize     = 20
)

// Building a Snake Game with Go and Ebiten
func main() {
	g := &Game{
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
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Go Snake")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
