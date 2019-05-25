// This file implements a line of sight algorithm.
//
// It works in a way that can remind of the Dijkstra algorithm, but within each
// cone between a diagonal and an orthogonal line, only movements along those
// two directions are allowed. This allows the algorithm to be a simple pass on
// squares around the player, starting from radius 1 until line of sight range.
//
// Going from a position from to a position pos has a cost, which depends
// essentially on the type of terrain in from. Some circumstances, such as
// being on top of a tree, can influence the cost of terrains.
//
// The obtained light rays are lines formed using at most two adjacent
// directions: a diagonal and an orthogonal one (for example north east and
// east).

package main

type raynode struct {
	Cost int
}

type rayMap map[position]raynode

func (g *game) BestParent(rm rayMap, from, pos position, rs raystyle) (position, int) {
	var parents [2]position
	p := parents[:0]
	p = pos.Parents(from, p)
	b := p[0]
	if len(p) > 1 && rm[p[1]].Cost+g.LOSCost(from, p[1], pos, rs) < rm[b].Cost+g.LOSCost(from, b, pos, rs) {
		b = p[1]
	}
	return b, rm[b].Cost + g.LOSCost(from, b, pos, rs)
}

func (g *game) DiagonalOpaque(from, to position, rs raystyle) bool {
	// The game uses cardinal movement only, so two diagonal walls should,
	// for example, block line of sight. This is in contrast with the main
	// mechanics of the line of sight algorithm, which for gameplay reasons
	// allows diagonals for light rays in normal circumstances.
	var cache [2]position
	p := cache[:0]
	switch to.Dir(from) {
	case NE:
		p = append(p, to.S(), to.W())
	case NW:
		p = append(p, to.S(), to.E())
	case SW:
		p = append(p, to.N(), to.E())
	case SE:
		p = append(p, to.N(), to.W())
	}
	count := 0
	for _, pos := range p {
		_, ok := g.Clouds[pos]
		if ok {
			count++
			continue
		}
		if !pos.valid() {
			continue
		}
		c := g.Dungeon.Cell(pos)
		switch c.T {
		case WallCell, HoledWallCell, WindowCell:
			count++
		case BarrierCell:
			if rs == TargetingRay {
				count++
			}
		}
	}
	return count > 1
}

func (g *game) DiagonalDifficult(from, to position) bool {
	// For reasons similar as in DiagonalOpaque, two diagonal foliage cells
	// should reduce range of line of sight in that diagonal direction.
	var cache [2]position
	p := cache[:0]
	switch to.Dir(from) {
	case NE:
		p = append(p, to.S(), to.W())
	case NW:
		p = append(p, to.S(), to.E())
	case SW:
		p = append(p, to.N(), to.E())
	case SE:
		p = append(p, to.N(), to.W())
	}
	count := 0
	for _, pos := range p {
		if !pos.valid() {
			continue
		}
		_, ok := g.Clouds[pos]
		if ok {
			count++
			continue
		}
		switch g.Dungeon.Cell(pos).T {
		case WallCell, FoliageCell, HoledWallCell:
			count++
		}
	}
	return count > 1
}

// LOSCost gives cost of expanding from 'pos' to 'to' light ray originated at
// 'from', for particular circumstances rs of light ray.
func (g *game) LOSCost(from, pos, to position, rs raystyle) int {
	var wallcost int
	switch rs {
	case TreePlayerRay:
		wallcost = TreeRange + 1
	case MonsterRay:
		wallcost = DefaultMonsterLOSRange + 1
	default:
		wallcost = g.LosRange() + 1
	}
	// diagonal costs
	if g.DiagonalOpaque(pos, to, rs) {
		return wallcost
	}
	// no terrain cost on origin
	if from == pos {
		if rs != TreePlayerRay && g.DiagonalDifficult(pos, to) {
			return wallcost - 1
		}
		return to.Distance(pos)
	}
	// pos terrain specific costs
	c := g.Dungeon.Cell(pos)
	if c.T == WallCell {
		return wallcost
	}
	if _, ok := g.Clouds[pos]; ok {
		return wallcost
	}
	if c.T == DoorCell {
		if pos != from {
			mons := g.MonsterAt(pos)
			if !mons.Exists() {
				return wallcost
			}
		}
	}
	if c.T == FoliageCell || c.T == HoledWallCell {
		switch rs {
		case TreePlayerRay:
			if c.T == FoliageCell {
				break
			}
			fallthrough
		default:
			return wallcost + to.Distance(pos) - 3
		}
	}
	if rs != TreePlayerRay && g.DiagonalDifficult(pos, to) {
		cost := wallcost - pos.Distance(from) - 1
		if cost < 1 {
			cost = 1
		}
		return cost
	}
	if rs == TreePlayerRay && c.T == WindowCell && from.Distance(pos) >= DefaultLOSRange {
		return wallcost - from.Distance(pos) - 1
	}
	return to.Distance(pos)
}

type raystyle int

const (
	NormalPlayerRay raystyle = iota
	MonsterRay
	TreePlayerRay
	TargetingRay
	LightRay
)

const LightRange = 6

func (g *game) BuildRayMap(from position, rs raystyle, rm rayMap) {
	var wallcost int
	switch rs {
	case TreePlayerRay:
		wallcost = TreeRange
	case MonsterRay:
		wallcost = DefaultMonsterLOSRange
	case LightRay:
		wallcost = LightRange
	default:
		wallcost = g.LosRange()
	}
	for k := range rm {
		delete(rm, k)
	}
	rm[from] = raynode{Cost: 0}
	for d := 1; d <= wallcost; d++ {
		for x := -d + from.X; x <= d+from.X; x++ {
			for _, pos := range []position{{x, from.Y + d}, {x, from.Y - d}} {
				if !pos.valid() {
					continue
				}
				_, c := g.BestParent(rm, from, pos, rs)
				rm[pos] = raynode{Cost: c}
			}
		}
		for y := -d + 1 + from.Y; y <= d-1+from.Y; y++ {
			for _, pos := range []position{{from.X + d, y}, {from.X - d, y}} {
				if !pos.valid() {
					continue
				}
				_, c := g.BestParent(rm, from, pos, rs)
				rm[pos] = raynode{Cost: c}
			}
		}
	}
}

const DefaultLOSRange = 12
const DefaultMonsterLOSRange = 12

func (g *game) LosRange() int {
	return DefaultLOSRange
}

func (g *game) StopAuto() {
	if g.Autoexploring && !g.AutoHalt {
		g.Print("You stop exploring.")
	} else if g.AutoDir != NoDir {
		g.Print("You stop.")
	} else if g.AutoTarget != InvalidPos {
		g.Print("You stop.")
	}
	g.AutoHalt = true
	g.AutoDir = NoDir
	g.AutoTarget = InvalidPos
}

const TreeRange = 50

func (g *game) ComputeLOS() {
	g.ComputeLights()
	for k := range g.Player.LOS {
		delete(g.Player.LOS, k)
	}
	c := g.Dungeon.Cell(g.Player.Pos)
	rs := NormalPlayerRay
	if c.T == TreeCell {
		rs = TreePlayerRay
	}
	g.BuildRayMap(g.Player.Pos, rs, g.Player.Rays)
	nb := make([]position, 8)
	for pos, n := range g.Player.Rays {
		if n.Cost <= DefaultLOSRange {
			g.Player.LOS[pos] = true
		} else if c.T == TreeCell && g.Illuminated[pos.idx()] && n.Cost <= TreeRange {
			if g.Dungeon.Cell(pos).T == WallCell {
				// this is just an approximation, but ok in practice
				nb = pos.Neighbors(nb, func(npos position) bool {
					if !npos.valid() || !g.Illuminated[npos.idx()] || g.Dungeon.Cell(npos).IsWall() {
						return false
					}
					node, ok := g.Player.Rays[npos]
					return ok && node.Cost < TreeRange
				})
				if len(nb) == 0 {
					continue
				}
			}
			g.Player.LOS[pos] = true
		}
	}
	for pos := range g.Player.LOS {
		if g.Player.Sees(pos) {
			g.SeePosition(pos)
		}
	}
	for _, mons := range g.Monsters {
		if mons.Exists() && g.Player.Sees(mons.Pos) {
			mons.ComputeLOS(g) // approximation of what the monster will see for player info purposes
			mons.UpdateKnowledge(g, mons.Pos)
			if mons.Seen {
				g.StopAuto()
				continue
			}
			mons.Seen = true
			g.Printf("You see %s (%v).", mons.Kind.Indefinite(false), mons.State)
			g.StopAuto()
		}
	}
}

func (m *monster) ComputeLOS(g *game) {
	if m.Kind.Peaceful() {
		return
	}
	for k := range m.LOS {
		delete(m.LOS, k)
	}
	losRange := DefaultMonsterLOSRange
	g.BuildRayMap(m.Pos, MonsterRay, g.RaysCache)
	for pos, n := range g.RaysCache {
		if pos == m.Pos {
			m.LOS[pos] = true
			continue
		}
		if n.Cost <= losRange && g.Dungeon.Cell(pos).T != BarrelCell {
			ppos, _ := g.BestParent(g.RaysCache, m.Pos, pos, MonsterRay)
			if !g.Dungeon.Cell(ppos).Hides() {
				m.LOS[pos] = true
			}
		}
	}
}

func (g *game) SeePosition(pos position) {
	c := g.Dungeon.Cell(pos)
	t, okT := g.TerrainKnowledge[pos]
	if !c.Explored {
		see := "see"
		if c.IsNotable() {
			g.Printf("You %s %s.", see, c.ShortDesc(g, pos))
			g.StopAuto()
		}
		g.Dungeon.SetExplored(pos)
		g.DijkstraMapRebuild = true
	} else {
		// XXX this can be improved to handle more terrain types changes
		if okT && t == WallCell && c.T != WallCell {
			g.Printf("There is no longer a wall there.")
			g.StopAuto()
			g.DijkstraMapRebuild = true
		}
		if cld, ok := g.Clouds[pos]; ok && cld == CloudFire && okT && (t == FoliageCell || t == DoorCell) {
			g.Printf("There are flames there.")
			g.StopAuto()
			g.DijkstraMapRebuild = true
		}
	}
	if okT {
		delete(g.TerrainKnowledge, pos)
		if c.T.IsPlayerPassable() {
			delete(g.MagicalBarriers, pos)
		}
	}
	if mons, ok := g.LastMonsterKnownAt[pos]; ok && (mons.Pos != pos || !mons.Exists()) {
		delete(g.LastMonsterKnownAt, pos)
		mons.LastKnownPos = InvalidPos
	}
	delete(g.NoiseIllusion, pos)
	if g.Objects.Story[pos] == StoryShaedra && !g.LiberatedShaedra && g.Player.Pos.Distance(pos) <= 2 && g.Player.Pos != g.Places.Marevor &&
		g.Player.Pos != g.Places.Monolith && g.Ev != nil {
		g.PushEvent(&simpleEvent{ERank: g.Ev.Rank(), EAction: ShaedraAnimation})
		g.LiberatedShaedra = true
	}
}

func (g *game) ComputeExclusion(pos position, toggle bool) {
	exclusionRange := g.LosRange()
	g.ExclusionsMap[pos] = toggle
	for d := 1; d <= exclusionRange; d++ {
		for x := -d + pos.X; x <= d+pos.X; x++ {
			for _, pos := range []position{{x, pos.Y + d}, {x, pos.Y - d}} {
				if !pos.valid() {
					continue
				}
				g.ExclusionsMap[pos] = toggle
			}
		}
		for y := -d + 1 + pos.Y; y <= d-1+pos.Y; y++ {
			for _, pos := range []position{{pos.X + d, y}, {pos.X - d, y}} {
				if !pos.valid() {
					continue
				}
				g.ExclusionsMap[pos] = toggle
			}
		}
	}
}

func (g *game) Ray(pos position) []position {
	if !g.Player.LOS[pos] {
		return nil
	}
	ray := []position{}
	for pos != g.Player.Pos {
		ray = append(ray, pos)
		pos, _ = g.BestParent(g.Player.Rays, g.Player.Pos, pos, TargetingRay)
	}
	return ray
}

func (g *game) ComputeRayHighlight(pos position) {
	g.Highlight = map[position]bool{}
	ray := g.Ray(pos)
	for _, p := range ray {
		g.Highlight[p] = true
	}
}

func (g *game) ComputeNoise() {
	dij := &noisePath{game: g}
	rg := DefaultLOSRange
	nm := Dijkstra(dij, []position{g.Player.Pos}, rg)
	count := 0
	for k := range g.Noise {
		delete(g.Noise, k)
	}
	rmax := 2
	if g.Player.Inventory.Body == CloakHear {
		rmax += 2
	}
	// TODO: maybe if they're close enough you could hear them breathe too, or something like that.
	nm.iter(g.Player.Pos, func(n *node) {
		pos := n.Pos
		if g.Player.Sees(pos) {
			return
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() && mons.State != Resting && mons.State != Watching && RandInt(rmax) > 0 {
			switch mons.Kind {
			case MonsMirrorSpecter, MonsSatowalgaPlant, MonsButterfly:
				if mons.Kind == MonsMirrorSpecter && g.Player.Inventory.Body == CloakHear {
					g.Noise[pos] = true
					g.Print("You hear an imperceptible air movement.")
					count++
				}
			case MonsWingedMilfid, MonsTinyHarpy:
				g.Noise[pos] = true
				g.Print("You hear the flapping of wings.")
				count++
			case MonsOricCelmist, MonsEarthDragon, MonsTreeMushroom:
				g.Noise[pos] = true
				g.Print("You hear heavy footsteps.")
				count++
			case MonsWorm:
				g.Noise[pos] = true
				g.Print("You hear a creep noise.")
				count++
			case MonsDog, MonsBlinkingFrog:
				g.Noise[pos] = true
				g.Print("You hear light footsteps.")
				count++
			default:
				g.Noise[pos] = true
				g.Print("You hear footsteps.")
				count++
			}
		}
	})
	if count > 0 {
		g.StopAuto()
	}
}

func (p *player) Sees(pos position) bool {
	//return pos == p.Pos || p.LOS[pos] && p.Dir.InViewCone(p.Pos, pos)
	return p.LOS[pos]
}

func (m *monster) SeesPlayer(g *game) bool {
	return m.Sees(g, g.Player.Pos) && g.Player.Sees(m.Pos)
}

func (m *monster) SeesLight(g *game, pos position) bool {
	if !(m.LOS[pos] && m.Dir.InViewCone(m.Pos, pos)) {
		return false
	}
	if m.State == Resting && m.Pos.Distance(pos) > 1 {
		return false
	}
	return true
}

func (m *monster) Sees(g *game, pos position) bool {
	var darkRange = 4
	if g.Player.Inventory.Body == CloakShadows {
		darkRange = 3
	}
	if g.Player.HasStatus(StatusShadows) {
		darkRange = 1
	}
	const tableRange = 1
	if !(m.LOS[pos] && m.Dir.InViewCone(m.Pos, pos)) {
		return false
	}
	if m.State == Resting && m.Pos.Distance(pos) > 1 {
		return false
	}
	c := g.Dungeon.Cell(pos)
	if (!g.Illuminated[pos.idx()] && !g.Player.HasStatus(StatusIlluminated) || !c.IsIlluminable()) && m.Pos.Distance(pos) > darkRange {
		return false
	}
	if c.T == TableCell && m.Pos.Distance(pos) > tableRange {
		return false
	}
	return true
}

func (g *game) ComputeMonsterLOS() {
	for k := range g.MonsterLOS {
		delete(g.MonsterLOS, k)
	}
	for _, mons := range g.Monsters {
		if !mons.Exists() || !g.Player.Sees(mons.Pos) {
			continue
		}
		for pos := range g.Player.LOS {
			if !g.Player.Sees(pos) {
				continue
			}
			if mons.Sees(g, pos) {
				g.MonsterLOS[pos] = true
			}
		}
	}
	if g.MonsterLOS[g.Player.Pos] {
		g.Player.Statuses[StatusUnhidden] = 1
		g.Player.Statuses[StatusHidden] = 0
	} else {
		g.Player.Statuses[StatusUnhidden] = 0
		g.Player.Statuses[StatusHidden] = 1
	}
	if g.Illuminated[g.Player.Pos.idx()] && g.Dungeon.Cell(g.Player.Pos).IsIlluminable() {
		g.Player.Statuses[StatusLight] = 1
	} else {
		g.Player.Statuses[StatusLight] = 0
	}
}

func (g *game) ComputeLights() {
	// XXX: could be optimized further to avoid unnecessary recalculations
	for i := 0; i < DungeonNCells; i++ {
		g.Illuminated[i] = false
	}
	for lpos, on := range g.Objects.Lights {
		if !on {
			continue
		}
		if lpos.Distance(g.Player.Pos) > DefaultLOSRange+LightRange && g.Dungeon.Cell(g.Player.Pos).T != TreeCell {
			continue
		}
		g.BuildRayMap(lpos, LightRay, g.RaysCache)
		for pos, n := range g.RaysCache {
			if n.Cost <= LightRange {
				g.Illuminated[pos.idx()] = true
			}
		}
	}
	for _, mons := range g.Monsters {
		if !mons.Exists() || mons.Kind != MonsButterfly {
			continue
		}
		if mons.Pos.Distance(g.Player.Pos) > DefaultLOSRange+LightRange && g.Dungeon.Cell(g.Player.Pos).T != TreeCell {
			continue
		}
		g.BuildRayMap(mons.Pos, LightRay, g.RaysCache)
		for pos, n := range g.RaysCache {
			if n.Cost <= LightRange {
				g.Illuminated[pos.idx()] = true
			}
		}
	}
}

func (g *game) ComputeMonsterCone(m *monster) {
	g.MonsterTargLOS = make(map[position]bool)
	for pos := range g.Player.LOS {
		if !g.Player.Sees(pos) {
			continue
		}
		if m.Sees(g, pos) {
			g.MonsterTargLOS[pos] = true
		}
	}
}

func (m *monster) UpdateKnowledge(g *game, pos position) {
	if mons, ok := g.LastMonsterKnownAt[pos]; ok {
		mons.LastKnownPos = InvalidPos
	}
	if m.LastKnownPos != InvalidPos {
		delete(g.LastMonsterKnownAt, m.LastKnownPos)
	}
	g.LastMonsterKnownAt[pos] = m
	m.LastSeenState = m.State
	m.LastKnownPos = pos
}
