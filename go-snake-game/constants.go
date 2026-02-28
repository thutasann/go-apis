package main

import "time"

var (
	dirUp    = Point{x: 0, y: -1}
	dirDown  = Point{x: 0, y: 1}
	dirLeft  = Point{x: -1, y: 0}
	dirRight = Point{x: 1, y: 0}
)

const (
	gameSpeed    = time.Second / 6
	screenWidth  = 640
	screenHeight = 480
	gridSize     = 20
)
