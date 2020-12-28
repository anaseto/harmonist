package main

import (
	"sort"
)

type dungeonPath struct {
	dungeon   *dungeon
	neighbors [8]gruid.Point
	wcost     int
}

func (dp *dungeonPath) Neighbors(pos gruid.Point) []gruid.Point {
	nb := dp.neighbors[:0]
	return pos.CardinalNeighbors(nb, func(npos gruid.Point) bool { return npos.valid() })
}

func (dp *dungeonPath) Cost(from, to gruid.Point) int {
	if dp.dungeon.Cell(to).T == WallCell {
		if dp.wcost > 0 {
			return dp.wcost
		}
		return 4
	}
	return 1
}

func (dp *dungeonPath) Estimation(from, to gruid.Point) int {
	return from.Distance(to)
}

type gridPath struct {
	dungeon   *dungeon
	neighbors [4]gruid.Point
}

func (gp *gridPath) Neighbors(pos gruid.Point) []gruid.Point {
	nb := gp.neighbors[:0]
	return pos.CardinalNeighbors(nb, func(npos gruid.Point) bool { return npos.valid() })
}

func (gp *gridPath) Cost(from, to gruid.Point) int {
	return 1
}

func (gp *gridPath) Estimation(from, to gruid.Point) int {
	return from.Distance(to)
}

type mappingPath struct {
	state      *state
	neighbors [8]gruid.Point
}

func (dp *mappingPath) Neighbors(pos gruid.Point) []gruid.Point {
	d := dp.state.Dungeon
	if d.Cell(pos).T == WallCell {
		return nil
	}
	nb := dp.neighbors[:0]
	keep := func(npos gruid.Point) bool {
		return npos.valid()
	}
	return pos.CardinalNeighbors(nb, keep)
}

func (dp *mappingPath) Cost(from, to gruid.Point) int {
	return 1
}

func (dp *mappingPath) Estimation(from, to gruid.Point) int {
	return from.Distance(to)
}

type tunnelPath struct {
	dg        *dgen
	neighbors [4]gruid.Point
	area      [9]gruid.Point
}

func (tp *tunnelPath) Neighbors(pos gruid.Point) []gruid.Point {
	nb := tp.neighbors[:0]
	return pos.CardinalNeighbors(nb, func(npos gruid.Point) bool { return npos.valid() })
}

func (tp *tunnelPath) Cost(from, to gruid.Point) int {
	if tp.dg.room[from] && !tp.dg.tunnel[from] {
		return 50
	}
	cost := 1
	c := tp.dg.d.Cell(from)
	if tp.dg.room[from] {
		cost += 7
	} else if !tp.dg.tunnel[from] && c.T != GroundCell {
		cost++
	}
	if c.IsPassable() {
		return cost
	}
	wc := tp.dg.WallAreaCount(tp.area[:0], from, 1)
	return cost + 8 - wc
}

func (tp *tunnelPath) Estimation(from, to gruid.Point) int {
	return from.Distance(to)
}

type playerPath struct {
	state      *state
	neighbors [8]gruid.Point
	goal      gruid.Point
}

func (pp *playerPath) Neighbors(pos gruid.Point) []gruid.Point {
	d := pp.state.Dungeon
	nb := pp.neighbors[:0]
	keep := func(npos gruid.Point) bool {
		t, okT := pp.state.TerrainKnowledge[npos]
		if cld, ok := pp.state.Clouds[npos]; ok && cld == CloudFire && (!okT || t != FoliageCell && t != DoorCell) {
			return false
		}
		return npos.valid() && d.Cell(npos).Explored && (d.Cell(npos).T.IsPlayerPassable() && !okT ||
			okT && t.IsPlayerPassable() ||
			pp.state.Player.HasStatus(StatusLevitation) && (t == BarrierCell || t == ChasmCell) ||
			pp.state.Player.HasStatus(StatusDig) && (d.Cell(npos).T.IsDiggable() && !okT || (okT && t.IsDiggable())))
	}
	nb = pos.CardinalNeighbors(nb, keep)
	sort.Slice(nb, func(i, j int) bool {
		return nb[i].MaxCardinalDist(pp.goal) <= nb[j].MaxCardinalDist(pp.goal)
	})
	return nb
}

func (pp *playerPath) Cost(from, to gruid.Point) int {
	if !pp.state.ExclusionsMap[from] && pp.state.ExclusionsMap[to] {
		return unreachable
	}
	return 1
}

func (pp *playerPath) Estimation(from, to gruid.Point) int {
	return from.Distance(to)
}

type jumpPath struct {
	state      *state
	neighbors [8]gruid.Point
}

func (jp *jumpPath) Neighbors(pos gruid.Point) []gruid.Point {
	nb := jp.neighbors[:0]
	keep := func(npos gruid.Point) bool {
		return jp.state.PlayerCanPass(npos)
	}
	nb = pos.CardinalNeighbors(nb, keep)
	nb = ShufflePos(nb)
	return nb
}

func (jp *jumpPath) Cost(from, to gruid.Point) int {
	return 1
}

func (jp *jumpPath) Estimation(from, to gruid.Point) int {
	return from.Distance(to)
}

type noisePath struct {
	state      *state
	neighbors [8]gruid.Point
}

func (fp *noisePath) Neighbors(pos gruid.Point) []gruid.Point {
	nb := fp.neighbors[:0]
	d := fp.state.Dungeon
	keep := func(npos gruid.Point) bool {
		return npos.valid() && d.Cell(npos).T != WallCell
	}
	return pos.CardinalNeighbors(nb, keep)
}

func (fp *noisePath) Cost(from, to gruid.Point) int {
	return 1
}

type autoexplorePath struct {
	state      *state
	neighbors [8]gruid.Point
}

func (ap *autoexplorePath) Neighbors(pos gruid.Point) []gruid.Point {
	if ap.state.ExclusionsMap[pos] {
		return nil
	}
	d := ap.state.Dungeon
	nb := ap.neighbors[:0]
	keep := func(npos gruid.Point) bool {
		t, okT := ap.state.TerrainKnowledge[npos]
		if cld, ok := ap.state.Clouds[npos]; ok && cld == CloudFire && (!okT || t != FoliageCell && t != DoorCell) {
			// XXX little info leak
			return false
		}
		return npos.valid() && (d.Cell(npos).T.IsPlayerPassable() && (!okT || t != WallCell)) &&
			!ap.state.ExclusionsMap[npos]
	}
	nb = pos.CardinalNeighbors(nb, keep)
	return nb
}

func (ap *autoexplorePath) Cost(from, to gruid.Point) int {
	return 1
}

type monPath struct {
	state      *state
	monster   *monster
	neighbors [8]gruid.Point
}

func ShufflePos(ps []gruid.Point) []gruid.Point {
	for i := 0; i < len(ps); i++ {
		j := i + RandInt(len(ps)-i)
		ps[i], ps[j] = ps[j], ps[i]
	}
	return ps
}

func (mp *monPath) Neighbors(pos gruid.Point) []gruid.Point {
	nb := mp.neighbors[:0]
	keep := func(npos gruid.Point) bool {
		return mp.monster.CanPassDestruct(mp.state, npos)
	}
	ret := pos.CardinalNeighbors(nb, keep)
	// shuffle so that monster movement is not unnaturally predictable
	ret = ShufflePos(ret)
	return ret
}

func (mp *monPath) Cost(from, to gruid.Point) int {
	g := mp.state
	mons := g.MonsterAt(to)
	if !mons.Exists() {
		c := g.Dungeon.Cell(to)
		if mp.monster.Kind == MonsEarthDragon && c.IsDestructible() && !mp.monster.Status(MonsConfused) {
			return 5
		}
		if to == g.Player.Pos && mp.monster.Kind.Peaceful() {
			switch mp.monster.Kind {
			case MonsEarthDragon:
				return 1
			default:
				return 4
			}
		}
		if mp.monster.Kind.Patrolling() && mp.monster.State != Hunting && !c.IsNormalPatrolWay() {
			return 4
		}
		return 1
	}
	if mons.Status(MonsLignified) {
		return 8
	}
	return 6
}

func (mp *monPath) Estimation(from, to gruid.Point) int {
	return from.Distance(to)
}

func (m *monster) APath(g *state, from, to gruid.Point) []gruid.Point {
	mp := &monPath{state: g, monster: m}
	path, _, found := AstarPath(mp, from, to)
	if !found {
		return nil
	}
	return path
}

func (g *state) PlayerPath(from, to gruid.Point) []gruid.Point {
	pp := &playerPath{state: g, goal: to}
	path, _, found := AstarPath(pp, from, to)
	if !found {
		return nil
	}
	return path
}

func (g *state) SortedNearestTo(cells []gruid.Point, to gruid.Point) []gruid.Point {
	ps := posSlice{}
	for _, pos := range cells {
		pp := &dungeonPath{dungeon: g.Dungeon, wcost: unreachable}
		_, cost, found := AstarPath(pp, pos, to)
		if found {
			ps = append(ps, posCost{pos, cost})
		}
	}
	sort.Sort(ps)
	sorted := []gruid.Point{}
	for _, pc := range ps {
		sorted = append(sorted, pc.pos)
	}
	return sorted
}

type posCost struct {
	pos  gruid.Point
	cost int
}

type posSlice []posCost

func (ps posSlice) Len() int           { return len(ps) }
func (ps posSlice) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }
func (ps posSlice) Less(i, j int) bool { return ps[i].cost < ps[j].cost }
