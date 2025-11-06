# Cat Fish Quest - Project 2
## Jordan DeAndrade

---

## Game Description
A cat-themed adventure game where you collect fish and avoid dangerous items while exploring two different levels! Navigate through beautiful tiled maps, dodge bad items, and make your way to the dimensional portal to progress.

---

## How to Run
```bash
go build -o catfishquest
./catfishquest
```

Or simply:
```bash
go run .
```

---

## Controls
- **Arrow Keys** or **WASD**: Move the cat player
  - ↑ / W: Move up
  - ↓ / S: Move down
  - ← / A: Move left
  - → / D: Move right
- **R**: Restart game (when game is over)

---

## Game Objectives

### Level 1
1. Collect fish (good items) scattered around the map
2. Avoid rusty cans (bad items) - they're dangerous for cats!
3. Collect **at least 9 fish** to unlock the portal
4. Enter the **Dimensional Portal** (purple swirling portal) to advance to Level 2

### Level 2
1. Collect more fish while avoiding the bad items
2. Watch out for the **2 animated NPCs** that move around the map:
   - **Drone Dog** (brown): Moves horizontally back and forth
   - **Angry Character** (pink-haired): Moves vertically up and down
3. Collect 9 fish to unlock the portal
4. Enter the portal to win the game!

---

## Items Guide

### Good Items (Collect These!) 🐟
- **Goldfish** (yellow/orange fish): Worth collecting!
- **Worm** (pink worm): Tasty cat snack!
- Collect **9 or more** to unlock the portal on each level

### Bad Items (AVOID These!) ⚠️
- **Rusty Can** (red can): Dangerous! Touching one ends the game immediately
- There are **5 bad items** on each level - be careful!

### Special Items
- **Dimensional Portal** (purple swirling portal):
  - Appears **locked** (grayed out) until you collect 9 fish
  - Once **unlocked** (fully colored), walk into it to progress
  - Located near the bottom-right corner of each map

---

## Game Features

✅ **Window Size**: 800x600 pixels  
✅ **Tiled Maps**: Two beautiful 20x20 tile maps  
✅ **TMX Files**: Maps loaded from level1.tmx and level2.tmx  
✅ **Camera System**: Follows the player smoothly  
✅ **Items**: 17+ good items and 5 bad items per level  
✅ **Collection Counter**: Always visible at the top of the screen  
✅ **Portal System**: Unlocks at 9 items, leads to next level  
✅ **Two Levels**: Progress from Level 1 to Level 2  
✅ **Animated NPCs**: 2 NPCs on Level 2 with movement  
✅ **Animated Player**: Cat sprite with directional animations  
✅ **Boundary Collision**: Player cannot move off the map  
✅ **Embedded Assets**: All assets loaded using go:embed  
✅ **Game States**: Playing, Game Over, and Game Won screens  

---

## NPC Recognition

### Level 1
- No NPCs (peaceful exploration level)

### Level 2
1. **Drone Dog** (Brown/tan colored)
   - Moves horizontally left and right
   - Has a mechanical/robotic appearance
   
2. **Murong Yi Character** (Pink-haired character)
   - Moves vertically up and down
   - Has an angry expression
   - From the Dragon Hero visual novel series

---

## Portal Recognition
- The **Dimensional Portal** is a purple/violet swirling energy portal
- Appears in the lower-right area of each map
- When **locked**: Appears grayed out/transparent (need more fish!)
- When **unlocked**: Appears bright and fully visible
- **UI indicator** at the top shows portal status at all times

---

## Game Over Conditions
- Touching any **bad item** (rusty can) = Instant Game Over
- Press **R** to restart from Level 1

---

## Winning the Game
1. Complete Level 1 by collecting 9+ fish and entering the portal
2. Complete Level 2 by collecting 9+ fish and entering the portal
3. See the victory screen with your total collection count!
4. Press **R** to play again

---

## Technical Details
- **Engine**: Ebitengine (v2)
- **Language**: Go 1.21
- **Map Format**: TMX (Tiled Map Editor)
- **Assets**: Embedded using `go:embed` directive
- **Asset Organization**: All assets in subfolders under `/assets/`
  - `/assets/background/` - Map files and tilesets
  - `/assets/sprites/` - Player and portal sprites
  - `/assets/items/` - Collectible items
  - `/assets/npc/` - NPC character sprites

---

## What's Implemented
✅ All rubric requirements completed
✅ 800x600 window (within 500-1000 requirement)
✅ Two 20x20 tiled maps from TMX files
✅ Camera system following player
✅ 17+ good items randomly distributed
✅ Collection counter displayed on screen
✅ 5 bad items that trigger game over
✅ Game over screen with restart option
✅ Portal unlocks at 9 items
✅ Level progression system
✅ 2 animated NPCs on Level 2
✅ Animated player with directional movement
✅ Map boundary collision
✅ go:embed for all assets
✅ Organized asset folder structure
✅ Complete README (this file!)

---

## Credits
- **Developer**: Jordan DeAndrade
- **Game Engine**: Ebitengine
- **Tiled Map Loader**: lafriks/go-tiled
- **Art Assets**: Various free pixel art resources
- **NPC Character**: Murong Yi from Dragon Hero series by linxuelian.itch.io

---

## Have Fun! 🐱🐟
Enjoy helping the cat collect fish and navigate through the dimensional portals!

