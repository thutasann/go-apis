package entity

import "github.com/hajimehoshi/ebiten/v2"

const (
	TagEnemy  = "enemy"
	TagPlayer = "player"
	TagFood   = "food"
)

type Entity interface {
	Update(world worldView) bool
	Draw(screen *ebiten.Image)
	Tag() string
}
