package main

import (
	"bytes"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 512
	screenHeight = 240
	//area where we play
	layoutCharacterWidth = 200
	startLayoutC2        = 312

	noteSize = 25

	lineMiddleY      = 190
	lineMiddleMargin = 25

	introFramesNbr             = 700
	firstTypingAttackFramesNbr = 300
)

var (
	arcadeFont font.Face
)

type Character struct {
	audioCharacter  AudioCharacter
	notes           []Note
	notesToFadeAway []NoteFadeAway
	characterSprite CharacterSprite
}

type CharacterSprite struct {
	img     *ebiten.Image
	sprites map[SpriteStance]Sprite
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

	character1         Character
	character2         Character
	currentPhaseStance PhaseStance
	mapNoteToPlay      map[int]int
	notesDisplayed     int
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
		if g.count >= introFramesNbr {
			g.currentPhaseStance = firstAttackC1
			g.count = 0
		}
	case firstAttackC1:
		if len(g.mapNoteToPlay) >= 3 {
			g.currentPhaseStance = defendC2
			g.count = 0
		}
		checkActionStartAttack(g)
	case defendC2:

		if g.notesDisplayed >= len(g.mapNoteToPlay) && len(g.character2.notes) <= 0 {
			g.currentPhaseStance = addNoteC2
			break
		}
		//DEFEND
		if line, ok := g.mapNoteToPlay[g.count]; ok {
			x := getPositionInLine(line, startLayoutC2)
			g.character2.notes = append(g.character2.notes, Note{
				x:         x,
				y:         20,
				line:      line,
				direction: down,
			})
			g.notesDisplayed++
		}
		checkActionC2(g)
	case addNoteC2:
		if noteWasAdded(g, false) {
			g.currentPhaseStance = defendC1
			g.count = 0
			g.notesDisplayed = 0
		}
	case defendC1:
		if g.notesDisplayed >= len(g.mapNoteToPlay) && len(g.character1.notes) <= 0 {
			g.currentPhaseStance = addNoteC1
			break
		}
		//DEFEND
		if line, ok := g.mapNoteToPlay[g.count]; ok {
			x := getPositionInLine(line, 0)
			g.character1.notes = append(g.character1.notes, Note{
				x:         x,
				y:         20,
				line:      line,
				direction: down,
			})
			g.notesDisplayed++
		}
		checkActionC1(g)
	case addNoteC1:
		if noteWasAdded(g, true) {
			g.currentPhaseStance = defendC2
			g.count = 0
			g.notesDisplayed = 0
		}
	case c1Lost, c2Lost:
		return nil

	default:
	}

	if g.character1.updateNotesAndCheckIfLost() {
		g.currentPhaseStance = c1Lost
	}
	if g.character2.updateNotesAndCheckIfLost() {
		g.currentPhaseStance = c2Lost
	}

	return nil
}

func (c *Character) updateNotesAndCheckIfLost() bool {
	notes := c.notes
	for i := 0; i < len(notes); i++ {
		//update position
		if notes[i].direction == up {
			notes[i].y -= 1
		} else {
			notes[i].y += 1
		}

		// if out of scope, it's lost
		if (notes[i].y < 0+10) || notes[i].y > screenHeight-10 {
			if notes[i].direction == down{
				return true
			}
			notes = removeNoteAnyOrder(notes, i)
			i--
		}
	}

	for i := 0; i < len(c.notesToFadeAway); i++ {
		c.notesToFadeAway[i].count++
		
		if c.notesToFadeAway[i].count >= 100 {
			c.notesToFadeAway = removeNoteToFadeAwayAnyOrder(c.notesToFadeAway, i)
			i--
		}
	}
	c.notes = notes
	return false
}

func removeNoteAnyOrder(s []Note, i int) []Note {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func removeNoteToFadeAwayAnyOrder(s []NoteFadeAway, i int) []NoteFadeAway {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Decode an image from the image file's byte slice.
	// Now the byte slice is generated with //go:generate for Go 1.15 or older.
	// If you use Go 1.16 or newer, it is strongly recommended to use //go:embed to embed the image file.
	// See https://pkg.go.dev/embed for more details.
	img, _, err := ebitenutil.NewImageFromFile("./res/sprite_dog.png")
	if err != nil {
		log.Fatal(err)
	}
	dogImage := ebiten.NewImageFromImage(img)

	imgKnight, _, err := ebitenutil.NewImageFromFile("./res/sprite_knight.png")
	if err != nil {
		log.Fatal(err)
	}
	knightImage := ebiten.NewImageFromImage(imgKnight)

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

	dogSprites := initDogSprites()
	knightSprites := initKnightSprites()
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("TODONAME")

	audioCtx := audio.NewContext(48000)
	initWillTellOverture()
	playWillTellOvertur(audioCtx)
	audioPlayer1 := initPlayer1(audioCtx)
	audioPlayer2 := initPlayer1(audioCtx)

	if err := ebiten.RunGame(&Game{
		audioContext: audioCtx,
		count:        0,
		character1: Character{
			audioCharacter:  audioPlayer1,
			notes:           []Note{},
			notesToFadeAway: []NoteFadeAway{},
			characterSprite: CharacterSprite{
				img:     dogImage,
				sprites: dogSprites,
			},
		},
		character2: Character{
			audioCharacter:  audioPlayer2,
			notes:           []Note{},
			notesToFadeAway: []NoteFadeAway{},
			characterSprite: CharacterSprite{
				img:     knightImage,
				sprites: knightSprites,
			},
		},
		currentPhaseStance: intro,
		mapNoteToPlay:      map[int]int{},
		notesDisplayed:     0,
	}); err != nil {
		log.Fatal(err)
	}

}
func playWillTellOvertur(audioCtx *audio.Context) {
	p := getPlayer("./res/william_tell_overture_8_bit.mp3", audioCtx)
	p.Play()
}
func initPlayer1(audioCtx *audio.Context) AudioCharacter {
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
