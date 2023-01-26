package main

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