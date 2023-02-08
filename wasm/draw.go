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
	//ebitenutil.DrawRect(screen, 2, 2, layoutCharacterWidth, screenHeight*0.9, ParseHexColorFast("#0074D9"))
	
	screen.DrawImage(c1Back, &ebiten.DrawImageOptions{})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenWidth-200, 0)
	screen.DrawImage(c2Back, op)

	//ebitenutil.DrawRect(screen, 2, 2, layoutCharacterWidth, screenHeight*0.9, ParseHexColorFast("#d35400"))

	ebitenutil.DrawLine(screen, 0, lineMiddleY, screenWidth, screenHeight-50, color.RGBA{200, 50, 150, 150})
	ebitenutil.DrawLine(screen, 0, lineMiddleY-lineMiddleMargin, screenWidth, lineMiddleY-lineMiddleMargin, color.RGBA{100, 80, 150, 150})
	ebitenutil.DrawLine(screen, 0, lineMiddleY+lineMiddleMargin, screenWidth, lineMiddleY+lineMiddleMargin, color.RGBA{220, 140, 90, 150})

	if g.currentPhaseStance == intro {
		drawIntro(screen, g)
		return
	}

	if g.currentPhaseStance == c1Lost || g.currentPhaseStance == c2Lost {
		drawLost(screen,g)
		return
	}

	if g.currentPhaseStance == addNoteC2 || g.currentPhaseStance == addNoteC1 {
		drawAddNote(screen, g)
	}

	drawCharacter(g.character1.characterSprite, Playing, g.count, screen, screenWidth/2, screenHeight/3)
	drawCharacter(g.character2.characterSprite, Playing, g.count, screen, screenWidth/2, screenHeight-screenHeight/3)

	for _, note := range g.character1.notes {
		ebitenutil.DrawRect(screen, float64(note.x), float64(note.y), noteSize, noteSize, ParseHexColorFast("#10ac84"))
	}
	for _, note := range g.character2.notes {
		ebitenutil.DrawRect(screen, float64(note.x), float64(note.y), noteSize, noteSize, ParseHexColorFast("#f368e0"))
	}

	for _, noteFadeAway := range g.character1.notesToFadeAway {
		ebitenutil.DrawRect(screen, float64(noteFadeAway.note.x), float64(noteFadeAway.note.y), noteSize, noteSize, color.RGBA{75, 205, 111, uint8(4)})
	}

	for _, noteFadeAway := range g.character2.notesToFadeAway {
		ebitenutil.DrawRect(screen, float64(noteFadeAway.note.x), float64(noteFadeAway.note.y), noteSize, noteSize, color.RGBA{75, 205, 111, uint8(6)})
	}

	//s := fmt.Sprintf("frame count : %d\n mapNoteToPlay size : %d\n c1Notes : %v", g.count, len(g.mapNoteToPlay), g.character1.notes)
	s := fmt.Sprintf("frame count : %d\n currentPhaseStance : %d\n", g.count, g.currentPhaseStance)
	ebitenutil.DebugPrint(screen, s)
	// 40 widht
	// 10 sprite
	// tous les 50
	/*if !g.typing {
		t := (g.count / 20) % 20
		if t > 9 {
			t = 20 - t - 1
		}
		s := fmt.Sprintf("score : %d\nmissed : %d\n frame count : %d\n test : %d", g.score, g.missed, g.count, t)
		ebitenutil.DebugPrint(screen, s)
	}*/
}

func drawIntro(screen *ebiten.Image, g *Game) {

	text.Draw(screen, "New Fight !", arcadeFont, screenWidth/2, screenHeight/4, color.White)

	if g.count > 200 {
		drawCharacter(g.character1.characterSprite, Playing, g.count, screen, screenWidth/2, screenHeight/4)
	}

	if g.count > 300 {
		text.Draw(screen, "Versus", arcadeFont, screenWidth/2, screenHeight/2, color.White)
	}
	if g.count > 400 {
		drawCharacter(g.character2.characterSprite, Playing, g.count, screen, screenWidth/2, screenHeight-(screenHeight/3))
	}
}

func drawLost(screen *ebiten.Image, g *Game) {

	getWinnerLooser := func() (string, string) {
		if g.currentPhaseStance == c1Lost {
			drawCharacter(g.character1.characterSprite, Lost, g.count, screen, screenWidth/2, screenHeight/4)
			drawCharacter(g.character2.characterSprite, Playing, g.count, screen, screenWidth/2, screenHeight-(screenHeight/3))
			return "2", "1"
		}
		drawCharacter(g.character1.characterSprite, Playing, g.count, screen, screenWidth/2, screenHeight/4)
		drawCharacter(g.character2.characterSprite, Lost, g.count, screen, screenWidth/2, screenHeight-(screenHeight/3))
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
