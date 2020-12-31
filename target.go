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

func (ui *model) CancelExamine() {
	ui.st.Targeting = InvalidPos
	ui.st.Highlight = nil
	ui.st.MonsterTargLOS = nil
	ui.HideCursor()
	ui.mp.targeting = false
}

func (ui *model) Examine(pos gruid.Point) {
	if !valid(pos) {
		return
	}
	if ui.mp.ex == nil {
		ui.mp.ex = &examination{
			pos:     pos,
			objects: []gruid.Point{},
		}
	}
	ui.SetCursor(pos)
	ui.ComputeHighlight()
	m := ui.st.MonsterAt(pos)
	if m.Exists() && ui.st.Player.Sees(pos) {
		ui.st.ComputeMonsterCone(m)
	} else {
		ui.st.MonsterTargLOS = nil
	}
	ui.updatePosInfo()
}

func (ui *model) StartExamine() {
	ui.mp.targeting = true
	g := ui.st
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
	ui.mp.ex = &examination{
		pos:     pos,
		objects: []gruid.Point{},
	}
	if pos == g.Player.Pos {
		ui.NextObject(InvalidPos, ui.mp.ex)
		if !valid(ui.mp.ex.pos) {
			ui.NextStair(ui.mp.ex)
		}
		if valid(ui.mp.ex.pos) && Distance(pos, ui.mp.ex.pos) < DefaultLOSRange+5 {
			pos = ui.mp.ex.pos
		}
	}
	ui.Examine(pos)
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

func (m *model) DrawPosInfo() {
	g := m.st
	p := gruid.Point{}
	if m.mp.ex.pos.X <= DungeonWidth/2 {
		p.X += DungeonWidth + 1
	}
	info := m.mp.ex.info

	y := 2
	formatBox := func(title, s string, fg gruid.Color) {
		m.label.Box = &ui.Box{Title: ui.NewStyledText(title).WithStyle(gruid.Style{}.WithFg(fg))}
		m.label.SetText(s)
		y += m.label.Draw(m.gd.Slice(gruid.NewRange(0, y, DungeonWidth/2, 2+DungeonHeight).Add(p))).Size().Y
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

func (m *monster) Color(gs *state) gruid.Color {
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

func (m *model) updatePosInfo() {
	g := m.st
	pi := posInfo{}
	pos := m.mp.ex.pos
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
	m.mp.ex.info = pi
}

func (m *model) ComputeHighlight() {
	m.st.ComputePathHighlight(m.mp.ex.pos)
}

func (g *state) ComputePathHighlight(pos gruid.Point) {
	path := g.PlayerPath(g.Player.Pos, pos)
	g.Highlight = map[gruid.Point]bool{}
	for _, p := range path {
		g.Highlight[p] = true
	}
}

func (m *model) Target() error {
	g := m.st
	pos := m.mp.ex.pos
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

func (ui *model) NextMonster(key gruid.Key, pos gruid.Point, data *examination) {
	g := ui.st
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
	ui.Examine(pos)
}

func (ui *model) NextStair(data *examination) {
	g := ui.st
	if data.sortedStairs == nil {
		stairs := g.StairsSlice()
		data.sortedStairs = g.SortedNearestTo(stairs, g.Player.Pos)
	}
	if data.stairIndex >= len(data.sortedStairs) {
		data.stairIndex = 0
	}
	if len(data.sortedStairs) > 0 {
		ui.Examine(data.sortedStairs[data.stairIndex])
		data.stairIndex++
	}
}

func (ui *model) NextObject(pos gruid.Point, data *examination) {
	g := ui.st
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
	ui.Examine(pos)
}

func (ui *model) ExcludeZone(pos gruid.Point) {
	g := ui.st
	if !g.Dungeon.Cell(pos).Explored {
		g.Print("You cannot choose an unexplored cell for exclusion.")
	} else {
		toggle := !g.ExclusionsMap[pos]
		g.ComputeExclusion(pos, toggle)
	}
}
