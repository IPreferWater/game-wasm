package main


type PhaseStance int64
const (
	intro PhaseStance = iota
	firstAttackC1
	attackC1
	defendC1
	attackC2
	defendC2
)


var (
	williamTellOverture map[int]int
)


func initWillTellOverture(){
	williamTellOverture= map[int]int{
		550: 0,
		560: 1,
		570: 2,
		580: 0,
		590: 1,
		600: 2,
	}
}

//introduction
//start typing-attack
// map j1 init 3 notes

//start typing-defense
// j2 play the map
// add timer for j2 to ad note