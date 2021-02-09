package main

import (
	//"log"
	"sort"
	"time"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
)

const (
	AnimDurShort       = 25 * time.Millisecond
	AnimDurShortMedium = 50 * time.Millisecond
	AnimDurMedium      = 75 * time.Millisecond
	AnimDurMediumLong  = 100 * time.Millisecond
	AnimDurLong        = 200 * time.Millisecond
	AnimDurExtraLong   = 300 * time.Millisecond
)

type Animations struct {
	frames []AnimFrame
	grid   gruid.Grid
	pgrid  gruid.Grid
	idx    int
	draw   bool
}

type AnimFrame struct {
	Cells    []gruid.FrameCell
	Duration time.Duration
}

type msgAnim int

func (md *model) initAnimations() {
	gd := md.gd.Slice(md.gd.Range().Shift(0, 2, 0, -1))
	max := gd.Size()
	md.anims.grid = gruid.NewGrid(max.X, max.Y)
	md.anims.pgrid = gruid.NewGrid(max.X, max.Y)
	md.anims.grid.Copy(gd)
	md.anims.pgrid.Copy(gd)
}

func (md *model) animNext() gruid.Cmd {
	d := md.anims.frames[0].Duration
	idx := md.anims.idx
	return func() gruid.Msg {
		t := time.NewTimer(d)
		<-t.C
		return msgAnim(idx)
	}
}

func (md *model) animCmd() gruid.Cmd {
	if len(md.anims.frames) == 0 {
		return nil
	}
	idx := md.anims.idx
	return func() gruid.Msg {
		return msgAnim(idx)
	}
}

func (md *model) startAnimSeq() {
	if md.anims.Done() {
		md.resetAnimations()
	}
	it := md.g.Dungeon.Grid.Iterator()
	for it.Next() {
		p := it.P()
		r, fg, bg := md.positionDrawing(p)
		attrs := AttrInMap
		if md.g.Highlight[p] || p == md.mp.ex.pos {
			attrs |= AttrReverse
		}
		md.anims.grid.Set(p, gruid.Cell{Rune: r, Style: gruid.Style{Fg: fg, Bg: bg, Attrs: attrs}})
	}
}

func (md *model) resetAnimations() {
	gd := md.gd.Slice(md.gd.Range().Shift(0, 2, 0, -1))
	md.anims.grid.Copy(gd)
	md.anims.pgrid.Copy(gd)
}

func (a *Animations) Finish() {
	a.idx++
	a.frames = nil
}

func (a *Animations) Done() bool {
	return len(a.frames) == 0
}

func (a *Animations) Draw(p gruid.Point, r rune, fg, bg gruid.Color) {
	c := a.grid.At(p)
	c.Rune = r
	c.Style.Fg = fg
	c.Style.Bg = bg
	a.grid.Set(p, c)
}

func (a *Animations) Frame(d time.Duration) {
	frame := AnimFrame{}
	frame.Duration = d
	it := a.grid.Iterator()
	itp := a.pgrid.Iterator()
	for it.Next() && itp.Next() {
		if it.Cell() == itp.Cell() {
			continue
		}
		frame.Cells = append(frame.Cells, gruid.FrameCell{P: it.P(), Cell: it.Cell()})
	}
	a.frames = append(a.frames, frame)
	a.pgrid.Copy(a.grid)
}

func (md *model) DrawAtPosition(p gruid.Point, targ bool, r rune, fg, bg gruid.Color) {
	// TODO
}

func (md *model) SwappingAnimation(mpos, ppos gruid.Point) {
	if DisableAnimations {
		return
	}
	md.startAnimSeq()
	md.anims.Frame(AnimDurShort)
	_, fgm, bgColorm := md.positionDrawing(mpos)
	_, _, bgColorp := md.positionDrawing(ppos)
	md.anims.Draw(mpos, 'Φ', fgm, bgColorp)
	md.anims.Draw(ppos, 'Φ', ColorFgPlayer, bgColorm)
	md.anims.Frame(AnimDurMedium)
	md.anims.Draw(mpos, 'Φ', ColorFgPlayer, bgColorp)
	md.anims.Draw(ppos, 'Φ', fgm, bgColorm)
	md.anims.Frame(AnimDurMedium)
}

func (md *model) TeleportAnimation(from, to gruid.Point, showto bool) {
	if DisableAnimations {
		return
	}
	md.startAnimSeq()
	_, _, bgColorf := md.positionDrawing(from)
	_, _, bgColort := md.positionDrawing(to)
	md.anims.Draw(from, 'Φ', ColorCyan, bgColorf)
	md.anims.Frame(AnimDurMediumLong)
	if showto {
		md.anims.Draw(from, 'Φ', ColorBlue, bgColorf)
		md.anims.Draw(to, 'Φ', ColorCyan, bgColort)
		md.anims.Frame(AnimDurMedium)
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
	md.startAnimSeq()
	md.anims.Frame(AnimDurShort)
	for i := 0; i < len(ray); i++ {
		pos := ray[i]
		or, fgColor, bgColor := md.positionDrawing(pos)
		md.anims.Draw(pos, r, fg, bgColor)
		md.anims.Frame(AnimDurShort)
		md.anims.Draw(pos, or, fgColor, bgColor)
	}
}

func (md *model) WaveDrawAt(pos gruid.Point, fg gruid.Color) {
	r, _, bgColor := md.positionDrawing(pos)
	md.anims.Draw(pos, r, bgColor, fg)
}

func (md *model) ExplosionDrawAt(pos gruid.Point, fg gruid.Color) {
	g := md.g
	_, _, bgColor := md.positionDrawing(pos)
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
	md.anims.Draw(pos, r, bgColor, fg)
}

func (md *model) NoiseAnimation(noises []gruid.Point) {
	if DisableAnimations {
		return
	}
	md.LOSWavesAnimation(DefaultLOSRange, WaveMagicNoise, md.g.Player.Pos)
	//md.startAnimSeq()
	colors := []gruid.Color{ColorFgSleepingMonster, ColorFgMagicPlace}
	for i := 0; i < 2; i++ {
		for _, pos := range noises {
			r := '♫'
			_, _, bgColor := md.positionDrawing(pos)
			md.anims.Draw(pos, r, bgColor, colors[i])
		}
		_, _, bgColor := md.positionDrawing(md.g.Player.Pos)
		md.anims.Draw(md.g.Player.Pos, '@', bgColor, colors[i])
		md.anims.Frame(AnimDurShortMedium)
	}

}

func (md *model) ExplosionAnimation(es explosionStyle, pos gruid.Point) {
	g := md.g
	if DisableAnimations {
		return
	}
	md.startAnimSeq()
	md.anims.Frame(AnimDurShort)
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
		md.anims.Frame(AnimDurMediumLong)
	}
}

func (g *game) Waves(maxCost int, ws wavestyle, center gruid.Point) (dists []int, cdists map[int][]int) {
	var dij paths.Dijkstra
	switch ws {
	case WaveMagicNoise:
		dij = &gridPath{dungeon: g.Dungeon}
	default:
		dij = &noisePath{state: g}
	}
	nodes := g.PR.DijkstraMap(dij, []gruid.Point{center}, maxCost)
	cdists = make(map[int][]int)
	for _, n := range nodes {
		cdists[n.Cost] = append(cdists[n.Cost], idx(n.P))
	}
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
	WaveSleeping // TODO: check if really useful
)

func (md *model) WaveAnimation(wave []int, ws wavestyle) {
	if DisableAnimations {
		return
	}
	md.startAnimSeq()
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
	md.anims.Frame(AnimDurShort)
}

func (md *model) WallExplosionAnimation(pos gruid.Point) {
	if DisableAnimations {
		return
	}
	md.startAnimSeq()
	colors := [2]gruid.Color{ColorFgExplosionWallStart, ColorFgExplosionWallEnd}
	for _, fg := range colors {
		_, _, bgColor := md.positionDrawing(pos)
		//md.anims.Draw(pos, '☼', fg, bgColor)
		md.anims.Draw(pos, '%', bgColor, fg)
		md.anims.Frame(AnimDurShort)
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
	md.startAnimSeq()
	md.anims.Frame(AnimDurShort)
	// change colors depending on effect
	var fg gruid.Color
	switch bs {
	//	case BeamSleeping:
	//		fg = ColorFgSleepingMonster
	case BeamLignification:
		fg = ColorFgLignifiedMonster
	case BeamObstruction:
		fg = ColorFgMagicPlace
	}
	for j := 0; j < 3; j++ {
		for i := len(ray) - 1; i >= 0; i-- {
			pos := ray[i]
			_, _, bgColor := md.positionDrawing(pos)
			r := '*'
			if RandInt(2) == 0 {
				r = '×'
			}
			md.anims.Draw(pos, r, bgColor, fg)
		}
		md.anims.Frame(AnimDurShortMedium)
	}
}

func (md *model) SlowingMagaraAnimation(ray []gruid.Point) {
	if DisableAnimations {
		return
	}
	md.startAnimSeq()
	md.anims.Frame(AnimDurShort)
	colors := [2]gruid.Color{ColorFgConfusedMonster, ColorFgMagicPlace}
	for j := 0; j < 3; j++ {
		for i := len(ray) - 1; i >= 0; i-- {
			fg := colors[RandInt(2)]
			pos := ray[i]
			_, _, bgColor := md.positionDrawing(pos)
			r := '*'
			if RandInt(2) == 0 {
				r = '×'
			}
			md.anims.Draw(pos, r, bgColor, fg)
		}
		md.anims.Frame(AnimDurShortMedium)
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
	md.startAnimSeq()
	md.anims.Frame(AnimDurShort)
	for i := 0; i < len(ray); i++ {
		pos := ray[i]
		r, fgColor, bgColor := md.positionDrawing(pos)
		md.anims.Draw(pos, md.ProjectileSymbol(Dir(g.Player.Pos, pos)), ColorFgMonster, bgColor)
		md.anims.Frame(AnimDurShort)
		md.anims.Draw(pos, r, fgColor, bgColor)
	}
	//md.anims.Frame(AnimDurShort)
}

func (md *model) WoundedAnimation() {
	g := md.g
	if DisableAnimations {
		return
	}
	md.startAnimSeq()
	r, _, bg := md.positionDrawing(g.Player.Pos)
	md.anims.Draw(g.Player.Pos, r, ColorFgHPwounded, bg)
	md.anims.Frame(AnimDurShortMedium)
	if g.Player.HP <= 15 {
		md.anims.Draw(g.Player.Pos, r, ColorFgHPcritical, bg)
		md.anims.Frame(AnimDurShortMedium)
	}
}

func (md *model) PlayerGoodEffectAnimation() {
	g := md.g
	if DisableAnimations {
		return
	}
	md.startAnimSeq()
	md.anims.Frame(AnimDurShort)
	r, _, bg := md.positionDrawing(g.Player.Pos)
	md.anims.Draw(g.Player.Pos, r, ColorGreen, bg)
	md.anims.Frame(AnimDurShortMedium)
	md.anims.Draw(g.Player.Pos, r, ColorYellow, bg)
	md.anims.Frame(AnimDurShortMedium)
	//md.anims.Draw(g.Player.Pos, r, fg, bg)
}

func (md *model) StatusEndAnimation() {
	g := md.g
	if DisableAnimations {
		return
	}
	md.startAnimSeq()
	r, _, bg := md.positionDrawing(g.Player.Pos)
	md.anims.Draw(g.Player.Pos, r, ColorViolet, bg)
	md.anims.Frame(AnimDurShortMedium)
}

func (md *model) FoundFakeStairsAnimation() {
	g := md.g
	if DisableAnimations {
		return
	}
	r, _, bg := md.positionDrawing(g.Player.Pos)
	md.anims.Draw(g.Player.Pos, r, ColorMagenta, bg)
	md.anims.Frame(AnimDurMediumLong)
	//md.DrawAtPosition(g.Player.Pos, false, r, fg, bg)
}

func (md *model) MusicAnimation(pos gruid.Point) {
	if DisableAnimations {
		return
	}
	// TODO: not convinced by this animation
	//r, fg, bg := ui.PositionDrawing(pos)
	//ui.DrawAtPosition(pos, false, '♪', ColorCyan, bg)
	//ui.Flush()
	//	//time.Sleep(AnimDurMediumLong)
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
	md.startAnimSeq()
	_, _, bg := md.positionDrawing(path[0])
	for _, pos := range path[:len(path)-1] {
		md.anims.Draw(pos, '×', ColorFgPlayer, bg)
	}
	md.anims.Draw(path[len(path)-1], '@', ColorFgPlayer, bg)
	md.anims.Frame(AnimDurMediumLong)
	//log.Print("anims: %+v", md.anims)
}

func (md *model) MagicMappingAnimation(border []int) {
	if DisableAnimations {
		return
	}
	for _, i := range border {
		pos := idxtopos(i)
		r, fg, bg := md.positionDrawing(pos)
		md.anims.Draw(pos, r, fg, bg)
	}
	md.anims.Frame(AnimDurShort)
}

func (md *model) Story() {
	switch md.g.Depth {
	case WinDepth:
		md.FreeingShaedraAnimation()
	case MaxDepth:
		md.TakingArtifactAnimation()
	default:
		md.mode = modeNormal // should not happen
	}
}

func (md *model) FreeingShaedraAnimation() {
	g := md.g
	switch md.story {
	case 0:
		md.mode = modeStory
		g.Print("You see Shaedra. She is wounded!")
		g.PrintStyled("Shaedra: “Oh, it's you, Syu! Let's flee with Marevor's magara!”", logSpecial)
		g.Print("[(x) to continue]")
	case 1:
		if !DisableAnimations {
			md.startAnimSeq()
			_, _, bg := md.positionDrawing(g.Places.Monolith)
			md.anims.Draw(g.Places.Monolith, 'Φ', ColorFgMagicPlace, bg)
			md.anims.Frame(AnimDurMediumLong)
		}
		g.Objects.Stairs[g.Places.Monolith] = WinStair
		g.Dungeon.SetCell(g.Places.Monolith, StairCell)
		if !DisableAnimations {
			md.startAnimSeq()
			md.anims.Frame(AnimDurMediumLong)
			_, _, bg := md.positionDrawing(g.Places.Marevor)
			md.anims.Draw(g.Places.Marevor, 'Φ', ColorFgMagicPlace, bg)
			md.anims.Frame(AnimDurMediumLong)
		}
		g.Objects.Story[g.Places.Marevor] = StoryMarevor
		g.PrintStyled("Marevor: “And what about the mission? Take that magara!”", logSpecial)
		g.PrintStyled("Shaedra: “Pff, don't be reckless!”", logSpecial)
		g.Print("[(x) to continue]")
	case 2:
		if !DisableAnimations {
			md.startAnimSeq()
			_, _, bg := md.positionDrawing(g.Places.Marevor)
			md.anims.Draw(g.Places.Marevor, 'Φ', ColorFgMagicPlace, bg)
			md.anims.Frame(AnimDurMediumLong)
		}
		g.Objects.Story[g.Places.Marevor] = StoryMarevor
		g.PrintStyled("Marevor: “And what about the mission? Take that magara!”", logSpecial)
		g.PrintStyled("Shaedra: “Pff, don't be reckless!”", logSpecial)
		g.Print("[(x) to continue]")
	case 3:
		if !DisableAnimations {
			md.startAnimSeq()
			_, _, bg := md.positionDrawing(g.Places.Marevor)
			md.anims.Draw(g.Places.Marevor, 'Φ', ColorFgMagicPlace, bg)
			md.anims.Draw(g.Places.Shaedra, 'Φ', ColorFgMagicPlace, bg)
			md.anims.Frame(AnimDurMediumLong)
		}
		//ui.Flush()
		//	Sleep(AnimDurMediumLong)
		g.Dungeon.SetCell(g.Places.Shaedra, GroundCell)
		g.Dungeon.SetCell(g.Places.Marevor, ScrollCell)
		g.Objects.Scrolls[g.Places.Marevor] = ScrollExtended
		g.RescuedShaedra()
		md.story = 0
		md.mode = modeNormal
		return
	}
	md.story++
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
	switch md.story {
	case 0:
		md.mode = modeStory
		g.PrintStyled("You take and use the artifact.", logSpecial)
		g.Print("[(x) to continue].")
	case 1:
		g.Dungeon.SetCell(g.Places.Artifact, GroundCell)
		if !DisableAnimations {
			md.startAnimSeq()
			_, _, bg := md.positionDrawing(g.Places.Monolith)
			md.DrawAtPosition(g.Places.Monolith, false, 'Φ', ColorFgMagicPlace, bg)
			md.anims.Frame(AnimDurMediumLong)
		}
		g.Objects.Stairs[g.Places.Monolith] = WinStair
		g.Dungeon.SetCell(g.Places.Monolith, StairCell)
		if !DisableAnimations {
			md.startAnimSeq()
			_, _, bg := md.positionDrawing(g.Places.Marevor)
			md.DrawAtPosition(g.Places.Marevor, false, 'Φ', ColorFgMagicPlace, bg)
			md.anims.Frame(AnimDurMediumLong)
		}
		g.Objects.Story[g.Places.Marevor] = StoryMarevor
		g.PrintStyled("Marevor: “Great! Let's escape and find some bones to celebrate!”", logSpecial)
		g.PrintStyled("Syu: “Sorry, but I prefer bananas!”", logSpecial)
		g.Print("[(x) to continue]")
	case 2:
		g.Dungeon.SetCell(g.Places.Marevor, GroundCell)
		AchRetrievedArtifact.Get(g)
		md.story = 0
		md.mode = modeNormal
		return
	}
	md.story++
}
