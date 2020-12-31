package main

import (
	"errors"

	"github.com/anaseto/gruid"
)

var InvalidPos = gruid.Point{-1, -1}

type examination struct {
	pos          gruid.Point
	nmonster     int
	objects      []gruid.Point
	nobject      int
	sortedStairs []gruid.Point
	stairIndex   int
}

func (ui *model) CancelExamine() {
	ui.st.Targeting = InvalidPos
	ui.st.Highlight = nil
	ui.st.MonsterTargLOS = nil
	ui.HideCursor()
	ui.mp.targeting = false
}

func (ui *model) Examine(pos gruid.Point) {
	if ui.mp.ex == nil {
		ui.mp.ex = &examination{
			pos:     pos,
			objects: []gruid.Point{},
		}
	}
	ui.mp.ex.pos = pos
	ui.ComputeHighlight()
	ui.SetCursor(pos)
	m := ui.st.MonsterAt(pos)
	if m.Exists() && ui.st.Player.Sees(pos) {
		ui.st.ComputeMonsterCone(m)
	} else {
		ui.st.MonsterTargLOS = nil
	}
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
	data.pos = pos
	data.nmonster = nmonster
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
		data.pos = data.sortedStairs[data.stairIndex]
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
	data.pos = pos
	data.nobject = nobject
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
