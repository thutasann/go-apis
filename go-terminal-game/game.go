package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type game struct {
	isRunning bool
	level     *level
	stats     *stats
	player    *player
	input     *input
	cleanup   func() error

	drawBuf *bytes.Buffer
}

func newGame(width, height int) *game {
	restoreTerminal, err := configureTerminal()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to configure terminal input: %v\n", err)
	}

	var (
		lvl   = newLevel(width, height)
		input = newInput()
	)
	return &game{
		level:   lvl,
		drawBuf: new(bytes.Buffer),
		stats:   newStats(),
		input:   input,
		cleanup: restoreTerminal,

		player: &player{
			level: lvl,
			input: input,
			pos:   position{x: 2, y: 5},
		},
	}
}

func configureTerminal() (func() error, error) {
	stateCmd := exec.Command("stty", "-g")
	stateCmd.Stdin = os.Stdin
	state, err := stateCmd.Output()
	if err != nil {
		return nil, err
	}

	rawCmd := exec.Command("stty", "cbreak", "min", "1", "-echo")
	rawCmd.Stdin = os.Stdin
	if err := rawCmd.Run(); err != nil {
		return nil, err
	}

	savedState := strings.TrimSpace(string(state))
	return func() error {
		restoreCmd := exec.Command("stty", savedState)
		restoreCmd.Stdin = os.Stdin
		return restoreCmd.Run()
	}, nil
}

func (g *game) start() {
	if g.cleanup != nil {
		defer func() {
			if err := g.cleanup(); err != nil {
				fmt.Fprintf(os.Stderr, "failed to restore terminal input: %v\n", err)
			}
		}()
	}

	// Switch to the alternate screen buffer so redraws replace the same frame.
	fmt.Fprint(os.Stdout, "\033[?1049h\033[?25l")
	defer fmt.Fprint(os.Stdout, "\033[?25h\033[?1049l")

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
