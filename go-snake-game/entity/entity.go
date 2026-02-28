package entity

import "github.com/hajimehoshi/ebiten/v2"

type Entity interface {
	Update(world worldView) bool
	Draw(screen *ebiten.Image)
	Tag() string
}
