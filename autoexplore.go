package main

import (
	"errors"
	"github.com/anaseto/gruid"
)

var DijkstraMapCache [DungeonNCells]int

func (g *game) Autoexplore() error {
	if mons := g.MonsterInLOS(); mons.Exists() {
		return errors.New("You cannot auto-explore while there are monsters in view.")
	}
	if g.ExclusionsMap[g.Player.Pos] {
		return errors.New("You cannot auto-explore while in an excluded area.")
	}
	if g.AllExplored() {
		return errors.New("Nothing left to explore.")
	}
	sources := g.AutoexploreSources()
	if len(sources) == 0 {
		return errors.New("Some excluded places remain unexplored.")
	}
	g.BuildAutoexploreMap(sources)
	n, finished := g.NextAuto()
	if finished || n == nil {
		return errors.New("You cannot reach safely some places.")
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

func (g *game) AutoexploreSources() []int {
	sources := []int{}
	np := &noisePath{state: g}
	it := g.Dungeon.Grid.Iterator()
	for i := 0; it.Next(); i++ {
		pos := it.P()
		c := cell(it.Cell())
		if c.IsWall() {
			if len(np.Neighbors(pos)) == 0 {
				continue
			}
		}
		if g.ExclusionsMap[pos] {
			continue
		}
		if !explored(c) {
			sources = append(sources, i)
		}

	}
	return sources
}

func (g *game) BuildAutoexploreMap(sources []int) {
	ap := &autoexplorePath{state: g}
	g.AutoExploreDijkstra(ap, sources)
	g.DijkstraMapRebuild = false
}

func (g *game) NextAuto() (next *gruid.Point, finished bool) {
	ap := &autoexplorePath{state: g}
	if DijkstraMapCache[idx(g.Player.Pos)] == unreachable {
		return nil, false
	}
	neighbors := ap.Neighbors(g.Player.Pos)
	if len(neighbors) == 0 {
		return nil, false
	}
	n := neighbors[0]
	ncost := DijkstraMapCache[idx(n)]
	for _, pos := range neighbors[1:] {
		cost := DijkstraMapCache[idx(pos)]
		if cost < ncost {
			n = pos
			ncost = cost
		}
	}
	if ncost >= DijkstraMapCache[idx(g.Player.Pos)] {
		finished = true
	}
	next = &n
	return next, finished
}
