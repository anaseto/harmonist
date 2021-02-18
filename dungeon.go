// many ideas here from articles found at http://www.roguebasin.com/

// TODO: some algorithms could use gruid's rl package, though it may not be
// worth the trouble.

package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
	"github.com/anaseto/gruid/rl"
)

type dungeon struct {
	Grid rl.Grid
}

func (d *dungeon) Cell(p gruid.Point) cell {
	return cell(d.Grid.At(p))
}

func (d *dungeon) Border(p gruid.Point) bool {
	return p.X == DungeonWidth-1 || p.Y == DungeonHeight-1 || p.X == 0 || p.Y == 0
}

func (d *dungeon) SetCell(p gruid.Point, c cell) {
	oc := d.Cell(p)
	d.Grid.Set(p, rl.Cell(c|oc&Explored))
}

func (d *dungeon) SetExplored(p gruid.Point) {
	oc := d.Cell(p)
	d.Grid.Set(p, rl.Cell(oc|Explored))
}

const maxIterations = 1000

func (dg *dgen) WallCell() gruid.Point {
	d := dg.d
	count := 0
	for {
		count++
		if count > maxIterations {
			panic("WallCell")
		}
		x := dg.rand.Intn(DungeonWidth)
		y := dg.rand.Intn(DungeonHeight)
		p := gruid.Point{x, y}
		c := d.Cell(p)
		if terrain(c) == WallCell {
			return p
		}
	}
}

func (d *dungeon) HasFreeNeighbor(nbs *paths.Neighbors, p gruid.Point) bool {
	nb := nbs.All(p, valid)
	for _, q := range nb {
		if d.Cell(q).IsPassable() {
			return true
		}
	}
	return false
}

func (dg *dgen) HasTooManyWallNeighbors(p gruid.Point) bool {
	d := dg.d
	nbs := dg.neighbors.All(p, valid)
	count := 8 - len(nbs)
	for _, q := range nbs {
		if !d.Cell(q).IsPassable() {
			count++
		}
	}
	return count > 1
}

func (g *game) HasNonWallExploredNeighbor(p gruid.Point) bool {
	d := g.Dungeon
	neighbors := g.cardinalNeighbors(p)
	for _, q := range neighbors {
		c := d.Cell(q)
		if t, ok := g.TerrainKnowledge[q]; ok {
			c = t
		}
		if !c.IsWall() && explored(c) {
			return true
		}
	}
	return false
}

type roomSlice []*room

func (rs roomSlice) Len() int      { return len(rs) }
func (rs roomSlice) Swap(i, j int) { rs[i], rs[j] = rs[j], rs[i] }
func (rs roomSlice) Less(i, j int) bool {
	center := gruid.Point{DungeonWidth / 2, DungeonHeight / 2}
	ipos := rs[i].p
	ipos.X += rs[i].w / 2
	ipos.Y += rs[i].h / 2
	jpos := rs[j].p
	jpos.X += rs[j].w / 2
	jpos.Y += rs[j].h / 2
	return rs[i].special || !rs[j].special && distance(ipos, center) <= distance(jpos, center)
}

type dgen struct {
	d         *dungeon
	tunnel    map[gruid.Point]bool
	room      map[gruid.Point]bool
	rooms     []*room
	spl       places
	special   specialRoom
	layout    maplayout
	PR        *paths.PathRange
	rand      *rand.Rand
	neighbors paths.Neighbors
}

func (dg *dgen) ConnectRoomsShortestPath(i, j int) bool {
	if i == j {
		return false
	}
	r1 := dg.rooms[i]
	r2 := dg.rooms[j]
	var e1pos, e2pos gruid.Point
	var e1i, e2i int
	e1i = r1.UnusedEntry()
	e1pos = r1.entries[e1i].p
	e2i = r2.UnusedEntry()
	e2pos = r2.entries[e2i].p
	tp := &tunnelPath{dg: dg}
	path := dg.PR.AstarPath(tp, e1pos, e2pos)
	if len(path) == 0 {
		log.Println(fmt.Sprintf("no path from %v to %v", e1pos, e2pos))
		return false
	}
	for _, p := range path {
		if !valid(p) {
			panic(fmt.Sprintf("gruid.Point %v from %v to %v", p, e1pos, e2pos))
		}
		t := terrain(dg.d.Cell(p))
		if t == WallCell || t == ChasmCell || t == GroundCell || t == FoliageCell {
			dg.d.SetCell(p, GroundCell)
			dg.tunnel[p] = true
		}
	}
	r1.entries[e1i].used = true
	r2.entries[e2i].used = true
	r1.tunnels++
	r2.tunnels++
	return true
}

func (dg *dgen) NewRoom(rpos gruid.Point, kind string) *room {
	r := &room{p: rpos, vault: &rl.Vault{}}
	err := r.vault.Parse(kind)
	if err != nil {
		log.Printf("bad vault:%v", err)
		return nil
	}
	r.w = r.vault.Size().X
	r.h = r.vault.Size().Y
	drev := 2
	if r.w > r.h+2 {
		drev += r.w - r.h - 2
		if drev > 4 {
			drev = 4
		}
	}
	if dg.rand.Intn(drev) == 0 {
		switch dg.rand.Intn(2) {
		case 0:
			r.vault.Reflect()
			r.vault.Rotate(1 + 2*dg.rand.Intn(2))
		default:
			r.vault.Rotate(1 + 2*dg.rand.Intn(2))
		}
	} else {
		switch dg.rand.Intn(2) {
		case 0:
			r.vault.Reflect()
			r.vault.Rotate(2 * dg.rand.Intn(2))
		default:
			r.vault.Rotate(2 * dg.rand.Intn(2))
		}
	}
	r.w = r.vault.Size().X
	r.h = r.vault.Size().Y
	if !r.HasSpace(dg) {
		return nil
	}
	r.Dig(dg)
	if r.w == 0 || r.h == 0 {
		log.Printf("bad vault size: %v", r.vault.Size())
	}
	return r
}

func (dg *dgen) nearestConnectedRoom(i int) (k int) {
	r := dg.rooms[i]
	d := unreachable
	for j, nextRoom := range dg.rooms[:i] {
		if j == i {
			continue
		}
		nd := roomDistance(r, nextRoom)
		if nd < d {
			d = nd
			k = j
		}
	}
	return k
}

func (dg *dgen) nearRoom(i int) (k int) {
	r := dg.rooms[i]
	d := unreachable
	for j, nextRoom := range dg.rooms {
		if j == i {
			continue
		}
		nd := roomDistance(r, nextRoom)
		if nd < d {
			n := dg.rand.Intn(5)
			if n > 0 {
				d = nd
				k = j
			}
		}
	}
	return k
}

func (dg *dgen) PutDoors(g *game) {
	for _, r := range dg.rooms {
		for _, e := range r.entries {
			if e.used && !e.virtual {
				r.places = append(r.places, place{p: e.p, kind: PlaceDoor})
			}
		}
		for _, pl := range r.places {
			if pl.kind != PlaceDoor {
				continue
			}
			dg.d.SetCell(pl.p, DoorCell)
			r.places = append(r.places, place{p: pl.p, kind: PlaceDoor})
		}
	}
}

func (dg *dgen) PutHoledWalls(g *game, n int) {
	candidates := []gruid.Point{}
	it := dg.d.Grid.Iterator()
	for it.Next() {
		p := it.P()
		if dg.room[p] && g.HoledWallCandidate(p) {
			candidates = append(candidates, p)
		}
	}
	if len(candidates) == 0 {
		return
	}
	for i := 0; i < n; i++ {
		p := candidates[dg.rand.Intn(len(candidates))]
		g.Dungeon.SetCell(p, HoledWallCell)
	}
}

func (dg *dgen) PutWindows(g *game, n int) {
	candidates := []gruid.Point{}
	it := dg.d.Grid.Iterator()
	for it.Next() {
		p := it.P()
		if dg.room[p] && g.HoledWallCandidate(p) {
			candidates = append(candidates, p)
		}
	}
	if len(candidates) == 0 {
		return
	}
	for i := 0; i < n; i++ {
		p := candidates[dg.rand.Intn(len(candidates))]
		g.Dungeon.SetCell(p, WindowCell)
	}
}

func (g *game) HoledWallCandidate(p gruid.Point) bool {
	d := g.Dungeon
	if !valid(p) || !d.Cell(p).IsWall() {
		return false
	}
	return valid(p.Add(gruid.Point{-1, 0})) && valid(p.Add(gruid.Point{1, 0})) &&
		d.Cell(p.Add(gruid.Point{-1, 0})).IsWall() && d.Cell(p.Add(gruid.Point{1, 0})).IsWall() &&
		valid(p.Add(gruid.Point{0, -1})) && d.Cell(p.Add(gruid.Point{0, -1})).IsPassable() &&
		valid(p.Add(gruid.Point{0, 1})) && d.Cell(p.Add(gruid.Point{0, 1})).IsPassable() ||
		(valid(p.Add(gruid.Point{-1, 0})) && valid(p.Add(gruid.Point{1, 0})) &&
			d.Cell(p.Add(gruid.Point{-1, 0})).IsPassable() && d.Cell(p.Add(gruid.Point{1, 0})).IsPassable() &&
			valid(p.Add(gruid.Point{0, -1})) && d.Cell(p.Add(gruid.Point{0, -1})).IsWall() &&
			valid(p.Add(gruid.Point{0, 1})) && d.Cell(p.Add(gruid.Point{0, 1})).IsWall())
}

type placement int

const (
	PlacementRandom placement = iota
	PlacementCenter
	PlacementEdge
)

func (dg *dgen) GenRooms(templates []string, n int, pl placement) (ps []gruid.Point, ok bool) {
	if len(templates) == 0 {
		return nil, false
	}
	ok = true
	for i := 0; i < n; i++ {
		var r *room
		count := 500
		var p gruid.Point
		var tpl string
		for r == nil && count > 0 {
			count--
			switch pl {
			case PlacementRandom:
				p = gruid.Point{dg.rand.Intn(DungeonWidth - 1), dg.rand.Intn(DungeonHeight - 1)}
			case PlacementCenter:
				p = gruid.Point{DungeonWidth/2 - 4 + dg.rand.Intn(5), DungeonHeight/2 - 3 + dg.rand.Intn(4)}
			case PlacementEdge:
				if dg.rand.Intn(2) == 0 {
					p = gruid.Point{dg.rand.Intn(DungeonWidth / 4), dg.rand.Intn(DungeonHeight - 1)}
				} else {
					p = gruid.Point{3*DungeonWidth/4 + dg.rand.Intn(DungeonWidth/4) - 1, dg.rand.Intn(DungeonHeight - 1)}
				}
			}
			tpl = templates[dg.rand.Intn(len(templates))]
			r = dg.NewRoom(p, tpl)
		}
		if r != nil {
			switch pl {
			case PlacementCenter, PlacementEdge:
				r.special = true
			}
			dg.rooms = append(dg.rooms, r)
			ps = append(ps, p)
		} else {
			ok = false
		}
	}
	return ps, ok
}

func (dg *dgen) ConnectRooms() {
	sort.Sort(roomSlice(dg.rooms))
	for i, r := range dg.rooms {
		if i == 0 {
			continue
		}
		if r.tunnels > 0 {
			continue
		}
		nearest := dg.nearestConnectedRoom(i)
		ok := dg.ConnectRoomsShortestPath(nearest, i)
		for !ok {
			panic("connect rooms")
		}
	}
	extraTunnels := 6
	switch dg.layout {
	case RandomSmallWalkCaveUrbanised:
		extraTunnels = 7
	case NaturalCave:
		extraTunnels = 4
	}
	count := 0
	for i, r := range dg.rooms {
		if count >= extraTunnels {
			break
		}
		if r.tunnels > 1 {
			continue
		}
		count++
		dg.ConnectRoomsShortestPath(i, dg.nearRoom(i))
	}
	if count >= extraTunnels {
		return
	}
	for i, r := range dg.rooms {
		if count >= extraTunnels {
			break
		}
		if r.tunnels > 2 {
			continue
		}
		count++
		dg.ConnectRoomsShortestPath(i, dg.nearRoom(i))
	}
}

type maplayout int

const (
	AutomataCave maplayout = iota
	RandomWalkCave
	RandomWalkTreeCave
	RandomSmallWalkCaveUrbanised
	NaturalCave
)

func (dg *dgen) GenShaedraCell(g *game) {
	g.Objects.Story = map[gruid.Point]story{}
	g.Places.Shaedra = dg.spl.Shaedra
	g.Objects.Story[g.Places.Shaedra] = StoryShaedra
	g.Places.Monolith = dg.spl.Monolith
	g.Objects.Story[g.Places.Monolith] = NoStory
	g.Places.Marevor = dg.spl.Marevor
	g.Objects.Story[g.Places.Marevor] = NoStory
}

func (dg *dgen) GenArtifactPlace(g *game) {
	g.Objects.Story = map[gruid.Point]story{}
	g.Places.Artifact = dg.spl.Artifact
	g.Objects.Story[g.Places.Artifact] = StoryArtifactSealed
	g.Places.Monolith = dg.spl.Monolith
	g.Objects.Story[g.Places.Monolith] = NoStory
	g.Places.Marevor = dg.spl.Marevor
	g.Objects.Story[g.Places.Marevor] = NoStory
}

func (g *game) GenRoomTunnels(ml maplayout) {
	dg := dgen{}
	dg.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	dg.PR = paths.NewPathRange(gruid.NewRange(0, 0, DungeonWidth, DungeonHeight))
	dg.layout = ml
	d := &dungeon{}
	d.Grid = rl.NewGrid(DungeonWidth, DungeonHeight)
	dg.d = d
	dg.tunnel = make(map[gruid.Point]bool)
	dg.room = make(map[gruid.Point]bool)
	dg.rooms = []*room{}
	switch ml {
	case AutomataCave:
		dg.GenCellularAutomataCaveMap()
	case RandomWalkCave:
		dg.GenCaveMap(21 * 40)
	case RandomWalkTreeCave:
		dg.GenTreeCaveMap()
	case RandomSmallWalkCaveUrbanised:
		dg.GenCaveMap(20 * 10)
	case NaturalCave:
		if dg.rand.Intn(3) == 0 {
			dg.GenCellularAutomataCaveMap()
		} else {
			dg.GenCaveMap(21 * 47)
		}
	}
	var places []gruid.Point
	var nspecial = 4
	if sr := g.Params.Special[g.Depth]; sr != noSpecialRoom {
		nspecial--
		pl := PlacementEdge
		if dg.rand.Intn(3) == 0 {
			pl = PlacementCenter
		}
		dg.special = sr
		var ok bool
		count := 0
		for {
			places, ok = dg.GenRooms(sr.Templates(), 1, pl)
			count++
			if count > 150 {
				if g.Depth == WinDepth || g.Depth == MaxDepth {
					panic("special room")
				}
				break
			}
			if ok {
				break
			}
		}
	}
	if g.Depth == WinDepth {
		dg.GenShaedraCell(g)
		nspecial--
	} else if g.Depth == MaxDepth {
		dg.GenArtifactPlace(g)
		nspecial--
	}
	normal := 5
	if g.Depth < 3 {
		nspecial--
		normal--
	}
	switch ml {
	case RandomWalkCave:
		dg.GenRooms(roomBigTemplates, nspecial-1, PlacementRandom)
		dg.GenRooms(roomNormalTemplates, normal, PlacementRandom)
	case RandomWalkTreeCave:
		dg.GenRooms(roomBigTemplates, nspecial+1, PlacementRandom)
		dg.GenRooms(roomNormalTemplates, normal+2, PlacementRandom)
	case RandomSmallWalkCaveUrbanised:
		nspecial += 3
		dg.GenRooms(roomBigTemplates, nspecial, PlacementRandom)
		dg.GenRooms(roomNormalTemplates, normal+5, PlacementRandom)
	case NaturalCave:
		nspecial++
		if g.Depth == WinDepth {
			nspecial += 2
		}
		dg.GenRooms(roomBigTemplates, nspecial, PlacementRandom)
		dg.GenRooms(roomNormalTemplates, normal-3, PlacementRandom)
	default:
		dg.GenRooms(roomBigTemplates, nspecial, PlacementRandom)
		dg.GenRooms(roomNormalTemplates, normal+2, PlacementRandom)
	}
	dg.ConnectRooms()
	g.Dungeon = d
	dg.PutDoors(g)
	dg.PlayerStartCell(g, places)
	dg.ClearUnconnected(g)
	if dg.rand.Intn(10) > 0 {
		var c cell
		if dg.rand.Intn(5) > 1 {
			c = ChasmCell
		} else {
			c = WaterCell
		}
		dg.GenLake(c)
		if dg.rand.Intn(5) == 0 {
			dg.GenLake(c)
		}
	}
	if g.Depth < MaxDepth {
		if g.Params.Blocked[g.Depth] {
			dg.GenStairs(g, BlockedStair)
		} else {
			dg.GenStairs(g, NormalStair)
		}
		dg.GenFakeStairs(g)
	}
	for i := 0; i < 4+dg.rand.Intn(2); i++ {
		dg.GenBarrel(g)
	}
	dg.AddSpecial(g, ml)
	dg.PR.CCMapAll(newPather(func(p gruid.Point) bool {
		return valid(p) && g.Dungeon.Cell(p).IsPassable()
	}))
	dg.GenMonsters(g)
	dg.PutCavernCells(g)
	if dg.rand.Intn(2) == 0 {
		dg.GenQueenRock()
	}
}

func (dg *dgen) PutCavernCells(g *game) {
	d := dg.d
	// TODO: improve handling and placement of this
	it := dg.d.Grid.Iterator()
	for it.Next() {
		p := it.P()
		if terrain(cell(it.Cell())) == GroundCell && !dg.room[p] && !dg.tunnel[p] {
			d.SetCell(p, CavernCell)
		}
	}
}

func (dg *dgen) ClearUnconnected(g *game) {
	d := dg.d
	sp := newPather(func(p gruid.Point) bool { return d.Cell(p).IsPlayerPassable() })
	dg.PR.CCMap(sp, g.Player.P)
	mg := rl.MapGen{Grid: dg.d.Grid}
	mg.KeepCC(dg.PR, g.Player.P, rl.Cell(WallCell))
}

func (dg *dgen) AddSpecial(g *game, ml maplayout) {
	g.Objects.Stones = map[gruid.Point]stone{}
	if g.Params.Blocked[g.Depth] || g.Depth == MaxDepth {
		dg.GenBarrierStone(g)
	}
	bananas := 1
	bananas += g.Params.ExtraBanana[g.Depth]
	for i := 0; i < bananas; i++ {
		dg.GenBanana(g)
	}
	if !g.Params.NoMagara[g.Depth] {
		dg.GenMagara(g)
	}
	dg.GenItem(g)
	dg.GenPotion(g, MagicPotion)
	if g.Params.HealthPotion[g.Depth] {
		dg.GenPotion(g, HealthPotion)
	}
	dg.GenStones(g)
	ntables := 4
	switch ml {
	case AutomataCave, RandomWalkCave, NaturalCave:
		if dg.rand.Intn(3) == 0 {
			ntables++
		} else if dg.rand.Intn(10) == 0 {
			ntables--
		}
	case RandomWalkTreeCave:
		if dg.rand.Intn(4) > 0 {
			ntables++
		}
		if dg.rand.Intn(4) > 0 {
			ntables++
		}
	case RandomSmallWalkCaveUrbanised:
		ntables += 2
		if dg.rand.Intn(4) > 0 {
			ntables++
		}
	}
	if g.Params.Tables[g.Depth] {
		ntables += 2 + dg.rand.Intn(2)
	}
	for i := 0; i < ntables; i++ {
		dg.GenTable(g)
	}
	dg.GenLight(g)
	ntrees := 1
	switch ml {
	case AutomataCave:
		if dg.rand.Intn(4) == 0 {
			ntrees++
		} else if dg.rand.Intn(8) == 0 {
			ntrees--
		}
	case RandomWalkCave:
		if dg.rand.Intn(4) > 0 {
			ntrees++
		}
		if dg.rand.Intn(8) == 0 {
			ntrees++
		}
	case NaturalCave:
		ntrees++
		if dg.rand.Intn(2) > 0 {
			ntrees++
		}
	case RandomWalkTreeCave, RandomSmallWalkCaveUrbanised:
		if dg.rand.Intn(2) == 0 {
			ntrees--
		}
	}
	if g.Params.Trees[g.Depth] {
		ntrees += 2 + dg.rand.Intn(2)
	}
	for i := 0; i < ntrees; i++ {
		dg.GenTree(g)
	}
	nhw := 1
	if dg.rand.Intn(3) > 0 {
		nhw++
	}
	if g.Params.Holes[g.Depth] {
		nhw += 3 + dg.rand.Intn(2)
	}
	switch ml {
	case RandomSmallWalkCaveUrbanised:
		if dg.rand.Intn(4) > 0 {
			nhw++
		}
	}
	dg.PutHoledWalls(g, nhw)
	nwin := 1
	if nhw == 1 {
		nwin++
	}
	if g.Params.Windows[g.Depth] {
		nwin += 4 + dg.rand.Intn(3)
	}
	switch ml {
	case RandomSmallWalkCaveUrbanised:
		if dg.rand.Intn(4) > 0 {
			nwin++
		}
	}
	dg.PutWindows(g, nwin)
	if g.Params.Lore[g.Depth] {
		dg.PutLore(g)
	}
}

func (dg *dgen) PutLore(g *game) {
	p := InvalidPos
	count := 0
	for p == InvalidPos {
		count++
		if count > 2000 {
			panic("PutLore1")
		}
		p = dg.rooms[RandInt(len(dg.rooms))].RandomPlace(PlaceItem)
	}
	count = 0
	for {
		count++
		if count > maxIterations {
			panic("PutLore2")
		}
		i := dg.rand.Intn(len(LoreMessages))
		if g.GeneratedLore[i] {
			continue
		}
		g.GeneratedLore[i] = true
		g.Objects.Lore[p] = i
		g.Objects.Scrolls[p] = ScrollLore
		g.Dungeon.SetCell(p, ScrollCell)
		break
	}
}

func (dg *dgen) GenLight(g *game) {
	no := 2
	ni := 8
	switch dg.layout {
	case NaturalCave:
		no += dg.rand.Intn(2)
		ni += dg.rand.Intn(3)
	case AutomataCave, RandomWalkCave:
		ni += dg.rand.Intn(4)
	case RandomWalkTreeCave:
		no--
		ni += dg.rand.Intn(4)
	case RandomSmallWalkCaveUrbanised:
		no--
		no -= dg.rand.Intn(2)
		ni += 2
		ni += dg.rand.Intn(4)
	}
	for i := 0; i < no; i++ {
		p := dg.OutsideGroundCell(g)
		g.Dungeon.SetCell(p, LightCell)
	}
	for i := 0; i < ni; i++ {
		p := dg.rooms[RandInt(len(dg.rooms))].RandomPlaces(PlaceSpecialOrStatic)
		if p != InvalidPos {
			g.Dungeon.SetCell(p, LightCell)
		} else if dg.rand.Intn(10) > 0 {
			i--
		}
	}
	it := g.Dungeon.Grid.Iterator()
	for it.Next() {
		if it.Cell() == rl.Cell(LightCell) {
			g.Objects.Lights[it.P()] = true
		}
	}
	g.ComputeLights()
}

func (dg *dgen) PlayerStartCell(g *game, places []gruid.Point) {
	const far = 30
	r := dg.rooms[len(dg.rooms)-1]
loop:
	for i := len(dg.rooms) - 2; i >= 0; i-- {
		for _, p := range places {
			if distance(r.p, p) < far && distance(dg.rooms[i].p, p) >= far {
				r = dg.rooms[i]
				dg.rooms[len(dg.rooms)-1], dg.rooms[i] = dg.rooms[i], dg.rooms[len(dg.rooms)-1]
				break loop
			}
		}
	}
	g.Player.P = r.RandomPlace(PlacePatrol)
	switch g.Depth {
	case 1, 4:
	default:
		return
	}
	itpos := InvalidPos
	neighbors := g.playerPassableNeighbors(g.Player.P)
	for i := 0; i < len(neighbors); i++ {
		j := RandInt(len(neighbors) - i)
		neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
	}
loopnb:
	for _, npos := range neighbors {
		c := g.Dungeon.Cell(npos)
		if c.IsGround() {
			for _, pl := range r.places {
				if npos == pl.p {
					continue loopnb
				}
			}
			itpos = npos
			break
		}
	}
	if itpos == InvalidPos {
		itpos = r.RandomPlace(PlaceItem)
	}
	if itpos == InvalidPos {
		itpos = r.RandomPlaces(PlaceSpecialOrStatic)
		if itpos == InvalidPos {
			panic("no item")
		}
	}
	g.Dungeon.SetCell(itpos, ScrollCell)
	switch g.Depth {
	case 1:
		g.Objects.Scrolls[itpos] = ScrollStory
	case 4:
		g.Objects.Scrolls[itpos] = ScrollDayoriahMessage
	}
}

func (dg *dgen) GenBanana(g *game) {
	count := 0
	for {
		count++
		if count > maxIterations {
			panic("GenBanana")
		}
		x := dg.rand.Intn(DungeonWidth)
		y := dg.rand.Intn(DungeonHeight)
		p := gruid.Point{x, y}
		c := dg.d.Cell(p)
		if terrain(c) == GroundCell && !dg.room[p] {
			dg.d.SetCell(p, BananaCell)
			g.Objects.Bananas[p] = true
			break
		}
	}
}

func (dg *dgen) GenPotion(g *game, ptn potion) {
	count := 0
	p := InvalidPos
	for p == InvalidPos {
		count++
		if count > maxIterations {
			return
		}
		p = dg.rooms[RandInt(len(dg.rooms))].RandomPlace(PlaceItem)
	}
	dg.d.SetCell(p, PotionCell)
	g.Objects.Potions[p] = ptn
}

func (dg *dgen) OutsideGroundCell(g *game) gruid.Point {
	count := 0
	for {
		count++
		if count > 1500 {
			panic("OutsideGroundCell")
		}
		x := dg.rand.Intn(DungeonWidth)
		y := dg.rand.Intn(DungeonHeight)
		p := gruid.Point{x, y}
		mons := g.MonsterAt(p)
		if mons.Exists() {
			continue
		}
		c := dg.d.Cell(p)
		if terrain(c) == GroundCell && !dg.room[p] {
			return p
		}
	}
}

func (dg *dgen) OutsideCavernMiddleCell(g *game) gruid.Point {
	count := 0
	for {
		count++
		if count > 2500 {
			return InvalidPos
		}
		x := dg.rand.Intn(DungeonWidth)
		y := dg.rand.Intn(DungeonHeight)
		p := gruid.Point{x, y}
		mons := g.MonsterAt(p)
		if mons.Exists() {
			continue
		}
		c := dg.d.Cell(p)
		if terrain(c) == GroundCell && count < 400 || terrain(c) == FoliageCell && count < 350 {
			continue
		}
		if (c.IsGround() || terrain(c) == FoliageCell) && !dg.room[p] && !dg.HasTooManyWallNeighbors(p) {
			return p
		}
	}
}

func (dg *dgen) SatowalgaCell(g *game) gruid.Point {
	count := 0
	for {
		count++
		if count > 2000 {
			return g.FreeCellForMonster()
		}
		x := dg.rand.Intn(DungeonWidth)
		y := dg.rand.Intn(DungeonHeight)
		p := gruid.Point{x, y}
		mons := g.MonsterAt(p)
		if mons.Exists() {
			continue
		}
		c := dg.d.Cell(p)
		if terrain(c) == GroundCell && count < 400 {
			continue
		}
		if c.IsGround() && !dg.room[p] && !dg.HasTooManyWallNeighbors(p) {
			return p
		}
	}
}

func (dg *dgen) FoliageCell(g *game) gruid.Point {
	count := 0
	for {
		count++
		if count > 1500 {
			return dg.OutsideCell(g)
		}
		x := dg.rand.Intn(DungeonWidth)
		y := dg.rand.Intn(DungeonHeight)
		p := gruid.Point{x, y}
		if distance(p, g.Player.P) < DefaultLOSRange {
			continue
		}
		mons := g.MonsterAt(p)
		if mons.Exists() {
			continue
		}
		c := dg.d.Cell(p)
		if terrain(c) == FoliageCell {
			return p
		}
	}
}

func (dg *dgen) OutsideCell(g *game) gruid.Point {
	count := 0
	for {
		count++
		if count > 1500 {
			panic("OutsideCell")
		}
		x := dg.rand.Intn(DungeonWidth)
		y := dg.rand.Intn(DungeonHeight)
		p := gruid.Point{x, y}
		mons := g.MonsterAt(p)
		if mons.Exists() {
			continue
		}
		c := dg.d.Cell(p)
		if !dg.room[p] && (terrain(c) == FoliageCell || terrain(c) == GroundCell) {
			return p
		}
	}
}

func (dg *dgen) InsideCell(g *game) gruid.Point {
	count := 0
	for {
		count++
		if count > 1500 {
			panic("InsideCell")
		}
		x := dg.rand.Intn(DungeonWidth)
		y := dg.rand.Intn(DungeonHeight)
		p := gruid.Point{x, y}
		mons := g.MonsterAt(p)
		if mons.Exists() {
			continue
		}
		if distance(p, g.Player.P) < DefaultLOSRange {
			continue
		}
		c := dg.d.Cell(p)
		if dg.room[p] && (terrain(c) == FoliageCell || terrain(c) == GroundCell) {
			return p
		}
	}
}

func (dg *dgen) GenItem(g *game) {
	plan := g.GenPlan[g.Depth]
	if plan != GenAmulet && plan != GenCloak {
		return
	}
	p := InvalidPos
	count := 0
	for p == InvalidPos {
		count++
		if count > maxIterations {
			panic("GenItem")
		}
		p = dg.rooms[RandInt(len(dg.rooms))].RandomPlace(PlaceItem)
	}
	g.Dungeon.SetCell(p, ItemCell)
	var it item
	switch plan {
	case GenCloak:
		it = g.RandomCloak()
		g.GeneratedCloaks = append(g.GeneratedCloaks, it)
	case GenAmulet:
		it = g.RandomAmulet()
		g.GeneratedAmulets = append(g.GeneratedAmulets, it)
	}
	g.Objects.Items[p] = it
}

func (dg *dgen) GenBarrierStone(g *game) {
	p := InvalidPos
	count := 0
	for p == InvalidPos {
		count++
		if count > maxIterations {
			panic("GenBarrierStone")
		}
		p = dg.rooms[RandInt(len(dg.rooms))].RandomPlaces(PlaceSpecialOrStatic)
	}
	g.Dungeon.SetCell(p, StoneCell)
	g.Objects.Stones[p] = SealStone
}

func (dg *dgen) GenMagara(g *game) {
	p := InvalidPos
	count := 0
	for p == InvalidPos {
		count++
		if count > maxIterations {
			panic("GenMagara")
		}
		p = dg.rooms[RandInt(len(dg.rooms))].RandomPlace(PlaceItem)
	}
	g.Dungeon.SetCell(p, MagaraCell)
	mag := g.RandomMagara()
	g.Objects.Magaras[p] = mag
	g.GeneratedMagaras = append(g.GeneratedMagaras, mag.Kind)
}

func (dg *dgen) GenStairs(g *game, st stair) {
	var ri, pj int
	best := 0
	for i, r := range dg.rooms {
		for j, pl := range r.places {
			score := distance(pl.p, g.Player.P) + dg.rand.Intn(20)
			if !pl.used && pl.kind == PlaceSpecialStatic && score > best {
				ri = i
				pj = j
				best = distance(pl.p, g.Player.P)
			}
		}
	}
	r := dg.rooms[ri]
	r.places[pj].used = true
	r.places[pj].used = true
	p := r.places[pj].p
	g.Dungeon.SetCell(p, StairCell)
	g.Objects.Stairs[p] = st
}

func (dg *dgen) GenFakeStairs(g *game) {
	if !g.Params.FakeStair[g.Depth] {
		return
	}
	var ri, pj int
	best := 0
loop:
	for i, r := range dg.rooms {
		for _, pl := range r.places {
			if terrain(dg.d.Cell(pl.p)) == StairCell {
				continue loop
			}
		}
		for j, pl := range r.places {
			score := distance(pl.p, g.Player.P) + dg.rand.Intn(20)
			if !pl.used && pl.kind == PlaceSpecialStatic && score > best {
				ri = i
				pj = j
				best = distance(pl.p, g.Player.P)
			}
		}
	}
	r := dg.rooms[ri]
	r.places[pj].used = true
	r.places[pj].used = true
	p := r.places[pj].p
	g.Dungeon.SetCell(p, FakeStairCell)
	g.Objects.FakeStairs[p] = true
}

func (dg *dgen) GenBarrel(g *game) {
	p := InvalidPos
	count := 0
	for p == InvalidPos {
		count++
		if count > 500 {
			return
		}
		p = dg.rooms[RandInt(len(dg.rooms))].RandomPlace(PlaceSpecialStatic)
	}
	g.Dungeon.SetCell(p, BarrelCell)
	g.Objects.Barrels[p] = true
}

func (dg *dgen) GenTable(g *game) {
	p := InvalidPos
	count := 0
	for p == InvalidPos {
		count++
		if count > 500 {
			return
		}
		p = dg.rooms[RandInt(len(dg.rooms))].RandomPlaces(PlaceSpecialOrStatic)
	}
	g.Dungeon.SetCell(p, TableCell)
}

func (dg *dgen) GenTree(g *game) {
	p := dg.OutsideCavernMiddleCell(g)
	if p != InvalidPos {
		g.Dungeon.SetCell(p, TreeCell)
	}
}

func (dg *dgen) CaveGroundCell(g *game) gruid.Point {
	count := 0
	for {
		count++
		if count > maxIterations {
			panic("CaveGroundCell")
		}
		x := dg.rand.Intn(DungeonWidth)
		y := dg.rand.Intn(DungeonHeight)
		p := gruid.Point{x, y}
		c := dg.d.Cell(p)
		if (terrain(c) == GroundCell || terrain(c) == CavernCell || terrain(c) == QueenRockCell) && !dg.room[p] {
			return p
		}
	}
}

func (dg *dgen) RandomInStone(g *game) stone {
	if g.Params.MappingStone[g.Depth] {
		g.Params.MappingStone[g.Depth] = false
		return MappingStone
	}
	instones := []stone{
		BarrelStone,
		QueenStone,
		TreeStone,
		TeleportStone,
		SensingStone,
	}
	if dg.rand.Intn(2) == 0 {
		// fog stone less often inside
		instones = append(instones, FogStone)
	}
	return instones[RandInt(len(instones))]
}

func (dg *dgen) RandomOutStone(g *game) stone {
	instones := []stone{
		BarrelStone,
		FogStone,
		QueenStone,
		NightStone,
		TreeStone,
		TeleportStone,
	}
	if dg.rand.Intn(2) == 0 {
		// sensing stone less often outside
		instones = append(instones, SensingStone)
	}
	return instones[dg.rand.Intn(len(instones))]
}

func (dg *dgen) GenStones(g *game) {
	// Magical Stones
	nstones := 3
	switch dg.rand.Intn(8) {
	case 1, 2, 3, 4, 5:
		nstones++
	case 6, 7:
		nstones += 2
	}
	inroom := 2
	if g.Params.Stones[g.Depth] {
		nstones += 4 + dg.rand.Intn(3)
		inroom += 2
	}
	if dg.layout == RandomSmallWalkCaveUrbanised {
		inroom++
	}
	for i := 0; i < nstones; i++ {
		p := InvalidPos
		var st stone
		if i < inroom {
			count := 0
			for p == InvalidPos {
				count++
				if count > 1500 {
					p = dg.CaveGroundCell(g)
					break
				}
				p = dg.rooms[RandInt(len(dg.rooms))].RandomPlace(PlaceStatic)
			}
			st = dg.RandomInStone(g)
		} else {
			p = dg.CaveGroundCell(g)
			st = dg.RandomOutStone(g)
		}
		g.Objects.Stones[p] = st
		g.Dungeon.SetCell(p, StoneCell)
	}
}

func (dg *dgen) GenCellularAutomataCaveMap() {
	dg.RunCellularAutomataCave()
	dg.Foliage(false)
}

func (dg *dgen) RunCellularAutomataCave() bool {
	rules := []rl.CellularAutomataRule{
		{WCutoff1: 5, WCutoff2: 2, Reps: 4, WallsOutOfRange: true},
		{WCutoff1: 5, WCutoff2: 25, Reps: 3, WallsOutOfRange: true},
	}
	mg := rl.MapGen{Rand: dg.rand, Grid: dg.d.Grid}
	if dg.rand.Intn(2) == 0 {
		mg.CellularAutomataCave(rl.Cell(WallCell), rl.Cell(GroundCell), 0.47, rules)
	} else {
		mg.CellularAutomataCave(rl.Cell(WallCell), rl.Cell(GroundCell), 0.43, rules)
	}
	return true
}

func (dg *dgen) GenLake(c cell) {
	walls := []gruid.Point{}
	xshift := 10 + dg.rand.Intn(5)
	yshift := 5 + dg.rand.Intn(3)
	it := dg.d.Grid.Iterator()
	for it.Next() {
		p := it.P()
		if p.X < xshift || p.Y < yshift || p.X > DungeonWidth-xshift || p.Y > DungeonHeight-yshift {
			continue
		}
		c := cell(it.Cell())
		if terrain(c) == WallCell && !dg.room[p] {
			walls = append(walls, p)
		}
	}
	count := 0
	var bestpos = walls[RandInt(len(walls))]
	var bestsize int
	d := dg.d
	passable := func(p gruid.Point) func(q gruid.Point) bool {
		return func(q gruid.Point) bool {
			return valid(q) && terrain(dg.d.Cell(q)) == WallCell && !dg.room[q] && distance(p, q) < 10+dg.rand.Intn(10)
		}
	}
	for {
		p := walls[RandInt(len(walls))]
		sp := newPather(passable(p))
		size := len(dg.PR.CCMap(sp, p))
		count++
		if abs(bestsize-90) > abs(size-90) {
			bestsize = size
			bestpos = p
		}
		if count > 15 || abs(size-90) < 25 {
			break
		}
	}
	sp := newPather(passable(bestpos))
	for _, p := range dg.PR.CCMap(sp, bestpos) {
		d.SetCell(p, c)
	}
}

func (dg *dgen) GenQueenRock() {
	cavern := []gruid.Point{}
	for i := 0; i < DungeonNCells; i++ {
		p := idxtopos(i)
		c := dg.d.Cell(p)
		if terrain(c) == CavernCell {
			cavern = append(cavern, p)
		}
	}
	if len(cavern) == 0 {
		return
	}
	for i := 0; i < 1+dg.rand.Intn(2); i++ {
		p := cavern[RandInt(len(cavern))]
		passable := func(q gruid.Point) bool {
			return valid(q) && terrain(dg.d.Cell(q)) == CavernCell && distance(q, p) < 15+dg.rand.Intn(5)
		}
		cp := newPather(passable)
		for _, q := range dg.PR.CCMap(cp, p) {
			dg.d.SetCell(q, QueenRockCell)
		}
	}

}

func (dg *dgen) Foliage(less bool) {
	gd := rl.NewGrid(DungeonWidth, DungeonHeight)
	rules := []rl.CellularAutomataRule{
		{WCutoff1: 5, WCutoff2: 2, Reps: 4, WallsOutOfRange: true},
		{WCutoff1: 5, WCutoff2: 25, Reps: 2, WallsOutOfRange: true},
	}
	mg := rl.MapGen{Rand: dg.rand, Grid: gd}
	winit := 0.55
	if less {
		winit = 0.53
	}
	mg.CellularAutomataCave(rl.Cell(WallCell), rl.Cell(FoliageCell), winit, rules)
	it := dg.d.Grid.Iterator()
	itgd := gd.Iterator()
	for it.Next() && itgd.Next() {
		if terrain(cell(it.Cell())) == GroundCell && cell(itgd.Cell()) == FoliageCell {
			it.SetCell(rl.Cell(FoliageCell))
		}
	}
}

// walker implements rl.RandomWalker.
type walker struct {
	rand *rand.Rand
}

// Neighbor returns a random neighbor position, favoring horizontal directions
// (because the maps we use are longer in that direction).
func (w walker) Neighbor(p gruid.Point) gruid.Point {
	switch w.rand.Intn(6) {
	case 0, 1:
		return p.Shift(1, 0)
	case 2, 3:
		return p.Shift(-1, 0)
	case 4:
		return p.Shift(0, 1)
	default:
		return p.Shift(0, -1)
	}
}

func (dg *dgen) GenCaveMap(size int) {
	mg := rl.MapGen{
		Rand: dg.rand,
		Grid: dg.d.Grid,
	}
	wlk := walker{rand: dg.rand}
	mg.RandomWalkCave(wlk, rl.Cell(GroundCell), float64(size)/float64(DungeonNCells), 8)
	dg.Foliage(false)
}

func (dg *dgen) DigBlock(block []gruid.Point) []gruid.Point {
	d := dg.d
	p := dg.WallCell()
	block = block[:0]
	count := 0
	for {
		count++
		if count > 3000 && count%500 == 0 {
			p = dg.WallCell()
			block = block[:0]
		}
		if count > 10000 {
			panic("DigBlock")
		}
		block = append(block, p)
		if d.HasFreeNeighbor(&dg.neighbors, p) {
			break
		}
		p = randomNeighbor(p)
		if !valid(p) {
			block = block[:0]
			p = dg.WallCell()
			continue
		}
		if !valid(p) {
			return nil
		}
	}
	return block
}

func (dg *dgen) GenTreeCaveMap() {
	d := dg.d
	center := gruid.Point{40, 10}
	d.SetCell(center, GroundCell)
	d.SetCell(center.Add(gruid.Point{1, 0}), GroundCell)
	d.SetCell(center.Add(gruid.Point{1, -1}), GroundCell)
	d.SetCell(center.Add(gruid.Point{0, 1}), GroundCell)
	d.SetCell(center.Add(gruid.Point{1, 1}), GroundCell)
	d.SetCell(center.Add(gruid.Point{0, -1}), GroundCell)
	d.SetCell(center.Add(gruid.Point{-1, -1}), GroundCell)
	d.SetCell(center.Add(gruid.Point{-1, 0}), GroundCell)
	d.SetCell(center.Add(gruid.Point{-1, 1}), GroundCell)
	max := 21 * 21
	cells := 1
	block := make([]gruid.Point, 0, 64)
loop:
	for cells < max {
		block = dg.DigBlock(block)
		if len(block) == 0 {
			continue loop
		}
		for _, p := range block {
			if terrain(d.Cell(p)) != GroundCell {
				d.SetCell(p, GroundCell)
				cells++
			}
		}
	}
	dg.Foliage(true)
}

// monster generation

func (g *game) GenBand(band monsterBand) []monsterKind {
	mbd := MonsBands[band]
	if !mbd.Band {
		return []monsterKind{mbd.Monster}
	}
	bandMonsters := []monsterKind{}
	for m, n := range mbd.Distribution {
		for i := 0; i < n; i++ {
			bandMonsters = append(bandMonsters, m)
		}
	}
	return bandMonsters
}

func (dg *dgen) BandInfoGuard(g *game, band monsterBand, pl placeKind) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	p := InvalidPos
	count := 0
loop:
	for p == InvalidPos {
		count++
		if count > maxIterations {
			p = dg.InsideCell(g)
			break
		}
		for i := 0; i < 20; i++ {
			r := dg.rooms[RandInt(len(dg.rooms)-1)]
			for _, e := range r.places {
				if e.kind == PlaceSpecialStatic {
					p = r.RandomPlace(pl)
					break
				}
			}
			if p != InvalidPos && !g.MonsterAt(p).Exists() {
				break loop
			}
		}
		r := dg.rooms[RandInt(len(dg.rooms)-1)]
		p = r.RandomPlace(pl)
	}
	bandinfo.Path = append(bandinfo.Path, p)
	bandinfo.Beh = BehGuard
	return bandinfo
}

func (dg *dgen) BandInfoGuardSpecial(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	p := InvalidPos
	count := 0
	for _, r := range dg.rooms {
		count++
		if count > 1 {
			log.Print("unavailable special guard position")
			p = dg.InsideCell(g)
			break
		}
		p = r.RandomPlace(PlacePatrolSpecial)
		if p != InvalidPos && !g.MonsterAt(p).Exists() {
			break
		}
	}
	bandinfo.Path = append(bandinfo.Path, p)
	bandinfo.Beh = BehGuard
	return bandinfo
}

func (dg *dgen) BandInfoPatrol(g *game, band monsterBand, pl placeKind) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	p := InvalidPos
	count := 0
	for p == InvalidPos || g.MonsterAt(p).Exists() {
		count++
		if count > 4000 {
			p = dg.InsideCell(g)
			break
		}
		p = dg.rooms[RandInt(len(dg.rooms)-1)].RandomPlace(pl)
	}
	target := InvalidPos
	count = 0
	for target == InvalidPos {
		// TODO: only find place in other room?
		count++
		if count > 4000 {
			target = dg.InsideCell(g)
			break
		}
		target = dg.rooms[RandInt(len(dg.rooms)-1)].RandomPlace(pl)
	}
	bandinfo.Path = append(bandinfo.Path, p)
	bandinfo.Path = append(bandinfo.Path, target)
	bandinfo.Beh = BehPatrol
	return bandinfo
}

func (dg *dgen) BandInfoPatrolSpecial(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	p := InvalidPos
	count := 0
	for _, r := range dg.rooms {
		count++
		if count > 1 {
			log.Print("unavailable special patrol position")
			p = dg.InsideCell(g)
			break
		}
		p = r.RandomPlace(PlacePatrolSpecial)
		if p != InvalidPos && !g.MonsterAt(p).Exists() {
			break
		}
	}
	target := InvalidPos
	count = 0
	for _, r := range dg.rooms {
		count++
		if count > 1 {
			log.Print("unavailable special second patrol position")
			p = dg.InsideCell(g)
			break
		}
		target = r.RandomPlace(PlacePatrolSpecial)
		if target != InvalidPos {
			break
		}
	}
	bandinfo.Path = append(bandinfo.Path, p)
	bandinfo.Path = append(bandinfo.Path, target)
	bandinfo.Beh = BehPatrol
	return bandinfo
}

func (dg *dgen) BandInfoOutsideGround(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	bandinfo.Path = append(bandinfo.Path, dg.OutsideGroundCell(g))
	bandinfo.Beh = BehWander
	return bandinfo
}

func (dg *dgen) BandInfoSatowalga(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	bandinfo.Path = append(bandinfo.Path, dg.SatowalgaCell(g))
	bandinfo.Beh = BehWander
	return bandinfo
}

func (dg *dgen) BandInfoOutside(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	bandinfo.Path = append(bandinfo.Path, dg.OutsideCell(g))
	bandinfo.Beh = BehWander
	return bandinfo
}

func (dg *dgen) BandInfoOutsideExplore(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	bandinfo.Path = append(bandinfo.Path, dg.OutsideCell(g))
	for i := 0; i < 4; i++ {
		for j := 0; j < 100; j++ {
			p := dg.OutsideCell(g)
			if dg.PR.CCMapAt(p) == dg.PR.CCMapAt(bandinfo.Path[0]) {
				bandinfo.Path = append(bandinfo.Path, p)
				break
			}
		}
	}
	bandinfo.Beh = BehExplore
	return bandinfo
}

func (dg *dgen) BandInfoOutsideExploreButterfly(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	bandinfo.Path = append(bandinfo.Path, dg.OutsideCell(g))
	for i := 0; i < 9; i++ {
		for j := 0; j < 100; j++ {
			p := dg.OutsideCell(g)
			if dg.PR.CCMapAt(p) == dg.PR.CCMapAt(bandinfo.Path[0]) {
				bandinfo.Path = append(bandinfo.Path, p)
				break
			}
		}
	}
	bandinfo.Beh = BehExplore
	return bandinfo
}

func (dg *dgen) BandInfoFoliage(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	bandinfo.Path = append(bandinfo.Path, dg.FoliageCell(g))
	bandinfo.Beh = BehWander
	return bandinfo
}

func (g *game) FreeCellForMonster() gruid.Point {
	d := g.Dungeon
	count := 0
	for {
		count++
		if count > maxIterations {
			panic("FreeCellForMonster")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		p := gruid.Point{x, y}
		c := d.Cell(p)
		if !c.IsPassable() {
			continue
		}
		if g.Player != nil && distance(g.Player.P, p) < 8 {
			continue
		}
		mons := g.MonsterAt(p)
		if mons.Exists() {
			continue
		}
		return p
	}
}

func (g *game) FreeCellForBandMonster(p gruid.Point) gruid.Point {
	count := 0
	for {
		count++
		if count > maxIterations {
			return g.FreeCellForMonster()
		}
		neighbors := g.playerPassableNeighbors(p)
		if len(neighbors) == 0 {
			// should not happen
			log.Printf("no neighbors for %v", p)
			count = maxIterations + 1
			continue
		}
		r := RandInt(len(neighbors))
		p = neighbors[r]
		if g.Player != nil && distance(g.Player.P, p) < 8 {
			continue
		}
		mons := g.MonsterAt(p)
		if mons.Exists() || !g.Dungeon.Cell(p).IsPassable() {
			continue
		}
		return p
	}
}

func (dg *dgen) PutMonsterBand(g *game, band monsterBand) bool {
	monsters := g.GenBand(band)
	if monsters == nil {
		return false
	}
	awake := dg.rand.Intn(5) > 0
	var bdinf bandInfo
	switch band {
	case LoneYack, LoneWorm, PairYack:
		bdinf = dg.BandInfoFoliage(g, band)
	case LoneDog, LoneHarpy:
		bdinf = dg.BandInfoOutsideGround(g, band)
	case LoneBlinkingFrog, LoneExplosiveNadre, PairExplosiveNadre:
		bdinf = dg.BandInfoOutside(g, band)
	case LoneMirrorSpecter, LoneWingedMilfid, LoneVampire, PairWingedMilfid, LoneEarthDragon, LoneHazeCat, LoneSpider:
		bdinf = dg.BandInfoOutsideExplore(g, band)
	case LoneButterfly:
		bdinf = dg.BandInfoOutsideExploreButterfly(g, band)
	case LoneTreeMushroom, LoneAcidMound:
		bdinf = dg.BandInfoOutside(g, band)
	case LoneHighGuard:
		bdinf = dg.BandInfoGuard(g, band, PlacePatrol)
	case LoneSatowalgaPlant:
		bdinf = dg.BandInfoSatowalga(g, band)
	case SpecialLoneVampire, SpecialLoneNixe, SpecialLoneMilfid, SpecialLoneOricCelmist, SpecialLoneHarmonicCelmist, SpecialLoneHighGuard,
		SpecialLoneHarpy, SpecialLoneTreeMushroom, SpecialLoneMirrorSpecter, SpecialLoneHazeCat, SpecialLoneSpider, SpecialLoneAcidMound, SpecialLoneDog, SpecialLoneExplosiveNadre, SpecialLoneYack, SpecialLoneBlinkingFrog:
		if dg.rand.Intn(5) > 0 {
			bdinf = dg.BandInfoPatrolSpecial(g, band)
		} else {
			bdinf = dg.BandInfoGuardSpecial(g, band)
		}
		if !awake && dg.rand.Intn(2) > 0 {
			awake = true
		}
	case UniqueCrazyImp:
		bdinf = dg.BandInfoOutside(g, band)
		bdinf.Beh = BehCrazyImp
	default:
		bdinf = dg.BandInfoPatrol(g, band, PlacePatrol)
	}
	g.Bands = append(g.Bands, bdinf)
	var p gruid.Point
	if len(bdinf.Path) == 0 {
		// should not happen now
		p = g.FreeCellForMonster()
	} else {
		p = bdinf.Path[0]
		if g.MonsterAt(p).Exists() {
			log.Printf("already monster at %v mons %v no place for %v", p, g.MonsterAt(p).Kind, monsters)
			p = g.FreeCellForMonster()
		}
	}
	for i, mk := range monsters {
		mons := &monster{Kind: mk}
		if awake {
			mons.State = Wandering
		}
		g.Monsters = append(g.Monsters, mons)
		mons.Init()
		mons.Index = len(g.Monsters) - 1
		mons.Band = len(g.Bands) - 1
		mons.PlaceAtStart(g, p)
		mons.Target = mons.NextTarget(g)
		if i < len(monsters)-1 {
			p = g.FreeCellForBandMonster(p)
		}
	}
	return true
}

func (dg *dgen) PutRandomBand(g *game, bands []monsterBand) bool {
	return dg.PutMonsterBand(g, bands[RandInt(len(bands))])
}

func (dg *dgen) PutRandomBandN(g *game, bands []monsterBand, n int) {
	for i := 0; i < n; i++ {
		dg.PutMonsterBand(g, bands[RandInt(len(bands))])
	}
}

func (dg *dgen) GenMonsters(g *game) {
	g.Monsters = []*monster{}
	g.Bands = []bandInfo{}
	// common bands
	bandsGuard := []monsterBand{LoneGuard}
	bandsButterfly := []monsterBand{LoneButterfly}
	bandsHighGuard := []monsterBand{LoneHighGuard}
	bandsAnimals := []monsterBand{LoneYack, LoneWorm, LoneDog, LoneBlinkingFrog, LoneExplosiveNadre, LoneHarpy, LoneAcidMound}
	bandsPlants := []monsterBand{LoneSatowalgaPlant}
	bandsBipeds := []monsterBand{LoneOricCelmist, LoneMirrorSpecter, LoneWingedMilfid, LoneMadNixe, LoneVampire, LoneHarmonicCelmist}
	bandsRare := []monsterBand{LoneTreeMushroom, LoneEarthDragon, LoneHazeCat, LoneSpider}
	// monster specific bands
	bandNadre := []monsterBand{LoneExplosiveNadre}
	bandFrog := []monsterBand{LoneBlinkingFrog}
	bandDog := []monsterBand{LoneDog}
	bandYack := []monsterBand{LoneYack}
	bandVampire := []monsterBand{LoneVampire}
	bandOricCelmist := []monsterBand{LoneOricCelmist}
	bandHarmonicCelmist := []monsterBand{LoneHarmonicCelmist}
	bandMadNixe := []monsterBand{LoneMadNixe}
	bandMirrorSpecter := []monsterBand{LoneMirrorSpecter}
	bandTreeMushroom := []monsterBand{LoneTreeMushroom}
	bandHazeCat := []monsterBand{LoneHazeCat}
	bandSpider := []monsterBand{LoneSpider}
	bandDragon := []monsterBand{LoneEarthDragon}
	bandGuardPair := []monsterBand{PairGuard}
	bandYackPair := []monsterBand{PairYack}
	bandExplosiveNadrePair := []monsterBand{PairExplosiveNadre}
	bandWingedMilfidPair := []monsterBand{PairWingedMilfid}
	bandNixePair := []monsterBand{PairNixe}
	bandVampirePair := []monsterBand{PairVampire}
	bandOricCelmistPair := []monsterBand{PairOricCelmist}
	bandHarmonicCelmistPair := []monsterBand{PairHarmonicCelmist}
	// special bands
	if g.Params.Special[g.Depth] != noSpecialRoom {
		switch dg.special {
		case roomVampires:
			dg.PutRandomBandN(g, []monsterBand{SpecialLoneVampire}, 2)
		case roomNixes:
			dg.PutRandomBandN(g, []monsterBand{SpecialLoneNixe}, 2)
		case roomFrogs:
			dg.PutRandomBandN(g, []monsterBand{SpecialLoneBlinkingFrog}, 2)
		case roomMilfids:
			switch dg.rand.Intn(6) {
			case 0:
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneYack}, 2)
			case 1:
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneDog}, 2)
			default:
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneMilfid}, 2)
			}
		case roomCelmists:
			switch dg.rand.Intn(3) {
			case 0:
				bandOricCelmists := []monsterBand{SpecialLoneOricCelmist}
				dg.PutRandomBandN(g, bandOricCelmists, 2)
			case 1:
				bandHarmonicCelmists := []monsterBand{SpecialLoneHarmonicCelmist}
				dg.PutRandomBandN(g, bandHarmonicCelmists, 2)
			case 2:
				bandOricCelmists := []monsterBand{SpecialLoneOricCelmist}
				bandHarmonicCelmists := []monsterBand{SpecialLoneHarmonicCelmist}
				dg.PutRandomBandN(g, bandHarmonicCelmists, 1)
				dg.PutRandomBandN(g, bandOricCelmists, 1)
			}
		case roomHarpies:
			if dg.rand.Intn(3) > 0 {
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneHarpy}, 2)
			} else {
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneSpider}, 2)
			}
		case roomTreeMushrooms:
			if dg.rand.Intn(3) > 0 {
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneTreeMushroom}, 2)
			} else {
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneHazeCat}, 2)
			}
		case roomMirrorSpecters:
			switch dg.rand.Intn(6) {
			case 0:
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneAcidMound}, 2)
			case 1:
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneExplosiveNadre}, 2)
			default:
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneMirrorSpecter}, 2)
			}
		case roomShaedra:
			if dg.rand.Intn(3) > 0 {
				dg.PutRandomBand(g, []monsterBand{SpecialLoneHighGuard})
			} else {
				dg.PutRandomBand(g, []monsterBand{SpecialLoneOricCelmist})
			}
		case roomArtifact:
			switch dg.rand.Intn(3) {
			case 0:
				dg.PutRandomBand(g, []monsterBand{SpecialLoneHarmonicCelmist})
			case 1:
				dg.PutRandomBand(g, []monsterBand{SpecialLoneOricCelmist})
			default:
				dg.PutRandomBand(g, []monsterBand{SpecialLoneHighGuard})
			}
		default:
			// XXX not used now
			bandOricCelmists := []monsterBand{SpecialLoneOricCelmist}
			dg.PutRandomBandN(g, bandOricCelmists, 2)
		}
	}
	if g.Depth == g.Params.CrazyImp {
		dg.PutRandomBand(g, []monsterBand{UniqueCrazyImp})
	}
	dg.PutRandomBandN(g, bandsButterfly, 2)
	if dg.layout == RandomSmallWalkCaveUrbanised {
		dg.PutRandomBandN(g, bandsGuard, 1+(g.Depth+1)/4)
	}
	switch g.Depth {
	case 1:
		// 8-9
		if dg.rand.Intn(2) == 0 {
			if dg.rand.Intn(2) == 0 {
				dg.PutRandomBandN(g, bandsGuard, 5)
			} else {
				dg.PutRandomBandN(g, bandsGuard, 4)
				dg.PutRandomBandN(g, bandsBipeds, 1)
			}
			dg.PutRandomBandN(g, bandsAnimals, 3)
		} else {
			dg.PutRandomBandN(g, bandsGuard, 4)
			if dg.rand.Intn(5) > 0 {
				dg.PutRandomBandN(g, bandsAnimals, 5)
			} else {
				dg.PutRandomBandN(g, bandsAnimals, 3)
				dg.PutRandomBandN(g, bandsRare, 1)
			}
		}
	case 2:
		// 10-11
		dg.PutRandomBandN(g, bandsGuard, 3)
		switch dg.rand.Intn(5) {
		case 0, 1:
			// 7
			dg.PutRandomBandN(g, bandsBipeds, 1)
			dg.PutRandomBandN(g, bandsAnimals, 4)
			dg.PutRandomBandN(g, bandsPlants, 2)
		case 2, 3:
			// 8
			dg.PutRandomBandN(g, bandsAnimals, 3)
			dg.PutRandomBandN(g, bandsRare, 1)
			dg.PutRandomBandN(g, bandsButterfly, 2)
			dg.PutRandomBandN(g, bandsPlants, 2)
		case 4:
			// 8
			dg.PutRandomBandN(g, bandsPlants, 3)
			if dg.rand.Intn(2) == 0 {
				dg.PutRandomBandN(g, bandFrog, 5)
			} else {
				dg.PutRandomBandN(g, bandYack, 5)
			}
		}
	case 3:
		// 11-12
		dg.PutRandomBandN(g, bandsHighGuard, 2)
		dg.PutRandomBandN(g, bandsGuard, 4)
		switch dg.rand.Intn(5) {
		case 0, 1:
			// 5
			if dg.rand.Intn(3) == 0 {
				dg.PutRandomBandN(g, bandDog, 3)
			} else {
				dg.PutRandomBandN(g, bandsAnimals, 3)
			}
			dg.PutRandomBandN(g, bandsPlants, 2)
		case 2, 3:
			// 5
			dg.PutRandomBandN(g, bandsAnimals, 2)
			dg.PutRandomBandN(g, bandsPlants, 1)
			dg.PutRandomBandN(g, bandsBipeds, 2)
		case 4:
			// 6
			dg.PutRandomBandN(g, bandsPlants, 1)
			dg.PutRandomBandN(g, bandNadre, 5)
		}
	case 4:
		// 12-13
		dg.PutRandomBandN(g, bandsHighGuard, 2)
		switch dg.rand.Intn(5) {
		case 0, 1:
			// 10
			dg.PutRandomBandN(g, bandsGuard, 4)
			dg.PutRandomBandN(g, bandsRare, 2)
			dg.PutRandomBandN(g, bandGuardPair, 1)
			dg.PutRandomBandN(g, bandsBipeds, 1)
			dg.PutRandomBandN(g, bandsPlants, 1)
		case 2, 3:
			// 11
			dg.PutRandomBandN(g, bandsGuard, 7)
			dg.PutRandomBandN(g, bandGuardPair, 1)
			dg.PutRandomBandN(g, bandsAnimals, 1)
			dg.PutRandomBandN(g, bandsPlants, 1)
		case 4:
			dg.PutRandomBandN(g, bandsGuard, 4)
			if dg.rand.Intn(2) == 0 {
				dg.PutRandomBandN(g, bandOricCelmist, 4)
			} else {
				dg.PutRandomBandN(g, bandHarmonicCelmist, 4)
			}
			dg.PutRandomBandN(g, bandsPlants, 1)
		}
	case 5:
		// 13-14
		dg.PutRandomBandN(g, bandsHighGuard, 2)
		if dg.rand.Intn(2) == 0 {
			// 11
			if dg.rand.Intn(2) == 0 {
				dg.PutRandomBandN(g, bandsGuard, 2)
				dg.PutRandomBandN(g, bandGuardPair, 1)
			} else {
				dg.PutRandomBandN(g, bandsGuard, 4)
			}
			dg.PutRandomBandN(g, bandsAnimals, 1)
			dg.PutRandomBandN(g, bandsRare, 2)
			dg.PutRandomBandN(g, bandsBipeds, 3)
			dg.PutRandomBandN(g, bandsPlants, 1)
		} else {
			// 12
			dg.PutRandomBandN(g, bandsGuard, 2)
			dg.PutRandomBandN(g, bandsAnimals, 3)
			dg.PutRandomBandN(g, bandsBipeds, 2)
			if dg.rand.Intn(2) == 0 {
				dg.PutRandomBandN(g, bandOricCelmistPair, 1)
			} else {
				dg.PutRandomBandN(g, bandHarmonicCelmistPair, 1)
			}
			dg.PutRandomBandN(g, bandsRare, 1)
			dg.PutRandomBandN(g, bandsPlants, 2)
		}
	case 6:
		// 15-17
		dg.PutRandomBandN(g, bandsHighGuard, 1)
		if dg.rand.Intn(2) == 0 {
			// 14
			dg.PutRandomBandN(g, bandsGuard, 3)
			dg.PutRandomBandN(g, bandsAnimals, 2)
			dg.PutRandomBandN(g, bandsRare, 3)
			dg.PutRandomBandN(g, bandsBipeds, 1)
			if dg.rand.Intn(2) == 0 {
				dg.PutRandomBandN(g, bandYackPair, 1)
			} else {
				dg.PutRandomBandN(g, bandWingedMilfidPair, 1)
			}
			dg.PutRandomBandN(g, bandsPlants, 3)
		} else {
			// 16
			dg.PutRandomBandN(g, bandsGuard, 2)
			if dg.rand.Intn(2) == 0 {
				if dg.rand.Intn(2) == 0 {
					dg.PutRandomBandN(g, bandYack, 8)
				} else {
					dg.PutRandomBandN(g, bandFrog, 8)
				}
			} else {
				dg.PutRandomBandN(g, bandsRare, 2)
				if dg.rand.Intn(3) == 0 {
					dg.PutRandomBandN(g, bandsAnimals, 4)
					dg.PutRandomBandN(g, []monsterBand{PairWorm}, 1)
				} else {
					dg.PutRandomBandN(g, bandsAnimals, 6)
				}
			}
			dg.PutRandomBandN(g, bandsButterfly, 1)
			dg.PutRandomBandN(g, bandsPlants, 5)
		}
	case 7:
		// 19
		dg.PutRandomBandN(g, bandsHighGuard, 1)
		if dg.rand.Intn(2) == 0 {
			// 18
			dg.PutRandomBandN(g, bandsGuard, 4)
			if dg.rand.Intn(3) == 0 {
				dg.PutRandomBandN(g, bandDog, 4)
				dg.PutRandomBandN(g, bandsAnimals, 2)
			} else {
				dg.PutRandomBandN(g, bandsAnimals, 6)
			}
			dg.PutRandomBandN(g, bandsButterfly, 1)
			dg.PutRandomBandN(g, bandsRare, 2)
			dg.PutRandomBandN(g, bandsBipeds, 3)
			dg.PutRandomBandN(g, bandsPlants, 2)
		} else {
			// 18
			dg.PutRandomBandN(g, bandsGuard, 1)
			dg.PutRandomBandN(g, bandsRare, 4)
			dg.PutRandomBandN(g, bandsButterfly, 1)
			if dg.rand.Intn(3) == 0 {
				dg.PutRandomBandN(g, bandNadre, 7)
			} else {
				dg.PutRandomBandN(g, bandsAnimals, 5)
				if dg.rand.Intn(2) == 0 {
					dg.PutRandomBandN(g, []monsterBand{PairFrog}, 1)
				} else {
					dg.PutRandomBandN(g, []monsterBand{PairDog}, 1)
				}
			}
			dg.PutRandomBandN(g, bandsPlants, 5)
		}
	case 8:
		// 18-19
		dg.PutRandomBandN(g, bandsHighGuard, 4)
		if dg.rand.Intn(2) == 0 {
			// 14
			dg.PutRandomBandN(g, bandsGuard, 5)
			dg.PutRandomBandN(g, bandsRare, 1)
			if dg.rand.Intn(3) == 0 {
				if dg.rand.Intn(2) == 0 {
					dg.PutRandomBandN(g, bandOricCelmist, 6)
				} else {
					dg.PutRandomBandN(g, bandMadNixe, 6)
				}
				dg.PutRandomBandN(g, bandsBipeds, 2)
			} else {
				dg.PutRandomBandN(g, bandsBipeds, 8)
			}
		} else {
			// 15
			dg.PutRandomBandN(g, bandsGuard, 5)
			dg.PutRandomBandN(g, bandsAnimals, 2)
			dg.PutRandomBandN(g, bandsRare, 1)
			dg.PutRandomBandN(g, bandHarmonicCelmistPair, 1)
			dg.PutRandomBandN(g, bandsBipeds, 4)
			dg.PutRandomBandN(g, bandsPlants, 1)
		}
	case 9:
		// 20-24
		dg.PutRandomBandN(g, bandsHighGuard, 2)
		if dg.rand.Intn(2) == 0 {
			// 18
			dg.PutRandomBandN(g, bandsGuard, 3)
			if dg.rand.Intn(2) == 0 {
				switch dg.rand.Intn(4) {
				case 0:
					dg.PutRandomBandN(g, bandTreeMushroom, 4)
					dg.PutRandomBandN(g, []monsterBand{PairTreeMushroom}, 1)
				case 1:
					dg.PutRandomBandN(g, bandDragon, 6)
					dg.PutRandomBandN(g, bandsAnimals, 2)
				case 2:
					dg.PutRandomBandN(g, bandHazeCat, 4)
					dg.PutRandomBandN(g, []monsterBand{PairHazeCat}, 1)
				case 3:
					dg.PutRandomBandN(g, bandSpider, 4)
					dg.PutRandomBandN(g, []monsterBand{PairSpider}, 1)
				}
			} else {
				dg.PutRandomBandN(g, bandsRare, 4)
				dg.PutRandomBandN(g, []monsterBand{PairTreeMushroom}, 1)
			}
			dg.PutRandomBandN(g, bandsRare, 4)
			dg.PutRandomBandN(g, bandsAnimals, 1)
			dg.PutRandomBandN(g, bandsBipeds, 2)
			dg.PutRandomBandN(g, bandsPlants, 2)
		} else {
			// 22+2
			dg.PutRandomBandN(g, bandsButterfly, 2)
			dg.PutRandomBandN(g, bandsGuard, 2)
			dg.PutRandomBandN(g, bandsAnimals, 8)
			if dg.rand.Intn(2) == 0 {
				dg.PutRandomBandN(g, bandExplosiveNadrePair, 2)
			} else {
				dg.PutRandomBandN(g, bandYackPair, 2)
			}
			dg.PutRandomBandN(g, bandsRare, 3)
			dg.PutRandomBandN(g, bandsPlants, 3)
			dg.PutRandomBandN(g, []monsterBand{PairSatowalga}, 1)
		}
	case 10:
		// 22
		dg.PutRandomBandN(g, bandsHighGuard, 3)
		if dg.rand.Intn(2) == 0 {
			// 19
			dg.PutRandomBandN(g, bandsGuard, 7)
			dg.PutRandomBandN(g, bandGuardPair, 1)
			dg.PutRandomBandN(g, bandsRare, 2)
			if dg.rand.Intn(2) == 0 {
				dg.PutRandomBandN(g, bandsBipeds, 8)
			} else {
				dg.PutRandomBandN(g, bandsBipeds, 4)
				dg.PutRandomBandN(g, bandMirrorSpecter, 4)
			}
		} else {
			// 19
			dg.PutRandomBandN(g, bandGuardPair, 1)
			if dg.rand.Intn(3) == 0 {
				dg.PutRandomBandN(g, bandsGuard, 4)
				dg.PutRandomBandN(g, bandVampire, 4)
				dg.PutRandomBandN(g, []monsterBand{PairVampire}, 1)
				dg.PutRandomBandN(g, bandsRare, 2)
			} else {
				dg.PutRandomBandN(g, bandsGuard, 6)
				dg.PutRandomBandN(g, bandsBipeds, 3)
				if dg.rand.Intn(2) == 0 {
					dg.PutRandomBandN(g, []monsterBand{PairNixe}, 1)
				} else {
					dg.PutRandomBandN(g, []monsterBand{PairOricCelmist}, 1)
				}
				dg.PutRandomBandN(g, bandsAnimals, 2)
			}
			dg.PutRandomBandN(g, bandsAnimals, 2)
			dg.PutRandomBandN(g, bandsPlants, 1)
			dg.PutRandomBandN(g, bandsRare, 1)
		}
	case 11:
		// 26
		dg.PutRandomBandN(g, bandsHighGuard, 5)
		if dg.rand.Intn(2) == 0 {
			// 21
			dg.PutRandomBandN(g, bandsGuard, 5)
			dg.PutRandomBandN(g, bandsRare, 2)
			if dg.rand.Intn(3) == 0 {
				if dg.rand.Intn(2) == 0 {
					dg.PutRandomBandN(g, bandOricCelmist, 5)
				} else {
					dg.PutRandomBandN(g, bandHarmonicCelmist, 5)
				}
				dg.PutRandomBandN(g, bandsBipeds, 5)
			} else {
				dg.PutRandomBandN(g, bandsBipeds, 10)
			}
			if dg.rand.Intn(2) == 0 {
				dg.PutRandomBandN(g, bandVampirePair, 1)
			} else {
				if dg.rand.Intn(2) == 0 {
					dg.PutRandomBandN(g, bandOricCelmistPair, 1)
				} else {
					dg.PutRandomBandN(g, bandHarmonicCelmistPair, 1)
				}
			}
			dg.PutRandomBandN(g, bandsAnimals, 2)
		} else {
			// 21
			dg.PutRandomBandN(g, bandsGuard, 5)
			dg.PutRandomBandN(g, bandOricCelmistPair, 1)
			dg.PutRandomBandN(g, []monsterBand{PairGuard}, 1)
			dg.PutRandomBandN(g, bandsRare, 1)
			dg.PutRandomBandN(g, bandsBipeds, 7)
			if dg.rand.Intn(2) == 0 {
				dg.PutRandomBandN(g, bandHarmonicCelmistPair, 1)
			} else {
				dg.PutRandomBandN(g, bandNixePair, 1)
			}
			dg.PutRandomBandN(g, bandsAnimals, 1)
			dg.PutRandomBandN(g, bandsPlants, 1)
		}
	}
}
