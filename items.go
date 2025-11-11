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
	animatedSprite *AnimatedSprite
	collected      bool
}

func NewItem(x, y float64, itemType ItemType, image *ebiten.Image) *Item {
	item := &Item{
		x:        x,
		y:        y,
		width:    64,
		height:   64,
		itemType: itemType,
		image:    image,
	}

	if itemType == ItemPortal {
		item.animatedSprite = NewAnimatedSprite(image, 32, 32)
		item.animatedSprite.frameCount = 6
	}

	return item
}

func (i *Item) Update() {
	if i.animatedSprite != nil {
		i.animatedSprite.Update()
	}
}

func (i *Item) Draw(screen *ebiten.Image, cameraX, cameraY float64) {
	i.DrawWithAlpha(screen, cameraX, cameraY, 1.0)
}

func (i *Item) DrawWithAlpha(screen *ebiten.Image, cameraX, cameraY, alpha float64) {
	if i.collected {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2.0, 2.0)
	op.GeoM.Translate(i.x-cameraX, i.y-cameraY)
	op.ColorScale.ScaleAlpha(float32(alpha))

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
