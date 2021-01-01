package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/ui"
)

var InvalidPos = gruid.Point{-1, -1}

type examination struct {
	pos          gruid.Point
	nmonster     int
	objects      []gruid.Point
	nobject      int
	sortedStairs []gruid.Point
	stairIndex   int
	info         posInfo
}

func (md *model) CancelExamine() {
	md.g.Targeting = InvalidPos
	md.g.Highlight = nil
	md.g.MonsterTargLOS = nil
	md.HideCursor()
	md.mp.targeting = false
}

func (md *model) Examine(pos gruid.Point) {
	if !valid(pos) {
		return
	}
	if md.mp.ex == nil {
		md.mp.ex = &examination{
			pos:     pos,
			objects: []gruid.Point{},
		}
	}
	md.SetCursor(pos)
	md.ComputeHighlight()
	m := md.g.MonsterAt(pos)
	if m.Exists() && md.g.Player.Sees(pos) {
		md.g.ComputeMonsterCone(m)
	} else {
		md.g.MonsterTargLOS = nil
	}
	md.updatePosInfo()
}

func (md *model) StartExamine() {
	md.mp.targeting = true
	g := md.g
	pos := g.Player.Pos
	minDist := 999
	for _, mons := range g.Monsters {
		if mons.Exists() && g.Player.LOS[mons.Pos] {
			dist := Distance(mons.Pos, g.Player.Pos)
			if minDist > dist {
				minDist = dist
				pos = mons.Pos
			}
		}
	}
	md.mp.ex = &examination{
		pos:     pos,
		objects: []gruid.Point{},
	}
	if pos == g.Player.Pos {
		md.NextObject(InvalidPos, md.mp.ex)
		if !valid(md.mp.ex.pos) {
			md.NextStair(md.mp.ex)
		}
		if valid(md.mp.ex.pos) && Distance(pos, md.mp.ex.pos) < DefaultLOSRange+5 {
			pos = md.mp.ex.pos
		}
	}
	md.Examine(pos)
}

type posInfo struct {
	Pos         gruid.Point
	Unknown     bool
	Noise       bool
	Unreachable bool
	Sees        bool
	Player      bool
	Monster     *monster
	Cell        string
	Cloud       string
	Lighted     bool
}

func (md *model) DrawPosInfo() {
	g := md.g
	p := gruid.Point{}
	if md.mp.ex.pos.X <= DungeonWidth/2 {
		p.X += DungeonWidth + 1
	}
	info := md.mp.ex.info

	y := 2
	formatBox := func(title, s string, fg gruid.Color) {
		md.label.Box = &ui.Box{Title: ui.NewStyledText(title).WithStyle(gruid.Style{}.WithFg(fg))}
		md.label.SetText(s)
		y += md.label.Draw(md.gd.Slice(gruid.NewRange(0, y, DungeonWidth/2, 2+DungeonHeight).Add(p))).Size().Y
	}

	features := []string{}
	if !info.Unknown {
		features = append(features, info.Cell)
		if info.Cloud != "" && info.Sees {
			features = append(features, info.Cloud)
		}
		if info.Lighted && info.Sees {
			features = append(features, "(lighted)")
		}
	} else {
		features = append(features, "unknown place")
	}
	if info.Noise {
		features = append(features, "noise")
	}
	if info.Unreachable {
		features = append(features, "unreachable")
	}
	t := " Terrain Features"
	if !info.Sees && !info.Unknown {
		t += " (seen)"
	} else if info.Unknown {
		t += " (unknown)"
	}
	fg := ColorFg
	if info.Unreachable {
		t += " - Unreachable"
		fg = ColorOrange
	}
	formatBox(t+" ", strings.Join(features, ", "), fg)

	if info.Player {
		formatBox("Player", "This is you.", ColorBlue)
	}

	mons := info.Monster
	if !mons.Exists() {
		return
	}
	title := fmt.Sprintf(" %s %s (%s %s) ", mons.Kind, mons.State, mons.Dir.String())
	fg = mons.Color(g)
	var mdesc []string

	statuses := mons.StatusesText()
	if statuses != "" {
		mdesc = append(mdesc, "Statuses: %s", statuses)
	}
	mdesc = append(mdesc, "Traits: "+mons.Traits())
	formatBox(title, strings.Join(mdesc, "\n"), fg)
}

func (m *monster) Color(gs *game) gruid.Color {
	var fg gruid.Color
	if m.Status(MonsLignified) {
		fg = ColorFgLignifiedMonster
	} else if m.Status(MonsConfused) {
		fg = ColorFgConfusedMonster
	} else if m.Status(MonsParalysed) {
		fg = ColorFgParalysedMonster
	} else if m.State == Resting {
		fg = ColorFgSleepingMonster
	} else if m.State == Hunting {
		fg = ColorFgMonster
	} else if m.Peaceful(gs) {
		fg = ColorFgPlayer
	} else {
		fg = ColorFgWanderingMonster
	}
	return fg
}

func (m *monster) StatusesText() string {
	infos := []string{}
	for st, i := range m.Statuses {
		if i > 0 {
			infos = append(infos, fmt.Sprintf("%s %d", monsterStatus(st), m.Statuses[monsterStatus(st)]))
		}
	}
	return strings.Join(infos, ", ")
}

func (m *monster) Traits() string {
	var info string
	info += fmt.Sprintf("Their size is %s.", m.Kind.Size())
	if m.Kind.Peaceful() {
		info += " " + fmt.Sprint("They are peaceful.")
	}
	if m.Kind.CanOpenDoors() {
		info += " " + fmt.Sprint("They can open doors.")
	}
	if m.Kind.CanFly() {
		info += " " + fmt.Sprint("They can fly.")
	}
	if m.Kind.CanSwim() {
		info += " " + fmt.Sprint("They can swim.")
	}
	if m.Kind.ShallowSleep() {
		info += " " + fmt.Sprint("They have very shallow sleep.")
	}
	if m.Kind.ResistsLignification() {
		info += " " + fmt.Sprint("They are unaffected by lignification.")
	}
	if m.Kind.ReflectsTeleport() {
		info += " " + fmt.Sprint("They partially reflect back oric teleport magic.")
	}
	if m.Kind.GoodFlair() {
		info += " " + fmt.Sprint("They have good flair.")
	}
	return info
}

func (md *model) updatePosInfo() {
	g := md.g
	pi := posInfo{}
	pos := md.mp.ex.pos
	pi.Pos = pos
	switch {
	case !g.Dungeon.Cell(pos).Explored:
		pi.Unknown = true
		if g.Noise[pos] || g.NoiseIllusion[pos] {
			pi.Noise = true
		}
		return
		//case !targ.Reachable(g, pos):
		//pi.Unreachable = true
		//return
	}
	mons := g.MonsterAt(pos)
	if pos == g.Player.Pos {
		pi.Player = true
	}
	if g.Player.Sees(pos) {
		pi.Sees = true
	}
	c := g.Dungeon.Cell(pos)
	if t, ok := g.TerrainKnowledge[pos]; ok {
		c.T = t
	}
	if mons.Exists() && g.Player.Sees(pos) {
		pi.Monster = mons
	}
	if cld, ok := g.Clouds[pos]; ok && g.Player.Sees(pos) {
		pi.Cloud = cld.String()
	}
	pi.Cell = c.ShortDesc(g, pos)
	if g.Illuminated[idx(pos)] && c.IsIlluminable() && g.Player.Sees(pos) {
		pi.Lighted = true
	}
	if g.Noise[pos] || g.NoiseIllusion[pos] {
		pi.Noise = true
	}
	md.mp.ex.info = pi
}

func (md *model) ComputeHighlight() {
	md.g.ComputePathHighlight(md.mp.ex.pos)
}

func (g *game) ComputePathHighlight(pos gruid.Point) {
	path := g.PlayerPath(g.Player.Pos, pos)
	g.Highlight = map[gruid.Point]bool{}
	for _, p := range path {
		g.Highlight[p] = true
	}
}

func (md *model) Target() error {
	g := md.g
	pos := md.mp.ex.pos
	if !g.Dungeon.Cell(pos).Explored {
		return errors.New("You do not know this place.")
	}
	if g.Dungeon.Cell(pos).T == WallCell && !g.Player.HasStatus(StatusDig) {
		return errors.New("You cannot travel into a wall.")
	}
	path := g.PlayerPath(g.Player.Pos, pos)
	if len(path) == 0 {
		return errors.New("There is no safe path to this place.")
	}
	if c := g.Dungeon.Cell(pos); c.Explored && c.T != WallCell {
		g.AutoTarget = pos
		g.Targeting = pos
		return nil
	}
	return errors.New("Invalid destination.")
}

func (md *model) NextMonster(key gruid.Key, pos gruid.Point, data *examination) {
	g := md.g
	nmonster := data.nmonster
	for i := 0; i < len(g.Monsters); i++ {
		if key == "+" {
			nmonster++
		} else {
			nmonster--
		}
		if nmonster > len(g.Monsters)-1 {
			nmonster = 0
		} else if nmonster < 0 {
			nmonster = len(g.Monsters) - 1
		}
		mons := g.Monsters[nmonster]
		if mons.Exists() && g.Player.LOS[mons.Pos] && pos != mons.Pos {
			pos = mons.Pos
			break
		}
	}
	data.nmonster = nmonster
	md.Examine(pos)
}

func (md *model) NextStair(data *examination) {
	g := md.g
	if data.sortedStairs == nil {
		stairs := g.StairsSlice()
		data.sortedStairs = g.SortedNearestTo(stairs, g.Player.Pos)
	}
	if data.stairIndex >= len(data.sortedStairs) {
		data.stairIndex = 0
	}
	if len(data.sortedStairs) > 0 {
		md.Examine(data.sortedStairs[data.stairIndex])
		data.stairIndex++
	}
}

func (md *model) NextObject(pos gruid.Point, data *examination) {
	g := md.g
	nobject := data.nobject
	if len(data.objects) == 0 {
		for p := range g.Objects.Stairs {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.FakeStairs {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.Stones {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.Barrels {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.Magaras {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.Bananas {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.Items {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.Scrolls {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.Potions {
			data.objects = append(data.objects, p)
		}
		data.objects = g.SortedNearestTo(data.objects, g.Player.Pos)
	}
	for i := 0; i < len(data.objects); i++ {
		p := data.objects[nobject]
		nobject++
		if nobject > len(data.objects)-1 {
			nobject = 0
		}
		if g.Dungeon.Cell(p).Explored {
			pos = p
			break
		}
	}
	data.nobject = nobject
	md.Examine(pos)
}

func (md *model) ExcludeZone(pos gruid.Point) {
	g := md.g
	if !g.Dungeon.Cell(pos).Explored {
		g.Print("You cannot choose an unexplored cell for exclusion.")
	} else {
		toggle := !g.ExclusionsMap[pos]
		g.ComputeExclusion(pos, toggle)
	}
}
