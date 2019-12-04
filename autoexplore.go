package main

import "errors"

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
	np := &noisePath{game: g}
	for i, c := range g.Dungeon.Cells {
		pos := idxtopos(i)
		if c.IsWall() {
			if len(np.Neighbors(pos)) == 0 {
				continue
			}
		}
		if !c.Explored {
			return false
		}
	}
	return true
}

func (g *game) AutoexploreSources() []int {
	sources := []int{}
	np := &noisePath{game: g}
	for i, c := range g.Dungeon.Cells {
		pos := idxtopos(i)
		if c.IsWall() {
			if len(np.Neighbors(pos)) == 0 {
				continue
			}
		}
		if g.ExclusionsMap[pos] {
			continue
		}
		if !c.Explored {
			sources = append(sources, i)
		}

	}
	return sources
}

func (g *game) BuildAutoexploreMap(sources []int) {
	ap := &autoexplorePath{game: g}
	g.AutoExploreDijkstra(ap, sources)
	g.DijkstraMapRebuild = false
}

func (g *game) NextAuto() (next *position, finished bool) {
	ap := &autoexplorePath{game: g}
	if DijkstraMapCache[g.Player.Pos.idx()] == unreachable {
		return nil, false
	}
	neighbors := ap.Neighbors(g.Player.Pos)
	if len(neighbors) == 0 {
		return nil, false
	}
	n := neighbors[0]
	ncost := DijkstraMapCache[n.idx()]
	for _, pos := range neighbors[1:] {
		cost := DijkstraMapCache[pos.idx()]
		if cost < ncost {
			n = pos
			ncost = cost
		}
	}
	if ncost >= DijkstraMapCache[g.Player.Pos.idx()] {
		finished = true
	}
	next = &n
	return next, finished
}
