package main

import (
	"sort"
	"time"
)

const (
	AnimDurShort       = 25
	AnimDurShortMedium = 50
	AnimDurMedium      = 75
	AnimDurMediumLong  = 100
	AnimDurLong        = 200
	AnimDurExtraLong   = 300
)

func (ui *gameui) SwappingAnimation(mpos, ppos position) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	ui.Sleep(AnimDurShort)
	_, fgm, bgColorm := ui.PositionDrawing(mpos)
	_, _, bgColorp := ui.PositionDrawing(ppos)
	ui.DrawAtPosition(mpos, true, 'Φ', fgm, bgColorp)
	ui.DrawAtPosition(ppos, true, 'Φ', ColorFgPlayer, bgColorm)
	ui.Flush()
	ui.Sleep(AnimDurMedium)
	ui.DrawAtPosition(mpos, true, 'Φ', ColorFgPlayer, bgColorp)
	ui.DrawAtPosition(ppos, true, 'Φ', fgm, bgColorm)
	ui.Flush()
	ui.Sleep(AnimDurMedium)
}

func (ui *gameui) TeleportAnimation(from, to position, showto bool) {
	if DisableAnimations {
		return
	}
	_, _, bgColorf := ui.PositionDrawing(from)
	_, _, bgColort := ui.PositionDrawing(to)
	ui.DrawAtPosition(from, true, 'Φ', ColorCyan, bgColorf)
	ui.Flush()
	ui.Sleep(AnimDurMediumLong)
	if showto {
		ui.DrawAtPosition(from, true, 'Φ', ColorBlue, bgColorf)
		ui.DrawAtPosition(to, true, 'Φ', ColorCyan, bgColort)
		ui.Flush()
		ui.Sleep(AnimDurMedium)
	}
}

type explosionStyle int

const (
	FireExplosion explosionStyle = iota
	WallExplosion
	AroundWallExplosion
)

func (ui *gameui) MonsterProjectileAnimation(ray []position, r rune, fg uicolor) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	ui.Sleep(AnimDurShort)
	for i := 0; i < len(ray); i++ {
		pos := ray[i]
		or, fgColor, bgColor := ui.PositionDrawing(pos)
		ui.DrawAtPosition(pos, true, r, fg, bgColor)
		ui.Flush()
		ui.Sleep(AnimDurShort)
		ui.DrawAtPosition(pos, true, or, fgColor, bgColor)
	}
}

func (ui *gameui) WaveDrawAt(pos position, fg uicolor) {
	r, _, bgColor := ui.PositionDrawing(pos)
	ui.DrawAtPosition(pos, true, r, bgColor, fg)
}

func (ui *gameui) ExplosionDrawAt(pos position, fg uicolor) {
	g := ui.g
	_, _, bgColor := ui.PositionDrawing(pos)
	mons := g.MonsterAt(pos)
	r := ';'
	switch RandInt(9) {
	case 0, 6:
		r = ','
	case 1:
		r = '}'
	case 2:
		r = '%'
	case 3, 7:
		r = ':'
	case 4:
		r = '\\'
	case 5:
		r = '~'
	}
	if mons.Exists() || g.Player.Pos == pos {
		r = '√'
	}
	ui.DrawAtPosition(pos, true, r, bgColor, fg)
}

func (ui *gameui) NoiseAnimation(noises []position) {
	if DisableAnimations {
		return
	}
	ui.LOSWavesAnimation(DefaultLOSRange, WaveMagicNoise, ui.g.Player.Pos)
	colors := []uicolor{ColorFgSleepingMonster, ColorFgMagicPlace}
	for i := 0; i < 2; i++ {
		for _, pos := range noises {
			r := '♫'
			_, _, bgColor := ui.PositionDrawing(pos)
			ui.DrawAtPosition(pos, false, r, bgColor, colors[i])
		}
		_, _, bgColor := ui.PositionDrawing(ui.g.Player.Pos)
		ui.DrawAtPosition(ui.g.Player.Pos, false, '@', bgColor, colors[i])
		ui.Flush()
		ui.Sleep(AnimDurShortMedium)
	}

}

func (ui *gameui) ExplosionAnimation(es explosionStyle, pos position) {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	ui.Sleep(AnimDurShort)
	colors := [2]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd}
	if es == WallExplosion || es == AroundWallExplosion {
		colors[0] = ColorFgExplosionWallStart
		colors[1] = ColorFgExplosionWallEnd
	}
	for i := 0; i < 3; i++ {
		nb := g.Dungeon.FreeNeighbors(pos)
		if es != AroundWallExplosion {
			nb = append(nb, pos)
		}
		for _, npos := range nb {
			fg := colors[RandInt(2)]
			if !g.Player.LOS[npos] {
				continue
			}
			ui.ExplosionDrawAt(npos, fg)
		}
		ui.Flush()
		ui.Sleep(AnimDurMediumLong)
	}
}

func (g *game) Waves(maxCost int, ws wavestyle, center position) (dists []int, cdists map[int][]int) {
	var dij Dijkstrer
	switch ws {
	case WaveMagicNoise:
		dij = &gridPath{dungeon: g.Dungeon}
	default:
		dij = &noisePath{game: g}
	}
	nm := Dijkstra(dij, []position{center}, maxCost)
	cdists = make(map[int][]int)
	nm.iter(g.Player.Pos, func(n *node) {
		pos := n.Pos
		cdists[n.Cost] = append(cdists[n.Cost], pos.idx())
	})
	for dist := range cdists {
		dists = append(dists, dist)
	}
	sort.Ints(dists)
	return dists, cdists
}

func (ui *gameui) LOSWavesAnimation(r int, ws wavestyle, center position) {
	dists, cdists := ui.g.Waves(r, ws, center)
	for _, d := range dists {
		wave := cdists[d]
		if len(wave) == 0 {
			break
		}
		ui.WaveAnimation(wave, ws)
	}
}

type wavestyle int

const (
	WaveMagicNoise wavestyle = iota
	WaveNoise
	WaveConfusion
	WaveSlowing
	WaveTree
	WaveSleeping
)

func (ui *gameui) WaveAnimation(wave []int, ws wavestyle) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	for _, i := range wave {
		pos := idxtopos(i)
		switch ws {
		case WaveConfusion:
			fg := ColorFgConfusedMonster
			if ui.g.Player.Sees(pos) {
				ui.WaveDrawAt(pos, fg)
			}
		case WaveSleeping:
			fg := ColorFgSleepingMonster
			if ui.g.Player.Sees(pos) {
				ui.WaveDrawAt(pos, fg)
			}
		case WaveSlowing:
			fg := ColorFgParalysedMonster
			if ui.g.Player.Sees(pos) {
				ui.WaveDrawAt(pos, fg)
			}
		case WaveTree:
			fg := ColorFgLignifiedMonster
			if ui.g.Player.Sees(pos) {
				ui.WaveDrawAt(pos, fg)
			}
		case WaveNoise:
			fg := ColorFgWanderingMonster
			if ui.g.Player.Sees(pos) {
				ui.WaveDrawAt(pos, fg)
			}
		case WaveMagicNoise:
			fg := ColorFgMagicPlace
			ui.WaveDrawAt(pos, fg)
		}
	}
	ui.Flush()
	ui.Sleep(AnimDurShort)
}

func (ui *gameui) WallExplosionAnimation(pos position) {
	if DisableAnimations {
		return
	}
	colors := [2]uicolor{ColorFgExplosionWallStart, ColorFgExplosionWallEnd}
	for _, fg := range colors {
		_, _, bgColor := ui.PositionDrawing(pos)
		//ui.DrawAtPosition(pos, true, '☼', fg, bgColor)
		ui.DrawAtPosition(pos, true, '%', bgColor, fg)
		ui.Flush()
		ui.Sleep(AnimDurShort)
	}
}

type beamstyle int

const (
	BeamSleeping beamstyle = iota
	BeamLignification
	BeamObstruction
)

func (ui *gameui) BeamsAnimation(ray []position, bs beamstyle) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	ui.Sleep(AnimDurShort)
	// change colors depending on effect
	var fg uicolor
	switch bs {
	case BeamSleeping:
		fg = ColorFgSleepingMonster
	case BeamLignification:
		fg = ColorFgLignifiedMonster
	case BeamObstruction:
		fg = ColorFgMagicPlace
	}
	for j := 0; j < 3; j++ {
		for i := len(ray) - 1; i >= 0; i-- {
			pos := ray[i]
			_, _, bgColor := ui.PositionDrawing(pos)
			r := '*'
			if RandInt(2) == 0 {
				r = '×'
			}
			ui.DrawAtPosition(pos, true, r, bgColor, fg)
		}
		ui.Flush()
		ui.Sleep(AnimDurShortMedium)
	}
}

func (ui *gameui) SlowingMagaraAnimation(ray []position) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	ui.Sleep(AnimDurShort)
	colors := [2]uicolor{ColorFgConfusedMonster, ColorFgMagicPlace}
	for j := 0; j < 3; j++ {
		for i := len(ray) - 1; i >= 0; i-- {
			fg := colors[RandInt(2)]
			pos := ray[i]
			_, _, bgColor := ui.PositionDrawing(pos)
			r := '*'
			if RandInt(2) == 0 {
				r = '×'
			}
			ui.DrawAtPosition(pos, true, r, bgColor, fg)
		}
		ui.Flush()
		ui.Sleep(AnimDurShortMedium)
	}
}

func (ui *gameui) ProjectileSymbol(dir direction) (r rune) {
	switch dir {
	case E, ENE, ESE, WNW, W, WSW:
		r = '—'
	case NE, SW:
		r = '/'
	case NNE, N, NNW, SSW, S, SSE:
		r = '|'
	case NW, SE:
		r = '\\'
	}
	return r
}

func (ui *gameui) MonsterJavelinAnimation(ray []position, hit bool) {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	ui.Sleep(AnimDurShort)
	for i := 0; i < len(ray); i++ {
		pos := ray[i]
		r, fgColor, bgColor := ui.PositionDrawing(pos)
		ui.DrawAtPosition(pos, true, ui.ProjectileSymbol(pos.Dir(g.Player.Pos)), ColorFgMonster, bgColor)
		ui.Flush()
		ui.Sleep(AnimDurShort)
		ui.DrawAtPosition(pos, true, r, fgColor, bgColor)
	}
	ui.Sleep(AnimDurShort)
}

func (ui *gameui) WoundedAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	r, _, bg := ui.PositionDrawing(g.Player.Pos)
	ui.DrawAtPosition(g.Player.Pos, false, r, ColorFgHPwounded, bg)
	ui.Flush()
	ui.Sleep(AnimDurShortMedium)
	if g.Player.HP <= 15 {
		ui.DrawAtPosition(g.Player.Pos, false, r, ColorFgHPcritical, bg)
		ui.Flush()
		ui.Sleep(AnimDurShortMedium)
	}
}

func (ui *gameui) PlayerGoodEffectAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.Sleep(AnimDurShortMedium)
	r, fg, bg := ui.PositionDrawing(g.Player.Pos)
	ui.DrawAtPosition(g.Player.Pos, false, r, ColorGreen, bg)
	ui.Flush()
	ui.Sleep(AnimDurMedium)
	ui.DrawAtPosition(g.Player.Pos, false, r, ColorYellow, bg)
	ui.Flush()
	ui.Sleep(AnimDurMedium)
	ui.DrawAtPosition(g.Player.Pos, false, r, fg, bg)
	ui.Flush()
}

func (ui *gameui) StatusEndAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	r, fg, bg := ui.PositionDrawing(g.Player.Pos)
	ui.DrawAtPosition(g.Player.Pos, false, r, ColorViolet, bg)
	ui.Flush()
	ui.Sleep(AnimDurMedium)
	ui.DrawAtPosition(g.Player.Pos, false, r, fg, bg)
	ui.Flush()
}

func (ui *gameui) FoundFakeStairsAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	r, fg, bg := ui.PositionDrawing(g.Player.Pos)
	ui.DrawAtPosition(g.Player.Pos, false, r, ColorMagenta, bg)
	ui.Flush()
	ui.Sleep(AnimDurMediumLong)
	ui.DrawAtPosition(g.Player.Pos, false, r, fg, bg)
	ui.Flush()
}

func (ui *gameui) MusicAnimation(pos position) {
	if DisableAnimations {
		return
	}
	// TODO: not convinced by this animation
	//r, fg, bg := ui.PositionDrawing(pos)
	//ui.DrawAtPosition(pos, false, '♪', ColorCyan, bg)
	//ui.Flush()
	//time.Sleep(AnimDurMediumLong)
	//ui.DrawAtPosition(pos, false, r, fg, bg)
	//ui.Flush()
}

func (ui *gameui) PushAnimation(path []position) {
	if DisableAnimations {
		return
	}
	if len(path) == 0 {
		// should not happen
		return
	}
	_, _, bg := ui.PositionDrawing(path[0])
	for _, pos := range path[:len(path)-1] {
		ui.DrawAtPosition(pos, false, '×', ColorFgPlayer, bg)
	}
	ui.DrawAtPosition(path[len(path)-1], false, '@', ColorFgPlayer, bg)
	ui.Flush()
	ui.Sleep(AnimDurMediumLong)
}

func (ui *gameui) MenuSelectedAnimation(m menu, ok bool) {
	if DisableAnimations {
		return
	}
	if !ui.Small() {
		var message string
		if m == MenuInteract {
			message = ui.UpdateInteractButton()
		} else {
			message = m.String()
		}
		if message == "" {
			return
		}
		if ok {
			ui.DrawColoredText(message, MenuCols[m][0], DungeonHeight, ColorCyan)
		} else {
			ui.DrawColoredText(message, MenuCols[m][0], DungeonHeight, ColorMagenta)
		}
		ui.Flush()
		var t time.Duration = 25
		if !ok {
			t += 25
		}
		ui.Sleep(t * time.Millisecond)
		ui.DrawColoredText(message, MenuCols[m][0], DungeonHeight, ColorViolet)
	}
}

func (ui *gameui) MagicMappingAnimation(border []int) {
	if DisableAnimations {
		return
	}
	for _, i := range border {
		pos := idxtopos(i)
		r, fg, bg := ui.PositionDrawing(pos)
		ui.DrawAtPosition(pos, false, r, fg, bg)
	}
	ui.Flush()
}

func (ui *gameui) FreeingShaedraAnimation() {
	g := ui.g
	//if DisableAnimations {
	// TODO this animation cannot be disabled as-is, because code is mixed with it...
	//return
	//}
	g.Print("You see Shaedra. She is wounded!")
	g.PrintStyled("Shaedra: “Oh, it's you, Syu! Let's flee with Marevor's magara!”", logSpecial)
	g.Print("[(x) to continue]")
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	ui.WaitForContinue(-1)
	_, _, bg := ui.PositionDrawing(g.Places.Monolith)
	ui.DrawAtPosition(g.Places.Monolith, false, 'Φ', ColorFgMagicPlace, bg)
	ui.Flush()
	ui.Sleep(AnimDurMediumLong)
	g.Objects.Stairs[g.Places.Monolith] = WinStair
	g.Dungeon.SetCell(g.Places.Monolith, StairCell)
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	ui.Sleep(AnimDurLong)
	_, _, bg = ui.PositionDrawing(g.Places.Marevor)
	ui.DrawAtPosition(g.Places.Marevor, false, 'Φ', ColorFgMagicPlace, bg)
	ui.Flush()
	ui.Sleep(AnimDurMediumLong)
	g.Objects.Story[g.Places.Marevor] = StoryMarevor
	g.PrintStyled("Marevor: “And what about the mission?”", logSpecial)
	g.PrintStyled("Shaedra: “Pff, don't be reckless!”", logSpecial)
	g.Print("[(x) to continue]")
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	ui.WaitForContinue(-1)
	ui.DrawDungeonView(NoFlushMode)
	ui.DrawAtPosition(g.Places.Marevor, false, 'Φ', ColorFgMagicPlace, bg)
	ui.DrawAtPosition(g.Places.Shaedra, false, 'Φ', ColorFgMagicPlace, bg)
	ui.Flush()
	ui.Sleep(AnimDurMediumLong)
	g.Dungeon.SetCell(g.Places.Shaedra, GroundCell)
	g.Dungeon.SetCell(g.Places.Marevor, ScrollCell)
	g.Objects.Scrolls[g.Places.Marevor] = ScrollExtended
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	g.Player.Magaras = append(g.Player.Magaras, magara{})
	g.Player.Inventory.Misc = NoItem
	g.PrintStyled("You have a new empty slot for a magara.", logSpecial)
	AchRescuedShaedra.Get(g)
}

func (ui *gameui) TakingArtifactAnimation() {
	g := ui.g
	//if DisableAnimations {
	// TODO this animation cannot be disabled as-is, because code is mixed with it...
	//return
	//}
	g.Print("You take and use the artifact.")
	g.Dungeon.SetCell(g.Places.Artifact, GroundCell)
	_, _, bg := ui.PositionDrawing(g.Places.Monolith)
	ui.DrawAtPosition(g.Places.Monolith, false, 'Φ', ColorFgMagicPlace, bg)
	ui.Flush()
	ui.Sleep(AnimDurMediumLong)
	g.Objects.Stairs[g.Places.Monolith] = WinStair
	g.Dungeon.SetCell(g.Places.Monolith, StairCell)
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	ui.Sleep(AnimDurLong)
	_, _, bg = ui.PositionDrawing(g.Places.Marevor)
	ui.DrawAtPosition(g.Places.Marevor, false, 'Φ', ColorFgMagicPlace, bg)
	ui.Flush()
	ui.Sleep(AnimDurMediumLong)
	g.Objects.Story[g.Places.Marevor] = StoryMarevor
	g.PrintStyled("Marevor: “Great! Let's escape and find some bones to celebrate!”", logSpecial)
	g.PrintStyled("Syu: “Sorry, but I prefer bananas!”", logSpecial)
	g.Print("[(x) to continue]")
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	ui.WaitForContinue(-1)
	ui.DrawDungeonView(NoFlushMode)
	ui.DrawAtPosition(g.Places.Marevor, false, 'Φ', ColorFgMagicPlace, bg)
	ui.Flush()
	ui.Sleep(AnimDurMediumLong)
	g.Dungeon.SetCell(g.Places.Marevor, GroundCell)
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	AchRetrievedArtifact.Get(g)
}
