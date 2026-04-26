package main

import (
	"bytes"
	"fmt"
)

func main() {
	width := 20
	height := 20
	level := make([][]byte, height)

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			level[h] = make([]byte, width)
		}
	}

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			level[h][w] = WALL
		}
	}

	buf := new(bytes.Buffer)
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			if level[h][w] == WALL {
				buf.WriteString("H")
			}
		}
		buf.WriteString("\n")
	}

	fmt.Println(buf.String())
}
