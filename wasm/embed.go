package main

import (
	_ "embed"
)

var (

	//go:embed  res/dog/back.png
	dog_back_png []byte

	//go:embed  res/dog/sprite.png
	dog_sprite_png []byte

	//go:embed res/dog/sound_0.mp3
	dog_sound_0 []byte

	//go:embed res/dog/sound_1.mp3
	dog_sound_1 []byte

	//go:embed res/dog/sound_2.mp3
	dog_sound_2 []byte

	//go:embed res/dog/sound_3.mp3
	dog_sound_3 []byte

	//go:embed  res/knight/back.png
	knight_back_png []byte

	//go:embed  res/knight/sprite.png
	knight_sprite_png []byte

	//go:embed res/knight/sound_0.mp3
	knight_sound_0 []byte

	//go:embed res/knight/sound_1.mp3
	knight_sound_1 []byte

	//go:embed res/knight/sound_2.mp3
	knight_sound_2 []byte

	//go:embed res/knight/sound_3.mp3
	knight_sound_3 []byte

	//go:embed res/william_tell_overture_8_bit.mp3
	william_tell_overture_8_bit []byte

	//go:embed res/note_sprite.png
	note_sprite_png []byte
)
