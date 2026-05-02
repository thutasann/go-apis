package main

type position struct {
	x, y int
}

type player struct {
	pos   position
	level *level
}

func (p *player) update() {
	p.level.data[p.pos.y][p.pos.x] = NOTHING
	p.pos.x += 1
}
