package main

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

type game struct {
	isRunning bool
	level     *level
	stats     *stats

	drawBuf *bytes.Buffer
}

func newGame(width, height int) *game {
	lvl := newLevel(width, height)
	return &game{
		level:   lvl,
		drawBuf: new(bytes.Buffer),
		stats:   newStats(),
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
		g.stats.update()
		time.Sleep(time.Millisecond * 13)
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

func (g *game) renderStats() {
	g.drawBuf.WriteString("-- STATS\n")
	g.drawBuf.WriteString(fmt.Sprintf("FPS: %.2f", g.stats.fps))
}

func (g *game) render() {
	g.drawBuf.Reset()
	fmt.Fprint(os.Stdout, "\033[2J\033[1;1H")

	g.renderLevel()
	g.renderStats()
	fmt.Fprint(os.Stdout, g.drawBuf.String())
}
