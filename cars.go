package main

import (
	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Car struct {
	spriteSheet  *ebiten.Image
	frames       []*ebiten.Image // Individual frames extracted from sheet
	currentFrame int
	frameCount   int
	animTimer    int
	x, y         float64
	speedX       float64
	speedY       float64
	width        int
	height       int
	changeTimer  int
	maxSpeed     float64
}

func NewCar(x, y float64, spriteSheet *ebiten.Image, maxSpeed float64) *Car {
	// Extract frames from sprite sheet
	// Assuming sprite sheets are grids - Blue Limo is 980x980, Police is 700x700
	sheetWidth := spriteSheet.Bounds().Dx()

	// Try to figure out frame size - likely 7x7 grid for Blue Limo (140x140 per frame)
	// or 10x10 for Police (70x70 per frame)
	var frameSize int
	var gridSize int

	if sheetWidth == 980 {
		// Blue Limo: 980/7 = 140
		frameSize = 140
		gridSize = 7
	} else {
		// Police: 700/10 = 70
		frameSize = 70
		gridSize = 10
	}

	// Extract just the first row of frames for animation
	frames := make([]*ebiten.Image, gridSize)
	for i := 0; i < gridSize; i++ {
		x := i * frameSize
		rect := image.Rect(x, 0, x+frameSize, frameSize)
		frames[i] = spriteSheet.SubImage(rect).(*ebiten.Image)
	}

	car := &Car{
		spriteSheet:  spriteSheet,
		frames:       frames,
		frameCount:   gridSize,
		currentFrame: 0,
		x:            x,
		y:            y,
		width:        80, // Display size
		height:       80,
		maxSpeed:     maxSpeed,
		changeTimer:  0,
		animTimer:    0,
	}

	// Set initial random direction
	car.changeDirection()

	return car
}

func (c *Car) changeDirection() {
	// Random direction with speed variation
	angle := rand.Float64() * 2 * math.Pi            // Random angle (0 to 2*PI)
	speed := c.maxSpeed * (0.5 + rand.Float64()*0.5) // 50-100% of max speed

	c.speedX = speed * math.Cos(angle)
	c.speedY = speed * math.Sin(angle)

	// Reset timer (change direction every 2-5 seconds)
	c.changeTimer = 120 + rand.Intn(180) // 2-5 seconds at 60fps
}

func (c *Car) Update(mapWidth, mapHeight int) {
	// Move the car
	c.x += c.speedX
	c.y += c.speedY

	// Bounce off walls
	if c.x < 0 {
		c.x = 0
		c.speedX = -c.speedX
		c.changeTimer = 60 // Change direction soon
	}
	if c.y < 0 {
		c.y = 0
		c.speedY = -c.speedY
		c.changeTimer = 60
	}
	if c.x+float64(c.width) > float64(mapWidth) {
		c.x = float64(mapWidth - c.width)
		c.speedX = -c.speedX
		c.changeTimer = 60
	}
	if c.y+float64(c.height) > float64(mapHeight) {
		c.y = float64(mapHeight - c.height)
		c.speedY = -c.speedY
		c.changeTimer = 60
	}

	// Animate the car (cycle through frames)
	c.animTimer++
	if c.animTimer >= 8 { // Change frame every 8 ticks
		c.animTimer = 0
		c.currentFrame = (c.currentFrame + 1) % c.frameCount
	}

	// Random direction changes
	c.changeTimer--
	if c.changeTimer <= 0 {
		c.changeDirection()
	}
}

func (c *Car) Draw(target *ebiten.Image, cameraX, cameraY float64) {
	if c.currentFrame >= len(c.frames) {
		return
	}

	frame := c.frames[c.currentFrame]

	op := &ebiten.DrawImageOptions{}

	// Scale the car frame to fit our display size
	frameWidth := float64(frame.Bounds().Dx())
	frameHeight := float64(frame.Bounds().Dy())
	scaleX := float64(c.width) / frameWidth
	scaleY := float64(c.height) / frameHeight
	op.GeoM.Scale(scaleX, scaleY)

	// Rotate car based on movement direction (optional for better visuals)
	angle := math.Atan2(c.speedY, c.speedX)
	op.GeoM.Translate(-float64(c.width)/2, -float64(c.height)/2) // Center for rotation
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(float64(c.width)/2, float64(c.height)/2) // Translate back

	op.GeoM.Translate(c.x-cameraX, c.y-cameraY)
	target.DrawImage(frame, op)
}

func (c *Car) CheckCollision(px, py, pw, ph float64) bool {
	return px < c.x+float64(c.width) &&
		px+pw > c.x &&
		py < c.y+float64(c.height) &&
		py+ph > c.y
}
