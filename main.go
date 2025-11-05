package main

import "embed"

var EmbeddedAssets embed.FS

type GameState int

const (
	StatePlaying GameState = iota
	StateGameOver
	StateNextLevel
)
