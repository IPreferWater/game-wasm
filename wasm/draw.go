package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (g *Game) Draw(screen *ebiten.Image) {
	if g.currentPhaseStance == gameError {
		text.Draw(screen, "error", nil, screenWidth/2, screenHeight/2, color.White)
		return
	}
	drawBackground(screen)

	if g.currentPhaseStance == intro {
		drawIntro(screen, g)
		return
	}

	if g.currentPhaseStance == c1Lost || g.currentPhaseStance == c2Lost {
		drawLost(screen, g)
		return
	}

	if g.currentPhaseStance == addNoteC2 || g.currentPhaseStance == addNoteC1 {
		drawAddNote(screen, g)
	}

	if g.currentPhaseStance == firstAttackC1 {
		drawBlinkingNote(screen, g)
	} else {
		drawCoolDowns(screen, g)
	}

	drawCharacter(g.character1.characterSprite, Playing, g.frameCount, screen)
	drawCharacter(g.character2.characterSprite, Playing, g.frameCount, screen)

	drawNotes(screen, g.character1.notes, true)
	drawNotes(screen, g.character2.notes, false)

	//FADEAWAY TODO
	for _, noteFadeAway := range g.character1.notesToFadeAway {
		ebitenutil.DrawRect(screen, float64(noteFadeAway.note.x), float64(noteFadeAway.note.y), noteSize, noteSize, color.RGBA{75, 205, 111, uint8(4)})
	}

	for _, noteFadeAway := range g.character2.notesToFadeAway {
		ebitenutil.DrawRect(screen, float64(noteFadeAway.note.x), float64(noteFadeAway.note.y), noteSize, noteSize, color.RGBA{75, 205, 111, uint8(6)})
	}
	
	//s := fmt.Sprintf("frame count : %d\n mapNoteToPlay size : %d\n c1Notes : %v", g.count, len(g.mapNoteToPlay), g.character1.notes)
	//ebitenutil.DebugPrint(screen, s)
}

func drawBackground(screen *ebiten.Image) {
	screen.DrawImage(dogBack, &ebiten.DrawImageOptions{})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenWidth-200, 0)
	screen.DrawImage(knightBack, op)

	ebitenutil.DrawLine(screen, 0, lineMiddleY, screenWidth, screenHeight-50, color.RGBA{200, 50, 150, 150})
}

func drawNotes(screen *ebiten.Image, notes []Note, isC1 bool) {

	xStartSprite := 0
	if !isC1 {
		xStartSprite += noteSize * 4
	}

	for _, note := range notes {
		opNotes := &ebiten.DrawImageOptions{}
		opNotes.GeoM.Translate(float64(note.x), float64(note.y))

		xStart := (note.line * noteSize) + xStartSprite
		subRectangle := image.Rect(xStart, 0, xStart+noteSize, noteSize)
		screen.DrawImage(notesSprite.SubImage(subRectangle).(*ebiten.Image), opNotes)
	}

}

func drawBlinkingNote(screen *ebiten.Image, g *Game) {
	if g.blink {
		for i := 0; i < 4; i++ {
			opNotes := &ebiten.DrawImageOptions{}
			x := getPositionInLine(i, 0)

			opNotes.GeoM.Translate(float64(x), screenHeight-noteSize)
			subRectangle := image.Rect(i*noteSize, 0, (i+1)*noteSize, noteSize)
			screen.DrawImage(notesSprite.SubImage(subRectangle).(*ebiten.Image), opNotes)
		}
	}

}
// get correct CD depending on the line 
func getCooldownInPercentage(line int, cd Cooldown, currentFrame int) float64 {
	switch line {
	case 0:
		return calculCooldownInPercentage(cd.line1, currentFrame)
	case 1:
		return calculCooldownInPercentage(cd.line2, currentFrame)
	case 2:
		return calculCooldownInPercentage(cd.line3, currentFrame)
	case 3:
		return calculCooldownInPercentage(cd.line4, currentFrame)
	}
	return 1
}

// calcul in % the cd
func calculCooldownInPercentage(cd int, currentFrame int) float64 {
	if cd+coolDownFrameForSameNote < currentFrame {
		return 1
	}

	did := currentFrame - cd
	return float64((did*100)/coolDownFrameForSameNote) / 100
}

func drawCoolDowns(screen *ebiten.Image, g *Game) {
	percentageOfCD := func(i int) float64 {
		if i < 4 {
			return getCooldownInPercentage(i, g.character1.cooldown, g.frameCount)
		}
		return getCooldownInPercentage(i-4, g.character2.cooldown, g.frameCount)
	}

	for i := 0; i < 8; i++ {
		opNotes := &ebiten.DrawImageOptions{}
		x := getPositionForCoolDowns(i)

		opNotes.GeoM.Translate(float64(x), screenHeight-noteSize)
		yToPut := noteSize * percentageOfCD(i)
		subRectangle := image.Rect(i*noteSize, 0, (i+1)*noteSize, int(yToPut))

		screen.DrawImage(notesSprite.SubImage(subRectangle).(*ebiten.Image), opNotes)
	}
}

func getPositionForCoolDowns(noteNumber int) float32 {
	// the first 4 notes are for character 1, the 4 lasts for character 2
	if noteNumber <= 3 {
		return getPositionInLine(noteNumber, 0)
	}
	return getPositionInLine(noteNumber-4, startLayoutC2)
}

func drawIntro(screen *ebiten.Image, g *Game) {

	text.Draw(screen, "New Fight !", arcadeFont, xMiddleTxt, yTopText, color.White)

	if g.frameCount > 200 {
		drawCharacter(g.character1.characterSprite, Playing, g.frameCount, screen)
	}

	if g.frameCount > 300 {
		text.Draw(screen, "Versus", arcadeFont, xMiddleTxt, yMiddleTxt, color.White)
	}
	if g.frameCount > 400 {
		drawCharacter(g.character2.characterSprite, Playing, g.frameCount, screen)
	}
}

func drawLost(screen *ebiten.Image, g *Game) {

	if g.currentPhaseStance == c1Lost {
		drawCharacter(g.character1.characterSprite, Lost, g.frameCount, screen)
		drawCharacter(g.character2.characterSprite, Playing, g.frameCount, screen)
		text.Draw(screen, "Knight won !", arcadeFont, xMiddleTxt, yKnightSprite+70, color.White)
	} else {
		drawCharacter(g.character1.characterSprite, Playing, g.frameCount, screen)
		drawCharacter(g.character2.characterSprite, Lost, g.frameCount, screen)
		text.Draw(screen, fmt.Sprintf("Dog won !\nWouaf !"), arcadeFont, xMiddleTxt, yMiddleTxt, color.White)
	}

	txtReplay := fmt.Sprintf("Tape space\nto replay")
	text.Draw(screen, txtReplay, arcadeFont, xMiddleTxt, yTopText, color.White)
}

func drawAddNote(screen *ebiten.Image, g *Game) {
	text.Draw(screen, "Add a note !", arcadeFont, xMiddleTxt, yMiddleTxt, color.White)
}

func drawCharacter(characterSprite CharacterSprite, spriteStance SpriteStance, frameCount int, screen *ebiten.Image) {
	sprite := characterSprite.sprites[spriteStance]

	op := &ebiten.DrawImageOptions{}
	xMidle := (screenWidth / 2) - (sprite.width / 2)
	op.GeoM.Translate(float64(xMidle), characterSprite.ySprite)

	spriteIdx := int(frameCount/sprite.changeSpriteAfterFrames) % (sprite.numberOfSprites * 2)
	if spriteIdx > sprite.numberOfSprites {
		spriteIdx = (sprite.numberOfSprites * 2) - spriteIdx - 1
	}

	x1 := sprite.width * spriteIdx
	x2 := sprite.width * (spriteIdx + 1)
	screen.DrawImage(characterSprite.img.SubImage(image.Rect(x1, sprite.yStar, x2, sprite.yStar+sprite.height)).(*ebiten.Image), op)
}
