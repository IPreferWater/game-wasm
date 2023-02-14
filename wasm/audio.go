package main

import (
	_ "embed"
)

var (
		//go:embed res/dog/sound_0.mp3
		dog_sound_0 []byte

		//go:embed res/dog/sound_1.mp3
		dog_sound_1 []byte

		//go:embed res/dog/sound_2.mp3
		dog_sound_2 []byte

		//go:embed res/dog/sound_3.mp3
		dog_sound_3 []byte

		//go:embed res/william_tell_overture_8_bit.mp3
		william_tell_overture_8_bit []byte
)