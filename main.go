package main

import (
	"bytes"
	"embed"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

// Embed all assets
//
//go:embed assets
var assetsFS embed.FS

type GameState int

const (
	StatePlaying GameState = iota
	StateGameOver
	StateLevelComplete
	StateGameWon
)

type Game struct {
	player         *Player
	npcs           []*NPC
	items          []*Item
	portal         *Item
	tileMap        *TileMap
	cameraX        float64
	cameraY        float64
	state          GameState
	currentLevel   int
	itemsCollected int
	portalUnlocked bool

	// Asset images
	goldfishImg *ebiten.Image
	wormImg     *ebiten.Image
	badItemImg  *ebiten.Image
	portalImg   *ebiten.Image
	dogNPCImg   *ebiten.Image
	npc2Img     *ebiten.Image
}

func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())

	g := &Game{
		state:        StatePlaying,
		currentLevel: 1,
	}

	// Load all images
	g.loadAssets()

	// Load level 1
	g.loadLevel(1)

	return g
}

func (g *Game) loadAssets() {
	g.goldfishImg = g.loadImageFromFS("assets/items/Goldfish.png")
	g.wormImg = g.loadImageFromFS("assets/items/Worm.png")
	g.badItemImg = g.loadImageFromFS("assets/items/Rusty Can.png")
	g.portalImg = g.loadImageFromFS("assets/sprites/Dimensional_Portal.png")
	g.dogNPCImg = g.loadImageFromFS("assets/npc/drone dog.png")
	g.npc2Img = g.loadImageFromFS("assets/npc/mry_dgh_angry.png")
}

func (g *Game) loadImageFromFS(path string) *ebiten.Image {
	data, err := assetsFS.ReadFile(path)
	if err != nil {
		log.Printf("Warning: Failed to load %s: %v", path, err)
		img := ebiten.NewImage(64, 64)
		img.Fill(color.RGBA{255, 0, 255, 255}) // Magenta placeholder
		return img
	}

	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))
	if err != nil {
		log.Printf("Warning: Failed to decode %s: %v", path, err)
		img := ebiten.NewImage(64, 64)
		img.Fill(color.RGBA{255, 0, 255, 255})
		return img
	}

	return img
}

func (g *Game) loadLevel(level int) {
	var tmxData []byte
	var tilesetImages map[string]*ebiten.Image
	var err error

	if level == 1 {
		tmxData, err = assetsFS.ReadFile("assets/background/level1.tmx")
		if err != nil {
			log.Fatal("Failed to load level1.tmx:", err)
		}

		// Load level 1 tileset
		img := g.loadImageFromFS("assets/background/orig_big copy.png")
		tilesetImages = map[string]*ebiten.Image{
			"orig_big copy.png": img,
		}

		// Create player for level 1
		playerImg := g.loadImageFromFS("assets/sprites/walk.png")
		g.player = NewPlayer(100, 100, playerImg)

		// No NPCs on level 1
		g.npcs = []*NPC{}

	} else {
		tmxData, err = assetsFS.ReadFile("assets/background/level2.tmx")
		if err != nil {
			log.Fatal("Failed to load level2.tmx:", err)
		}

		// Load level 2 tilesets
		tilesetImages = map[string]*ebiten.Image{
			"clay_tile_64_01.png":      g.loadImageFromFS("assets/background/clay_tile_64_01.png"),
			"grass_01_tile_64_01.png":  g.loadImageFromFS("assets/background/grass_01_tile_64_01.png"),
			"ice_tile_64_01.png":       g.loadImageFromFS("assets/background/ice_tile_64_01.png"),
			"paving_01_tile_64_02.png": g.loadImageFromFS("assets/background/paving_01_tile_64_02.png"),
			"paving_02_tile_64_01.png": g.loadImageFromFS("assets/background/paving_02_tile_64_01.png"),
			"sand_01_tile_64_01.png":   g.loadImageFromFS("assets/background/sand_01_tile_64_01.png"),
		}

		// Reset player position for level 2
		g.player.x = 100
		g.player.y = 100

		// Add NPCs to level 2
		g.npcs = []*NPC{
			NewNPC(400, 300, g.dogNPCImg, 150, true), // Horizontal movement
			NewNPC(800, 200, g.npc2Img, 100, false),  // Vertical movement
		}
	}

	// Load tilemap
	g.tileMap, err = NewTileMap(tmxData, tilesetImages)
	if err != nil {
		log.Fatal("Failed to load tilemap:", err)
	}

	// Spawn items
	g.spawnItems()

	// Reset camera
	g.cameraX = 0
	g.cameraY = 0
}

func (g *Game) spawnItems() {
	g.items = []*Item{}
	g.portalUnlocked = false

	mapWidth := g.tileMap.Width()
	mapHeight := g.tileMap.Height()

	// Spawn 15+ good items (mix of goldfish and worms)
	for i := 0; i < 17; i++ {
		x := float64(rand.Intn(mapWidth-100) + 50)
		y := float64(rand.Intn(mapHeight-100) + 50)

		var img *ebiten.Image
		if i%2 == 0 {
			img = g.goldfishImg
		} else {
			img = g.wormImg
		}

		g.items = append(g.items, NewItem(x, y, ItemGood, img))
	}

	// Spawn 5 bad items (rusty cans - dangerous for cats!)
	for i := 0; i < 5; i++ {
		x := float64(rand.Intn(mapWidth-100) + 50)
		y := float64(rand.Intn(mapHeight-100) + 50)

		g.items = append(g.items, NewItem(x, y, ItemBad, g.badItemImg))
	}

	// Spawn portal near the end of the map
	portalX := float64(mapWidth - 150)
	portalY := float64(mapHeight - 150)
	g.portal = NewItem(portalX, portalY, ItemPortal, g.portalImg)
}

func (g *Game) Update() error {
	if g.state == StatePlaying {
		// Update player
		g.player.Update(g.tileMap.Width(), g.tileMap.Height())

		// Update NPCs (only on level 2)
		for _, npc := range g.npcs {
			npc.Update()
		}

		// Update camera to follow player
		g.cameraX = g.player.x - float64(screenWidth)/2 + float64(g.player.width)/2
		g.cameraY = g.player.y - float64(screenHeight)/2 + float64(g.player.height)/2

		// Clamp camera to map bounds
		if g.cameraX < 0 {
			g.cameraX = 0
		}
		if g.cameraY < 0 {
			g.cameraY = 0
		}
		if g.cameraX > float64(g.tileMap.Width()-screenWidth) {
			g.cameraX = float64(g.tileMap.Width() - screenWidth)
		}
		if g.cameraY > float64(g.tileMap.Height()-screenHeight) {
			g.cameraY = float64(g.tileMap.Height() - screenHeight)
		}

		// Check item collisions
		px, py, pw, ph := g.player.GetBounds()
		for _, item := range g.items {
			if item.CheckCollision(px, py, pw, ph) {
				item.collected = true

				if item.itemType == ItemGood {
					g.itemsCollected++

					// Unlock portal at 9 items
					if g.itemsCollected >= 9 {
						g.portalUnlocked = true
					}
				} else if item.itemType == ItemBad {
					// Game over!
					g.state = StateGameOver
				}
			}
		}

		// Check portal collision if unlocked
		if g.portalUnlocked && g.portal.CheckCollision(px, py, pw, ph) {
			if g.currentLevel == 1 {
				// Move to level 2
				g.currentLevel = 2
				g.itemsCollected = 0
				g.loadLevel(2)
				g.state = StatePlaying
			} else {
				// Game won!
				g.state = StateGameWon
			}
		}

	} else if g.state == StateGameOver {
		// Press R to restart
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.state = StatePlaying
			g.currentLevel = 1
			g.itemsCollected = 0
			g.loadLevel(1)
		}
	} else if g.state == StateGameWon {
		// Press R to play again
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.state = StatePlaying
			g.currentLevel = 1
			g.itemsCollected = 0
			g.loadLevel(1)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 50, 50, 255})

	if g.state == StatePlaying {
		// Draw tilemap
		g.tileMap.Draw(screen, g.cameraX, g.cameraY)

		// Draw items
		for _, item := range g.items {
			item.Draw(screen, g.cameraX, g.cameraY)
		}

		// Draw portal (only if unlocked or always visible)
		if g.portalUnlocked {
			g.portal.Draw(screen, g.cameraX, g.cameraY)
		} else {
			// Draw grayed out portal
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(g.portal.x-g.cameraX, g.portal.y-g.cameraY)
			op.ColorScale.ScaleAlpha(0.3)
			screen.DrawImage(g.portal.image, op)
		}

		// Draw NPCs
		for _, npc := range g.npcs {
			npc.Draw(screen, g.cameraX, g.cameraY)
		}

		// Draw player
		g.player.Draw(screen, g.cameraX, g.cameraY)

		// Draw UI
		g.drawUI(screen)

	} else if g.state == StateGameOver {
		text.Draw(screen, "GAME OVER!", basicfont.Face7x13, screenWidth/2-50, screenHeight/2, color.White)
		text.Draw(screen, "You touched a bad item!", basicfont.Face7x13, screenWidth/2-90, screenHeight/2+20, color.White)
		text.Draw(screen, "Press R to restart", basicfont.Face7x13, screenWidth/2-80, screenHeight/2+40, color.White)

	} else if g.state == StateGameWon {
		text.Draw(screen, "CONGRATULATIONS!", basicfont.Face7x13, screenWidth/2-70, screenHeight/2, color.White)
		text.Draw(screen, "You completed both levels!", basicfont.Face7x13, screenWidth/2-100, screenHeight/2+20, color.White)
		text.Draw(screen, fmt.Sprintf("Total items collected: %d", g.itemsCollected), basicfont.Face7x13, screenWidth/2-100, screenHeight/2+40, color.White)
		text.Draw(screen, "Press R to play again", basicfont.Face7x13, screenWidth/2-90, screenHeight/2+60, color.White)
	}
}

func (g *Game) drawUI(screen *ebiten.Image) {
	// Draw semi-transparent background for UI
	uiRect := ebiten.NewImage(screenWidth, 40)
	uiRect.Fill(color.RGBA{0, 0, 0, 180})
	screen.DrawImage(uiRect, &ebiten.DrawImageOptions{})

	// Draw level and collection info
	levelText := fmt.Sprintf("Level: %d", g.currentLevel)
	text.Draw(screen, levelText, basicfont.Face7x13, 10, 20, color.White)

	collectionText := fmt.Sprintf("Fish Collected: %d", g.itemsCollected)
	text.Draw(screen, collectionText, basicfont.Face7x13, 10, 35, color.White)

	// Portal status
	portalText := "Portal: Locked (Need 9 fish)"
	if g.portalUnlocked {
		portalText = "Portal: UNLOCKED! Go to portal!"
	}
	text.Draw(screen, portalText, basicfont.Face7x13, screenWidth-250, 20, color.RGBA{255, 215, 0, 255})

	// Controls
	controlsText := "Arrow Keys/WASD: Move"
	text.Draw(screen, controlsText, basicfont.Face7x13, screenWidth-200, 35, color.RGBA{200, 200, 200, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Cat Fish Quest - Project 2 - Jordan DeAndrade")
	ebiten.SetWindowResizable(false)

	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
