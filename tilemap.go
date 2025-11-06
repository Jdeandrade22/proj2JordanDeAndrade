package main

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
)

type TileMap struct {
	tiledMap   *tiled.Map
	tileImages map[uint32]*ebiten.Image
	width      int
	height     int
}

func NewTileMap(tmxData []byte, tilesetImages map[string]*ebiten.Image) (*TileMap, error) {
	tiledMap, err := tiled.LoadReader("map.tmx", bytes.NewReader(tmxData))
	if err != nil {
		return nil, err
	}

	tm := &TileMap{
		tiledMap:   tiledMap,
		tileImages: make(map[uint32]*ebiten.Image),
		width:      tiledMap.Width * tiledMap.TileWidth,
		height:     tiledMap.Height * tiledMap.TileHeight,
	}

	// Load tile images from tilesets
	for _, tileset := range tiledMap.Tilesets {
		if tileset.Image != nil {
			// Tileset with a single image
			tilesetImg := tilesetImages[tileset.Image.Source]
			if tilesetImg == nil {
				log.Printf("Warning: tileset image not found: %s", tileset.Image.Source)
				continue
			}

			columns := tileset.Columns
			if columns == 0 {
				columns = tileset.TileCount
			}

			for i := 0; i < tileset.TileCount; i++ {
				gid := tileset.FirstGID + uint32(i)
				x := int(i % columns)
				y := int(i / columns)

				sx := x * tileset.TileWidth
				sy := y * tileset.TileHeight

				rect := image.Rect(sx, sy, sx+tileset.TileWidth, sy+tileset.TileHeight)
				tm.tileImages[gid] = tilesetImg.SubImage(rect).(*ebiten.Image)
			}
		} else {
			// Collection of individual tiles
			for _, tile := range tileset.Tiles {
				if tile.Image != nil {
					gid := tileset.FirstGID + tile.ID
					tileImg := tilesetImages[tile.Image.Source]
					if tileImg != nil {
						tm.tileImages[gid] = tileImg
					}
				}
			}
		}
	}

	return tm, nil
}

func (tm *TileMap) Draw(screen *ebiten.Image, cameraX, cameraY float64) {
	for _, layer := range tm.tiledMap.Layers {
		if !layer.Visible {
			continue
		}

		for tileY := 0; tileY < tm.tiledMap.Height; tileY++ {
			for tileX := 0; tileX < tm.tiledMap.Width; tileX++ {
				tileIndex := tileY*tm.tiledMap.Width + tileX
				if tileIndex >= len(layer.Tiles) {
					continue
				}

				tile := layer.Tiles[tileIndex]
				if tile.IsNil() {
					continue
				}

				tileImage := tm.tileImages[tile.ID]
				if tileImage == nil {
					continue
				}

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(
					float64(tileX*tm.tiledMap.TileWidth)-cameraX,
					float64(tileY*tm.tiledMap.TileHeight)-cameraY,
				)

				// Apply layer opacity
				if layer.Opacity < 1.0 {
					op.ColorScale.ScaleAlpha(float32(layer.Opacity))
				}

				screen.DrawImage(tileImage, op)
			}
		}
	}
}

func (tm *TileMap) Width() int {
	return tm.width
}

func (tm *TileMap) Height() int {
	return tm.height
}

