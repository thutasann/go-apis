package main

type position struct {
	x, y int
}

type player struct {
	pos   position
	level *level

	reverse bool
}

func (p *player) update() {
	if p.reverse {
		p.pos.x -= 1
		if p.pos.x == 2 {
			p.pos.x += 1
			p.reverse = false
		}
		return
	}

	p.pos.x += 1
	if p.pos.x == p.level.width-2 {
		p.pos.x -= 1
		p.reverse = true
	}
}
