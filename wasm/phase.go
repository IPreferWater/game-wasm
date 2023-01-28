package main


type PhaseStance int64
const (
	intro PhaseStance = iota
	attackC1
	defendC1
	attackC2
	defendC2
)
type phase struct {
	introFramesNbr int
	firstTypingAttackFramesNbr int
}
type bla struct {
	line int
}

var (
	williamTellOverture map[int]bla
)


func initWillTellOverture(){
	williamTellOverture= map[int]bla{
		550: {0},
		560: {1},
		570: {2},
		580: {0},
		590: {1},
		600: {2},
	}
}

//introduction
//start typing-attack
// map j1 init 3 notes

//start typing-defense
// j2 play the map
// add timer for j2 to ad note