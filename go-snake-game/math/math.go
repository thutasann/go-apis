package math

import (
	"math/rand"
	"slices"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/thuta/go-snake/common"
)

var (
	DirUp           = Point{X: 0, Y: -1}
	DirDown         = Point{X: 0, Y: 1}
	DirLeft         = Point{X: -1, Y: 0}
	DirRight        = Point{X: 1, Y: 0}
	MplusFaceSource *text.GoTextFaceSource
)

type Point struct {
	X, Y int
}

func (p Point) Equals(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

func (p Point) Add(other Point) Point {
	return Point{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

func (p Point) IsBadCollision(
	points []Point,
) bool {
	if p.X < 0 ||
		p.Y < 0 ||
		p.X >= common.ScreenWidth/common.GridSize ||
		p.Y >= common.ScreenHeight/common.GridSize {
		return true
	}

	return slices.Contains(points, p)
}

func RandomPosition() Point {
	return Point{
		X: rand.Intn(common.ScreenWidth / common.GridSize),
		Y: rand.Intn(common.ScreenHeight / common.GridSize),
	}
}
