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
	x, y           float64
	width          int
	height         int
	itemType       ItemType
	image          *ebiten.Image
	animatedSprite *AnimatedSprite // For animated items like portal
	collected      bool
}

func NewItem(x, y float64, itemType ItemType, image *ebiten.Image) *Item {
	item := &Item{
		x:        x,
		y:        y,
		width:    64, // Doubled from 32 to 64
		height:   64, // Doubled from 32 to 64
		itemType: itemType,
		image:    image,
	}

	// If it's a portal, make it animated (3x2 grid = 6 frames at 32x32 each)
	if itemType == ItemPortal {
		item.animatedSprite = NewAnimatedSprite(image, 32, 32) // 96x64 sprite sheet / 3x2 = 32x32 per frame
		item.animatedSprite.frameCount = 6                     // Portal has 6 frames (3 columns x 2 rows)
	}

	return item
}

func (i *Item) Update() {
	// Update animation for animated items (like portal)
	if i.animatedSprite != nil {
		i.animatedSprite.Update()
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

	// Use animated sprite if available (for portal)
	if i.animatedSprite != nil {
		i.animatedSprite.Draw(screen, op)
	} else {
		screen.DrawImage(i.image, op)
	}
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
