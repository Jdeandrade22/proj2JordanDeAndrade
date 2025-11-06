package main

import (
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	FRAME_COUNT     = 4
	ANIMATION_SPEED = 200 // milliseconds per frame
)

type AnimatedSprite struct {
	sheet        *ebiten.Image
	frameWidth   int
	frameHeight  int
	currentFrame int
	lastUpdate   time.Time
	row          int // which row in the sprite sheet (for different animations)
}

func NewAnimatedSprite(sheet *ebiten.Image, frameWidth, frameHeight int) *AnimatedSprite {
	return &AnimatedSprite{
		sheet:       sheet,
		frameWidth:  frameWidth,
		frameHeight: frameHeight,
		lastUpdate:  time.Now(),
		row:         0, // default to first row (down animation)
	}
}

func (a *AnimatedSprite) Update() {
	now := time.Now()
	if now.Sub(a.lastUpdate) > ANIMATION_SPEED*time.Millisecond {
		a.currentFrame = (a.currentFrame + 1) % FRAME_COUNT
		a.lastUpdate = now
	}
}

func (a *AnimatedSprite) SetAnimationRow(row int) {
	a.row = row
}

func (a *AnimatedSprite) Draw(target *ebiten.Image, op *ebiten.DrawImageOptions) {
	frameX := a.currentFrame * a.frameWidth
	frameY := a.row * a.frameHeight

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
