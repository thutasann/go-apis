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
	player    *player
	input     *input

	drawBuf *bytes.Buffer
}

func newGame(width, height int) *game {
	var (
		lvl   = newLevel(width, height)
		input = &input{}
	)
	return &game{
		level:   lvl,
		drawBuf: new(bytes.Buffer),
		stats:   newStats(),
		input:   input,

		player: &player{
			level: lvl,
			input: input,
			pos:   position{x: 2, y: 5},
		},
	}
}

func (g *game) start() {
	g.isRunning = true
	g.loop()
}

func (g *game) loop() {
	for g.isRunning {
		g.input.update()
		g.update()
		g.render()
		g.stats.update()
		time.Sleep(time.Millisecond * 16) // limit FPS
	}
}

func (g *game) update() {
	g.level.set(g.player.pos, NOTHING)
	g.player.update()
	g.level.set(g.player.pos, PLAYER)
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
}

func (g *game) render() {
	g.drawBuf.Reset()
	fmt.Fprint(os.Stdout, "\033[2J\033[1;1H")

	g.renderLevel()
	g.renderStats()
	fmt.Fprint(os.Stdout, g.drawBuf.String())
}
