package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	sprite    *AnimatedSprite
	x, y      float64
	width     int
	height    int
	direction int // 0=down, 1=left, 2=right, 3=up
	isMoving  bool
	speed     float64
}

func NewPlayer(x, y float64, spriteSheet *ebiten.Image) *Player {
	return &Player{
		sprite:    NewAnimatedSprite(spriteSheet, 48, 48), // Adjust size based on your sprite
		x:         x,
		y:         y,
		width:     48,
		height:    48,
		direction: 0, // facing down initially
		speed:     3.0,
	}
}

func (p *Player) Update(mapWidth, mapHeight int) {
	p.isMoving = false

	// Handle movement and set animation row based on direction
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		p.x -= p.speed
		p.direction = 1
		p.isMoving = true
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		p.x += p.speed
		p.direction = 2
		p.isMoving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		p.y -= p.speed
		p.direction = 3
		p.isMoving = true
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		p.y += p.speed
		p.direction = 0
		p.isMoving = true
	}

	// Keep player within map bounds
	if p.x < 0 {
		p.x = 0
	}
	if p.y < 0 {
		p.y = 0
	}
	if p.x+float64(p.width) > float64(mapWidth) {
		p.x = float64(mapWidth - p.width)
	}
	if p.y+float64(p.height) > float64(mapHeight) {
		p.y = float64(mapHeight - p.height)
	}

	// Update animation
	if p.isMoving {
		p.sprite.SetAnimationRow(p.direction)
		p.sprite.Update()
	} else {
		// When not moving, show first frame of current direction
		p.sprite.SetAnimationRow(p.direction)
		p.sprite.currentFrame = 0 // Reset to first frame when idle
	}
}

func (p *Player) Draw(target *ebiten.Image, cameraX, cameraY float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.x-cameraX, p.y-cameraY)
	p.sprite.Draw(target, op)
}

func (p *Player) GetBounds() (float64, float64, float64, float64) {
	return p.x, p.y, float64(p.width), float64(p.height)
}
