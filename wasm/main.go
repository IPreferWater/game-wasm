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

	// Where we display the playing are of 1 character
	layoutCharacterWidth = 200
	// Where the x playing area of character 2 si starting
	startLayoutC2 = 312

	// Size of the square in pixel of 1 note
	noteSize = 25

	// Where the note should be typed on the playing area
	lineMiddleY = 190
	// Hitbox for the notes
	lineMiddleMargin = 25

	// How many frame for the introduction stance
	introFramesNbr = 700

	// How many frame we wait before be able to type the same note on attack stance
	coolDownFrameForSameNote = 40
)

var (
	arcadeFont  font.Face
	c1Back      *ebiten.Image
	c2Back      *ebiten.Image
	notesSprite *ebiten.Image
)

type Game struct {
	audioContext *audio.Context
	count        int

	character1         Character
	character2         Character
	currentPhaseStance PhaseStance
	mapNoteToPlay      map[int]int
	notesDisplayed     int
	williamTellPlayer  *audio.Player
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

	imgC1Back, _, err := ebitenutil.NewImageFromFile("./res/dog/back.png")
	if err != nil {
		log.Fatal(err)
	}
	c1Back = imgC1Back

	imgC2Back, _, err := ebitenutil.NewImageFromFile("./res/knight/back.png")
	if err != nil {
		log.Fatal(err)
	}
	c2Back = imgC2Back

	imgNoteSprite, _, err := ebitenutil.NewImageFromFile("./res/note_sprite.png")
	if err != nil {
		log.Fatal(err)
	}
	notesSprite = ebiten.NewImageFromImage(imgNoteSprite)

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
	//playWillTellOvertur(audioCtx)
	audioPlayer1 := initAudioCharacter(audioCtx, "dog")
	audioPlayer2 := initAudioCharacter(audioCtx, "knight")

	if err := ebiten.RunGame(&Game{
		audioContext: audioCtx,
		count:        700,
		character1: Character{
			audioCharacter:  audioPlayer1,
			notes:           []Note{},
			notesToFadeAway: []NoteFadeAway{},
			characterSprite: CharacterSprite{
				img:     dogImage,
				sprites: dogSprites,
			},
			cooldown: Cooldown{
				line1: -coolDownFrameForSameNote,
				line2: -coolDownFrameForSameNote,
				line3: -coolDownFrameForSameNote,
				line4: -coolDownFrameForSameNote,
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
			cooldown: Cooldown{
				line1: 0,
				line2: 0,
				line3: 0,
				line4: 0,
			},
		},
		currentPhaseStance: intro,
		mapNoteToPlay:      map[int]int{},
		notesDisplayed:     0,
		williamTellPlayer:  getPlayer("./res/william_tell_overture_8_bit.mp3", audioCtx),
	}); err != nil {
		log.Fatal(err)
	}

}

func initAudioCharacter(audioCtx *audio.Context, folderName string) AudioCharacter {
	mapPath := make(map[int]string)

	for i := 0; i <= 3; i++ {
		path := fmt.Sprintf("./res/%s/sound_%d.mp3", folderName, i)
		mapPath[i] = path
	}

	return AudioCharacter{
		sound0: getPlayer(mapPath[0], audioCtx),
		sound1: getPlayer(mapPath[1], audioCtx),
		sound2: getPlayer(mapPath[2], audioCtx),
		sound3: getPlayer(mapPath[3], audioCtx),
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
