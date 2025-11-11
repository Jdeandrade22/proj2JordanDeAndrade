package main

import (
	"bytes"
	"embed"
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

//go:embed assets
var assetsFS embed.FS

type GameState int

const (
	StatePlaying GameState = iota
	StateGameOver
	StateCarDeath // Special state for car collisions
	StateGameWon
)

type Game struct {
	player         *Player
	npcs           []*NPC
	cars           []*Car
	items          []*Item
	portal         *Item
	tileMap        *TileMap
	camera         *Camera
	world          *ebiten.Image
	state          GameState
	currentLevel   int
	itemsCollected int
	portalUnlocked bool

	goldfishImg       *ebiten.Image
	rainbowTroutImg   *ebiten.Image
	angelfishImg      *ebiten.Image
	bassImg           *ebiten.Image
	catfishImg        *ebiten.Image
	wormImg           *ebiten.Image
	badItemImg        *ebiten.Image
	portalImg         *ebiten.Image
	femalePortraitImg *ebiten.Image
	femaleWalkImg     *ebiten.Image
	blueCarImg        *ebiten.Image
	policeCarImg      *ebiten.Image
}

func NewGame() *Game {
	g := &Game{
		state:        StatePlaying,
		currentLevel: 1,
		camera:       Init(screenWidth, screenHeight),
	}

	g.loadAssets()
	g.loadLevel(1)

	return g
}

func (g *Game) loadAssets() {
	g.goldfishImg = g.loadImageFromFS("assets/items/Goldfish.png")
	g.rainbowTroutImg = g.loadImageFromFS("assets/items/Rainbow Trout.png")
	g.angelfishImg = g.loadImageFromFS("assets/items/Angelfish.png")
	g.bassImg = g.loadImageFromFS("assets/items/Bass.png")
	g.catfishImg = g.loadImageFromFS("assets/items/Catfish.png")
	g.wormImg = g.loadImageFromFS("assets/items/Worm.png")
	g.badItemImg = g.loadImageFromFS("assets/items/Rusty Can.png")
	g.portalImg = g.loadImageFromFS("assets/items/Dimensional_Portal.png")
	g.femalePortraitImg = g.loadImageFromFS("assets/npc/portrait female.png")
	g.femaleWalkImg = g.loadImageFromFS("assets/npc/walk and idle.png")
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

		tilesetImages = map[string]*ebiten.Image{
			"clay_tile_64_01.png":      g.loadImageFromFS("assets/background/clay_tile_64_01.png"),
			"grass_01_tile_64_01.png":  g.loadImageFromFS("assets/background/grass_01_tile_64_01.png"),
			"ice_tile_64_01.png":       g.loadImageFromFS("assets/background/ice_tile_64_01.png"),
			"paving_01_tile_64_02.png": g.loadImageFromFS("assets/background/paving_01_tile_64_02.png"),
			"paving_02_tile_64_01.png": g.loadImageFromFS("assets/background/paving_02_tile_64_01.png"),
			"sand_01_tile_64_01.png":   g.loadImageFromFS("assets/background/sand_01_tile_64_01.png"),
		}

		g.player.x = 100
		g.player.y = 100

		// Add NPCs to level 2
		g.npcs = []*NPC{
			NewAnimatedNPC(400, 300, g.femaleWalkImg, 24, 24, 8, 150, true),
			NewAnimatedNPC(800, 200, g.femaleWalkImg, 24, 24, 8, 100, false),
			NewStaticNPC(600, 500, g.femalePortraitImg, 80, false),
			NewStaticNPC(300, 600, g.femalePortraitImg, 120, true),
		}

		g.cars = []*Car{
			NewCar(500, 400, g.blueCarImg, 2.0),
		}

	} else if level == 3 {
		tmxData, err = assetsFS.ReadFile("assets/background/level3.tmx")
		if err != nil {
			log.Fatal("Failed to load level3.tmx:", err)
		}

		img := g.loadImageFromFS("assets/background/orig_big1.png")
		tilesetImages = map[string]*ebiten.Image{
			"orig_big1.png": img,
		}

		g.player.x = 100
		g.player.y = 100

		g.npcs = []*NPC{
			NewAnimatedNPC(300, 250, g.femaleWalkImg, 24, 24, 8, 200, true),
			NewAnimatedNPC(900, 300, g.femaleWalkImg, 24, 24, 8, 180, true),
			NewAnimatedNPC(600, 400, g.femaleWalkImg, 24, 24, 8, 150, false),
			NewAnimatedNPC(450, 600, g.femaleWalkImg, 24, 24, 8, 120, false),
			NewStaticNPC(750, 150, g.femalePortraitImg, 100, true),
			NewStaticNPC(200, 500, g.femalePortraitImg, 130, false),
		}

		g.cars = []*Car{
			NewCar(400, 200, g.blueCarImg, 2.5),
			NewCar(700, 500, g.policeCarImg, 3.0),
		}
	}

	g.tileMap, err = NewTileMap(tmxData, tilesetImages)
	if err != nil {
		log.Fatal("Failed to load tilemap:", err)
	}

	g.world = ebiten.NewImage(g.tileMap.Width(), g.tileMap.Height())
	g.spawnItems()
}

func (g *Game) spawnItems() {
	g.items = []*Item{}
	g.portalUnlocked = false

	mapWidth := g.tileMap.Width()
	mapHeight := g.tileMap.Height()

	goodItems := []*ebiten.Image{
		g.goldfishImg,
		g.rainbowTroutImg,
		g.angelfishImg,
		g.bassImg,
		g.catfishImg,
	}

	for i := 0; i < 17; i++ {
		x := float64(rand.Intn(mapWidth-100) + 50)
		y := float64(rand.Intn(mapHeight-100) + 50)
		img := goodItems[rand.Intn(len(goodItems))]
		g.items = append(g.items, NewItem(x, y, ItemGood, img))
	}

	for i := 0; i < 3; i++ {
		x := float64(rand.Intn(mapWidth-100) + 50)
		y := float64(rand.Intn(mapHeight-100) + 50)

		g.items = append(g.items, NewItem(x, y, ItemBad, g.badItemImg))
	}

	for i := 0; i < 2; i++ {
		x := float64(rand.Intn(mapWidth-100) + 50)
		y := float64(rand.Intn(mapHeight-100) + 50)

		g.items = append(g.items, NewItem(x, y, ItemBad, g.wormImg))
	}

	portalX := float64(mapWidth - 150)
	portalY := float64(mapHeight - 150)
	g.portal = NewItem(portalX, portalY, ItemPortal, g.portalImg)
}

func (g *Game) Update() error {
	if g.state == StatePlaying {
		g.player.Update(g.tileMap.Width(), g.tileMap.Height())

		for _, npc := range g.npcs {
			npc.Update()
		}

		for _, car := range g.cars {
			car.Update(g.tileMap.Width(), g.tileMap.Height())
		}

		g.portal.Update()

		g.camera.Follow.W = int(g.player.x + float64(g.player.width)/2)
		g.camera.Follow.H = int(g.player.y + float64(g.player.height)/2)

		px, py, pw, ph := g.player.GetBounds()
		for _, item := range g.items {
			if item.CheckCollision(px, py, pw, ph) {
				item.collected = true

				if item.itemType == ItemGood {
					g.itemsCollected++

					if g.itemsCollected >= 9 {
						g.portalUnlocked = true
					}
				} else if item.itemType == ItemBad {
					g.state = StateGameOver
				}
			}
		}

		for _, car := range g.cars {
			if car.CheckCollision(px, py, pw, ph) {
				g.state = StateCarDeath
			}
		}

		if g.portalUnlocked && g.portal.CheckCollision(px, py, pw, ph) {
			if g.currentLevel == 1 {
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
				g.state = StateGameWon
			}
		}
	} else if g.state == StateGameOver || g.state == StateCarDeath {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.state = StatePlaying
			g.currentLevel = 1
			g.itemsCollected = 0
			g.loadLevel(1)
		}
	} else if g.state == StateGameWon {
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
		g.world.Clear()

		g.tileMap.Draw(g.world, 0, 0)

		for _, item := range g.items {
			item.Draw(g.world, 0, 0)
		}

		if g.portalUnlocked {
			g.portal.DrawWithAlpha(g.world, 0, 0, 1.0)
		} else {
			g.portal.DrawWithAlpha(g.world, 0, 0, 0.3)
		}

		for _, npc := range g.npcs {
			npc.Draw(g.world, 0, 0)
		}

		for _, car := range g.cars {
			car.Draw(g.world, 0, 0)
		}

		g.player.Draw(g.world, 0, 0)
		g.camera.Draw(g.world, screen)
		g.drawUI(screen)

	} else if g.state == StateGameOver {
		text.Draw(screen, "GAME OVER!", basicfont.Face7x13, screenWidth/2-50, screenHeight/2, color.White)
		text.Draw(screen, "You touched a bad item!", basicfont.Face7x13, screenWidth/2-90, screenHeight/2+20, color.White)
		text.Draw(screen, "Press R to restart", basicfont.Face7x13, screenWidth/2-80, screenHeight/2+40, color.White)

	} else if g.state == StateCarDeath {
		text.Draw(screen, "GAME OVER!", basicfont.Face7x13, screenWidth/2-50, screenHeight/2-20, color.White)
		text.Draw(screen, "Cats only have 1 life around here!", basicfont.Face7x13, screenWidth/2-120, screenHeight/2, color.White)
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
	uiRect := ebiten.NewImage(screenWidth, 40)
	uiRect.Fill(color.RGBA{0, 0, 0, 180})
	screen.DrawImage(uiRect, &ebiten.DrawImageOptions{})

	levelText := fmt.Sprintf("Level: %d", g.currentLevel)
	text.Draw(screen, levelText, basicfont.Face7x13, 10, 20, color.White)

	collectionText := fmt.Sprintf("Fish Collected: %d", g.itemsCollected)
	text.Draw(screen, collectionText, basicfont.Face7x13, 10, 35, color.White)

	portalText := "Portal: Locked (Need 9 fish)"
	if g.portalUnlocked {
		portalText = "Portal: UNLOCKED! Go to portal!"
	}
	text.Draw(screen, portalText, basicfont.Face7x13, screenWidth-250, 20, color.RGBA{255, 215, 0, 255})

	controlsText := "WASD/Arrows: Move"
	text.Draw(screen, controlsText, basicfont.Face7x13, screenWidth-250, 35, color.RGBA{200, 200, 200, 255})
}

func (g *Game) Layout(_ int, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Cat's Quest - Project 2 - Jordan DeAndrade")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
