package main

import (
	"errors"
	"github.com/anaseto/gruid"
)

func (g *game) Autoexplore() (again bool, err error) {
	if mons := g.MonsterInLOS(); mons.Exists() {
		return again, errors.New("You cannot auto-explore while there are monsters in view.")
	}
	if g.ExclusionsMap[g.Player.P] {
		return again, errors.New("You cannot auto-explore while in an excluded area.")
	}
	if g.AllExplored() {
		return again, errors.New("Nothing left to explore.")
	}
	sources := g.AutoexploreSources()
	if len(sources) == 0 {
		return again, errors.New("Some excluded places remain unexplored.")
	}
	g.BuildAutoexploreMap(sources)
	n, finished := g.NextAuto()
	if finished || n == nil {
		return again, errors.New("You cannot reach safely some places.")
	}
	g.Autoexploring = true
	g.AutoHalt = false
	return g.PlayerBump(*n)
}

func (g *game) AllExplored() bool {
	np := &noisePath{state: g}
	it := g.Dungeon.Grid.Iterator()
	for it.Next() {
		pos := it.P()
		c := cell(it.Cell())
		if c.IsWall() {
			if len(np.Neighbors(pos)) == 0 {
				continue
			}
		}
		if !explored(c) {
			return false
		}
	}
	return true
}

func (g *game) AutoexploreSources() []gruid.Point {
	g.autosources = g.autosources[:0]
	np := &noisePath{state: g}
	it := g.Dungeon.Grid.Iterator()
	for it.Next() {
		p := it.P()
		c := cell(it.Cell())
		if c.IsWall() {
			if len(np.Neighbors(p)) == 0 {
				continue
			}
		}
		if g.ExclusionsMap[p] {
			continue
		}
		if !explored(c) {
			g.autosources = append(g.autosources, p)
		}

	}
	return g.autosources
}

const unreachable = 9999

func (g *game) BuildAutoexploreMap(sources []gruid.Point) {
	ap := &autoexplorePath{state: g}
	g.PR.BreadthFirstMap(ap, sources, unreachable)
	g.AutoexploreMapRebuild = false
}

func (g *game) NextAuto() (next *gruid.Point, finished bool) {
	ap := &autoexplorePath{state: g}
	if g.PR.BreadthFirstMapAt(g.Player.P) > unreachable {
		return nil, false
	}
	neighbors := ap.Neighbors(g.Player.P)
	if len(neighbors) == 0 {
		return nil, false
	}
	n := neighbors[0]
	ncost := g.PR.BreadthFirstMapAt(n)
	for _, pos := range neighbors[1:] {
		cost := g.PR.BreadthFirstMapAt(pos)
		if cost < ncost {
			n = pos
			ncost = cost
		}
	}
	if ncost >= g.PR.BreadthFirstMapAt(g.Player.P) {
		finished = true
	}
	next = &n
	return next, finished
}
