package main

import (
	"bytes"
	"fmt"
)

func main() {
	width := 80
	height := 18
	level := make([][]byte, height)

	for h := range height {
		for range width {
			level[h] = make([]byte, width)
		}
	}

	for h := range height {
		for w := range width {
			if h == 0 {
				level[h][w] = WALL
			}

			if w == 0 {
				level[h][w] = WALL
			}

			if w == width-1 {
				level[h][w] = WALL
			}

			if h == height-1 {
				level[h][w] = WALL
			}
		}
	}

	buf := new(bytes.Buffer)
	for h := range height {
		for w := range width {
			if level[h][w] == NOTHING {
				buf.WriteString(" ")
			}
			if level[h][w] == WALL {
				buf.WriteString("H")
			}
		}
		buf.WriteString("\n")
	}

	fmt.Println(buf.String())
}
