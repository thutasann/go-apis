package main

import (
	"fmt"
	"os"
)

func (g *game) update() {
	g.level.set(g.player.pos, NOTHING)
	g.player.update()
	g.level.set(g.player.pos, PLAYER)
}

func (g *game) render() {
	g.drawBuf.Reset()
	fmt.Fprint(os.Stdout, "\033[H")

	g.renderLevel()
	g.renderStats()
	fmt.Fprint(os.Stdout, g.drawBuf.String())
}

func (g *game) renderLevel() {
	for h := range g.level.height {
		for w := range g.level.width {
			if g.level.data[h][w] == NOTHING {
				g.drawBuf.WriteString(" ")
			}
			if g.level.data[h][w] == WALL {
				g.drawBuf.WriteString("☐")
			}
			if g.level.data[h][w] == PLAYER {
				g.drawBuf.WriteString("త")
			}
		}
		g.drawBuf.WriteString("\n")
	}
}

func (g *game) renderStats() {
	g.drawBuf.WriteString("-- STATS\n")
	g.drawBuf.WriteString(fmt.Sprintf("FPS: %.2f\n", g.stats.fps))
	g.drawBuf.WriteString(fmt.Sprintf("KEYPRESS: %v\n", g.input.key()))
}
