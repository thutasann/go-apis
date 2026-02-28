package main

func (g *Game) updateSnake(snake *[]Point, direction Point) {
	head := (*snake)[0]

	newHead := Point{
		x: head.x + direction.x,
		y: head.y + direction.y,
	}

	if g.isBadCollision(newHead, *snake) {
		g.gameOver = true
		return
	}

	if newHead == g.food {
		*snake = append(
			[]Point{newHead},
			*snake...,
		)
		g.spwanFood()
	} else {
		*snake = append(
			[]Point{newHead},
			(*snake)[:len(*snake)-1]...,
		)
	}
}
