package main

import (
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	FRAME_COUNT     = 8  // Your sprites have 8 frames
	ANIMATION_SPEED = 80 // milliseconds per frame (fast smooth animation)
)

type AnimatedSprite struct {
	sheet        *ebiten.Image
	frameWidth   int
	frameHeight  int
	currentFrame int
	lastUpdate   time.Time
	row          int // which row in the sprite sheet (for different animations)
	frameCount   int // number of frames to cycle through
}

func NewAnimatedSprite(sheet *ebiten.Image, frameWidth, frameHeight int) *AnimatedSprite {
	return &AnimatedSprite{
		sheet:       sheet,
		frameWidth:  frameWidth,
		frameHeight: frameHeight,
		lastUpdate:  time.Now(),
		row:         0,           // default to first row (down animation)
		frameCount:  FRAME_COUNT, // default to 8 frames for player sprites
	}
}

func (a *AnimatedSprite) Update() {
	now := time.Now()
	if now.Sub(a.lastUpdate) > ANIMATION_SPEED*time.Millisecond {
		a.currentFrame = (a.currentFrame + 1) % a.frameCount
		a.lastUpdate = now
	}
}

func (a *AnimatedSprite) SetAnimationRow(row int) {
	a.row = row
}

func (a *AnimatedSprite) Draw(target *ebiten.Image, op *ebiten.DrawImageOptions) {
	// Calculate frame position in the sprite sheet
	// For multi-row sprites (like portal 3x2), calculate row and column
	sheetCols := a.sheet.Bounds().Dx() / a.frameWidth
	frameCol := a.currentFrame % sheetCols
	frameRow := a.currentFrame / sheetCols

	frameX := frameCol * a.frameWidth
	frameY := frameRow * a.frameHeight

	frame := a.sheet.SubImage(image.Rect(
		frameX, frameY,
		frameX+a.frameWidth,
		frameY+a.frameHeight,
	)).(*ebiten.Image)

	target.DrawImage(frame, op)
}

func (a *AnimatedSprite) Reset() {
	a.currentFrame = 0
	a.lastUpdate = time.Now()
}
