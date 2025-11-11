package main

import (
	"bytes"
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const (
	sampleRate = 48000
)

type AudioManager struct {
	audioContext  *audio.Context
	bgmPlayer     *audio.Player
	eatSoundData  []byte
	carHonkData   []byte
	ouchSoundData []byte
}

func NewAudioManager() *AudioManager {
	am := &AudioManager{
		audioContext: audio.NewContext(sampleRate),
	}

	// Load background music
	am.loadBackgroundMusic()

	// Load eat sound effect
	am.loadEatSound()

	// Load car honk sound
	am.loadCarHonkSound()

	// Load ouch sound
	am.loadOuchSound()

	return am
}

func (am *AudioManager) loadBackgroundMusic() {
	// Load the MP3 from embedded filesystem
	data, err := assetsFS.ReadFile("assets/sounds/cottagecore-17463.mp3")
	if err != nil {
		log.Printf("Failed to load background music: %v", err)
		return
	}

	// Decode MP3
	stream, err := mp3.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		log.Printf("Failed to decode background music: %v", err)
		return
	}

	// Create an infinite loop
	loopStream := audio.NewInfiniteLoop(stream, stream.Length())

	// Create player
	player, err := am.audioContext.NewPlayer(loopStream)
	if err != nil {
		log.Printf("Failed to create music player: %v", err)
		return
	}

	am.bgmPlayer = player
}

func (am *AudioManager) loadEatSound() {
	// Load the MP3 from embedded filesystem
	data, err := assetsFS.ReadFile("assets/sounds/wet-squelchy-impact-352302.mp3")
	if err != nil {
		log.Printf("Failed to load eat sound: %v", err)
		return
	}

	// Decode MP3
	stream, err := mp3.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		log.Printf("Failed to decode eat sound: %v", err)
		return
	}

	// Read all data into memory for quick playback
	eatData, err := io.ReadAll(stream)
	if err != nil {
		log.Printf("Failed to read eat sound data: %v", err)
		return
	}

	am.eatSoundData = eatData
}

func (am *AudioManager) PlayBackgroundMusic() {
	if am.bgmPlayer != nil && !am.bgmPlayer.IsPlaying() {
		am.bgmPlayer.Play()
	}
}

func (am *AudioManager) StopBackgroundMusic() {
	if am.bgmPlayer != nil && am.bgmPlayer.IsPlaying() {
		am.bgmPlayer.Pause()
	}
}

func (am *AudioManager) loadCarHonkSound() {
	// Load the MP3 from embedded filesystem
	data, err := assetsFS.ReadFile("assets/sounds/car-honk-386166.mp3")
	if err != nil {
		log.Printf("Failed to load car honk sound: %v", err)
		return
	}

	// Decode MP3
	stream, err := mp3.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		log.Printf("Failed to decode car honk sound: %v", err)
		return
	}

	// Read all data into memory for quick playback
	carHonkData, err := io.ReadAll(stream)
	if err != nil {
		log.Printf("Failed to read car honk sound data: %v", err)
		return
	}

	am.carHonkData = carHonkData
}

func (am *AudioManager) loadOuchSound() {
	// Load the MP3 from embedded filesystem
	data, err := assetsFS.ReadFile("assets/sounds/ouchnoise-96832.mp3")
	if err != nil {
		log.Printf("Failed to load ouch sound: %v", err)
		return
	}

	// Decode MP3
	stream, err := mp3.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		log.Printf("Failed to decode ouch sound: %v", err)
		return
	}

	// Read all data into memory for quick playback
	ouchData, err := io.ReadAll(stream)
	if err != nil {
		log.Printf("Failed to read ouch sound data: %v", err)
		return
	}

	am.ouchSoundData = ouchData
}

func (am *AudioManager) PlayEatSound() {
	if am.eatSoundData == nil {
		return
	}

	// Create a new player each time for the sound effect
	// This allows multiple sounds to play simultaneously if needed
	stream := bytes.NewReader(am.eatSoundData)
	player, err := am.audioContext.NewPlayer(stream)
	if err != nil {
		log.Printf("Failed to create eat sound player: %v", err)
		return
	}

	// Play the sound
	player.Play()

	// Note: The player will be garbage collected after the sound finishes
}

func (am *AudioManager) PlayCarHonkSound() {
	if am.carHonkData == nil {
		return
	}

	stream := bytes.NewReader(am.carHonkData)
	player, err := am.audioContext.NewPlayer(stream)
	if err != nil {
		log.Printf("Failed to create car honk player: %v", err)
		return
	}

	player.Play()
}

func (am *AudioManager) PlayOuchSound() {
	if am.ouchSoundData == nil {
		return
	}

	stream := bytes.NewReader(am.ouchSoundData)
	player, err := am.audioContext.NewPlayer(stream)
	if err != nil {
		log.Printf("Failed to create ouch sound player: %v", err)
		return
	}

	player.Play()
}

//audio functions implemented with DeepseekR1
