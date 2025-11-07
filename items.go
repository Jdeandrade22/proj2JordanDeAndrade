package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type ItemType int

const (
	ItemGood ItemType = iota
	ItemBad
	ItemPortal
)

type Item struct {
	x, y      float64
	width     int
	height    int
	itemType  ItemType
	image     *ebiten.Image
	collected bool
}

func NewItem(x, y float64, itemType ItemType, image *ebiten.Image) *Item {
	return &Item{
		x:        x,
		y:        y,
		width:    64, // Doubled from 32 to 64
		height:   64, // Doubled from 32 to 64
		itemType: itemType,
		image:    image,
	}
}

func (i *Item) Draw(screen *ebiten.Image, cameraX, cameraY float64) {
	if i.collected {
		return
	}

	op := &ebiten.DrawImageOptions{}

	// Scale the item to be bigger (2x size)
	op.GeoM.Scale(2.0, 2.0)
	op.GeoM.Translate(i.x-cameraX, i.y-cameraY)

	screen.DrawImage(i.image, op)
}

func (i *Item) CheckCollision(px, py, pw, ph float64) bool {
	if i.collected {
		return false
	}

	return px < i.x+float64(i.width) &&
		px+pw > i.x &&
		py < i.y+float64(i.height) &&
		py+ph > i.y
}
