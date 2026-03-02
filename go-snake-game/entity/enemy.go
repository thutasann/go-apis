package entity

import (
	"image/color"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/thuta/go-snake/common"
	"github.com/thuta/go-snake/math"
)

var _ Entity = (*Enemy)(nil)

type Enemy struct {
	body      []math.Point
	direction math.Point
}

func NewEnemy(start, dir math.Point) *Enemy {
	return &Enemy{
		body:      []math.Point{start},
		direction: dir,
	}
}

func (e *Enemy) Update(worldView worldView) bool {
	if len(e.body) == 0 {
		return true
	}

	head := e.body[0]

	possibleDirections := []math.Point{
		math.DirRight,
		math.DirLeft,
		math.DirDown,
		math.DirUp,
	}

	bestDirection := e.direction
	bestScore := 1

	for _, dir := range possibleDirections {
		newHead := math.Point{X: head.X + dir.X, Y: head.Y + dir.Y}
		if e.isCollision(newHead, worldView) {
			continue
		}

		score := e.evaluateMove(newHead, worldView)
		if score > bestScore {
			bestScore = score
			bestDirection = dir
		}
	}

	e.direction = bestDirection
	newHead := math.Point{X: head.X + e.direction.X, Y: head.Y + e.direction.Y}
	if e.isCollision(newHead, worldView) {
		return true
	}

	grow := false
	for _, entity := range worldView.GetEntities(TagFood) {
		food := entity.(*Food)
		if newHead.Equals(food.position) {
			grow = true
			food.Respawn()
			break
		}
	}

	if grow {
		e.body = append(
			[]math.Point{newHead},
			e.body...,
		)
	} else {
		e.body = append(
			[]math.Point{newHead},
			e.body[:len(e.body)-1]...,
		)
	}

	return false
}

func (e *Enemy) evaluateMove(newHead math.Point, worldView worldView) int {
	score := 0

	bestFood, distanceToFood, hasFood := e.findBestFood(newHead, worldView)

	if hasFood {
		score += 2000 - distanceToFood*20
		score += e.evaluatePathToFood(newHead, bestFood, worldView)
	}

	playerRaw, ok := worldView.GetFirstEntity(TagPlayer)
	if ok {
		player := playerRaw.(*Player)
		for _, p := range player.body {
			distanceToPlayer := newHead.Distance(p)
			if distanceToPlayer < 5 {
				score -= (5 - distanceToPlayer) * 150
			}
		}
	}

	nearByFoodCount := 0
	for _, entity := range worldView.GetEntities(TagFood) {
		food := entity.(*Food)
		if newHead.Distance(food.position) < 10 {
			nearByFoodCount++
		}
	}

	score += nearByFoodCount * 50

	return score
}

func (e *Enemy) evaluatePathToFood(
	newHead, targetFood math.Point, worldView worldView,
) int {
	if newHead.X == targetFood.X || newHead.Y == targetFood.Y {
		if e.isClearPath(newHead, targetFood, worldView) {
			return 500
		}
		return -300
	}
	if e.body[0].X != targetFood.X && newHead.X == targetFood.X {
		return 200
	}

	if e.body[0].Y != targetFood.Y && newHead.Y == targetFood.Y {
		return 200
	}
	return 0
}

func (e *Enemy) findBestFood(
	position math.Point, worldView worldView,
) (math.Point, int, bool) {
	bestFoodPos := math.Point{}
	bestScore := -1
	bestDistance := 0

	entities := worldView.GetEntities(TagFood)
	if len(entities) == 0 {
		return bestFoodPos, bestDistance, false
	}

	for _, entity := range entities {
		food := entity.(*Food)
		distance := position.Distance(food.position)
		score := 1000 - distance*10

		if e.isClearPath(position, food.position, worldView) {
			score += 500
		}

		if score > bestScore || (score == bestScore && distance < bestDistance) {
			bestScore = score
			bestDistance = distance
			bestFoodPos = food.position
		}
	}
	return bestFoodPos, bestDistance, true
}

func (e *Enemy) isClearPath(from, to math.Point, worldView worldView) bool {
	if from.X == to.X {
		start, end := from.Y, to.Y
		if start > end {
			start, end = end, start
		}

		for y := start; y <= end; y++ {
			if e.isCollision(math.Point{X: from.X, Y: y}, worldView) {
				return false
			}
		}
	} else if from.Y == to.Y {
		start, end := from.X, to.X
		if start > end {
			start, end = end, start
		}

		for x := start; x <= end; x++ {
			if e.isCollision(math.Point{X: x, Y: from.Y}, worldView) {
				return false
			}
		}
	}
	return true
}

func (e *Enemy) isCollision(p math.Point, worldView worldView) bool {
	if p.IsBadCollision(e.body) {
		return true
	}

	playerRaw, ok := worldView.GetFirstEntity(TagPlayer)
	if ok {
		player := playerRaw.(*Player)
		if slices.Contains(player.body, p) {
			return true
		}
	}
	return false
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	for _, pt := range e.body {
		vector.DrawFilledRect(
			screen,
			float32(pt.X*common.GridSize),
			float32(pt.Y*common.GridSize),
			common.GridSize,
			common.GridSize,
			color.White,
			true,
		)
	}
}

func (e Enemy) Tag() string {
	return TagEnemy
}
