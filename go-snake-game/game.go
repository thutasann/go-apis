package main

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Point struct {
	x, y int
}

type Game struct {
	snake      []Point
	direction  Point
	lastUpdate time.Time
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.direction = dirUp
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.direction = dirDown
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.direction = dirLeft
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.direction = dirRight
	}

	if time.Since(g.lastUpdate) < gameSpeed {
		return nil
	}
	g.lastUpdate = time.Now()

	g.updateSnake(&g.snake, g.direction)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.snake {
		vector.FillRect(
			screen,
			float32(p.x*gridSize),
			float32(p.y*gridSize),
			gridSize,
			gridSize,
			color.White,
			true,
		)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
