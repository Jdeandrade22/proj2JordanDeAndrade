package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Camerastruct
type Camera struct {
	ViewportWidth  int
	ViewportHeight int
	Follow         Follow
}

//positioning
type Follow struct {
	W int // X position
	H int // Y position
}

// Init creates a new camera with the given viewport dimensions
func Init(width, height int) *Camera {
	return &Camera{
		ViewportWidth:  width,
		ViewportHeight: height,
		Follow: Follow{
			W: 0,
			H: 0,
		},
	}
}

// Draw renders the world image to the screen, centered on the Follow position
func (c *Camera) Draw(world, screen *ebiten.Image) {
	cameraX := c.Follow.W - c.ViewportWidth/2
	cameraY := c.Follow.H - c.ViewportHeight/2

	// Clamp camera to world bounds
	worldWidth := world.Bounds().Dx()
	worldHeight := world.Bounds().Dy()

	if cameraX < 0 {
		cameraX = 0
	}
	if cameraY < 0 {
		cameraY = 0
	}
	if cameraX > worldWidth-c.ViewportWidth {
		cameraX = worldWidth - c.ViewportWidth
	}
	if cameraY > worldHeight-c.ViewportHeight {
		cameraY = worldHeight - c.ViewportHeight
	}

	//subimage
	sx := cameraX
	sy := cameraY
	sw := c.ViewportWidth
	sh := c.ViewportHeight

	// Ensure we don't go out of bounds
	if sx < 0 {
		sx = 0
	}
	if sy < 0 {
		sy = 0
	}
	if sx+sw > worldWidth {
		sw = worldWidth - sx
	}
	if sy+sh > worldHeight {
		sh = worldHeight - sy
	}

	// Draw the visible portion of the world to the screen
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(world.SubImage(image.Rect(sx, sy, sx+sw, sy+sh)).(*ebiten.Image), op)
}
