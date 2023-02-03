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

func getPositionInLine(line int,from int)float32{
	return (layoutCharacterWidth/4)*float32(line) + float32(from) // 20 as layout
}

func checkActionC2(g *Game) {
	keysPressed := inpututil.PressedKeys()
	if len(keysPressed) == 0 {
		return
	}

	if contains(keysPressed, ebiten.KeyH) {
		checkIfNoteHit(g, g.character2, 0, false)
	}

	if contains(keysPressed, ebiten.KeyJ) {
		checkIfNoteHit(g, g.character2, 1, false)
	}

	if contains(keysPressed, ebiten.KeyK) {
		checkIfNoteHit(g, g.character2, 2, false)
	}

	if contains(keysPressed, ebiten.KeyL) {
		checkIfNoteHit(g, g.character2, 3, false)
	}
}

func checkActionC1(g *Game) {
	keysPressed := inpututil.PressedKeys()
	if len(keysPressed) == 0 {
		return
	}

	if contains(keysPressed, ebiten.KeyQ) {
		checkIfNoteHit(g, g.character1, 0, true)
	}

	if contains(keysPressed, ebiten.KeyW) {
		checkIfNoteHit(g, g.character1, 1, true)
	}

	if contains(keysPressed, ebiten.KeyE) {
		checkIfNoteHit(g, g.character1, 2, true)
	}

	if contains(keysPressed, ebiten.KeyR) {
		checkIfNoteHit(g, g.character1, 3, true)
	}
}

func noteWasAdded(g *Game, isC1 bool) bool{
	// 180 is aprox the time a note reach the line
	count := g.count-160
	keysPressed := inpututil.PressedKeys()
	if len(keysPressed) == 0 {
		return false
	}
	correctKeyPressed := false
	if isC1 {
		if contains(keysPressed, ebiten.KeyQ) {
			g.mapNoteToPlay[count] = 0
			correctKeyPressed=true
		}
	
		if contains(keysPressed, ebiten.KeyW) {
			g.mapNoteToPlay[count] = 1
			correctKeyPressed=true
		}
	
		if contains(keysPressed, ebiten.KeyE) {
			g.mapNoteToPlay[count] = 2
			correctKeyPressed=true
		}
	
		if contains(keysPressed, ebiten.KeyR) {
			g.mapNoteToPlay[count] = 3
			correctKeyPressed=true
		}
		return correctKeyPressed
	}
	// it's c2

	if contains(keysPressed, ebiten.KeyH) {
		g.mapNoteToPlay[count] = 0
			correctKeyPressed=true
	}

	if contains(keysPressed, ebiten.KeyJ) {
		g.mapNoteToPlay[count] = 1
			correctKeyPressed=true
	}

	if contains(keysPressed, ebiten.KeyK) {
		g.mapNoteToPlay[count] = 2
			correctKeyPressed=true
	}

	if contains(keysPressed, ebiten.KeyL) {
		g.mapNoteToPlay[count] = 3
			correctKeyPressed=true
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
			x:    getPositionInLine(0,0),
			y:    screenHeight - 20,
			line: 0,
			direction: up,
		})
	}

	if contains(keysPressed, ebiten.KeyW) {
		g.mapNoteToPlay[g.count] = 1
		g.character1.notes = append(g.character1.notes, Note{
			x:    getPositionInLine(1,0),
			y:    screenHeight - 20,
			line: 1,
			direction: up,
		})
	}

	if contains(keysPressed, ebiten.KeyE) {
		g.mapNoteToPlay[g.count] = 2
		g.character1.notes = append(g.character1.notes, Note{
			x:    getPositionInLine(2,0),
			y:    screenHeight - 20,
			line: 2,
			direction: up,
		})
	}

	if contains(keysPressed, ebiten.KeyR) {
		g.mapNoteToPlay[g.count] = 3
		g.character1.notes = append(g.character1.notes, Note{
			x:    getPositionInLine(3,0),
			y:    screenHeight - 20,
			line: 3,
			direction: up,
		})
	}
}

func checkIfNoteHit(g *Game, c Character, line int, isC1 bool){

	for i, note := range c.notes {
		if note.line == line {
			if note.y+noteSize > lineMiddleY && note.y < lineMiddleY {			
				c.notes = append(c.notes[:i], c.notes[i+1:]...)		
				fmt.Println(len(c.notes))	
				c.notesToFadeAway = append(c.notesToFadeAway, NoteFadeAway{
					note:    note,
					success: true,
					count:   100,
				})

				p := g.character1.audioCharacter
				switch line {
				case 0:
					rewindAndPlay(p.sound0)
				case 1:
					rewindAndPlay(p.sound1)
				case 2:
					rewindAndPlay(p.sound2)
				case 3:
					rewindAndPlay(p.sound3)
				}
				break
				//return
			}
		}
	}

	//TODO I think it's better to return the object character
	if isC1 {
		g.character1 = c
		return
	}
	g.character2 = c
}

func rewindAndPlay(p *audio.Player) {
	p.Rewind()
	p.Play()
}
