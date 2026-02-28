package main

func (g *Game) updateSnake(snake *[]Point, direction Point) {
	head := (*snake)[0]

	newHead := Point{
		x: head.x + direction.x,
		y: head.y + direction.y,
	}

	*snake = append(
		[]Point{newHead},
		(*snake)[:len(*snake)-1]...,
	)
}
