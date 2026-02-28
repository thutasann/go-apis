package main

func (g *Game) isBadCollision(
	p Point,
	snake []Point,
) bool {

	if p.x < 0 || p.y < 0 || p.x >= screenWidth/gridSize || p.y >= screenHeight/gridSize {
		return true
	}

	for _, sp := range snake {
		if sp == p {
			return true
		}
	}

	return false
}
