package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	walkSprites [8]*AnimatedSprite //8 imgaes
	x, y        float64
	width       int
	height      int
	direction   int
	isMoving    bool
	speed       float64
}

func NewPlayer(x, y float64, walkSprites [8]*ebiten.Image) *Player {
	p := &Player{
		x:         x,
		y:         y,
		width:     64,
		height:    64,
		direction: 0,
		speed:     3.0,
	}

	for i := 0; i < 8; i++ {
		p.walkSprites[i] = NewAnimatedSprite(walkSprites[i], 64, 64)
	}

	return p
}

func (p *Player) Update(mapWidth, mapHeight int) {
	p.isMoving = false

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

	if moveX != 0 || moveY != 0 {
		p.isMoving = true

		//rnad movement (ai)
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

		// run cycle
		if moveX != 0 && moveY != 0 {
			moveX *= 0.707 //diag
			moveY *= 0.707
		}

		p.x += moveX * p.speed
		p.y += moveY * p.speed
	}

	//boundaries
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
		//initialize frame
		p.walkSprites[p.direction].currentFrame = 0
	}
}

func (p *Player) Draw(target *ebiten.Image, cameraX, cameraY float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.x-cameraX, p.y-cameraY)
	p.walkSprites[p.direction].Draw(target, op)
}

func (p *Player) GetBounds() (float64, float64, float64, float64) {
	hitboxPadding := 16.0 // hitbox increase
	return p.x + hitboxPadding,
		p.y + hitboxPadding,
		float64(p.width) - (hitboxPadding * 2),
		float64(p.height) - (hitboxPadding * 2)
}
