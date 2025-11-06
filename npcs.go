package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type NPC struct {
	image     *ebiten.Image
	x, y      float64
	direction float64
	moveRange float64
	startX    float64
	startY    float64
	width     int
	height    int
	moveHorz  bool
}

func NewNPC(x, y float64, image *ebiten.Image, moveRange float64, moveHorizontal bool) *NPC {
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
	}
}

func (npc *NPC) Update() {
	// Simple back and forth movement
	if npc.moveHorz {
		npc.x += npc.direction
		// Reverse direction when hitting movement boundaries
		if npc.x >= npc.startX+npc.moveRange || npc.x <= npc.startX-npc.moveRange {
			npc.direction = -npc.direction
		}
	} else {
		npc.y += npc.direction
		if npc.y >= npc.startY+npc.moveRange || npc.y <= npc.startY-npc.moveRange {
			npc.direction = -npc.direction
		}
	}
}

func (npc *NPC) Draw(target *ebiten.Image, cameraX, cameraY float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(npc.x-cameraX, npc.y-cameraY)
	target.DrawImage(npc.image, op)
}
