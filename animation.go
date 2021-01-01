package main

import (
	"sort"
	//"time"

	"github.com/anaseto/gruid"
)

const (
	AnimDurShort       = 25
	AnimDurShortMedium = 50
	AnimDurMedium      = 75
	AnimDurMediumLong  = 100
	AnimDurLong        = 200
	AnimDurExtraLong   = 300
)

func (md *model) DrawAtPosition(p gruid.Point, targ bool, r rune, fg, bg gruid.Color) {
	// TODO
}

func (md *model) SwappingAnimation(mpos, ppos gruid.Point) {
	if DisableAnimations {
		return
	}
	md.DrawDungeonView(NormalMode)
	Sleep(AnimDurShort)
	_, fgm, bgColorm := md.PositionDrawing(mpos)
	_, _, bgColorp := md.PositionDrawing(ppos)
	md.DrawAtPosition(mpos, true, 'Φ', fgm, bgColorp)
	md.DrawAtPosition(ppos, true, 'Φ', ColorFgPlayer, bgColorm)
	//ui.Flush()
	Sleep(AnimDurMedium)
	md.DrawAtPosition(mpos, true, 'Φ', ColorFgPlayer, bgColorp)
	md.DrawAtPosition(ppos, true, 'Φ', fgm, bgColorm)
	//ui.Flush()
	Sleep(AnimDurMedium)
}

func (md *model) TeleportAnimation(from, to gruid.Point, showto bool) {
	if DisableAnimations {
		return
	}
	_, _, bgColorf := md.PositionDrawing(from)
	_, _, bgColort := md.PositionDrawing(to)
	md.DrawAtPosition(from, true, 'Φ', ColorCyan, bgColorf)
	//ui.Flush()
	Sleep(AnimDurMediumLong)
	if showto {
		md.DrawAtPosition(from, true, 'Φ', ColorBlue, bgColorf)
		md.DrawAtPosition(to, true, 'Φ', ColorCyan, bgColort)
		//ui.Flush()
		Sleep(AnimDurMedium)
	}
}

type explosionStyle int

const (
	FireExplosion explosionStyle = iota
	WallExplosion
	AroundWallExplosion
)

func (md *model) MonsterProjectileAnimation(ray []gruid.Point, r rune, fg gruid.Color) {
	if DisableAnimations {
		return
	}
	md.DrawDungeonView(NormalMode)
	Sleep(AnimDurShort)
	for i := 0; i < len(ray); i++ {
		pos := ray[i]
		or, fgColor, bgColor := md.PositionDrawing(pos)
		md.DrawAtPosition(pos, true, r, fg, bgColor)
		//ui.Flush()
		Sleep(AnimDurShort)
		md.DrawAtPosition(pos, true, or, fgColor, bgColor)
	}
}

func (md *model) WaveDrawAt(pos gruid.Point, fg gruid.Color) {
	r, _, bgColor := md.PositionDrawing(pos)
	md.DrawAtPosition(pos, true, r, bgColor, fg)
}

func (md *model) ExplosionDrawAt(pos gruid.Point, fg gruid.Color) {
	g := md.g
	_, _, bgColor := md.PositionDrawing(pos)
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
	md.DrawAtPosition(pos, true, r, bgColor, fg)
}

func (md *model) NoiseAnimation(noises []gruid.Point) {
	if DisableAnimations {
		return
	}
	md.LOSWavesAnimation(DefaultLOSRange, WaveMagicNoise, md.g.Player.Pos)
	colors := []gruid.Color{ColorFgSleepingMonster, ColorFgMagicPlace}
	for i := 0; i < 2; i++ {
		for _, pos := range noises {
			r := '♫'
			_, _, bgColor := md.PositionDrawing(pos)
			md.DrawAtPosition(pos, false, r, bgColor, colors[i])
		}
		_, _, bgColor := md.PositionDrawing(md.g.Player.Pos)
		md.DrawAtPosition(md.g.Player.Pos, false, '@', bgColor, colors[i])
		//ui.Flush()
		Sleep(AnimDurShortMedium)
	}

}

func (md *model) ExplosionAnimation(es explosionStyle, pos gruid.Point) {
	g := md.g
	if DisableAnimations {
		return
	}
	md.DrawDungeonView(NormalMode)
	Sleep(AnimDurShort)
	colors := [2]gruid.Color{ColorFgExplosionStart, ColorFgExplosionEnd}
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
			md.ExplosionDrawAt(npos, fg)
		}
		//ui.Flush()
		Sleep(AnimDurMediumLong)
	}
}

func (g *game) Waves(maxCost int, ws wavestyle, center gruid.Point) (dists []int, cdists map[int][]int) {
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

func (md *model) LOSWavesAnimation(r int, ws wavestyle, center gruid.Point) {
	dists, cdists := md.g.Waves(r, ws, center)
	for _, d := range dists {
		wave := cdists[d]
		if len(wave) == 0 {
			break
		}
		md.WaveAnimation(wave, ws)
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

func (md *model) WaveAnimation(wave []int, ws wavestyle) {
	if DisableAnimations {
		return
	}
	md.DrawDungeonView(NormalMode)
	for _, i := range wave {
		pos := idxtopos(i)
		switch ws {
		case WaveConfusion:
			fg := ColorFgConfusedMonster
			if md.g.Player.Sees(pos) {
				md.WaveDrawAt(pos, fg)
			}
		case WaveSleeping:
			fg := ColorFgSleepingMonster
			if md.g.Player.Sees(pos) {
				md.WaveDrawAt(pos, fg)
			}
		case WaveSlowing:
			fg := ColorFgParalysedMonster
			if md.g.Player.Sees(pos) {
				md.WaveDrawAt(pos, fg)
			}
		case WaveTree:
			fg := ColorFgLignifiedMonster
			if md.g.Player.Sees(pos) {
				md.WaveDrawAt(pos, fg)
			}
		case WaveNoise:
			fg := ColorFgWanderingMonster
			if md.g.Player.Sees(pos) {
				md.WaveDrawAt(pos, fg)
			}
		case WaveMagicNoise:
			fg := ColorFgMagicPlace
			md.WaveDrawAt(pos, fg)
		}
	}
	//ui.Flush()
	Sleep(AnimDurShort)
}

func (md *model) WallExplosionAnimation(pos gruid.Point) {
	if DisableAnimations {
		return
	}
	colors := [2]gruid.Color{ColorFgExplosionWallStart, ColorFgExplosionWallEnd}
	for _, fg := range colors {
		_, _, bgColor := md.PositionDrawing(pos)
		//ui.DrawAtPosition(pos, true, '☼', fg, bgColor)
		md.DrawAtPosition(pos, true, '%', bgColor, fg)
		//ui.Flush()
		Sleep(AnimDurShort)
	}
}

type beamstyle int

const (
	BeamSleeping beamstyle = iota
	BeamLignification
	BeamObstruction
)

func (md *model) BeamsAnimation(ray []gruid.Point, bs beamstyle) {
	if DisableAnimations {
		return
	}
	md.DrawDungeonView(NormalMode)
	Sleep(AnimDurShort)
	// change colors depending on effect
	var fg gruid.Color
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
			_, _, bgColor := md.PositionDrawing(pos)
			r := '*'
			if RandInt(2) == 0 {
				r = '×'
			}
			md.DrawAtPosition(pos, true, r, bgColor, fg)
		}
		//ui.Flush()
		Sleep(AnimDurShortMedium)
	}
}

func (md *model) SlowingMagaraAnimation(ray []gruid.Point) {
	if DisableAnimations {
		return
	}
	md.DrawDungeonView(NormalMode)
	Sleep(AnimDurShort)
	colors := [2]gruid.Color{ColorFgConfusedMonster, ColorFgMagicPlace}
	for j := 0; j < 3; j++ {
		for i := len(ray) - 1; i >= 0; i-- {
			fg := colors[RandInt(2)]
			pos := ray[i]
			_, _, bgColor := md.PositionDrawing(pos)
			r := '*'
			if RandInt(2) == 0 {
				r = '×'
			}
			md.DrawAtPosition(pos, true, r, bgColor, fg)
		}
		//ui.Flush()
		Sleep(AnimDurShortMedium)
	}
}

func (md *model) ProjectileSymbol(dir direction) (r rune) {
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

func (md *model) MonsterJavelinAnimation(ray []gruid.Point, hit bool) {
	g := md.g
	if DisableAnimations {
		return
	}
	md.DrawDungeonView(NormalMode)
	Sleep(AnimDurShort)
	for i := 0; i < len(ray); i++ {
		pos := ray[i]
		r, fgColor, bgColor := md.PositionDrawing(pos)
		md.DrawAtPosition(pos, true, md.ProjectileSymbol(Dir(g.Player.Pos, pos)), ColorFgMonster, bgColor)
		//ui.Flush()
		Sleep(AnimDurShort)
		md.DrawAtPosition(pos, true, r, fgColor, bgColor)
	}
	Sleep(AnimDurShort)
}

func (md *model) WoundedAnimation() {
	g := md.g
	if DisableAnimations {
		return
	}
	md.DrawDungeonView(AnimationMode)
	r, _, bg := md.PositionDrawing(g.Player.Pos)
	md.DrawAtPosition(g.Player.Pos, false, r, ColorFgHPwounded, bg)
	//ui.Flush()
	Sleep(AnimDurShortMedium)
	if g.Player.HP <= 15 {
		md.DrawAtPosition(g.Player.Pos, false, r, ColorFgHPcritical, bg)
		//ui.Flush()
		Sleep(AnimDurShortMedium)
	}
}

func (md *model) PlayerGoodEffectAnimation() {
	g := md.g
	if DisableAnimations {
		return
	}
	md.DrawDungeonView(AnimationMode)
	Sleep(AnimDurShort)
	r, fg, bg := md.PositionDrawing(g.Player.Pos)
	md.DrawAtPosition(g.Player.Pos, false, r, ColorGreen, bg)
	//ui.Flush()
	Sleep(AnimDurShortMedium)
	md.DrawAtPosition(g.Player.Pos, false, r, ColorYellow, bg)
	//ui.Flush()
	Sleep(AnimDurShortMedium)
	md.DrawAtPosition(g.Player.Pos, false, r, fg, bg)
	//ui.Flush()
}

func (md *model) StatusEndAnimation() {
	g := md.g
	if DisableAnimations {
		return
	}
	md.DrawDungeonView(AnimationMode)
	r, _, bg := md.PositionDrawing(g.Player.Pos)
	md.DrawAtPosition(g.Player.Pos, false, r, ColorViolet, bg)
	//ui.Flush()
	Sleep(AnimDurShortMedium)
}

func (md *model) FoundFakeStairsAnimation() {
	g := md.g
	if DisableAnimations {
		return
	}
	r, fg, bg := md.PositionDrawing(g.Player.Pos)
	md.DrawAtPosition(g.Player.Pos, false, r, ColorMagenta, bg)
	//ui.Flush()
	Sleep(AnimDurMediumLong)
	md.DrawAtPosition(g.Player.Pos, false, r, fg, bg)
	//ui.Flush()
}

func (md *model) MusicAnimation(pos gruid.Point) {
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

func (md *model) PushAnimation(path []gruid.Point) {
	if DisableAnimations {
		return
	}
	if len(path) == 0 {
		// should not happen
		return
	}
	md.DrawDungeonView(AnimationMode)
	_, _, bg := md.PositionDrawing(path[0])
	for _, pos := range path[:len(path)-1] {
		md.DrawAtPosition(pos, false, '×', ColorFgPlayer, bg)
	}
	md.DrawAtPosition(path[len(path)-1], false, '@', ColorFgPlayer, bg)
	//ui.Flush()
	Sleep(AnimDurMediumLong)
}

func (md *model) MagicMappingAnimation(border []int) {
	if DisableAnimations {
		return
	}
	for _, i := range border {
		pos := idxtopos(i)
		r, fg, bg := md.PositionDrawing(pos)
		md.DrawAtPosition(pos, false, r, fg, bg)
	}
	//ui.Flush()
}

func (md *model) FreeingShaedraAnimation() {
	g := md.g
	//if DisableAnimations {
	// TODO this animation cannot be disabled as-is, because code is mixed with it...
	//return
	//}
	g.Print("You see Shaedra. She is wounded!")
	g.PrintStyled("Shaedra: “Oh, it's you, Syu! Let's flee with Marevor's magara!”", logSpecial)
	g.Print("[(x) to continue]")
	//ui.DrawDungeonView(NoFlushMode)
	//ui.Flush()
	//ui.WaitForContinue(-1)
	_, _, bg := md.PositionDrawing(g.Places.Monolith)
	md.DrawAtPosition(g.Places.Monolith, false, 'Φ', ColorFgMagicPlace, bg)
	//ui.Flush()
	Sleep(AnimDurMediumLong)
	g.Objects.Stairs[g.Places.Monolith] = WinStair
	g.Dungeon.SetCell(g.Places.Monolith, StairCell)
	//ui.DrawDungeonView(NoFlushMode)
	//ui.Flush()
	Sleep(AnimDurLong)
	_, _, bg = md.PositionDrawing(g.Places.Marevor)
	md.DrawAtPosition(g.Places.Marevor, false, 'Φ', ColorFgMagicPlace, bg)
	//ui.Flush()
	Sleep(AnimDurMediumLong)
	g.Objects.Story[g.Places.Marevor] = StoryMarevor
	g.PrintStyled("Marevor: “And what about the mission? Take that magara!”", logSpecial)
	g.PrintStyled("Shaedra: “Pff, don't be reckless!”", logSpecial)
	g.Print("[(x) to continue]")
	//ui.DrawDungeonView(NoFlushMode)
	//ui.Flush()
	//ui.WaitForContinue(-1)
	//ui.DrawDungeonView(NoFlushMode)
	md.DrawAtPosition(g.Places.Marevor, false, 'Φ', ColorFgMagicPlace, bg)
	md.DrawAtPosition(g.Places.Shaedra, false, 'Φ', ColorFgMagicPlace, bg)
	//ui.Flush()
	Sleep(AnimDurMediumLong)
	g.Dungeon.SetCell(g.Places.Shaedra, GroundCell)
	g.Dungeon.SetCell(g.Places.Marevor, ScrollCell)
	g.Objects.Scrolls[g.Places.Marevor] = ScrollExtended
	//ui.DrawDungeonView(NoFlushMode)
	//ui.Flush()
	g.RescuedShaedra()
}

func (g *game) RescuedShaedra() {
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

func (md *model) TakingArtifactAnimation() {
	g := md.g
	//if DisableAnimations {
	// TODO this animation cannot be disabled as-is, because code is mixed with it...
	//return
	//}
	g.PrintStyled("You take and use the artifact.", logSpecial)
	g.Print("[(x) to continue].")
	//ui.DrawDungeonView(NoFlushMode)
	//ui.Flush()
	//ui.WaitForContinue(-1)
	g.Dungeon.SetCell(g.Places.Artifact, GroundCell)
	_, _, bg := md.PositionDrawing(g.Places.Monolith)
	md.DrawAtPosition(g.Places.Monolith, false, 'Φ', ColorFgMagicPlace, bg)
	//ui.Flush()
	Sleep(AnimDurMediumLong)
	g.Objects.Stairs[g.Places.Monolith] = WinStair
	g.Dungeon.SetCell(g.Places.Monolith, StairCell)
	//ui.DrawDungeonView(NoFlushMode)
	//ui.Flush()
	Sleep(AnimDurLong)
	_, _, bg = md.PositionDrawing(g.Places.Marevor)
	md.DrawAtPosition(g.Places.Marevor, false, 'Φ', ColorFgMagicPlace, bg)
	//ui.Flush()
	Sleep(AnimDurMediumLong)
	g.Objects.Story[g.Places.Marevor] = StoryMarevor
	g.PrintStyled("Marevor: “Great! Let's escape and find some bones to celebrate!”", logSpecial)
	g.PrintStyled("Syu: “Sorry, but I prefer bananas!”", logSpecial)
	g.Print("[(x) to continue]")
	//ui.DrawDungeonView(NoFlushMode)
	//ui.Flush()
	//ui.WaitForContinue(-1)
	//ui.DrawDungeonView(NoFlushMode)
	md.DrawAtPosition(g.Places.Marevor, false, 'Φ', ColorFgMagicPlace, bg)
	//ui.Flush()
	Sleep(AnimDurMediumLong)
	g.Dungeon.SetCell(g.Places.Marevor, GroundCell)
	//ui.DrawDungeonView(NoFlushMode)
	//ui.Flush()
	AchRetrievedArtifact.Get(g)
}
