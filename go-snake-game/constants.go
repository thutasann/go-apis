package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	dirUp           = Point{x: 0, y: -1}
	dirDown         = Point{x: 0, y: 1}
	dirLeft         = Point{x: -1, y: 0}
	dirRight        = Point{x: 1, y: 0}
	mplusFaceSource *text.GoTextFaceSource
)

const (
	gameSpeed    = time.Second / 6
	screenWidth  = 640
	screenHeight = 480
	gridSize     = 20
)
