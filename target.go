package main

import "errors"

type Targeter interface {
	ComputeHighlight(*state, gruid.Point)
	Action(*state, gruid.Point) error
	Reachable(*state, gruid.Point) bool
	Done() bool
}

type examiner struct {
	done   bool
	stairs bool
}

func (ex *examiner) ComputeHighlight(g *state, pos gruid.Point) {
	g.ComputePathHighlight(pos)
}

func (g *state) ComputePathHighlight(pos gruid.Point) {
	path := g.PlayerPath(g.Player.Pos, pos)
	g.Highlight = map[gruid.Point]bool{}
	for _, p := range path {
		g.Highlight[p] = true
	}
}

func (ex *examiner) Action(g *state, pos gruid.Point) error {
	if !g.Dungeon.Cell(pos).Explored {
		return errors.New("You do not know this place.")
	}
	if g.Dungeon.Cell(pos).T == WallCell && !g.Player.HasStatus(StatusDig) {
		return errors.New("You cannot travel into a wall.")
	}
	path := g.PlayerPath(g.Player.Pos, pos)
	if len(path) == 0 {
		if ex.stairs {
			return errors.New("There is no safe path to the nearest stairs.")
		}
		return errors.New("There is no safe path to this place.")
	}
	if c := g.Dungeon.Cell(pos); c.Explored && c.T != WallCell {
		g.AutoTarget = pos
		g.Targeting = pos
		ex.done = true
		return nil
	}
	return errors.New("Invalid destination.")
}

func (ex *examiner) Reachable(g *state, pos gruid.Point) bool {
	return true
}

func (ex *examiner) Done() bool {
	return ex.done
}
