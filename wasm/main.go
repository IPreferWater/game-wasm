package main

import (
	"bytes"
	"fmt"
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
	frameCount   int

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

func getEbitenImageFromRes(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return img
}
func main() {

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("TODONAME")

	dogImage := getEbitenImageFromRes("./res/sprite_dog.png")
	knightImage := getEbitenImageFromRes("./res/sprite_knight.png")

	c1Back = getEbitenImageFromRes("./res/dog/back.png")
	c2Back = getEbitenImageFromRes("./res/knight/back.png")
	notesSprite = getEbitenImageFromRes("./res/note_sprite.png")

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

	audioCtx := audio.NewContext(48000)
	initWillTellOverture()
	audioPlayer1 := initAudioCharacter(audioCtx, "dog")
	audioPlayer2 := initAudioCharacter(audioCtx, "knight")

	if err := ebiten.RunGame(&Game{
		audioContext:       audioCtx,
		frameCount:         700,
		character1:         initNewCharacter(audioPlayer1, dogImage, dogSprites),
		character2:         initNewCharacter(audioPlayer2, knightImage, knightSprites),
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
