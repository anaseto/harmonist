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

func (ui *model) SwappingAnimation(mpos, ppos gruid.Point) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	Sleep(AnimDurShort)
	_, fgm, bgColorm := ui.PositionDrawing(mpos)
	_, _, bgColorp := ui.PositionDrawing(ppos)
	ui.DrawAtPosition(mpos, true, 'Φ', fgm, bgColorp)
	ui.DrawAtPosition(ppos, true, 'Φ', ColorFgPlayer, bgColorm)
	ui.Flush()
	Sleep(AnimDurMedium)
	ui.DrawAtPosition(mpos, true, 'Φ', ColorFgPlayer, bgColorp)
	ui.DrawAtPosition(ppos, true, 'Φ', fgm, bgColorm)
	ui.Flush()
	Sleep(AnimDurMedium)
}

func (ui *model) TeleportAnimation(from, to gruid.Point, showto bool) {
	if DisableAnimations {
		return
	}
	_, _, bgColorf := ui.PositionDrawing(from)
	_, _, bgColort := ui.PositionDrawing(to)
	ui.DrawAtPosition(from, true, 'Φ', ColorCyan, bgColorf)
	ui.Flush()
	Sleep(AnimDurMediumLong)
	if showto {
		ui.DrawAtPosition(from, true, 'Φ', ColorBlue, bgColorf)
		ui.DrawAtPosition(to, true, 'Φ', ColorCyan, bgColort)
		ui.Flush()
		Sleep(AnimDurMedium)
	}
}

type explosionStyle int

const (
	FireExplosion explosionStyle = iota
	WallExplosion
	AroundWallExplosion
)

func (ui *model) MonsterProjectileAnimation(ray []gruid.Point, r rune, fg uicolor) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	Sleep(AnimDurShort)
	for i := 0; i < len(ray); i++ {
		pos := ray[i]
		or, fgColor, bgColor := ui.PositionDrawing(pos)
		ui.DrawAtPosition(pos, true, r, fg, bgColor)
		ui.Flush()
		Sleep(AnimDurShort)
		ui.DrawAtPosition(pos, true, or, fgColor, bgColor)
	}
}

func (ui *model) WaveDrawAt(pos gruid.Point, fg uicolor) {
	r, _, bgColor := ui.PositionDrawing(pos)
	ui.DrawAtPosition(pos, true, r, bgColor, fg)
}

func (ui *model) ExplosionDrawAt(pos gruid.Point, fg uicolor) {
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

func (ui *model) NoiseAnimation(noises []gruid.Point) {
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
		Sleep(AnimDurShortMedium)
	}

}

func (ui *model) ExplosionAnimation(es explosionStyle, pos gruid.Point) {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	Sleep(AnimDurShort)
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
		Sleep(AnimDurMediumLong)
	}
}

func (g *state) Waves(maxCost int, ws wavestyle, center gruid.Point) (dists []int, cdists map[int][]int) {
	var dij Dijkstrer
	switch ws {
	case WaveMagicNoise:
		dij = &gridPath{dungeon: g.Dungeon}
	default:
		dij = &noisePath{state: g}
	}
	nm := Dijkstra(dij, []gruid.Point{center}, maxCost)
	cdists = make(map[int][]int)
	nm.iter(g.Player.Pos, func(n *node) {
		pos := n.Pos
		cdists[n.Cost] = append(cdists[n.Cost], idx(pos))
	})
	for dist := range cdists {
		dists = append(dists, dist)
	}
	sort.Ints(dists)
	return dists, cdists
}

func (ui *model) LOSWavesAnimation(r int, ws wavestyle, center gruid.Point) {
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

func (ui *model) WaveAnimation(wave []int, ws wavestyle) {
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
	Sleep(AnimDurShort)
}

func (ui *model) WallExplosionAnimation(pos gruid.Point) {
	if DisableAnimations {
		return
	}
	colors := [2]uicolor{ColorFgExplosionWallStart, ColorFgExplosionWallEnd}
	for _, fg := range colors {
		_, _, bgColor := ui.PositionDrawing(pos)
		//ui.DrawAtPosition(pos, true, '☼', fg, bgColor)
		ui.DrawAtPosition(pos, true, '%', bgColor, fg)
		ui.Flush()
		Sleep(AnimDurShort)
	}
}

type beamstyle int

const (
	BeamSleeping beamstyle = iota
	BeamLignification
	BeamObstruction
)

func (ui *model) BeamsAnimation(ray []gruid.Point, bs beamstyle) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	Sleep(AnimDurShort)
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
		Sleep(AnimDurShortMedium)
	}
}

func (ui *model) SlowingMagaraAnimation(ray []gruid.Point) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	Sleep(AnimDurShort)
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
		Sleep(AnimDurShortMedium)
	}
}

func (ui *model) ProjectileSymbol(dir direction) (r rune) {
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

func (ui *model) MonsterJavelinAnimation(ray []gruid.Point, hit bool) {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	Sleep(AnimDurShort)
	for i := 0; i < len(ray); i++ {
		pos := ray[i]
		r, fgColor, bgColor := ui.PositionDrawing(pos)
		ui.DrawAtPosition(pos, true, ui.ProjectileSymbol(Dir(g.Player.Pos, pos)), ColorFgMonster, bgColor)
		ui.Flush()
		Sleep(AnimDurShort)
		ui.DrawAtPosition(pos, true, r, fgColor, bgColor)
	}
	Sleep(AnimDurShort)
}

func (ui *model) WoundedAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(AnimationMode)
	r, _, bg := ui.PositionDrawing(g.Player.Pos)
	ui.DrawAtPosition(g.Player.Pos, false, r, ColorFgHPwounded, bg)
	ui.Flush()
	Sleep(AnimDurShortMedium)
	if g.Player.HP <= 15 {
		ui.DrawAtPosition(g.Player.Pos, false, r, ColorFgHPcritical, bg)
		ui.Flush()
		Sleep(AnimDurShortMedium)
	}
}

func (ui *model) PlayerGoodEffectAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(AnimationMode)
	Sleep(AnimDurShort)
	r, fg, bg := ui.PositionDrawing(g.Player.Pos)
	ui.DrawAtPosition(g.Player.Pos, false, r, ColorGreen, bg)
	ui.Flush()
	Sleep(AnimDurShortMedium)
	ui.DrawAtPosition(g.Player.Pos, false, r, ColorYellow, bg)
	ui.Flush()
	Sleep(AnimDurShortMedium)
	ui.DrawAtPosition(g.Player.Pos, false, r, fg, bg)
	ui.Flush()
}

func (ui *model) StatusEndAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(AnimationMode)
	r, _, bg := ui.PositionDrawing(g.Player.Pos)
	ui.DrawAtPosition(g.Player.Pos, false, r, ColorViolet, bg)
	ui.Flush()
	Sleep(AnimDurShortMedium)
}

func (ui *model) FoundFakeStairsAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	r, fg, bg := ui.PositionDrawing(g.Player.Pos)
	ui.DrawAtPosition(g.Player.Pos, false, r, ColorMagenta, bg)
	ui.Flush()
	Sleep(AnimDurMediumLong)
	ui.DrawAtPosition(g.Player.Pos, false, r, fg, bg)
	ui.Flush()
}

func (ui *model) MusicAnimation(pos gruid.Point) {
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

func (ui *model) PushAnimation(path []gruid.Point) {
	if DisableAnimations {
		return
	}
	if len(path) == 0 {
		// should not happen
		return
	}
	ui.DrawDungeonView(AnimationMode)
	_, _, bg := ui.PositionDrawing(path[0])
	for _, pos := range path[:len(path)-1] {
		ui.DrawAtPosition(pos, false, '×', ColorFgPlayer, bg)
	}
	ui.DrawAtPosition(path[len(path)-1], false, '@', ColorFgPlayer, bg)
	ui.Flush()
	Sleep(AnimDurMediumLong)
}

func (ui *model) MenuSelectedAnimation(m menu, ok bool) {
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
		Sleep(t)
		ui.DrawColoredText(message, MenuCols[m][0], DungeonHeight, ColorViolet)
	}
}

func (ui *model) MagicMappingAnimation(border []int) {
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

func (ui *model) FreeingShaedraAnimation() {
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
	Sleep(AnimDurMediumLong)
	g.Objects.Stairs[g.Places.Monolith] = WinStair
	g.Dungeon.SetCell(g.Places.Monolith, StairCell)
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	Sleep(AnimDurLong)
	_, _, bg = ui.PositionDrawing(g.Places.Marevor)
	ui.DrawAtPosition(g.Places.Marevor, false, 'Φ', ColorFgMagicPlace, bg)
	ui.Flush()
	Sleep(AnimDurMediumLong)
	g.Objects.Story[g.Places.Marevor] = StoryMarevor
	g.PrintStyled("Marevor: “And what about the mission? Take that magara!”", logSpecial)
	g.PrintStyled("Shaedra: “Pff, don't be reckless!”", logSpecial)
	g.Print("[(x) to continue]")
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	ui.WaitForContinue(-1)
	ui.DrawDungeonView(NoFlushMode)
	ui.DrawAtPosition(g.Places.Marevor, false, 'Φ', ColorFgMagicPlace, bg)
	ui.DrawAtPosition(g.Places.Shaedra, false, 'Φ', ColorFgMagicPlace, bg)
	ui.Flush()
	Sleep(AnimDurMediumLong)
	g.Dungeon.SetCell(g.Places.Shaedra, GroundCell)
	g.Dungeon.SetCell(g.Places.Marevor, ScrollCell)
	g.Objects.Scrolls[g.Places.Marevor] = ScrollExtended
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	g.RescuedShaedra()
}

func (g *state) RescuedShaedra() {
	g.Player.Magaras = append(g.Player.Magaras, magara{})
	g.Player.Inventory.Misc = NoItem
	g.PrintStyled("You equip the new magara in the artifact's old place.", logSpecial)
	if RandInt(2) == 0 {
		g.Player.Magaras[len(g.Player.Magaras)-1] = magara{Kind: DispersalMagara, Charges: DispersalMagara.DefaultCharges()}
	} else {
		g.Player.Magaras[len(g.Player.Magaras)-1] = magara{Kind: DelayedOricExplosionMagara, Charges: DelayedOricExplosionMagara.DefaultCharges()}
	}
	AchRescuedShaedra.Get(g)
}

func (ui *model) TakingArtifactAnimation() {
	g := ui.g
	//if DisableAnimations {
	// TODO this animation cannot be disabled as-is, because code is mixed with it...
	//return
	//}
	g.PrintStyled("You take and use the artifact.", logSpecial)
	g.Print("[(x) to continue].")
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	ui.WaitForContinue(-1)
	g.Dungeon.SetCell(g.Places.Artifact, GroundCell)
	_, _, bg := ui.PositionDrawing(g.Places.Monolith)
	ui.DrawAtPosition(g.Places.Monolith, false, 'Φ', ColorFgMagicPlace, bg)
	ui.Flush()
	Sleep(AnimDurMediumLong)
	g.Objects.Stairs[g.Places.Monolith] = WinStair
	g.Dungeon.SetCell(g.Places.Monolith, StairCell)
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	Sleep(AnimDurLong)
	_, _, bg = ui.PositionDrawing(g.Places.Marevor)
	ui.DrawAtPosition(g.Places.Marevor, false, 'Φ', ColorFgMagicPlace, bg)
	ui.Flush()
	Sleep(AnimDurMediumLong)
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
	Sleep(AnimDurMediumLong)
	g.Dungeon.SetCell(g.Places.Marevor, GroundCell)
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	AchRetrievedArtifact.Get(g)
}
