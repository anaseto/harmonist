package main

import "errors"

type Targeter interface {
	ComputeHighlight(*game, position)
	Action(*game, position) error
	Reachable(*game, position) bool
	Done() bool
}

type examiner struct {
	done   bool
	stairs bool
}

func (ex *examiner) ComputeHighlight(g *game, pos position) {
	g.ComputePathHighlight(pos)
}

func (g *game) ComputePathHighlight(pos position) {
	path := g.PlayerPath(g.Player.Pos, pos)
	g.Highlight = map[position]bool{}
	for _, p := range path {
		g.Highlight[p] = true
	}
}

func (ex *examiner) Action(g *game, pos position) error {
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

func (ex *examiner) Reachable(g *game, pos position) bool {
	return true
}

func (ex *examiner) Done() bool {
	return ex.done
}
