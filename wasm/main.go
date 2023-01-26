package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 240

	frameCount = 8

	noteSize = 30

	lineMiddleY      = 190
	lineMiddleMargin = 25
)

var (
	runnerImage *ebiten.Image
	dogSprites  map[SpriteStance]Sprite
)

type Game struct {
	audioContext *audio.Context
	count           int
	notes           []Note
	notesToFadeAway []NoteFadeAway
	typing          bool
	missed          int
	score           int
}

type NoteFadeAway struct {
	note    Note
	success bool
	count   int
}

type Note struct {
	x    float32
	y    float32
	line int
}

func (g *Game) Update() error {

	secondUpdate := false
	if g.count%60 == 0 {
		secondUpdate = true
	}

	g.count++

	// turn to play the notes
	//TODO typing is not a good word
	if g.typing {
		checkAction(g)
		for i := 0; i < len(g.notes); i++ {
			g.notes[i].y -= 1

			if g.notes[i].y < 0+20 { // 20 as layout
				g.notes = append(g.notes[:i], g.notes[i+1:]...)
				i--
			}
		}

		return nil
	}

	// turn to defeat the notes
	if secondUpdate {
		g.notes = append(g.notes, Note{
			x:    0,
			y:    20,
			line: rand.Intn(4),
		})
	}

	checkActionTaping(g)
	for i := 0; i < len(g.notes); i++ {
		g.notes[i].y += 1
	}

	for i := 0; i < len(g.notesToFadeAway); i++ {
		g.notesToFadeAway[i].count -= 1

		if g.notesToFadeAway[i].count <= 0 {
			g.notesToFadeAway = append(g.notesToFadeAway[:i], g.notesToFadeAway[i+1:]...)
			i--
		}
	}

	return nil
}

func drawCharacter(sprite Sprite, frameCount int, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(sprite.width)/2, -float64(sprite.height)/2)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	//i := (sprite.frameCount / 5) %

	//sx, sy := sprite.width+i*frameWidth, frameOY
	spriteIdx := int(frameCount/sprite.changeSpriteAfterFrames) % (sprite.numberOfSprites * 2)
	if spriteIdx > sprite.numberOfSprites {
		spriteIdx = (sprite.numberOfSprites * 2) - spriteIdx - 1
	}

	//fmt.Printf("x1 %d y1 %d x2 %d y2 %d\n",sprite.width*spriteNumber, 0,sprite.width*(spriteNumber+1), sprite.height)

	//sx := (fff/sprite.width)%float64(sprite.spriteNumber)
	x1 := sprite.width * spriteIdx
	x2 := sprite.width * (spriteIdx + 1)
	s := fmt.Sprintf("spriteIdx %d \n x1 %d x2 %d", spriteIdx, x1, x2)
	ebitenutil.DebugPrint(screen, s)
	screen.DrawImage(runnerImage.SubImage(image.Rect(x1, sprite.yStar, x2, sprite.yStar+sprite.height)).(*ebiten.Image), op)
}
func (g *Game) Draw(screen *ebiten.Image) {

	ebitenutil.DrawRect(screen, 2, 2, 30, 30, color.RGBA{200, 50, 150, 150})
	ebitenutil.DrawLine(screen, 0, lineMiddleY, screenWidth, screenHeight-50, color.RGBA{200, 50, 150, 150})
	ebitenutil.DrawLine(screen, 0, lineMiddleY-lineMiddleMargin, screenWidth, lineMiddleY-lineMiddleMargin, color.RGBA{100, 80, 150, 150})
	ebitenutil.DrawLine(screen, 0, lineMiddleY+lineMiddleMargin, screenWidth, lineMiddleY+lineMiddleMargin, color.RGBA{220, 140, 90, 150})

	drawCharacter(dogSprites[Playing], g.count, screen)

	for _, note := range g.notes {
		x := (screenWidth/4)*note.line + 20 // 20 as layout
		if g.typing {
			ebitenutil.DrawRect(screen, float64(x), float64(note.y), 30, 30, color.RGBA{200, 50, 150, 150})
		} else {

			ebitenutil.DrawRect(screen, float64(x), float64(note.y), noteSize, noteSize, color.NRGBA{250, 177, 160, 200})
		}

	}

	for _, noteFadeAway := range g.notesToFadeAway {
		x := (screenWidth/4)*noteFadeAway.note.line + 20 // 20 as layout

		ebitenutil.DrawRect(screen, float64(x), float64(noteFadeAway.note.y), noteSize, noteSize, color.RGBA{75, 205, 111, uint8(noteFadeAway.count)})
	}

	// 40 widht
	// 10 sprite
	// tous les 50
	if !g.typing {
		t := (g.count / 20) % 20
		if t > 9 {
			t = 20 - t - 1
		}
		/*s := fmt.Sprintf("score : %d\nmissed : %d\n frame count : %d\n test : %d", g.score, g.missed, g.count, t)
		ebitenutil.DebugPrint(screen, s)*/
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

/*func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}*/

func main() {
	// Decode an image from the image file's byte slice.
	// Now the byte slice is generated with //go:generate for Go 1.15 or older.
	// If you use Go 1.16 or newer, it is strongly recommended to use //go:embed to embed the image file.
	// See https://pkg.go.dev/embed for more details.
	img, _, err := ebitenutil.NewImageFromFile("./res/sprite_dog.png")
	if err != nil {
		log.Fatal(err)
	}
	runnerImage = ebiten.NewImageFromImage(img)
	dogSprites = initDogSprites()
	fmt.Println("allo?")
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebiten Demo)")
	if err := ebiten.RunGame(&Game{
		audioContext:    audio.NewContext(48000),
		count:           frameCount,
		notes:           []Note{},
		notesToFadeAway: []NoteFadeAway{},
		typing:          false,
		missed:          0,
		score:           0,
	}); err != nil {
		log.Fatal(err)
	}

}
