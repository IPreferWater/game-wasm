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
		/*indexNoteHit := checkIfNoteHit(g.character2.notes, 0)

		if indexNoteHit == -1 {
			g.currentPhaseStance = c2Lost
			return
		}

		noteToRemove := g.character2.notes[indexNoteHit]
		g.character2.notes = removeNoteAnyOrder(g.character2.notes, indexNoteHit)
		g.character2.notesToFadeAway = append(g.character2.notesToFadeAway, NoteFadeAway{
			note:    noteToRemove,
			success: true,
			count:   100,
		})

		rewindAndPlay(g.character2.audioCharacter.sound0)*/
	}

	if contains(keysPressed, ebiten.KeyJ) {
		indexNoteHit := checkIfNoteHit(g.character2.notes, 1)

		if indexNoteHit == -1 {
			g.currentPhaseStance = c2Lost
			return
		}

		noteToRemove := g.character2.notes[indexNoteHit]
		g.character2.notes = removeNoteAnyOrder(g.character2.notes, indexNoteHit)
		g.character2.notesToFadeAway = append(g.character2.notesToFadeAway, NoteFadeAway{
			note:    noteToRemove,
			success: true,
			count:   100,
		})

		rewindAndPlay(g.character2.audioCharacter.sound0)
	}

	if contains(keysPressed, ebiten.KeyK) {
		indexNoteHit := checkIfNoteHit(g.character2.notes, 2)

		if indexNoteHit == -1 {
			g.currentPhaseStance = c2Lost
			return
		}

		noteToRemove := g.character2.notes[indexNoteHit]
		g.character2.notes = removeNoteAnyOrder(g.character2.notes, indexNoteHit)
		g.character2.notesToFadeAway = append(g.character2.notesToFadeAway, NoteFadeAway{
			note:    noteToRemove,
			success: true,
			count:   100,
		})

		rewindAndPlay(g.character2.audioCharacter.sound0)
	}

	if contains(keysPressed, ebiten.KeyL) {
		indexNoteHit := checkIfNoteHit(g.character2.notes, 3)

		if indexNoteHit == -1 {
			g.currentPhaseStance = c2Lost
			return
		}

		noteToRemove := g.character2.notes[indexNoteHit]
		g.character2.notes = removeNoteAnyOrder(g.character2.notes, indexNoteHit)
		g.character2.notesToFadeAway = append(g.character2.notesToFadeAway, NoteFadeAway{
			note:    noteToRemove,
			success: true,
			count:   100,
		})

		rewindAndPlay(g.character2.audioCharacter.sound0)
	}

}

func checkActionC1(g *Game) {
	keysPressed := inpututil.PressedKeys()
	if len(keysPressed) == 0 {
		return
	}

	if contains(keysPressed, ebiten.KeyQ) {
		indexNoteHit := checkIfNoteHit(g.character1.notes, 0)

		if indexNoteHit == -1 {
			g.currentPhaseStance = c1Lost
			return
		}

		noteToRemove := g.character1.notes[indexNoteHit]
		g.character1.notes = removeNoteAnyOrder(g.character1.notes, indexNoteHit)
		g.character1.notesToFadeAway = append(g.character1.notesToFadeAway, NoteFadeAway{
			note:    noteToRemove,
			success: true,
			count:   100,
		})

		rewindAndPlay(g.character1.audioCharacter.sound0)
	}

	if contains(keysPressed, ebiten.KeyW) {
		indexNoteHit := checkIfNoteHit(g.character1.notes, 1)

		if indexNoteHit == -1 {
			g.currentPhaseStance = c1Lost
			return
		}

		noteToRemove := g.character1.notes[indexNoteHit]
		g.character1.notes = removeNoteAnyOrder(g.character1.notes, indexNoteHit)
		g.character1.notesToFadeAway = append(g.character1.notesToFadeAway, NoteFadeAway{
			note:    noteToRemove,
			success: true,
			count:   100,
		})

		rewindAndPlay(g.character1.audioCharacter.sound0)
	}

	if contains(keysPressed, ebiten.KeyE) {
		indexNoteHit := checkIfNoteHit(g.character1.notes, 2)

		if indexNoteHit == -1 {
			g.currentPhaseStance = c1Lost
			return
		}

		noteToRemove := g.character1.notes[indexNoteHit]
		g.character1.notes = removeNoteAnyOrder(g.character1.notes, indexNoteHit)
		g.character1.notesToFadeAway = append(g.character1.notesToFadeAway, NoteFadeAway{
			note:    noteToRemove,
			success: true,
			count:   100,
		})

		rewindAndPlay(g.character1.audioCharacter.sound0)
	}

	if contains(keysPressed, ebiten.KeyR) {
		indexNoteHit := checkIfNoteHit(g.character1.notes, 3)

		if indexNoteHit == -1 {
			g.currentPhaseStance = c1Lost
			return
		}

		noteToRemove := g.character1.notes[indexNoteHit]
		g.character1.notes = removeNoteAnyOrder(g.character1.notes, indexNoteHit)
		g.character1.notesToFadeAway = append(g.character1.notesToFadeAway, NoteFadeAway{
			note:    noteToRemove,
			success: true,
			count:   100,
		})

		rewindAndPlay(g.character1.audioCharacter.sound0)
	}
}

//TODO try character with pointer
func todoName(g *Game, character *Character, line int, isC1 bool) {
	indexNoteHit := checkIfNoteHit(character.notes, line)

	if indexNoteHit == -1 && isC1 {
		g.currentPhaseStance = c1Lost
		return
	} else if indexNoteHit == -1 && isC1 {
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

	//TODO I think it's better to return the object character
	/*if isC1 {
		g.character1 = c
		return
	}
	g.character2 = c*/
}

func rewindAndPlay(p *audio.Player) {
	p.Rewind()
	p.Play()
}
