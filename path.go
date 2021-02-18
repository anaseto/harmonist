package main

import (
	"sort"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
	"github.com/anaseto/gruid/rl"
)

type simplePath struct {
	passable  func(gruid.Point) bool
	neighbors paths.Neighbors
}

func newPather(passable func(gruid.Point) bool) *simplePath {
	return &simplePath{
		passable: passable,
	}
}

func (sp *simplePath) Neighbors(p gruid.Point) []gruid.Point {
	if !sp.passable(p) {
		return nil
	}
	return sp.neighbors.Cardinal(p, sp.passable)
}

type dungeonPath struct {
	dungeon *dungeon
	nbs     paths.Neighbors
	wcost   int
}

func (dp *dungeonPath) Neighbors(p gruid.Point) []gruid.Point {
	return dp.nbs.Cardinal(p, valid)
}

func (dp *dungeonPath) Cost(from, to gruid.Point) int {
	if terrain(dp.dungeon.Cell(to)) == WallCell {
		if dp.wcost > 0 {
			return dp.wcost
		}
		return 4
	}
	return 1
}

func (dp *dungeonPath) Estimation(from, to gruid.Point) int {
	return distance(from, to)
}

type gridPath struct {
	dungeon *dungeon
	nbs     paths.Neighbors
}

func (gp *gridPath) Neighbors(p gruid.Point) []gruid.Point {
	return gp.nbs.Cardinal(p, valid)
}

func (gp *gridPath) Cost(from, to gruid.Point) int {
	return 1
}

func (gp *gridPath) Estimation(from, to gruid.Point) int {
	return distance(from, to)
}

type mappingPath struct {
	g   *game
	nbs paths.Neighbors
}

func (dp *mappingPath) Neighbors(p gruid.Point) []gruid.Point {
	d := dp.g.Dungeon
	if terrain(d.Cell(p)) == WallCell {
		return nil
	}
	return dp.nbs.Cardinal(p, valid)
}

func (dp *mappingPath) Cost(from, to gruid.Point) int {
	return 1
}

func (dp *mappingPath) Estimation(from, to gruid.Point) int {
	return distance(from, to)
}

type tunnelPath struct {
	dg  *dgen
	nbs paths.Neighbors
}

func (tp *tunnelPath) Neighbors(p gruid.Point) []gruid.Point {
	return tp.nbs.Cardinal(p, valid)
}

func (tp *tunnelPath) Cost(from, to gruid.Point) int {
	if tp.dg.room[from] && !tp.dg.tunnel[from] {
		return 50
	}
	cost := 1
	c := tp.dg.d.Cell(from)
	if tp.dg.room[from] {
		cost += 7
	} else if !tp.dg.tunnel[from] && terrain(c) != GroundCell {
		cost++
	}
	if c.IsPassable() {
		return cost
	}
	wc := countWalls(tp.dg.d.Grid, from, 1, true)
	return cost + 8 - wc
}

func countWalls(gd rl.Grid, p gruid.Point, radius int, countOut bool) int {
	count := 0
	rg := gruid.Range{
		gruid.Point{p.X - radius, p.Y - radius},
		gruid.Point{p.X + radius + 1, p.Y + radius + 1},
	}
	if countOut {
		osize := rg.Size()
		rg = rg.Intersect(gd.Range())
		size := rg.Size()
		count += osize.X*osize.Y - size.X*size.Y
	} else {
		rg = rg.Intersect(gd.Range())
	}
	gd = gd.Slice(rg)
	count += gd.Count(rl.Cell(WallCell))
	return count
}

func (tp *tunnelPath) Estimation(from, to gruid.Point) int {
	return distance(from, to)
}

type playerPath struct {
	g    *game
	nbs  paths.Neighbors
	goal gruid.Point
}

func (g *game) ppPassable(p gruid.Point) bool {
	d := g.Dungeon
	t, okT := g.TerrainKnowledge[p]
	if cld, ok := g.Clouds[p]; ok && cld == CloudFire && (!okT || t != FoliageCell && t != DoorCell) {
		return false
	}
	return valid(p) && explored(d.Cell(p)) && (d.Cell(p).IsPlayerPassable() && !okT ||
		okT && t.IsPlayerPassable() ||
		g.Player.HasStatus(StatusLevitation) && (t == BarrierCell || t == ChasmCell) ||
		g.Player.HasStatus(StatusDig) && (d.Cell(p).IsDiggable() && !okT || (okT && t.IsDiggable())))
}

func (pp *playerPath) Neighbors(p gruid.Point) []gruid.Point {
	nbs := pp.nbs.Cardinal(p, pp.g.ppPassable)
	sort.Slice(nbs, func(i, j int) bool {
		return maxCardinalDist(nbs[i], pp.goal) <= maxCardinalDist(nbs[j], pp.goal)
	})
	return nbs
}

func (pp *playerPath) Cost(from, to gruid.Point) int {
	if !pp.g.ExclusionsMap[from] && pp.g.ExclusionsMap[to] {
		return unreachable
	}
	return 1
}

func (pp *playerPath) Estimation(from, to gruid.Point) int {
	return distance(from, to)
}

type jumpPath struct {
	g   *game
	nbs paths.Neighbors
}

func (jp *jumpPath) Neighbors(p gruid.Point) []gruid.Point {
	keep := func(q gruid.Point) bool {
		return jp.g.PlayerCanPass(q)
	}
	nbs := jp.nbs.Cardinal(p, keep)
	nbs = ShufflePos(nbs)
	return nbs
}

func (jp *jumpPath) Cost(from, to gruid.Point) int {
	return 1
}

func (jp *jumpPath) Estimation(from, to gruid.Point) int {
	return distance(from, to)
}

type noisePath struct {
	g   *game
	nbs paths.Neighbors
}

func (fp *noisePath) Neighbors(p gruid.Point) []gruid.Point {
	d := fp.g.Dungeon
	keep := func(q gruid.Point) bool {
		return valid(q) && terrain(d.Cell(q)) != WallCell
	}
	return fp.nbs.Cardinal(p, keep)
}

func (fp *noisePath) Cost(from, to gruid.Point) int {
	return 1
}

type autoexplorePath struct {
	g   *game
	nbs paths.Neighbors
}

func (ap *autoexplorePath) Neighbors(p gruid.Point) []gruid.Point {
	if ap.g.ExclusionsMap[p] {
		return nil
	}
	d := ap.g.Dungeon
	keep := func(q gruid.Point) bool {
		t, okT := ap.g.TerrainKnowledge[q]
		if cld, ok := ap.g.Clouds[q]; ok && cld == CloudFire && (!okT || t != FoliageCell && t != DoorCell) {
			// XXX little info leak
			return false
		}
		if !valid(q) {
			return false
		}
		c := d.Cell(q)
		return c.IsPlayerPassable() && (!okT && !c.IsWall() || !t.IsWall()) &&
			!ap.g.ExclusionsMap[q]
	}
	nbs := ap.nbs.Cardinal(p, keep)
	return nbs
}

func (ap *autoexplorePath) Cost(from, to gruid.Point) int {
	return 1
}

type monPath struct {
	g       *game
	monster *monster
	nbs     paths.Neighbors
}

func ShufflePos(ps []gruid.Point) []gruid.Point {
	for i := 0; i < len(ps); i++ {
		j := i + RandInt(len(ps)-i)
		ps[i], ps[j] = ps[j], ps[i]
	}
	return ps
}

func (mp *monPath) Neighbors(p gruid.Point) []gruid.Point {
	keep := func(q gruid.Point) bool {
		return mp.monster.CanPassDestruct(mp.g, q)
	}
	nbs := mp.nbs.Cardinal(p, keep)
	// shuffle so that monster movement is not unnaturally predictable
	nbs = ShufflePos(nbs)
	return nbs
}

func (mp *monPath) Cost(from, to gruid.Point) int {
	g := mp.g
	mons := g.MonsterAt(to)
	if !mons.Exists() {
		c := g.Dungeon.Cell(to)
		if mp.monster.Kind == MonsEarthDragon && c.IsDestructible() && !mp.monster.Status(MonsConfused) {
			return 5
		}
		if to == g.Player.P && mp.monster.Kind.Peaceful() {
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
	return distance(from, to)
}

func (m *monster) APath(g *game, from, to gruid.Point) []gruid.Point {
	mp := &monPath{g: g, monster: m}
	path := g.PR.AstarPath(mp, from, to)
	if len(path) == 0 {
		return nil
	}
	return path
}

func (g *game) PlayerPath(from, to gruid.Point) []gruid.Point {
	path := []gruid.Point{}
	if !g.ExclusionsMap[from] && g.ExclusionsMap[to] {
		pp := &playerPath{g: g, goal: to}
		path = g.PR.AstarPath(pp, from, to)
	} else {
		path = g.PR.JPSPath(path, from, to, g.ppPassable, false)
	}
	if len(path) == 0 {
		return nil
	}
	return path
}

func (g *game) SortedNearestTo(cells []gruid.Point, to gruid.Point) []gruid.Point {
	ps := posSlice{}
	for _, p := range cells {
		pp := &dungeonPath{dungeon: g.Dungeon, wcost: unreachable}
		path := g.PR.AstarPath(pp, p, to) // TODO: use JPS?
		if len(path) > 0 {
			ps = append(ps, posCost{p, len(path)})
		}
	}
	sort.Sort(ps)
	sorted := []gruid.Point{}
	for _, pc := range ps {
		sorted = append(sorted, pc.p)
	}
	return sorted
}

type posCost struct {
	p    gruid.Point
	cost int
}

type posSlice []posCost

func (ps posSlice) Len() int           { return len(ps) }
func (ps posSlice) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }
func (ps posSlice) Less(i, j int) bool { return ps[i].cost < ps[j].cost }
