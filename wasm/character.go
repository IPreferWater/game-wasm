package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type Character struct {
	audioCharacter  AudioCharacter
	notes           []Note
	notesToFadeAway []NoteFadeAway
	characterSprite CharacterSprite
	cooldown        Cooldown
}

type Note struct {
	x         float32
	y         float32
	line      int
	direction direction
}

type direction int64

const (
	up direction = iota
	down
)

type AudioCharacter struct {
	sound0 *audio.Player
	sound1 *audio.Player
	sound2 *audio.Player
	sound3 *audio.Player
}

type NoteFadeAway struct {
	note    Note
	success bool
	count   int
}

type CharacterSprite struct {
	img     *ebiten.Image
	sprites map[SpriteStance]Sprite
}

type Cooldown struct {
	line1 int
	line2 int
	line3 int
	line4 int
}

func (c *Character) updateNotesAndCheckIfLost() bool {
	notes := c.notes
	for i := 0; i < len(notes); i++ {
		//update position
		if notes[i].direction == up {
			notes[i].y -= 1
		} else {
			notes[i].y += 1
		}

		// if out of scope, it's lost
		if (notes[i].y < 0+10) || notes[i].y > screenHeight-10 {
			if notes[i].direction == down {
				return true
			}
			notes = removeNoteAnyOrder(notes, i)
			i--
		}
	}

	for i := 0; i < len(c.notesToFadeAway); i++ {
		c.notesToFadeAway[i].count++

		if c.notesToFadeAway[i].count >= 100 {
			c.notesToFadeAway = removeNoteToFadeAwayAnyOrder(c.notesToFadeAway, i)
			i--
		}
	}
	c.notes = notes
	return false
}

func removeNoteAnyOrder(s []Note, i int) []Note {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func removeNoteToFadeAwayAnyOrder(s []NoteFadeAway, i int) []NoteFadeAway {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
