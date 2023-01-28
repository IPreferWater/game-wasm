package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth          = 512
	screenHeight         = 240
	layoutCharacterWidth = 200
	startLayoutC2        = 312

	frameCount = 8

	noteSize = 25

	lineMiddleY      = 190
	lineMiddleMargin = 25
)

var (
	runnerImage *ebiten.Image
	dogSprites  map[SpriteStance]Sprite
	arcadeFont  font.Face
)

type Character struct {
	audioCharacter AudioCharacter
	notes          []Note
	m              map[int]bla
}

type AudioCharacter struct {
	sound0 *audio.Player
	sound1 *audio.Player
	sound2 *audio.Player
	sound3 *audio.Player
}

type Game struct {
	audioContext *audio.Context
	count        int
	//notesUpC1          []Note
	notesDownC2        []Note
	notesToFadeAway    []NoteFadeAway
	typing             bool
	missed             int
	score              int
	character1         Character
	character2         Character
	phase              phase
	currentPhaseStance PhaseStance
	//notesTyping        map[int]bla
}

type NoteFadeAway struct {
	note    Note
	success bool
	count   int
}

type Note struct {
	x         float32
	y         float32
	line      int
	direction direction
}

type direction int64

const (
	up direction = iota
	down
)

func (g *Game) Update() error {

	/*secondUpdate := false
	if g.count%20 == 0 {
		secondUpdate = true
	}*/

	g.count++

	switch g.currentPhaseStance {
	case intro:
		if g.count >= g.phase.introFramesNbr {
			g.currentPhaseStance = firstAttackC1
			g.count = 0
		}
	case firstAttackC1:
		if len(g.character1.m) >= 3 {
			g.currentPhaseStance = defendC2
			g.character2.m = g.character1.m
			g.count = 0
		}
		checkAction(g)
	case defendC2:
		if val, ok := g.character2.m[g.count]; ok {
			x := getPositionInLine(val.line, startLayoutC2)
			g.character2.notes = append(g.character2.notes, Note{
				x:         x,
				y:         20,
				line:      val.line,
				direction: down,
			})
		}
		checkAction(g)

	default:
	}

	for i := 0; i < len(g.character1.notes); i++ {
		if g.character1.notes[i].direction == up {
			g.character1.notes[i].y -= 1
		} else {
			g.character1.notes[i].y += 1
		}

		if g.character1.notes[i].y < 0+10 { // 20 as layout
			g.character1.notes = append(g.character1.notes[:i], g.character1.notes[i+1:]...)
			i--
		}
	}

	//TODO factorize
	for i := 0; i < len(g.character2.notes); i++ {
		if g.character2.notes[i].direction == up {
			g.character2.notes[i].y -= 1
		} else {
			g.character2.notes[i].y += 1
		}

		if g.character2.notes[i].y < 0+10 { // 20 as layout
			g.character2.notes = append(g.character2.notes[:i], g.character2.notes[i+1:]...)
			i--
		}
	}

	/*if g.currentPhaseStance == attackC1 {
		if len(g.character1.m) >= 3 {
			g.currentPhaseStance = defendC1
		}
		checkAction(g)
		for i := 0; i < len(g.notesUpC1); i++ {
			g.notesUpC1[i].y -= 1

			if g.notesUpC1[i].y < 0+20 { // 20 as layout
				g.notesUpC1 = append(g.notesUpC1[:i], g.notesUpC1[i+1:]...)
				i--
			}
		}
	}*/

	// turn to defeat the notes
	/*if secondUpdate {
		g.notes = append(g.notes, Note{
			x:    0,
			y:    20,
			line: rand.Intn(4),
		})
	}*/
	/*if val, ok := williamTellOverture[g.count]; ok {
		g.notes = append(g.notes, Note{
			x:    0,
			y:    20,
			line: val.line,
		})
	}*/

	/*checkActionTaping(g)
	for i := 0; i < len(g.notesUpC1); i++ {
		g.notesUpC1[i].y += 1
	}

	for i := 0; i < len(g.notesToFadeAway); i++ {
		g.notesToFadeAway[i].count -= 1

		if g.notesToFadeAway[i].count <= 0 {
			g.notesToFadeAway = append(g.notesToFadeAway[:i], g.notesToFadeAway[i+1:]...)
			i--
		}
	}*/

	return nil
}
func drawIntro(screen *ebiten.Image, g *Game) {

	text.Draw(screen, "New Fight !", arcadeFont, screenWidth/2, screenHeight/4, color.White)

	if g.count > 200 {
		drawCharacter(dogSprites[Playing], g.count, screen, screenWidth/2, screenHeight/4)
	}

	if g.count > 300 {
		text.Draw(screen, "Versus", arcadeFont, screenWidth/2, screenHeight/2, color.White)
	}

	if g.count > 400 {
		drawCharacter(dogSprites[Playing], g.count, screen, screenWidth/2, screenHeight-(screenHeight/3))
	}
}

func drawCharacter(sprite Sprite, frameCount int, screen *ebiten.Image, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)

	spriteIdx := int(frameCount/sprite.changeSpriteAfterFrames) % (sprite.numberOfSprites * 2)
	if spriteIdx > sprite.numberOfSprites {
		spriteIdx = (sprite.numberOfSprites * 2) - spriteIdx - 1
	}

	x1 := sprite.width * spriteIdx
	x2 := sprite.width * (spriteIdx + 1)
	screen.DrawImage(runnerImage.SubImage(image.Rect(x1, sprite.yStar, x2, sprite.yStar+sprite.height)).(*ebiten.Image), op)
}

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

	drawCharacter(dogSprites[Playing], g.count, screen, screenWidth/2, screenHeight/3)

	if g.currentPhaseStance == firstAttackC1 || g.currentPhaseStance == attackC1 {

	}

	for _, note := range g.character1.notes {
		//x := ((screenWidth/3)/4)*note.line + 20 // 20 as layout
		ebitenutil.DrawRect(screen, float64(note.x), float64(note.y), noteSize, noteSize, ParseHexColorFast("#10ac84"))
	}
	for _, note := range g.character2.notes {

		ebitenutil.DrawRect(screen, float64(note.x), float64(note.y), noteSize, noteSize, ParseHexColorFast("#f368e0"))
	}

	/*for _, note := range g.notesDownC2 {
		x := ((screenWidth-(screenWidth/3))/4)*note.line + 20 // 20 as layout
		ebitenutil.DrawRect(screen, float64(x), float64(note.y), noteSize, noteSize, ParseHexColorFast("#192a56"))
	}*/

	for _, noteFadeAway := range g.notesToFadeAway {
		x := ((screenWidth/3)/4)*noteFadeAway.note.line + 20 // 20 as layout
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
		s := fmt.Sprintf("score : %d\nmissed : %d\n frame count : %d\n test : %d", g.score, g.missed, g.count, t)
		ebitenutil.DebugPrint(screen, s)
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

	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const (
		arcadeFontSize = 8
		dpi            = 72
	)
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    arcadeFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	runnerImage = ebiten.NewImageFromImage(img)
	dogSprites = initDogSprites()
	fmt.Println("allo?")
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebiten Demo)")

	audioCtx := audio.NewContext(48000)
	initWillTellOverture()
	//playWillTellOvertur(audioCtx)
	player1 := initPlayer1(audioCtx)
	player2 := initPlayer1(audioCtx)

	if err := ebiten.RunGame(&Game{
		audioContext: audioCtx,
		count:        600,
		//notesUpC1:       []Note{},
		notesToFadeAway: []NoteFadeAway{},
		typing:          false,
		missed:          0,
		score:           0,
		character1: Character{
			audioCharacter: player1,
			notes:          []Note{},
			m:              map[int]bla{},
		},
		character2: Character{
			audioCharacter: player2,
			notes:          []Note{},
			m:              map[int]bla{},
		},
		phase:              phase{introFramesNbr: 700, firstTypingAttackFramesNbr: 300},
		currentPhaseStance: intro,
		//notesTyping:        map[int]bla{},
	}); err != nil {
		log.Fatal(err)
	}

}
func playWillTellOvertur(audioCtx *audio.Context) {
	p := getPlayer("./res/william_tell_overture_8_bit.mp3", audioCtx)
	p.Play()
}
func initPlayer1(audioCtx *audio.Context) AudioCharacter {
	//sound TODO https://ebitengine.org/en/examples/audio.html

	return AudioCharacter{
		sound0: getPlayer("./res/bark_0.mp3", audioCtx),
		sound1: getPlayer("./res/bark_1.mp3", audioCtx),
		sound2: getPlayer("./res/bark_2.mp3", audioCtx),
		sound3: getPlayer("./res/bark_3.mp3", audioCtx),
	}
}

func getPlayer(fileName string, audioCtx *audio.Context) *audio.Player {
	b, err := os.ReadFile(fileName) // just pass the file name
	if err != nil {
		fmt.Print("readfile")
		panic(err)
	}
	s, err := mp3.DecodeWithoutResampling(bytes.NewReader(b))
	if err != nil {
		fmt.Print("decode")
		panic(err)
	}
	p, err := audioCtx.NewPlayer(s)
	if err != nil {
		fmt.Print("new player")
		panic(err)
	}
	return p
}

func ParseHexColorFast(s string) (c color.RGBA) {
	c.A = 0xff

	if s[0] != '#' {
		return c
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	}
	return
}
