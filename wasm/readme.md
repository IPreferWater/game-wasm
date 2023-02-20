GOOS=js GOARCH=wasm go build -o  ../assets/game.wasm  

TODO : 
✔️ if a player strike a bad note, he loose
✔️ add a cooldown on the same line
✔️ stop the musique when someone lost
✔️ touch space to play again
✔️ add specifics sounds for the knight
✔️ switch squares to a "buble" sprite
✔️ make a blinking "call to action" for c1 to make him start the game
❌ can't increase g.frameCount more than one or it break the notes to display[increase speed over time] 
✔️ remove the wait before start to type the first note
✔️ refactor draw.go
✔️ Should include the sprite and font go:embeded like mp3
✔️ displayed notes are not good for knight
- add horizontal lines only when it's your time to play
- display the cooldown
- fade away is not working

