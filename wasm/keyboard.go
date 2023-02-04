package main

import (
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
		todoName(g,&g.character2,0,false)
	}

	if contains(keysPressed, ebiten.KeyJ) {
		todoName(g,&g.character2,1,false)
	}

	if contains(keysPressed, ebiten.KeyK) {
		todoName(g,&g.character2,2,false)
	}

	if contains(keysPressed, ebiten.KeyL) {
		todoName(g,&g.character2,3,false)
	}

}

func checkActionC1(g *Game) {
	keysPressed := inpututil.PressedKeys()
	if len(keysPressed) == 0 {
		return
	}

	if contains(keysPressed, ebiten.KeyQ) {
		todoName(g,&g.character1,0,true)
	}

	if contains(keysPressed, ebiten.KeyW) {
		todoName(g,&g.character1,1,true)
	}

	if contains(keysPressed, ebiten.KeyE) {
		todoName(g,&g.character1,2,true)
	}

	if contains(keysPressed, ebiten.KeyR) {
		todoName(g,&g.character1,3,true)
	}
}

//TODO try character with pointer
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

	rewindAndPlay(character.audioCharacter.sound0)
}

func noteWasAdded(g *Game, isC1 bool) bool {
	// 180 is aprox the time a note reach the line
	count := g.count - 160
	keysPressed := inpututil.PressedKeys()
	if len(keysPressed) == 0 {
		return false
	}
	correctKeyPressed := false
	if isC1 {
		if contains(keysPressed, ebiten.KeyQ) {
			g.mapNoteToPlay[count] = 0
			correctKeyPressed = true
		}

		if contains(keysPressed, ebiten.KeyW) {
			g.mapNoteToPlay[count] = 1
			correctKeyPressed = true
		}

		if contains(keysPressed, ebiten.KeyE) {
			g.mapNoteToPlay[count] = 2
			correctKeyPressed = true
		}

		if contains(keysPressed, ebiten.KeyR) {
			g.mapNoteToPlay[count] = 3
			correctKeyPressed = true
		}
		return correctKeyPressed
	}
	// it's c2

	if contains(keysPressed, ebiten.KeyH) {
		g.mapNoteToPlay[count] = 0
		correctKeyPressed = true
	}

	if contains(keysPressed, ebiten.KeyJ) {
		g.mapNoteToPlay[count] = 1
		correctKeyPressed = true
	}

	if contains(keysPressed, ebiten.KeyK) {
		g.mapNoteToPlay[count] = 2
		correctKeyPressed = true
	}

	if contains(keysPressed, ebiten.KeyL) {
		g.mapNoteToPlay[count] = 3
		correctKeyPressed = true
	}
	return correctKeyPressed
}

func checkActionStartAttack(g *Game) {
	keysPressed := inpututil.PressedKeys()
	if len(keysPressed) == 0 {
		return
	}

	if contains(keysPressed, ebiten.KeyQ) {
		g.mapNoteToPlay[g.count] = 0

		g.character1.notes = append(g.character1.notes, Note{
			x:         getPositionInLine(0, 0),
			y:         screenHeight - 20,
			line:      0,
			direction: up,
		})
	}

	if contains(keysPressed, ebiten.KeyW) {
		g.mapNoteToPlay[g.count] = 1
		g.character1.notes = append(g.character1.notes, Note{
			x:         getPositionInLine(1, 0),
			y:         screenHeight - 20,
			line:      1,
			direction: up,
		})
	}

	if contains(keysPressed, ebiten.KeyE) {
		g.mapNoteToPlay[g.count] = 2
		g.character1.notes = append(g.character1.notes, Note{
			x:         getPositionInLine(2, 0),
			y:         screenHeight - 20,
			line:      2,
			direction: up,
		})
	}

	if contains(keysPressed, ebiten.KeyR) {
		g.mapNoteToPlay[g.count] = 3
		g.character1.notes = append(g.character1.notes, Note{
			x:         getPositionInLine(3, 0),
			y:         screenHeight - 20,
			line:      3,
			direction: up,
		})
	}
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
