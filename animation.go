package main

import (
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	FRAME_COUNT     = 8
	ANIMATION_SPEED = 80
)

type AnimatedSprite struct {
	sheet          *ebiten.Image
	frameWidth     int
	frameHeight    int
	currentFrame   int
	lastUpdate     time.Time
	row            int // which row in the sprite sheet (for different animations)
	frameCount     int // number of frames to cycle through
	animationSpeed int // milliseconds per frame (can be customized per sprite)
}

func NewAnimatedSprite(sheet *ebiten.Image, frameWidth, frameHeight int) *AnimatedSprite {
	return &AnimatedSprite{
		sheet:          sheet,
		frameWidth:     frameWidth,
		frameHeight:    frameHeight,
		lastUpdate:     time.Now(),
		row:            0,
		frameCount:     FRAME_COUNT,
		animationSpeed: ANIMATION_SPEED,
	}
}

func (a *AnimatedSprite) Update() {
	now := time.Now()
	if now.Sub(a.lastUpdate) > time.Duration(a.animationSpeed)*time.Millisecond {
		a.currentFrame = (a.currentFrame + 1) % a.frameCount
		a.lastUpdate = now
	}
}

func (a *AnimatedSprite) Draw(target *ebiten.Image, op *ebiten.DrawImageOptions) {
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
