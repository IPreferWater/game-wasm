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

func checkAction(g *Game) {
	keysPressed := inpututil.PressedKeys()
	if len(keysPressed) == 0 {
		return
	}
	//key w = wait
	if contains(keysPressed, ebiten.KeyQ) {
		g.notesTyping[g.count] = bla{
			line: 0,
		}

		g.notesUpC1 = append(g.notesUpC1, Note{
			x:    0,
			y:    screenHeight - 20,
			line: 0,
		})
	}

	if contains(keysPressed, ebiten.KeyW) {
		g.notesTyping[g.count] = bla{
			line: 1,
		}
		g.notesUpC1 = append(g.notesUpC1, Note{
			x:    0,
			y:    screenHeight - 20,
			line: 1,
		})
	}

	if contains(keysPressed, ebiten.KeyE) {
		g.notesTyping[g.count] = bla{
			line: 2,
		}
		g.notesUpC1 = append(g.notesUpC1, Note{
			x:    0,
			y:    screenHeight - 20,
			line: 2,
		})
	}

	if contains(keysPressed, ebiten.KeyR) {
		g.notesTyping[g.count] = bla{
			line: 3,
		}
		g.notesUpC1 = append(g.notesUpC1, Note{
			x:    0,
			y:    screenHeight - 20,
			line: 3,
		})
	}
}

func checkIfNoteIsHit(g *Game, line int) {
	for i, note := range g.notesUpC1 {
		if note.line == line {
			if note.y+noteSize > lineMiddleY && note.y < lineMiddleY {
				g.notesToFadeAway = append(g.notesToFadeAway, NoteFadeAway{
					note:    note,
					success: true,
					count:   100,
				})

				g.notesUpC1 = append(g.notesUpC1[:i], g.notesUpC1[i+1:]...)
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
