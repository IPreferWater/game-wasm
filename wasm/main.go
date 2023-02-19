package main

import (
	"bytes"
	"fmt"
	_ "image/png"
	"io"
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
	xMiddleTxt = screenWidth-startLayoutC2 + 10
	yMiddleTxt = 90
	yTopText = 24

	// Size of the square in pixel of 1 note
	noteSize = 25

	// Where the note should be typed on the playing area
	lineMiddleY = 190
	// Hitbox for the notes
	lineMiddleMargin = 25

	// How many frame for the introduction stance
	introFramesNbr = 700
	blinkFrameNbr  = 35

	// How many frame we wait before be able to type the same note on attack stance
	coolDownFrameForSameNote = 40

	yDogSprite    = 28
	yKnightSprite = 100
)

var (
	arcadeFont  font.Face
	c1Back      *ebiten.Image
	c2Back      *ebiten.Image
	notesSprite *ebiten.Image
	errors      []string
)

type Game struct {
	audioContext *audio.Context
	frameCount   int

	character1         Character
	character2         Character
	currentPhaseStance PhaseStance
	mapNoteToPlay      map[int]int
	notesDisplayed     int
	blink              bool
	blinkCount         int
	williamTellPlayer  *audio.Player
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func getEbitenImageFromRes(path string) (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil, err
	}

	return img, nil
}
func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("IPreferWater Bubble Game")

	dogImage, err := getEbitenImageFromRes("./res/dog/sprite.png")
	if err != nil {
		addError(err)
	}

	knightImage, err := getEbitenImageFromRes("./res/knight/sprite.png")
	if err != nil {
		addError(err)
	}

	c1Back, err = getEbitenImageFromRes("./res/dog/back.png")
	if err != nil {
		addError(err)
	}

	c2Back, err = getEbitenImageFromRes("./res/knight/back.png")
	if err != nil {
		addError(err)
	}

	notesSprite, err = getEbitenImageFromRes("./res/note_sprite.png")
	if err != nil {
		addError(err)
	}

	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		addError(err)
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
	audioPlayer2:= initAudioCharacter(audioCtx, "knight")

	williamTellPlayer :=  newFuncNewPlayer(william_tell_overture_8_bit, audioCtx)

	if len(errors) > 0 {
		fmt.Println("yes?")
		for i, err := range errors {
			fmt.Printf("error %d => %s\n", i, err)
		}
		fmt.Println(errors)
		if err := ebiten.RunGame(&Game{
			currentPhaseStance: gameError,
		}); err != nil {
			fmt.Println(err)
		}
	}

	if err := ebiten.RunGame(&Game{
		audioContext:       audioCtx,
		frameCount:         0,
		character1:         initNewCharacter(audioPlayer1, dogImage, dogSprites, yDogSprite),
		character2:         initNewCharacter(audioPlayer2, knightImage, knightSprites, yKnightSprite),
		currentPhaseStance: intro,
		mapNoteToPlay:      map[int]int{},
		notesDisplayed:     0,
		williamTellPlayer:  williamTellPlayer,
		blink:              false,
		blinkCount:         0,
	}); err != nil {
		fmt.Println(err)
	}

}

func addError(err error) {
	errors = append(errors, err.Error())
}
func initAudioCharacter(audioCtx *audio.Context, characterName string) AudioCharacter {

	if characterName == "dog" {

		return AudioCharacter{
			sound0: newFuncNewPlayer(dog_sound_0, audioCtx),
			sound1: newFuncNewPlayer(dog_sound_1, audioCtx),
			sound2: newFuncNewPlayer(dog_sound_2, audioCtx),
			sound3: newFuncNewPlayer(dog_sound_3, audioCtx),
		}
	}

	return AudioCharacter{
		sound0: newFuncNewPlayer(dog_sound_0, audioCtx),
		sound1: newFuncNewPlayer(dog_sound_1, audioCtx),
		sound2: newFuncNewPlayer(dog_sound_2, audioCtx),
		sound3: newFuncNewPlayer(dog_sound_3, audioCtx),
	}
}

func newFuncNewPlayer(b []byte, audioContext *audio.Context) *audio.Player {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}

	const bytesPerSample = 4 // TODO: This should be defined in audio package

	s, err := mp3.DecodeWithoutResampling(bytes.NewReader(b))
	if err != nil {
		fmt.Printf("error new player => %s\n", s)
		return nil
	}

	p, err := audioContext.NewPlayer(s)
	if err != nil {
		fmt.Printf("error new player => %s\n", s)
		return nil
	}
	return p
}

func getPlayer(fileName string, audioCtx *audio.Context) (*audio.Player, error) {
	b, err := os.ReadFile(fileName) // just pass the file name
	if err != nil {
		return nil, fmt.Errorf("error on readFile %s => %s\n", fileName, err)
	}
	s, err := mp3.DecodeWithoutResampling(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("error on decode %s => %s\n", fileName, err)
	}
	p, err := audioCtx.NewPlayer(s)
	if err != nil {
		return nil, fmt.Errorf("error creating newPlayer %s => %s\n", fileName, err)
	}
	return p, nil
}
