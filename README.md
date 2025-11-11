# Cat's Quest - Project 2
**Jordan DeAndrade**

## Overview

A 2D proto-game built with Ebitengine v2 featuring an animated cat protagonist navigating three distinct levels. Players collect fish while avoiding obstacles, utilizing a dynamic camera system that follows the player across large tiled maps.

## Installation and Running

**Prerequisites:** Go 1.21 or higher

```bash
git clone https://github.com/Jdeandrade22/proj2JordanDeAndrade.git
cd proj2JordanDeAndrade
go run .
```

Or build and run:
```bash
go build
./project2_jordandeandrade
```

## Controls

- **WASD / Arrow Keys** - Move the cat in 8 directions (including diagonals)
- **R** - Restart after game over or winning

## Gameplay

### Objective
Collect 9 fish to unlock the portal in each level. Complete all three levels to win. Avoid hazards and moving obstacles.

### Good Items (17 per level)
- **Goldfish** - Yellow/orange fish
- **Rainbow Trout** - Colorful fish
- **Angelfish** - Elegant fish
- **Bass** - Strong fish
- **Catfish** - Whiskered fish

Collect 9 to unlock the portal and advance to the next level.

### Bad Items (5 per level)
- **Rusty Can** (3x) - Red can sprite, instant game over
- **Worm** (2x) - Pink worm sprite, instant game over

### Vehicle Hazards (Levels 2 & 3)
- **Blue Limo** - Cyan animated car with random movement
- **Police Car** - Police vehicle with random patrol patterns

Collision with vehicles triggers: "Cats only have 1 life around here!"

### Portal
- **Location:** Near bottom-right of each map
- **Locked State:** 30% opacity until 9 items collected
- **Unlocked State:** Full opacity with 6-frame animation
- **Function:** Advances to next level (or wins game on Level 3)

## Levels

**Level 1 - Introduction**  
Basic level with only fish and stationary hazards. No NPCs or vehicles.

**Level 2 - NPCs and Vehicles**  
Introduces animated Female Walking Characters and Female Portrait NPCs with patrol patterns. Blue Limo appears as the first moving hazard.

**Level 3 - Full Challenge**  
Four animated walking NPCs with extended patrol ranges. Both Blue Limo and Police Car move at high speeds with random patterns.

## NPCs

- **Female Walking Character** - 8-frame animated sprite, patrols horizontally or vertically
- **Female Portrait** - Static sprite with predictable movement patterns
- **Blue Limo** - Animated cyan vehicle with random movement
- **Police Car** - Animated police vehicle with high-speed random patrol

## Technical Features

### Window and Display
800x600 pixel window with fixed resolution.

### Camera System
Custom camera implementation that follows the player, keeping them centered while clamping to map boundaries. Matches the camera library interface from class with `Init()`, `Follow`, and `Draw()` methods.

### Tiled Maps
Three 20x20 tile maps loaded from TMX files using the `go-tiled` library. Each level is a separate map file stored in `/assets/background/`.

### Animation System
Custom sprite animation supporting multi-frame sheets with variable frame counts:
- **Player:** 8 directional animations with unique sprites for each direction
- **NPCs:** Varied frame counts and animation speeds
- **Portal:** 6-frame looping animation

### Asset Management
All assets embedded using `go:embed` directive. Organized in subfolder structure:
- `/assets/background/` - TMX files and tilesets
- `/assets/sprites/` - Player directional sprites
- `/assets/items/` - Collectibles and hazards
- `/assets/npc/` - NPC and vehicle sprites

### Collision Detection
AABB collision system for item pickup, hazard contact, and portal entry. Player hitbox is 32x32 (smaller than visual sprite) for better gameplay feel.

## Project Structure

```
proj2JordanDeAndrade/
├── main.go          - Game loop, state management, level loading
├── player.go        - Player movement and animation (8 directions)
├── npcs.go          - NPC behavior and rendering
├── cars.go          - Vehicle hazards with random movement
├── items.go         - Collectibles, hazards, and portal
├── tilemap.go       - TMX map loading and rendering
├── animation.go     - Sprite animation system
├── camera.go        - Camera (Init, Follow, Draw)
├── go.mod           - Dependencies
└── assets/          - Embedded game assets
```

## Dependencies

- `github.com/hajimehoshi/ebiten/v2` - 2D game engine
- `github.com/lafriks/go-tiled` - TMX map parser
- `golang.org/x/image` - Image processing

## Author

**Jordan DeAndrade**  
Email: j2deandrade@student.bridgew.edu  
GitHub: https://github.com/Jdeandrade22

some readme documentation provided through ai!
