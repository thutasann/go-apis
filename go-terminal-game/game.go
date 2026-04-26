package main

import (
	"bytes"
	"fmt"
	"os"
)

type game struct {
	isRunning bool
	level     *level

	drawBuf *bytes.Buffer
}

func newGame(width, height int) *game {
	lvl := newLevel(width, height)
	return &game{
		level:   lvl,
		drawBuf: new(bytes.Buffer),
	}
}

func (g *game) start() {
	g.isRunning = true
	g.loop()
}

func (g *game) loop() {
	for g.isRunning {
		g.update()
		g.render()
	}
}

func (g *game) update() {}

func (g *game) renderLevel() {
	for h := range g.level.height {
		for w := range g.level.width {
			if g.level.data[h][w] == NOTHING {
				g.drawBuf.WriteString(" ")
			}
			if g.level.data[h][w] == WALL {
				g.drawBuf.WriteString("☐")
			}
		}
		g.drawBuf.WriteString("\n")
	}
}

func (g *game) render() {
	g.renderLevel()
	fmt.Fprint(os.Stdout, g.drawBuf.String())
}
