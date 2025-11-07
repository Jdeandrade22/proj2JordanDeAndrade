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
	cars           []*Car // New: Cars for level 3!
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
	goldfishImg       *ebiten.Image
	rainbowTroutImg   *ebiten.Image
	angelfishImg      *ebiten.Image
	bassImg           *ebiten.Image
	catfishImg        *ebiten.Image
	wormImg           *ebiten.Image
	badItemImg        *ebiten.Image
	portalImg         *ebiten.Image
	femalePortraitImg *ebiten.Image // Female portrait
	femaleWalkImg     *ebiten.Image // Female walk and idle sprite sheet
	blueCarImg        *ebiten.Image // Blue Limo
	policeCarImg      *ebiten.Image // Police Car
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
	// Load all fish items
	g.goldfishImg = g.loadImageFromFS("assets/items/Goldfish.png")
	g.rainbowTroutImg = g.loadImageFromFS("assets/items/Rainbow Trout.png")
	g.angelfishImg = g.loadImageFromFS("assets/items/Angelfish.png")
	g.bassImg = g.loadImageFromFS("assets/items/Bass.png")
	g.catfishImg = g.loadImageFromFS("assets/items/Catfish.png")
	g.wormImg = g.loadImageFromFS("assets/items/Worm.png")

	// Load bad items and other assets
	g.badItemImg = g.loadImageFromFS("assets/items/Rusty Can.png")
	g.portalImg = g.loadImageFromFS("assets/items/Dimensional_Portal.png")

	// Load NPC sprites
	g.femalePortraitImg = g.loadImageFromFS("assets/npc/portrait female.png")
	g.femaleWalkImg = g.loadImageFromFS("assets/npc/walk and idle.png")

	// Load car sprite sheets
	g.blueCarImg = g.loadImageFromFS("assets/npc/Blue_LIMO_CLEAN_All_000-sheet.png")
	g.policeCarImg = g.loadImageFromFS("assets/npc/POLICE_CLEAN_ALLD0000-sheet.png")
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

		// Create player for level 1 - Cat character with all 8 directions
		var walkSprites [8]*ebiten.Image

		// Load all 8 walk sprites
		for i := 0; i < 8; i++ {
			walkSprites[i] = g.loadImageFromFS(fmt.Sprintf("assets/sprites/walk_%d.png", i+1))
		}

		g.player = NewPlayer(100, 100, walkSprites)

		// No NPCs on level 1
		g.npcs = []*NPC{}
		g.cars = []*Car{} // No cars either

	} else if level == 2 {
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
			NewAnimatedNPC(400, 300, g.femaleWalkImg, 24, 24, 8, 3, 150, true),  // Female walking - horizontal
			NewAnimatedNPC(800, 200, g.femaleWalkImg, 24, 24, 8, 3, 100, false), // Female walking - vertical
			NewStaticNPC(600, 500, g.femalePortraitImg, 80, false),              // Female portrait - vertical
			NewStaticNPC(300, 600, g.femalePortraitImg, 120, true),              // Female portrait - horizontal
		}
		g.cars = []*Car{} // No cars on level 2

	} else if level == 3 {
		// LEVEL 3 - FINAL BOSS LEVEL!
		tmxData, err = assetsFS.ReadFile("assets/background/level3.tmx")
		if err != nil {
			log.Fatal("Failed to load level3.tmx:", err)
		}

		// Load level 3 tileset
		img := g.loadImageFromFS("assets/background/orig_big1.png")
		tilesetImages = map[string]*ebiten.Image{
			"orig_big1.png": img,
		}

		// Reset player position for level 3
		g.player.x = 100
		g.player.y = 100

		// Add MANY NPCs to level 3 - FINAL CHALLENGE!
		g.npcs = []*NPC{
			NewAnimatedNPC(300, 250, g.femaleWalkImg, 24, 24, 8, 3, 200, true),  // Female walking 1 - horizontal
			NewAnimatedNPC(900, 300, g.femaleWalkImg, 24, 24, 8, 3, 180, true),  // Female walking 2 - horizontal
			NewAnimatedNPC(600, 400, g.femaleWalkImg, 24, 24, 8, 3, 150, false), // Female walking 3 - vertical
			NewAnimatedNPC(450, 600, g.femaleWalkImg, 24, 24, 8, 3, 120, false), // Female walking 4 - vertical
			NewStaticNPC(750, 150, g.femalePortraitImg, 100, true),              // Female portrait 1 - horizontal
			NewStaticNPC(200, 500, g.femalePortraitImg, 130, false),             // Female portrait 2 - vertical
		}

		// Add CARS that move randomly - AVOID THEM!
		g.cars = []*Car{
			NewCar(400, 200, g.blueCarImg, 2.5),   // Blue Limo - faster
			NewCar(700, 500, g.policeCarImg, 3.0), // Police Car - even faster!
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

	// Spawn 15+ good items (random mix of all fish types!)
	goodItems := []*ebiten.Image{
		g.goldfishImg,
		g.rainbowTroutImg,
		g.angelfishImg,
		g.bassImg,
		g.catfishImg,
		g.wormImg,
	}

	for i := 0; i < 17; i++ {
		x := float64(rand.Intn(mapWidth-100) + 50)
		y := float64(rand.Intn(mapHeight-100) + 50)

		// Randomly select from all available fish/worm images
		img := goodItems[rand.Intn(len(goodItems))]

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

		// Update NPCs
		for _, npc := range g.npcs {
			npc.Update()
		}

		// Update Cars (only on level 3)
		for _, car := range g.cars {
			car.Update(g.tileMap.Width(), g.tileMap.Height())
		}

		// Update portal animation
		g.portal.Update()

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

		// Check car collisions (level 3 only) - GAME OVER if hit!
		for _, car := range g.cars {
			if car.CheckCollision(px, py, pw, ph) {
				g.state = StateGameOver
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
			} else if g.currentLevel == 2 {
				// Move to level 3 - FINAL LEVEL!
				g.currentLevel = 3
				g.itemsCollected = 0
				g.loadLevel(3)
				g.state = StatePlaying
			} else {
				// Beat level 3 - GAME WON!
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

		// Draw portal - animate always, but with different opacity
		if g.portalUnlocked {
			// Full opacity when unlocked
			g.portal.DrawWithAlpha(screen, g.cameraX, g.cameraY, 1.0)
		} else {
			// 30% opacity when locked (still animates!)
			g.portal.DrawWithAlpha(screen, g.cameraX, g.cameraY, 0.3)
		}

		// Draw NPCs
		for _, npc := range g.npcs {
			npc.Draw(screen, g.cameraX, g.cameraY)
		}

		// Draw Cars (level 3)
		for _, car := range g.cars {
			car.Draw(screen, g.cameraX, g.cameraY)
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
		text.Draw(screen, "CONGRATULATIONS!", basicfont.Face7x13, screenWidth/2-70, screenHeight/2-20, color.White)
		text.Draw(screen, "YOU BEAT ALL 3 LEVELS!", basicfont.Face7x13, screenWidth/2-100, screenHeight/2, color.White)
		text.Draw(screen, "You are a true Cat Champion!", basicfont.Face7x13, screenWidth/2-110, screenHeight/2+20, color.White)
		text.Draw(screen, fmt.Sprintf("Final score: %d fish collected", g.itemsCollected), basicfont.Face7x13, screenWidth/2-120, screenHeight/2+40, color.White)
		text.Draw(screen, "Press R to play again", basicfont.Face7x13, screenWidth/2-90, screenHeight/2+70, color.White)
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
	controlsText := "WASD/Arrows: Move"
	text.Draw(screen, controlsText, basicfont.Face7x13, screenWidth-250, 35, color.RGBA{200, 200, 200, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Cat's Quest - Project 2 - Jordan DeAndrade")
	ebiten.SetWindowResizable(false)

	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
