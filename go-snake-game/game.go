package main

import (
	"errors"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/thuta/go-snake/common"
	"github.com/thuta/go-snake/entity"
	"github.com/thuta/go-snake/game"
	"github.com/thuta/go-snake/math"
)

type Game struct {
	world      *game.World
	lastUpdate time.Time
	gameOver   bool
}

func (g *Game) Update() error {
	if g.gameOver {
		return nil
	}

	playerRaw, ok := g.world.GetFirstEntity(entity.TagPlayer)
	if !ok {
		return errors.New("entity player was not found")
	}
	player := playerRaw.(*entity.Player)

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		player.SetDirection(math.DirUp)
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		player.SetDirection(math.DirDown)
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		player.SetDirection(math.DirLeft)
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		player.SetDirection(math.DirRight)
	}

	if time.Since(g.lastUpdate) < common.GameSpeed {
		return nil
	}
	g.lastUpdate = time.Now()

	for _, entity := range g.world.Entities() {
		if entity.Update(g.world) {
			g.gameOver = true
			return nil
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	for _, entity := range g.world.Entities() {
		entity.Draw(screen)
	}

	if g.gameOver {
		face := &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   48,
		}

		t := "Game Over!"
		w, h := text.Measure(
			t,
			face,
			face.Size,
		)

		op := &text.DrawOptions{}
		op.GeoM.Translate(
			common.ScreenWidth/2-w/2,
			common.ScreenHeight/2-h/2,
		)
		op.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, t, face, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.ScreenWidth, common.ScreenHeight
}
