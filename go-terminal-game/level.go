package main

type level struct {
	width, height int
	data          [][]int
}

func newLevel(width, height int) *level {
	data := make([][]int, height)

	for h := range height {
		for range width {
			data[h] = make([]int, width)
		}
	}

	for h := range height {
		for w := range width {
			if h == 0 {
				data[h][w] = WALL
			}

			if w == 0 {
				data[h][w] = WALL
			}

			if w == width-1 {
				data[h][w] = WALL
			}

			if h == height-1 {
				data[h][w] = WALL
			}
		}
	}

	return &level{
		width:  width,
		height: height,
		data:   data,
	}
}

func (l *level) set(pos position, v int) {
	l.data[pos.y][pos.x] = v
}
