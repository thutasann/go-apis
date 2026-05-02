package main

type position struct {
	x, y int
}

type player struct {
	pos   position
	level *level
	input *input

	reverse bool
}

func (p *player) update() {
	switch p.input.consumeFrameKey() {
	case 'a', 'A':
		p.move(-1, 0)
	case 'd', 'D':
		p.move(1, 0)
	case 'w', 'W':
		p.move(0, -1)
	case 's', 'S':
		p.move(0, 1)
	}
}

func (p *player) move(dx, dy int) {
	next := position{
		x: p.pos.x + dx,
		y: p.pos.y + dy,
	}

	if next.x <= 0 || next.x >= p.level.width-1 {
		return
	}

	if next.y <= 0 || next.y >= p.level.height-1 {
		return
	}

	p.pos = next
}
