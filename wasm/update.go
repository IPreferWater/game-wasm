package main

func (g *Game) Update() error {

	g.frameCount++

	switch g.currentPhaseStance {
	case intro:
		g.williamTellPlayer.Play()
		if g.frameCount >= introFramesNbr {
			g.currentPhaseStance = firstAttackC1
			g.frameCount = 0
		}

	case firstAttackC1:
		if !g.williamTellPlayer.IsPlaying() {
			g.williamTellPlayer.Rewind()
			g.williamTellPlayer.Play()
		}
		if len(g.mapNoteToPlay) >= 3 {
			g.shortenNotesToPlay()
			g.character1.cooldown.resetCooldowns()
			g.currentPhaseStance = defendC2
			g.frameCount = 0
		}

		//blink for call to action
		g.blinkCount++
		if g.blinkCount > blinkFrameNbr {
			g.blinkCount = 0
			g.blink = !g.blink
		}

		line := checkActionStartAttack(g)
		if line <= -1 {
			break
		}

		g.mapNoteToPlay[g.frameCount] = line
		g.character1.notes = append(g.character1.notes, Note{
			x:         getPositionInLine(line, 0),
			y:         screenHeight - 20,
			line:      line,
			direction: up,
		})

	case defendC2:

		if g.notesDisplayed >= len(g.mapNoteToPlay) && len(g.character2.notes) <= 0 {
			g.currentPhaseStance = addNoteC2
			break
		}
		//DEFEND
		if line, ok := g.mapNoteToPlay[g.frameCount]; ok {
			x := getPositionInLine(line, startLayoutC2)
			g.character2.notes = append(g.character2.notes, Note{
				x:         x,
				y:         20,
				line:      line,
				direction: down,
			})
			g.notesDisplayed++
		}
		checkActionC2(g)
	case addNoteC2:
		lineToAddNote := getLineOfnoteAdded(g, false)
		if lineToAddNote > -1 {

			// 160 is aprox the time a note reach the line
			count := g.frameCount - 160
			g.mapNoteToPlay[count] = lineToAddNote
			g.character2.cooldown.resetCooldowns()

			g.currentPhaseStance = defendC1
			g.frameCount = 0
			g.notesDisplayed = 0
		}
	case defendC1:
		if g.notesDisplayed >= len(g.mapNoteToPlay) && len(g.character1.notes) <= 0 {
			g.currentPhaseStance = addNoteC1
			break
		}
		//DEFEND
		if line, ok := g.mapNoteToPlay[g.frameCount]; ok {
			x := getPositionInLine(line, 0)
			g.character1.notes = append(g.character1.notes, Note{
				x:         x,
				y:         20,
				line:      line,
				direction: down,
			})
			g.notesDisplayed++
		}
		checkActionC1(g)
	case addNoteC1:
		lineToAddNote := getLineOfnoteAdded(g, true)
		if lineToAddNote > -1 {

			count := g.frameCount - 160
			g.mapNoteToPlay[count] = lineToAddNote
			g.character1.cooldown.resetCooldowns()

			g.currentPhaseStance = defendC2
			g.frameCount = 0
			g.notesDisplayed = 0
		}
	case c1Lost, c2Lost:
		// stop musique
		if g.williamTellPlayer.IsPlaying() {
			g.williamTellPlayer.Pause()
		}
		// replay
		if checkActionResetGame(g) {
			g.frameCount = 0
			g.currentPhaseStance = firstAttackC1
			g.mapNoteToPlay = make(map[int]int)
			g.character1.notes = []Note{}
			g.character2.notes = []Note{}
			g.character2.cooldown.resetCooldowns()
			g.character1.cooldown.resetCooldowns()
			g.notesDisplayed = 0
		}
		return nil

	default:
	}

	if g.character1.updateNotesAndCheckIfLost() {
		g.currentPhaseStance = c1Lost
	}
	if g.character2.updateNotesAndCheckIfLost() {
		g.currentPhaseStance = c2Lost
	}

	return nil
}

func (cd *Cooldown) resetCooldowns() {
	cd.line1 = -50
	cd.line2 = -50
	cd.line3 = -50
	cd.line4 = -50
}

func (g *Game) shortenNotesToPlay() {
	min := -1

	for frame := range g.mapNoteToPlay {
		// min == -1 mean it was not set yet
		if min == -1 {
			min = frame
			continue
		}
		if frame < min {
			min = frame
		}
	}

	newMap := make(map[int]int)

	//don't shorten too much or character 2 will have to play to quick
	frameToShorten := min - 20
	for frame, line := range g.mapNoteToPlay {
		newMap[frame-frameToShorten] = line
	}

	g.mapNoteToPlay = newMap
}
