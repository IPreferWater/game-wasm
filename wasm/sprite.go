package main

type Sprite struct {
	width                   int
	height                  int
	yStar                   int
	numberOfSprites         int
	changeSpriteAfterFrames int
}
type SpriteStance int64

const (
	Playing SpriteStance = iota
	Lost
)

func initDogSprites() map[SpriteStance]Sprite {
	return map[SpriteStance]Sprite{
		Playing: {
			width:                   64,
			height:                  43,
			yStar: 0,
			numberOfSprites:         8,
			changeSpriteAfterFrames: 10,
		},
		Lost: {
			width:                   64,
			height:                  58,
			yStar: 44,
			numberOfSprites:         3,
			changeSpriteAfterFrames: 10,
		},
	}
}

func initKnightSprites() map[SpriteStance]Sprite {
	return map[SpriteStance]Sprite{
		Playing: {
			width:                   56,
			height:                  47,
			yStar: 0,
			numberOfSprites:         9,
			changeSpriteAfterFrames: 10,
		},
		Lost: {
			width:                   56,
			height:                  47,
			yStar: 49,
			numberOfSprites:         3,
			changeSpriteAfterFrames: 10,
		},
	}
}