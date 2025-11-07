package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	walkSprites   [8]*AnimatedSprite // All 8 walk animations
	attackSprites [8]*AnimatedSprite // All 8 attack animations
	x, y          float64
	width         int
	height        int
	direction     int // 0-7 for 8 directions
	isMoving      bool
	isAttacking   bool
	attackTimer   int
	speed         float64
}

func NewPlayer(x, y float64, walkSprites [8]*ebiten.Image, attackSprites [8]*ebiten.Image) *Player {
	p := &Player{
		x:         x,
		y:         y,
		width:     64,
		height:    64,
		direction: 0, // facing down initially
		speed:     3.0,
	}

	// Initialize walk sprites
	for i := 0; i < 8; i++ {
		p.walkSprites[i] = NewAnimatedSprite(walkSprites[i], 64, 64)
	}

	// Initialize attack sprites
	for i := 0; i < 8; i++ {
		p.attackSprites[i] = NewAnimatedSprite(attackSprites[i], 64, 64)
	}

	return p
}

func (p *Player) Update(mapWidth, mapHeight int) {
	// Handle attack timer
	if p.isAttacking {
		p.attackTimer--
		if p.attackTimer <= 0 {
			p.isAttacking = false
		}
	}

	// Check for attack input (Space or X key)
	if !p.isAttacking && (ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyX)) {
		p.isAttacking = true
		p.attackTimer = 30                            // Attack lasts 30 frames (~0.5 seconds)
		p.attackSprites[p.direction].currentFrame = 0 // Reset attack animation
	}

	// Don't move while attacking
	if p.isAttacking {
		p.attackSprites[p.direction].Update()
		return
	}

	p.isMoving = false

	// Handle 8-directional movement
	moveX := 0.0
	moveY := 0.0

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		moveX = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		moveX = 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		moveY = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		moveY = 1
	}

	// Calculate direction (0-7) based on movement
	if moveX != 0 || moveY != 0 {
		p.isMoving = true

		// Determine direction index (0-7)
		if moveY < 0 && moveX == 0 {
			p.direction = 1 // Up
		} else if moveY < 0 && moveX > 0 {
			p.direction = 2 // Up-Right
		} else if moveY == 0 && moveX > 0 {
			p.direction = 3 // Right
		} else if moveY > 0 && moveX > 0 {
			p.direction = 4 // Down-Right
		} else if moveY > 0 && moveX == 0 {
			p.direction = 5 // Down
		} else if moveY > 0 && moveX < 0 {
			p.direction = 6 // Down-Left
		} else if moveY == 0 && moveX < 0 {
			p.direction = 7 // Left
		} else if moveY < 0 && moveX < 0 {
			p.direction = 0 // Up-Left
		}

		// Normalize diagonal movement
		if moveX != 0 && moveY != 0 {
			moveX *= 0.707 // sqrt(2)/2 for diagonal
			moveY *= 0.707
		}

		p.x += moveX * p.speed
		p.y += moveY * p.speed
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
		p.walkSprites[p.direction].Update()
	} else {
		// When not moving, show first frame of current direction
		p.walkSprites[p.direction].currentFrame = 0
	}
}

func (p *Player) Draw(target *ebiten.Image, cameraX, cameraY float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.x-cameraX, p.y-cameraY)

	if p.isAttacking {
		p.attackSprites[p.direction].Draw(target, op)
	} else {
		p.walkSprites[p.direction].Draw(target, op)
	}
}

func (p *Player) GetBounds() (float64, float64, float64, float64) {
	return p.x, p.y, float64(p.width), float64(p.height)
}
