package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func contains(keys []ebiten.Key, key ebiten.Key) bool {
	for _, v := range keys {
		if v == key && inpututil.IsKeyJustPressed(key) {
			return true
		}
	}
	return false
}

func getPositionInLine(line int, from int) float32 {
	return (layoutCharacterWidth/4)*float32(line) + float32(from) // 20 as layout
}

func checkActionC2(g *Game) {
	keysPressed := inpututil.PressedKeys()
	if len(keysPressed) == 0 {
		return
	}

	if contains(keysPressed, ebiten.KeyH) {
		todoName(g, &g.character2, 0, false)
	}

	if contains(keysPressed, ebiten.KeyJ) {
		todoName(g, &g.character2, 1, false)
	}

	if contains(keysPressed, ebiten.KeyK) {
		todoName(g, &g.character2, 2, false)
	}

	if contains(keysPressed, ebiten.KeyL) {
		todoName(g, &g.character2, 3, false)
	}

}

func checkActionC1(g *Game) {
	keysPressed := inpututil.PressedKeys()
	if len(keysPressed) == 0 {
		return
	}

	if contains(keysPressed, ebiten.KeyQ) {
		todoName(g, &g.character1, 0, true)
	}

	if contains(keysPressed, ebiten.KeyW) {
		todoName(g, &g.character1, 1, true)
	}

	if contains(keysPressed, ebiten.KeyE) {
		todoName(g, &g.character1, 2, true)
	}

	if contains(keysPressed, ebiten.KeyR) {
		todoName(g, &g.character1, 3, true)
	}
}

func todoName(g *Game, character *Character, line int, isC1 bool) {
	indexNoteHit := checkIfNoteHit(character.notes, line)

	if indexNoteHit == -1 && isC1 {
		g.currentPhaseStance = c1Lost
		return
	} else if indexNoteHit == -1 && !isC1 {
		g.currentPhaseStance = c2Lost
		return
	}

	noteToRemove := character.notes[indexNoteHit]
	character.notes = removeNoteAnyOrder(character.notes, indexNoteHit)
	character.notesToFadeAway = append(character.notesToFadeAway, NoteFadeAway{
		note:    noteToRemove,
		success: true,
		count:   100,
	})

	setCoolDown(character, line, g.frameCount)
	//TODO it would be better to keep the sounds in a map
	switch line {
	case 0:
		rewindAndPlay(character.audioCharacter.sound0)
	case 1:
		rewindAndPlay(character.audioCharacter.sound1)
	case 2:
		rewindAndPlay(character.audioCharacter.sound2)
	case 3:
		rewindAndPlay(character.audioCharacter.sound3)
	}

}

func setCoolDown(c *Character, line int, frameCount int) {
	switch line {
	case 0:
		c.cooldown.line1 = frameCount
	case 1:
		c.cooldown.line2 = frameCount
	case 2:
		c.cooldown.line3 = frameCount
	case 3:
		c.cooldown.line4 = frameCount
	default:
		panic(fmt.Sprintf("should have a line between 1 & 4 but got %d", line))
	}
}

func isLineNotInCoolDown(countFrame int, coolDown int) bool {
	if (countFrame - coolDown) < coolDownFrameForSameNote {
		return false
	}
	return true
}
func getLineOfnoteAdded(g *Game, isC1 bool) int {
	// 180 is aprox the time a note reach the line

	keysPressed := inpututil.PressedKeys()
	if len(keysPressed) == 0 {
		return -1
	}

	if isC1 {
		cooldowns := g.character1.cooldown
		if contains(keysPressed, ebiten.KeyQ) && isLineNotInCoolDown(g.frameCount, cooldowns.line1) {
			return 0
		}

		if contains(keysPressed, ebiten.KeyW) && isLineNotInCoolDown(g.frameCount, cooldowns.line2) {
			return 1
		}

		if contains(keysPressed, ebiten.KeyE) && isLineNotInCoolDown(g.frameCount, cooldowns.line3) {
			return 2
		}

		if contains(keysPressed, ebiten.KeyR) && isLineNotInCoolDown(g.frameCount, cooldowns.line4) {
			return 3
		}
		return -1
	}
	// it's c2
	cooldowns := g.character2.cooldown
	if contains(keysPressed, ebiten.KeyH) && isLineNotInCoolDown(g.frameCount, cooldowns.line1) {
		return 0
	}

	if contains(keysPressed, ebiten.KeyJ) && isLineNotInCoolDown(g.frameCount, cooldowns.line2) {
		return 1
	}

	if contains(keysPressed, ebiten.KeyK) && isLineNotInCoolDown(g.frameCount, cooldowns.line3) {
		return 2
	}

	if contains(keysPressed, ebiten.KeyL) && isLineNotInCoolDown(g.frameCount, cooldowns.line4) {
		return 3
	}
	return -1
}

func checkActionResetGame(g *Game) bool {
	keysPressed := inpututil.PressedKeys()
	if len(keysPressed) == 0 {
		return false
	}

	if contains(keysPressed, ebiten.KeySpace) {
		return true
	}
	return false
}

func checkActionStartAttack(g *Game) int {
	keysPressed := inpututil.PressedKeys()
	cooldowns := &g.character1.cooldown
	if len(keysPressed) == 0 {
		return -1
	}

	if contains(keysPressed, ebiten.KeyQ) && isLineNotInCoolDown(g.frameCount, cooldowns.line1) {
		cooldowns.line1 = g.frameCount
		return 0
	}

	if contains(keysPressed, ebiten.KeyW) && isLineNotInCoolDown(g.frameCount, cooldowns.line2) {
		cooldowns.line2 = g.frameCount
		return 1
	}

	if contains(keysPressed, ebiten.KeyE) && isLineNotInCoolDown(g.frameCount, cooldowns.line3) {
		cooldowns.line3 = g.frameCount
		return 2
	}

	if contains(keysPressed, ebiten.KeyR) && isLineNotInCoolDown(g.frameCount, cooldowns.line4) {
		cooldowns.line4 = g.frameCount
		return 3
	}
	return -1
}

func checkIfNoteHit(notes []Note, line int) int {

	for i, note := range notes {
		if note.line == line {
			if note.y+noteSize > lineMiddleY && note.y < lineMiddleY {
				return i
			}
		}
	}
	return -1
}

func rewindAndPlay(p *audio.Player) {
	p.Rewind()
	p.Play()
}
