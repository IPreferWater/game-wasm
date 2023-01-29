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

func getPositionInLine(line int,from int)float32{
	//((screenWidth/3)/4)*note.line + 20 // 20 as layout
	return (layoutCharacterWidth/4)*float32(line) + float32(from) // 20 as layout
}

func checkActionC2(g *Game) {
	keysPressed := inpututil.PressedKeys()
	if len(keysPressed) == 0 {
		return
	}

	if contains(keysPressed, ebiten.KeyJ) {
		checkIfNoteIsHit(g, 0)
	}

	if contains(keysPressed, ebiten.KeyK) {
		checkIfNoteIsHit(g, 1)
	}

	if contains(keysPressed, ebiten.KeyL) {
		checkIfNoteIsHit(g, 2)
	}

	if contains(keysPressed, ebiten.KeyM) {
		checkIfNoteIsHit(g, 3)
	}
}

func checkAction(g *Game) {
	keysPressed := inpututil.PressedKeys()
	if len(keysPressed) == 0 {
		return
	}

	if contains(keysPressed, ebiten.KeyQ) {
		g.character1.m[g.count] = 0

		g.character1.notes = append(g.character1.notes, Note{
			x:    getPositionInLine(0,0),
			y:    screenHeight - 20,
			line: 0,
			direction: up,
		})
	}

	if contains(keysPressed, ebiten.KeyW) {
		g.character1.m[g.count] = 1
		g.character1.notes = append(g.character1.notes, Note{
			x:    getPositionInLine(1,0),
			y:    screenHeight - 20,
			line: 1,
			direction: up,
		})
	}

	if contains(keysPressed, ebiten.KeyE) {
		g.character1.m[g.count] = 2
		g.character1.notes = append(g.character1.notes, Note{
			x:    getPositionInLine(2,0),
			y:    screenHeight - 20,
			line: 2,
			direction: up,
		})
	}

	if contains(keysPressed, ebiten.KeyR) {
		g.character1.m[g.count] = 3
		g.character1.notes = append(g.character1.notes, Note{
			x:    getPositionInLine(3,0),
			y:    screenHeight - 20,
			line: 3,
			direction: up,
		})
	}
}

func checkIfNoteIsHit(g *Game, line int) {
	for i, note := range g.character1.notes {
		if note.line == line {
			if note.y+noteSize > lineMiddleY && note.y < lineMiddleY {
				g.notesToFadeAway = append(g.notesToFadeAway, NoteFadeAway{
					note:    note,
					success: true,
					count:   100,
				})

				g.character1.notes = append(g.character1.notes[:i], g.character1.notes[i+1:]...)
				g.score++

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

				return
			}
		}
	}
	p := g.audioContext.NewPlayerFromBytes([]byte{20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20})
	p.Play()
	g.missed++
}

func rewindAndPlay(p *audio.Player) {
	p.Rewind()
	p.Play()
}

func checkActionTaping(g *Game) {
	keysPressed := inpututil.PressedKeys()
	if len(keysPressed) == 0 {
		return
	}
	//key w = wait
	if contains(keysPressed, ebiten.KeyQ) {
		checkIfNoteIsHit(g, 0)
	}

	if contains(keysPressed, ebiten.KeyW) {
		checkIfNoteIsHit(g, 1)
	}

	if contains(keysPressed, ebiten.KeyE) {
		checkIfNoteIsHit(g, 2)
	}

	if contains(keysPressed, ebiten.KeyR) {
		checkIfNoteIsHit(g, 3)
	}
}
