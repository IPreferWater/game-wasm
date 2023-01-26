package main

import (
	"github.com/hajimehoshi/ebiten/v2"
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
		g.notes = append(g.notes, Note{
			x:    0,
			y:    screenHeight-20,
			line: 0,
		})
	}

	if contains(keysPressed, ebiten.KeyW) {
		g.notes = append(g.notes, Note{
			x:    0,
			y:    screenHeight-20,
			line: 1,
		})
	}

	if contains(keysPressed, ebiten.KeyE) {
		g.notes = append(g.notes, Note{
			x:    0,
			y:    screenHeight-20,
			line: 2,
		})
	}

	if contains(keysPressed, ebiten.KeyR) {
		g.notes = append(g.notes, Note{
			x:    0,
			y:    screenHeight-20,
			line: 3,
		})
	}
}

func checkIfNoteIsHit(g *Game, line int){
	for i, note := range g.notes {
		if note.line == line{
			if note.y+noteSize > lineMiddleY && note.y<lineMiddleY {
				g.notesToFadeAway = append(g.notesToFadeAway, NoteFadeAway{
					note:    note,
					success: true,
					count: 100,
				})
				g.notes = append(g.notes[:i], g.notes[i+1:]...)
				g.score++
				return
			}
		}
	}

	g.missed++
}

func checkActionTaping(g *Game) {
	keysPressed := inpututil.PressedKeys()
	if len(keysPressed) == 0 {
		return
	}
	//key w = wait
	if contains(keysPressed, ebiten.KeyQ) {
		checkIfNoteIsHit(g,0)
	}

	if contains(keysPressed, ebiten.KeyW) {
		checkIfNoteIsHit(g,1)
	}

	if contains(keysPressed, ebiten.KeyE) {
		checkIfNoteIsHit(g,2)
	}

	if contains(keysPressed, ebiten.KeyR) {
		checkIfNoteIsHit(g,3)
	}
}
