package entity

import (
	"image/color"

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
