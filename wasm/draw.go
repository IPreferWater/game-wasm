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

	//layouts


	screen.DrawImage(c1Back, &ebiten.DrawImageOptions{})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenWidth-200, 0)
	screen.DrawImage(c2Back, op)



	ebitenutil.DrawLine(screen, 0, lineMiddleY, screenWidth, screenHeight-50, color.RGBA{200, 50, 150, 150})
	//ebitenutil.DrawLine(screen, 0, lineMiddleY-lineMiddleMargin, screenWidth, lineMiddleY-lineMiddleMargin, color.RGBA{100, 80, 150, 150})
	//ebitenutil.DrawLine(screen, 0, lineMiddleY+lineMiddleMargin, screenWidth, lineMiddleY+lineMiddleMargin, color.RGBA{220, 140, 90, 150})

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
	}
	
	drawCharacter(g.character1.characterSprite, Playing, g.frameCount, screen, screenWidth/2, screenHeight/3)
	drawCharacter(g.character2.characterSprite, Playing, g.frameCount, screen, screenWidth/2, screenHeight-screenHeight/3)

	//TODO refactor this method
	for _, note := range g.character1.notes {
		opNotes := &ebiten.DrawImageOptions{}
		opNotes.GeoM.Translate(float64(note.x), float64(note.y))
		xStart := note.line * noteSize
		subRectangle := image.Rect(xStart, 0, xStart+noteSize, noteSize)
		screen.DrawImage(notesSprite.SubImage(subRectangle).(*ebiten.Image), opNotes)
	}
	for _, note := range g.character2.notes {
		opNotes := &ebiten.DrawImageOptions{}
		opNotes.GeoM.Translate(float64(note.x), float64(note.y))
		xStart := (note.line * noteSize) + (noteSize * 4)
		subRectangle := image.Rect(xStart, 0, xStart+noteSize, noteSize)
		screen.DrawImage(notesSprite.SubImage(subRectangle).(*ebiten.Image), opNotes)
	}

	//FADEAWAY TODO
	for _, noteFadeAway := range g.character1.notesToFadeAway {
		ebitenutil.DrawRect(screen, float64(noteFadeAway.note.x), float64(noteFadeAway.note.y), noteSize, noteSize, color.RGBA{75, 205, 111, uint8(4)})
	}

	for _, noteFadeAway := range g.character2.notesToFadeAway {
		ebitenutil.DrawRect(screen, float64(noteFadeAway.note.x), float64(noteFadeAway.note.y), noteSize, noteSize, color.RGBA{75, 205, 111, uint8(6)})
	}

	//s := fmt.Sprintf("frame count : %d\n mapNoteToPlay size : %d\n c1Notes : %v", g.count, len(g.mapNoteToPlay), g.character1.notes)
	s := fmt.Sprintf("frame count : %d\n currentPhaseStance : %d\n", g.frameCount, g.currentPhaseStance)
	ebitenutil.DebugPrint(screen, s)
}

func drawBlinkingNote(screen *ebiten.Image, g *Game) {
	if g.blink {
		for i:=0; i<4;i++{
			opNotes := &ebiten.DrawImageOptions{}	
			x :=  getPositionInLine(i,0)
			
			opNotes.GeoM.Translate(float64(x), screenHeight-noteSize)
			subRectangle := image.Rect(i*noteSize, 0, (i+1)*noteSize, noteSize)
			screen.DrawImage(notesSprite.SubImage(subRectangle).(*ebiten.Image), opNotes)
		}
	}
	
}

/*func drawCoolDowns(screen *ebiten.Image, g *Game) {

	for i:=0; i<8;i++{
		opNotes := &ebiten.DrawImageOptions{}	
		x := getPositionForCoolDowns(i)
		
		opNotes.GeoM.Translate(float64(x), screenHeight-noteSize)
		subRectangle := image.Rect(i*noteSize, 0, (i+1)*noteSize, noteSize)
		screen.DrawImage(notesSprite.SubImage(subRectangle).(*ebiten.Image), opNotes)
	}
}*/

/*func getPositionForCoolDowns(noteNumber int) float32 {
	// the first 4 notes are for character 1, the 4 lasts for character 2
	if noteNumber<=3 {
		return getPositionInLine(noteNumber,0)
	}
	return getPositionInLine(noteNumber-4,startLayoutC2)
}*/

func drawIntro(screen *ebiten.Image, g *Game) {

	text.Draw(screen, "New Fight !", arcadeFont, screenWidth/2, screenHeight/4, color.White)

	if g.frameCount > 200 {
		drawCharacter(g.character1.characterSprite, Playing, g.frameCount, screen, screenWidth/2, screenHeight/4)
	}

	if g.frameCount > 300 {
		text.Draw(screen, "Versus", arcadeFont, screenWidth/2, screenHeight/2, color.White)
	}
	if g.frameCount > 400 {
		drawCharacter(g.character2.characterSprite, Playing, g.frameCount, screen, screenWidth/2, screenHeight-(screenHeight/3))
	}
}

func drawLost(screen *ebiten.Image, g *Game) {

	getWinnerLooser := func() (string, string) {
		if g.currentPhaseStance == c1Lost {
			drawCharacter(g.character1.characterSprite, Lost, g.frameCount, screen, screenWidth/2, screenHeight/4)
			drawCharacter(g.character2.characterSprite, Playing, g.frameCount, screen, screenWidth/2, screenHeight-(screenHeight/3))
			return "2", "1"
		}
		drawCharacter(g.character1.characterSprite, Playing, g.frameCount, screen, screenWidth/2, screenHeight/4)
		drawCharacter(g.character2.characterSprite, Lost, g.frameCount, screen, screenWidth/2, screenHeight-(screenHeight/3))
		return "1", "2"
	}
	winner, looser := getWinnerLooser()
	txt := fmt.Sprintf("Player %s win !\n Player %s is such a looser ...", winner, looser)
	text.Draw(screen, txt, arcadeFont, screenWidth/2, screenHeight/4, color.White)

	txtReplay := fmt.Sprintf("Tape space to replay")
	text.Draw(screen, txtReplay, arcadeFont, screenWidth/2, screenHeight/2, color.White)
}

func drawAddNote(screen *ebiten.Image, g *Game) {
	text.Draw(screen, "Add a note !!! !", arcadeFont, screenWidth/2, screenHeight/4, color.White)
}

func drawCharacter(characterSprite CharacterSprite, spriteStance SpriteStance, frameCount int, screen *ebiten.Image, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)

	sprite := characterSprite.sprites[spriteStance]
	spriteIdx := int(frameCount/sprite.changeSpriteAfterFrames) % (sprite.numberOfSprites * 2)
	if spriteIdx > sprite.numberOfSprites {
		spriteIdx = (sprite.numberOfSprites * 2) - spriteIdx - 1
	}

	x1 := sprite.width * spriteIdx
	x2 := sprite.width * (spriteIdx + 1)
	screen.DrawImage(characterSprite.img.SubImage(image.Rect(x1, sprite.yStar, x2, sprite.yStar+sprite.height)).(*ebiten.Image), op)
}
