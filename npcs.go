package main

import (
	image2 "image"

	"github.com/hajimehoshi/ebiten/v2"
)

type NPC struct {
	image        *ebiten.Image
	frames       []*ebiten.Image // Pre-extracted frames for animations
	currentFrame int
	frameDelay   int
	frameCounter int
	x, y         float64
	direction    float64
	moveRange    float64
	startX       float64
	startY       float64
	width        int
	height       int
	moveHorz     bool
	scale        float64 // Scale factor for drawing
}

// NewStaticNPC creates a non-animated NPC (like the portrait)
func NewStaticNPC(x, y float64, image *ebiten.Image, moveRange float64, moveHorizontal bool) *NPC {
	return &NPC{
		image:     image,
		x:         x,
		y:         y,
		direction: 1,
		moveRange: moveRange,
		startX:    x,
		startY:    y,
		width:     64,
		height:    64,
		moveHorz:  moveHorizontal,
		scale:     2.0, // Default scale for static NPCs
	}
}

// NewAnimatedNPC creates an NPC with sprite sheet animation
func NewAnimatedNPC(x, y float64, image *ebiten.Image, frameWidth, frameHeight, cols, rows int, moveRange float64, moveHorizontal bool) *NPC {
	npc := &NPC{
		image:     image,
		x:         x,
		y:         y,
		direction: 1,
		moveRange: moveRange,
		startX:    x,
		startY:    y,
		width:     64,
		height:    64,
		moveHorz:  moveHorizontal,
		scale:     3.0, // Bigger scale for animated NPCs
	}

	// Extract frames from the second row (walking left animation)
	npc.frames = make([]*ebiten.Image, cols)
	for i := 0; i < cols; i++ {
		x1 := i * frameWidth
		x2 := x1 + frameWidth
		y1 := frameHeight
		y2 := frameHeight * 2
		npc.frames[i] = image.SubImage(image2.Rect(x1, y1, x2, y2)).(*ebiten.Image)
	}

	npc.currentFrame = 0
	npc.frameDelay = 10 // Frames between animation changes
	npc.frameCounter = 0

	return npc
}

// Update handles NPC movement and animation
func (npc *NPC) Update() {
	// Movement logic
	if npc.moveHorz {
		npc.x += npc.direction
		if npc.x >= npc.startX+npc.moveRange || npc.x <= npc.startX-npc.moveRange {
			npc.direction = -npc.direction
		}
	} else {
		npc.y += npc.direction
		if npc.y >= npc.startY+npc.moveRange || npc.y <= npc.startY-npc.moveRange {
			npc.direction = -npc.direction
		}
	}

	// Animation logic (if frames exist)
	if len(npc.frames) > 0 {
		npc.frameCounter++
		if npc.frameCounter >= npc.frameDelay {
			npc.frameCounter = 0
			npc.currentFrame = (npc.currentFrame + 1) % len(npc.frames)
		}
	}
}

// Draw renders the NPC to the screen
func (npc *NPC) Draw(target *ebiten.Image, cameraX, cameraY float64) {
	op := &ebiten.DrawImageOptions{}

	// Apply scale
	op.GeoM.Scale(npc.scale, npc.scale)
	op.GeoM.Translate(npc.x-cameraX, npc.y-cameraY)

	// Draw animated frames or static image
	if len(npc.frames) > 0 {
		target.DrawImage(npc.frames[npc.currentFrame], op)
	} else {
		target.DrawImage(npc.image, op)
	}
}
