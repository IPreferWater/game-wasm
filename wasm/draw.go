package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (g *Game) Draw(screen *ebiten.Image) {

	ebitenutil.DrawRect(screen, 2, 2, 30, 30, color.RGBA{200, 50, 150, 150})

	//layouts
	ebitenutil.DrawRect(screen, 2, 2, layoutCharacterWidth, screenHeight*0.9, ParseHexColorFast("#0074D9"))
	ebitenutil.DrawRect(screen, startLayoutC2, 2, layoutCharacterWidth, screenHeight*0.9, ParseHexColorFast("#d35400"))

	ebitenutil.DrawLine(screen, 0, lineMiddleY, screenWidth, screenHeight-50, color.RGBA{200, 50, 150, 150})
	ebitenutil.DrawLine(screen, 0, lineMiddleY-lineMiddleMargin, screenWidth, lineMiddleY-lineMiddleMargin, color.RGBA{100, 80, 150, 150})
	ebitenutil.DrawLine(screen, 0, lineMiddleY+lineMiddleMargin, screenWidth, lineMiddleY+lineMiddleMargin, color.RGBA{220, 140, 90, 150})

	if g.currentPhaseStance == intro {
		drawIntro(screen, g)
		return
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
		x := ((screenWidth/3)/4)*noteFadeAway.note.line + 20 // 20 as layout
		ebitenutil.DrawRect(screen, float64(x), float64(noteFadeAway.note.y), noteSize, noteSize, color.RGBA{75, 205, 111, uint8(noteFadeAway.count)})
	}

	for _, noteFadeAway := range g.character2.notesToFadeAway {
		x := ((screenWidth/3)/4)*noteFadeAway.note.line + 20 // 20 as layout
		ebitenutil.DrawRect(screen, float64(x), float64(noteFadeAway.note.y), noteSize, noteSize, color.RGBA{75, 205, 111, uint8(noteFadeAway.count)})
	}

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