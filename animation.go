package main

import (
	//"log"
	"sort"
	"time"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
	"github.com/anaseto/gruid/rl"
)

const (
	AnimDurShort       = 25 * time.Millisecond
	AnimDurShortMedium = 50 * time.Millisecond
	AnimDurMedium      = 75 * time.Millisecond
	AnimDurMediumLong  = 100 * time.Millisecond
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
		if md.g.Highlight[p] || p == md.targ.ex.p {
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

func (a *Animations) DrawReverse(p gruid.Point, r rune, fg, bg gruid.Color) {
	c := a.grid.At(p)
	c.Rune = r
	c.Style.Fg = fg
	c.Style.Bg = bg
	c.Style.Attrs |= AttrReverse
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

func (md *model) SwappingAnimation(mp, pp gruid.Point) {
	if DisableAnimations {
		return
	}
	md.startAnimSeq()
	md.anims.Frame(AnimDurShort)
	_, fgm, bgColorm := md.positionDrawing(mp)
	_, _, bgColorp := md.positionDrawing(pp)
	md.anims.Draw(mp, 'Φ', fgm, bgColorp)
	md.anims.Draw(pp, 'Φ', ColorFgPlayer, bgColorm)
	md.anims.Frame(AnimDurMedium)
	md.anims.Draw(mp, 'Φ', ColorFgPlayer, bgColorp)
	md.anims.Draw(pp, 'Φ', fgm, bgColorm)
	md.anims.Frame(AnimDurMedium)
}

func (md *model) TeleportAnimation(from, to gruid.Point, showto bool) {
	if DisableAnimations {
		return
	}
	md.startAnimSeq()
	_, _, bgColorf := md.positionDrawing(from)
	_, _, bgColort := md.positionDrawing(to)
	if showto {
		md.anims.Draw(from, 'Φ', ColorBlue, bgColorf)
		md.anims.Draw(to, 'Φ', ColorCyan, bgColort)
		md.anims.Frame(AnimDurMediumLong)
	} else {
		md.anims.Draw(from, 'Φ', ColorCyan, bgColorf)
		md.anims.Frame(AnimDurMediumLong)
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
	for i := 0; i < len(ray); i++ {
		p := ray[i]
		or, fgColor, bgColor := md.positionDrawing(p)
		md.anims.Draw(p, r, fg, bgColor)
		md.anims.Frame(AnimDurShort)
		md.anims.Draw(p, or, fgColor, bgColor)
	}
}

func (md *model) WaveDrawAt(p gruid.Point, fg gruid.Color) {
	r, _, bgColor := md.positionDrawing(p)
	md.anims.DrawReverse(p, r, fg, bgColor)
}

func (md *model) ExplosionDrawAt(p gruid.Point, fg gruid.Color) {
	g := md.g
	_, _, bgColor := md.positionDrawing(p)
	mons := g.MonsterAt(p)
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
	if mons.Exists() || g.Player.P == p {
		r = '√'
	}
	md.anims.DrawReverse(p, r, fg, bgColor)
}

func (md *model) NoiseAnimation(noises []gruid.Point) {
	if DisableAnimations {
		return
	}
	md.LOSWavesAnimation(DefaultLOSRange, WaveMagicNoise, md.g.Player.P)
	//md.startAnimSeq()
	colors := []gruid.Color{ColorFgSleepingMonster, ColorFgMagicPlace}
	for i := 0; i < 2; i++ {
		for _, p := range noises {
			r := '♫'
			_, _, bgColor := md.positionDrawing(p)
			md.anims.DrawReverse(p, r, colors[i], bgColor)
		}
		_, _, bgColor := md.positionDrawing(md.g.Player.P)
		md.anims.DrawReverse(md.g.Player.P, '@', colors[i], bgColor)
		md.anims.Frame(AnimDurShortMedium)
	}

}

func (md *model) ExplosionAnimation(es explosionStyle, p gruid.Point) {
	if DisableAnimations {
		return
	}
	g := md.g
	md.startAnimSeq()
	md.anims.Frame(AnimDurShort)
	colors := [2]gruid.Color{ColorFgExplosionStart, ColorFgExplosionEnd}
	if es == WallExplosion || es == AroundWallExplosion {
		colors[0] = ColorFgExplosionWallStart
		colors[1] = ColorFgExplosionWallEnd
	}
	for i := 0; i < 3; i++ {
		nb := g.playerPassableNeighbors(p)
		if es != AroundWallExplosion {
			nb = append(nb, p)
		}
		for _, q := range nb {
			fg := colors[RandInt(2)]
			if !g.Player.LOS[q] {
				continue
			}
			md.ExplosionDrawAt(q, fg)
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
		dij = &noisePath{g: g}
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
		p := idxtopos(i)
		switch ws {
		case WaveConfusion:
			fg := ColorFgConfusedMonster
			if md.g.Player.Sees(p) {
				md.WaveDrawAt(p, fg)
			}
		case WaveSleeping:
			fg := ColorFgSleepingMonster
			if md.g.Player.Sees(p) {
				md.WaveDrawAt(p, fg)
			}
		case WaveSlowing:
			fg := ColorFgParalysedMonster
			if md.g.Player.Sees(p) {
				md.WaveDrawAt(p, fg)
			}
		case WaveTree:
			fg := ColorFgLignifiedMonster
			if md.g.Player.Sees(p) {
				md.WaveDrawAt(p, fg)
			}
		case WaveNoise:
			fg := ColorFgWanderingMonster
			if md.g.Player.Sees(p) {
				md.WaveDrawAt(p, fg)
			}
		case WaveMagicNoise:
			fg := ColorFgMagicPlace
			md.WaveDrawAt(p, fg)
		}
	}
	md.anims.Frame(AnimDurShort)
}

func (md *model) WallExplosionAnimation(p gruid.Point) {
	if DisableAnimations {
		return
	}
	md.startAnimSeq()
	colors := [2]gruid.Color{ColorFgExplosionWallStart, ColorFgExplosionWallEnd}
	for _, fg := range colors {
		_, _, bgColor := md.positionDrawing(p)
		//md.anims.Draw(pos, '☼', fg, bgColor)
		md.anims.DrawReverse(p, '%', fg, bgColor)
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
			p := ray[i]
			_, _, bgColor := md.positionDrawing(p)
			r := '*'
			if RandInt(2) == 0 {
				r = '×'
			}
			md.anims.DrawReverse(p, r, fg, bgColor)
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
			p := ray[i]
			_, _, bgColor := md.positionDrawing(p)
			r := '*'
			if RandInt(2) == 0 {
				r = '×'
			}
			md.anims.DrawReverse(p, r, fg, bgColor)
		}
		md.anims.Frame(AnimDurShortMedium)
	}
}

func (md *model) ProjectileSymbol(dir gruid.Point) (r rune) {
	switch dir {
	case gruid.Point{1, 0}, gruid.Point{-1, 0}:
		r = '—'
	case gruid.Point{1, -1}, gruid.Point{-1, 1}:
		r = '/'
	case gruid.Point{0, 1}, gruid.Point{0, -1}:
		r = '|'
	case gruid.Point{1, 1}, gruid.Point{-1, -1}:
		r = '\\'
	}
	return r
}

func (md *model) MonsterJavelinAnimation(ray []gruid.Point, hit bool) {
	if DisableAnimations {
		return
	}
	g := md.g
	md.startAnimSeq()
	md.anims.Frame(AnimDurShort)
	for i := 0; i < len(ray); i++ {
		p := ray[i]
		r, fgColor, bgColor := md.positionDrawing(p)
		md.anims.Draw(p, md.ProjectileSymbol(dirnorm(g.Player.P, p)), ColorFgMonster, bgColor)
		md.anims.Frame(AnimDurShort)
		md.anims.Draw(p, r, fgColor, bgColor)
	}
	//md.anims.Frame(AnimDurShort)
}

func (md *model) WoundedAnimation() {
	if DisableAnimations {
		return
	}
	g := md.g
	md.startAnimSeq()
	r, _, bg := md.positionDrawing(g.Player.P)
	md.anims.Draw(g.Player.P, r, ColorFgHPwounded, bg)
	md.anims.Frame(AnimDurShortMedium)
	if g.Player.HP <= 15 {
		md.anims.Draw(g.Player.P, r, ColorFgHPcritical, bg)
		md.anims.Frame(AnimDurShortMedium)
	}
}

func (md *model) PlayerGoodEffectAnimation() {
	if DisableAnimations {
		return
	}
	g := md.g
	md.startAnimSeq()
	md.anims.Frame(AnimDurShort)
	r, _, bg := md.positionDrawing(g.Player.P)
	md.anims.Draw(g.Player.P, r, ColorGreen, bg)
	md.anims.Frame(AnimDurShortMedium)
	md.anims.Draw(g.Player.P, r, ColorYellow, bg)
	md.anims.Frame(AnimDurShortMedium)
	//md.anims.Draw(g.Player.Pos, r, fg, bg)
}

func (md *model) StatusEndAnimation() {
	if DisableAnimations {
		return
	}
	g := md.g
	md.startAnimSeq()
	r, _, bg := md.positionDrawing(g.Player.P)
	md.anims.Draw(g.Player.P, r, ColorViolet, bg)
	md.anims.Frame(AnimDurShortMedium)
}

func (md *model) EffectAtPPAnimation() {
	if DisableAnimations {
		return
	}
	g := md.g
	md.startAnimSeq()
	r, _, bg := md.positionDrawing(g.Player.P)
	md.anims.Draw(g.Player.P, r, ColorCyan, bg)
	md.anims.Frame(AnimDurShortMedium)
}

func (md *model) FoundFakeStairsAnimation() {
	if DisableAnimations {
		return
	}
	g := md.g
	r, _, bg := md.positionDrawing(g.Player.P)
	md.anims.Draw(g.Player.P, r, ColorMagenta, bg)
	md.anims.Frame(AnimDurMediumLong)
}

func (md *model) MusicAnimation(p gruid.Point) {
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
	for _, p := range path[:len(path)-1] {
		md.anims.Draw(p, '×', ColorFgPlayer, bg)
	}
	md.anims.Draw(path[len(path)-1], '@', ColorFgPlayer, bg)
	md.anims.Frame(AnimDurMediumLong)
	//log.Print("anims: %+v", md.anims)
}

func (md *model) MagicMappingAnimation() {
	if DisableAnimations {
		return
	}
	md.startAnimSeq()
	md.anims.Frame(AnimDurShort)
}

func (md *model) AbyssFallAnimation() {
	if DisableAnimations {
		return
	}
	gd := rl.NewGrid(DungeonWidth, DungeonHeight)
	max := 5
	gd.FillFunc(func() rl.Cell {
		return rl.Cell(RandInt(max))
	})
	for i := 0; i < max; i++ {
		it := gd.Iterator()
		for it.Next() {
			if it.Cell() == rl.Cell(i) {
				md.anims.Draw(it.P(), ' ', ColorFg, ColorBgDark)
			}
		}
		md.anims.Frame(AnimDurShort)
	}
}

func (md *model) Story() {
	switch md.g.Depth {
	case WinDepth:
		md.FreeingShaedra()
	case MaxDepth:
		md.TakingArtifact()
	default:
		md.mode = modeNormal // should not happen
	}
}

func (md *model) FreeingShaedra() {
	g := md.g
	switch md.story {
	case 0:
		md.mode = modeStory
		g.Print("You see Shaedra. She is wounded!")
		g.PrintStyled("Shaedra: “Oh, it's you, Syu! Let's flee with Marevor's magara!”", logSpecial)
		md.confirm = true
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
		g.PrintStyled("[(x) to continue]", logConfirm)
	case 2:
		if !DisableAnimations {
			md.startAnimSeq()
			_, _, bg := md.positionDrawing(g.Places.Marevor)
			md.anims.Draw(g.Places.Marevor, 'Φ', ColorFgMagicPlace, bg)
			md.anims.Draw(g.Places.Shaedra, 'Φ', ColorFgMagicPlace, bg)
			md.anims.Frame(AnimDurMediumLong)
		}
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

func (md *model) TakingArtifact() {
	g := md.g
	switch md.story {
	case 0:
		md.mode = modeStory
		g.PrintStyled("You take and use the artifact.", logSpecial)
		md.confirm = true
	case 1:
		g.Dungeon.SetCell(g.Places.Artifact, GroundCell)
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
			_, _, bg := md.positionDrawing(g.Places.Marevor)
			md.anims.Draw(g.Places.Marevor, 'Φ', ColorFgMagicPlace, bg)
			md.anims.Frame(AnimDurMediumLong)
		}
		g.Objects.Story[g.Places.Marevor] = StoryMarevor
		g.PrintStyled("Marevor: “Great! Let's escape and find some bones to celebrate!”", logSpecial)
		g.PrintStyled("Syu: “Sorry, but I prefer bananas!”", logSpecial)
		g.PrintStyled("[(x) to continue]", logConfirm)
	case 2:
		g.Dungeon.SetCell(g.Places.Marevor, GroundCell)
		AchRetrievedArtifact.Get(g)
		md.story = 0
		md.mode = modeNormal
		return
	}
	md.story++
}
