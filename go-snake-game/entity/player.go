package entity

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/thuta/go-snake/common"
	"github.com/thuta/go-snake/math"
)

var _ Entity = (*Player)(nil)

type Player struct {
	body      []math.Point
	direction math.Point
}

func NewPlayer(start, dir math.Point) *Player {
	return &Player{
		body:      []math.Point{start},
		direction: dir,
	}
}

func (p *Player) Update(worldview worldView) bool {
	newHead := p.body[0].Add(p.direction)

	if newHead.IsBadCollision(p.body) {
		return true
	}

	grow := false
	for _, entity := range worldview.GetEntities("food") {
		food := entity.(*Food)
		if newHead.Equals(food.position) {
			grow = true
			food.Respawn()
			break
		}
	}

	if grow {
		p.body = append(
			[]math.Point{newHead},
			p.body...,
		)
	} else {
		p.body = append(
			[]math.Point{newHead},
			p.body[:len(p.body)-1]...,
		)
	}

	return false
}

func (p *Player) Draw(screen *ebiten.Image) {
	for _, pt := range p.body {
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

func (p *Player) SetDirection(dir math.Point) {
	p.direction = dir
}

func (p Player) Tag() string {
	return "player"
}
