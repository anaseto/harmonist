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
	p          gruid.Point
	nmonster     int
	objects      []gruid.Point
	nobject      int
	sortedStairs []gruid.Point
	stairIndex   int
	info         posInfo
	scroll       bool
}

// HideCursor hides the target cursor.
func (md *model) HideCursor() {
	md.mp.ex.p = InvalidPos
}

// SetCursor sets the target cursor.
func (md *model) SetCursor(pos gruid.Point) {
	md.mp.ex.p = pos
}

// CancelExamine cancels current targeting.
func (md *model) CancelExamine() {
	md.g.Highlight = nil
	md.g.MonsterTargLOS = nil
	md.HideCursor()
	md.mp.kbTargeting = false
	md.mp.ex.scroll = false
}

// Examine targets a given position with the cursor.
func (md *model) Examine(p gruid.Point) {
	if md.mp.ex.p == p {
		return
	}
	md.examine(p)
}

func (md *model) examine(p gruid.Point) {
	if !valid(p) {
		return
	}
	md.SetCursor(p)
	md.computeHighlight()
	m := md.g.MonsterAt(p)
	if m.Exists() && md.g.Player.Sees(p) {
		md.g.ComputeMonsterCone(m)
	} else {
		md.g.MonsterTargLOS = nil
	}
	md.updatePosInfo()
	md.mp.ex.scroll = false
}

// KeyboardExamine starts keyboard examination mode, with a sensible default
// target.
func (md *model) KeyboardExamine() {
	md.mp.kbTargeting = true
	g := md.g
	pos := g.Player.P
	minDist := 999
	for _, mons := range g.Monsters {
		if mons.Exists() && g.Player.LOS[mons.P] {
			dist := Distance(mons.P, g.Player.P)
			if minDist > dist {
				minDist = dist
				pos = mons.P
			}
		}
	}
	md.mp.ex = &examination{
		p:     pos,
		objects: []gruid.Point{},
	}
	if pos == g.Player.P {
		md.nextObject(InvalidPos, md.mp.ex)
		if !valid(md.mp.ex.p) {
			md.nextStair(md.mp.ex)
		}
		if valid(md.mp.ex.p) && Distance(pos, md.mp.ex.p) < DefaultLOSRange+5 {
			pos = md.mp.ex.p
		}
	}
	md.examine(pos)
}

type posInfo struct {
	Pos         gruid.Point
	Unknown     bool
	Noise       bool
	Unreachable bool
	Sees        bool
	Player      bool
	Monster     *monster
	Cell        cell
	Cloud       string
	Lighted     bool
}

func (md *model) drawPosInfo() {
	g := md.g
	p := gruid.Point{}
	if md.mp.ex.p.X <= DungeonWidth/2 {
		p.X += DungeonWidth/2 + 1
	}
	info := md.mp.ex.info

	y := 2
	formatBox := func(title, s string, fg gruid.Color) {
		md.description.Content = md.description.Content.WithText(s).Format(DungeonWidth/2 - 2)
		if md.description.Content.Size().Y+2 > 2+DungeonHeight {
			md.description.Box = &ui.Box{Title: ui.NewStyledText(title, gruid.Style{}.WithFg(fg)),
				Footer: ui.Text("scroll/page down for more...")}
		} else {
			md.description.Box = &ui.Box{Title: ui.NewStyledText(title, gruid.Style{}.WithFg(fg))}
		}
		y += md.description.Draw(md.gd.Slice(gruid.NewRange(0, y, DungeonWidth/2, 2+DungeonHeight).Add(p))).Size().Y
	}

	features := []string{}
	if !info.Unknown {
		features = append(features, info.Cell.ShortString(g, info.Pos))
		if info.Cloud != "" && info.Sees {
			features = append(features, info.Cloud)
		}
		if info.Lighted && info.Sees {
			features = append(features, "lighted")
		}
	} else {
		features = append(features, "unknown")
	}
	if info.Noise {
		features = append(features, "noise")
	}
	if info.Unreachable {
		features = append(features, "unreachable")
	}
	if !info.Sees && !info.Unknown {
		features = append(features, "seen")
	} else if info.Unknown {
		features = append(features, "unexplored")
	}
	t := features[0]
	if len(features) > 1 {
		t += " (" + strings.Join(features[1:], ", ") + ")"
	}
	fg := ColorFg
	if info.Unreachable {
		fg = ColorOrange
	}
	desc := ""
	if info.Unknown {
		desc = "You do not know what is in there."
	} else {
		desc = info.Cell.Desc(g, info.Pos)
	}

	if info.Player {
		if !md.mp.ex.scroll {
			formatBox(t+" ", desc, fg)
		}
		formatBox("Player", "This is you.", ColorBlue)
		return
	}

	mons := info.Monster
	if !mons.Exists() {
		formatBox(t+" ", desc, fg)
		return
	}
	title := fmt.Sprintf("%s (%s %s)", mons.Kind, mons.State, mons.Dir.String())
	fg = mons.color(g)
	var mdesc []string

	statuses := mons.statusesText()
	if statuses != "" {
		mdesc = append(mdesc, "Statuses: %s", statuses)
	}
	mdesc = append(mdesc, "Traits: "+mons.traits())
	if !md.mp.ex.scroll {
		formatBox(t+" ", desc, fg)
	}
	formatBox(title, strings.Join(mdesc, "\n"), fg)
}

func (m *monster) color(gs *game) gruid.Color {
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

func (m *monster) statusesText() string {
	infos := []string{}
	for st, i := range m.Statuses {
		if i > 0 {
			infos = append(infos, fmt.Sprintf("%s %d", monsterStatus(st), m.Statuses[monsterStatus(st)]))
		}
	}
	return strings.Join(infos, ", ")
}

func (m *monster) traits() string {
	var info string
	info += fmt.Sprintf("Their size is %s.", m.Kind.Size())
	if m.Kind.Peaceful() {
		info += "They are peaceful."
	}
	if m.Kind.CanOpenDoors() {
		info += " " + "They can open doors."
	}
	if m.Kind.CanFly() {
		info += " " + "They can fly."
	}
	if m.Kind.CanSwim() {
		info += " " + "They can swim."
	}
	if m.Kind.ShallowSleep() {
		info += " " + "They have very shallow sleep."
	}
	if m.Kind.ResistsLignification() {
		info += " " + "They are unaffected by lignification."
	}
	if m.Kind.ReflectsTeleport() {
		info += " " + "They partially reflect back oric teleport magic."
	}
	if m.Kind.GoodFlair() {
		info += " " + "They have good flair."
	}
	return info
}

func (md *model) updatePosInfo() {
	g := md.g
	pi := posInfo{}
	pos := md.mp.ex.p
	pi.Pos = pos
	switch {
	case !explored(g.Dungeon.Cell(pos)):
		pi.Unknown = true
		if g.Noise[pos] || g.NoiseIllusion[pos] {
			pi.Noise = true
		}
		md.mp.ex.info = pi
		return
		//case !targ.Reachable(g, pos):
		//pi.Unreachable = true
		//return
	}
	mons := g.MonsterAt(pos)
	if pos == g.Player.P {
		pi.Player = true
	}
	if g.Player.Sees(pos) {
		pi.Sees = true
	}
	c := g.Dungeon.Cell(pos)
	if t, ok := g.TerrainKnowledge[pos]; ok {
		c = t | c&Explored
	}
	if mons.Exists() && g.Player.Sees(pos) {
		pi.Monster = mons
	}
	if cld, ok := g.Clouds[pos]; ok && g.Player.Sees(pos) {
		pi.Cloud = cld.String()
	}
	pi.Cell = c
	if g.Illuminated(pos) && c.IsIlluminable() && g.Player.Sees(pos) {
		pi.Lighted = true
	}
	if g.Noise[pos] || g.NoiseIllusion[pos] {
		pi.Noise = true
	}
	md.mp.ex.info = pi
}

func (md *model) computeHighlight() {
	md.g.computePathHighlight(md.mp.ex.p)
}

func (g *game) computePathHighlight(pos gruid.Point) {
	path := g.PlayerPath(g.Player.P, pos)
	g.Highlight = map[gruid.Point]bool{}
	for _, p := range path {
		g.Highlight[p] = true
	}
}

func (md *model) target() error {
	g := md.g
	pos := md.mp.ex.p
	if !explored(g.Dungeon.Cell(pos)) {
		return errors.New("You do not know this place.")
	}
	if terrain(g.Dungeon.Cell(pos)) == WallCell && !g.Player.HasStatus(StatusDig) {
		return errors.New("You cannot travel into a wall.")
	}
	path := g.PlayerPath(g.Player.P, pos)
	if len(path) == 0 {
		return errors.New("There is no safe path to this place.")
	}
	if c := g.Dungeon.Cell(pos); explored(c) && terrain(c) != WallCell {
		g.AutoTarget = pos
		return nil
	}
	return errors.New("Invalid destination.")
}

func (md *model) nextMonster(key gruid.Key, pos gruid.Point, data *examination) {
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
		if mons.Exists() && g.Player.LOS[mons.P] && pos != mons.P {
			pos = mons.P
			break
		}
	}
	data.nmonster = nmonster
	md.Examine(pos)
}

func (md *model) nextStair(data *examination) {
	g := md.g
	if data.sortedStairs == nil {
		stairs := g.StairsSlice()
		data.sortedStairs = g.SortedNearestTo(stairs, g.Player.P)
	}
	if data.stairIndex >= len(data.sortedStairs) {
		data.stairIndex = 0
	}
	if len(data.sortedStairs) > 0 {
		md.Examine(data.sortedStairs[data.stairIndex])
		data.stairIndex++
	}
}

func (md *model) nextObject(pos gruid.Point, data *examination) {
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
		data.objects = g.SortedNearestTo(data.objects, g.Player.P)
	}
	for i := 0; i < len(data.objects); i++ {
		p := data.objects[nobject]
		nobject++
		if nobject > len(data.objects)-1 {
			nobject = 0
		}
		if explored(g.Dungeon.Cell(p)) {
			pos = p
			break
		}
	}
	data.nobject = nobject
	md.Examine(pos)
}

func (md *model) excludeZone(pos gruid.Point) {
	g := md.g
	if !explored(g.Dungeon.Cell(pos)) {
		g.Print("You cannot choose an unexplored cell for exclusion.")
	} else {
		toggle := !g.ExclusionsMap[pos]
		g.ComputeExclusion(pos, toggle)
	}
}
