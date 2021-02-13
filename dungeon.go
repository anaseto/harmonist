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

func (d *dungeon) Cell(pos gruid.Point) cell {
	return cell(d.Grid.At(pos))
}

func (d *dungeon) Border(pos gruid.Point) bool {
	return pos.X == DungeonWidth-1 || pos.Y == DungeonHeight-1 || pos.X == 0 || pos.Y == 0
}

func (d *dungeon) SetCell(pos gruid.Point, c cell) {
	oc := d.Cell(pos)
	d.Grid.Set(pos, rl.Cell(c|oc&Explored))
}

func (d *dungeon) SetExplored(pos gruid.Point) {
	oc := d.Cell(pos)
	d.Grid.Set(pos, rl.Cell(oc|Explored))
}

func (d *dungeon) FreePassableCell() gruid.Point {
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeCell")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := gruid.Point{x, y}
		c := d.Cell(pos)
		if c.IsPassable() {
			return pos
		}
	}
}

func (d *dungeon) WallCell() gruid.Point {
	count := 0
	for {
		count++
		if count > 1000 {
			panic("WallCell")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := gruid.Point{x, y}
		c := d.Cell(pos)
		if terrain(c) == WallCell {
			return pos
		}
	}
}

func (d *dungeon) HasFreeNeighbor(pos gruid.Point) bool {
	neighbors := ValidCardinalNeighbors(pos)
	for _, pos := range neighbors {
		if d.Cell(pos).IsPassable() {
			return true
		}
	}
	return false
}

func (d *dungeon) HasTooManyWallNeighbors(pos gruid.Point) bool {
	neighbors := ValidNeighbors(pos)
	count := 8 - len(neighbors)
	for _, pos := range neighbors {
		if !d.Cell(pos).IsPassable() {
			count++
		}
	}
	return count > 1
}

func (g *game) HasNonWallExploredNeighbor(pos gruid.Point) bool {
	d := g.Dungeon
	neighbors := ValidCardinalNeighbors(pos)
	for _, pos := range neighbors {
		c := d.Cell(pos)
		if t, ok := g.TerrainKnowledge[pos]; ok {
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
	//return rs[i].pos.Y < rs[j].pos.Y || rs[i].pos.Y == rs[j].pos.Y && rs[i].pos.X < rs[j].pos.X
	center := gruid.Point{DungeonWidth / 2, DungeonHeight / 2}
	ipos := rs[i].pos
	ipos.X += rs[i].w / 2
	ipos.Y += rs[i].h / 2
	jpos := rs[j].pos
	jpos.X += rs[j].w / 2
	jpos.Y += rs[j].h / 2
	return rs[i].special || !rs[j].special && Distance(ipos, center) <= Distance(jpos, center)
}

type dgen struct {
	d       *dungeon
	tunnel  map[gruid.Point]bool
	room    map[gruid.Point]bool
	rooms   []*room
	spl     places
	special specialRoom
	layout  maplayout
	cc      []int
	PR      *paths.PathRange
	rand    *rand.Rand
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
	e1pos = r1.entries[e1i].pos
	e2i = r2.UnusedEntry()
	e2pos = r2.entries[e2i].pos
	tp := &tunnelPath{dg: dg}
	path := dg.PR.AstarPath(tp, e1pos, e2pos)
	if len(path) == 0 {
		log.Println(fmt.Sprintf("no path from %v to %v", e1pos, e2pos))
		return false
	}
	for _, pos := range path {
		if !valid(pos) {
			panic(fmt.Sprintf("gruid.Point %v from %v to %v", pos, e1pos, e2pos))
		}
		t := terrain(dg.d.Cell(pos))
		if t == WallCell || t == ChasmCell || t == GroundCell || t == FoliageCell {
			dg.d.SetCell(pos, GroundCell)
			dg.tunnel[pos] = true
		}
	}
	r1.entries[e1i].used = true
	r2.entries[e2i].used = true
	r1.tunnels++
	r2.tunnels++
	return true
}

func (dg *dgen) NewRoom(rpos gruid.Point, kind string) *room {
	r := &room{pos: rpos, vault: &rl.Vault{}}
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
	if RandInt(drev) == 0 {
		switch RandInt(2) {
		case 0:
			r.vault.Reflect()
			r.vault.Rotate(1 + 2*RandInt(2))
		default:
			r.vault.Rotate(1 + 2*RandInt(2))
		}
	} else {
		switch RandInt(2) {
		case 0:
			r.vault.Reflect()
			r.vault.Rotate(2 * RandInt(2))
		default:
			r.vault.Rotate(2 * RandInt(2))
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
			n := RandInt(5)
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
			//if e.used && g.DoorCandidate(e.pos) {
			if e.used && !e.virtual {
				r.places = append(r.places, place{pos: e.pos, kind: PlaceDoor})
			}
		}
		for _, pl := range r.places {
			if pl.kind != PlaceDoor {
				continue
			}
			dg.d.SetCell(pl.pos, DoorCell)
			r.places = append(r.places, place{pos: pl.pos, kind: PlaceDoor})
		}
	}
}

func (g *game) DoorCandidate(pos gruid.Point) bool {
	d := g.Dungeon
	if !valid(pos) || d.Cell(pos).IsPassable() {
		return false
	}
	return valid(pos.Add(gruid.Point{-1, 0})) && valid(pos.Add(gruid.Point{1, 0})) &&
		d.Cell(pos.Add(gruid.Point{-1, 0})).IsGround() && d.Cell(pos.Add(gruid.Point{1, 0})).IsGround() &&
		(!valid(pos.Add(gruid.Point{0, -1})) || terrain(d.Cell(pos.Add(gruid.Point{0, -1}))) == WallCell) &&
		(!valid(pos.Add(gruid.Point{0, 1})) || terrain(d.Cell(pos.Add(gruid.Point{0, 1}))) == WallCell) &&
		((valid(pos.Add(gruid.Point{-1, -1})) && d.Cell(pos.Add(gruid.Point{-1, -1})).IsPassable()) ||
			(valid(pos.Add(gruid.Point{-1, 1})) && d.Cell(pos.Add(gruid.Point{-1, 1})).IsPassable()) ||
			(valid(pos.Add(gruid.Point{1, -1})) && d.Cell(pos.Add(gruid.Point{1, -1})).IsPassable()) ||
			(valid(pos.Add(gruid.Point{1, 1})) && d.Cell(pos.Add(gruid.Point{1, 1})).IsPassable())) ||
		valid(pos.Add(gruid.Point{0, -1})) && valid(pos.Add(gruid.Point{0, 1})) &&
			d.Cell(pos.Add(gruid.Point{0, -1})).IsGround() && d.Cell(pos.Add(gruid.Point{0, 1})).IsGround() &&
			(!valid(pos.Add(gruid.Point{1, 0})) || terrain(d.Cell(pos.Add(gruid.Point{1, 0}))) == WallCell) &&
			(!valid(pos.Add(gruid.Point{-1, 0})) || terrain(d.Cell(pos.Add(gruid.Point{-1, 0}))) == WallCell) &&
			((valid(pos.Add(gruid.Point{-1, -1})) && d.Cell(pos.Add(gruid.Point{-1, -1})).IsPassable()) ||
				(valid(pos.Add(gruid.Point{-1, 1})) && d.Cell(pos.Add(gruid.Point{-1, 1})).IsPassable()) ||
				(valid(pos.Add(gruid.Point{1, -1})) && d.Cell(pos.Add(gruid.Point{1, -1})).IsPassable()) ||
				(valid(pos.Add(gruid.Point{1, 1})) && d.Cell(pos.Add(gruid.Point{1, 1})).IsPassable()))
}

func (dg *dgen) PutHoledWalls(g *game, n int) {
	candidates := []gruid.Point{}
	it := dg.d.Grid.Iterator()
	for it.Next() {
		pos := it.P()
		if dg.room[pos] && g.HoledWallCandidate(pos) {
			candidates = append(candidates, pos)
		}
	}
	if len(candidates) == 0 {
		return
	}
	for i := 0; i < n; i++ {
		pos := candidates[RandInt(len(candidates))]
		g.Dungeon.SetCell(pos, HoledWallCell)
	}
}

func (dg *dgen) PutWindows(g *game, n int) {
	candidates := []gruid.Point{}
	it := dg.d.Grid.Iterator()
	for it.Next() {
		pos := it.P()
		if dg.room[pos] && g.HoledWallCandidate(pos) {
			candidates = append(candidates, pos)
		}
	}
	if len(candidates) == 0 {
		return
	}
	for i := 0; i < n; i++ {
		pos := candidates[RandInt(len(candidates))]
		g.Dungeon.SetCell(pos, WindowCell)
	}
}

func (g *game) HoledWallCandidate(pos gruid.Point) bool {
	d := g.Dungeon
	if !valid(pos) || !d.Cell(pos).IsWall() {
		return false
	}
	return valid(pos.Add(gruid.Point{-1, 0})) && valid(pos.Add(gruid.Point{1, 0})) &&
		d.Cell(pos.Add(gruid.Point{-1, 0})).IsWall() && d.Cell(pos.Add(gruid.Point{1, 0})).IsWall() &&
		valid(pos.Add(gruid.Point{0, -1})) && d.Cell(pos.Add(gruid.Point{0, -1})).IsPassable() &&
		valid(pos.Add(gruid.Point{0, 1})) && d.Cell(pos.Add(gruid.Point{0, 1})).IsPassable() ||
		(valid(pos.Add(gruid.Point{-1, 0})) && valid(pos.Add(gruid.Point{1, 0})) &&
			d.Cell(pos.Add(gruid.Point{-1, 0})).IsPassable() && d.Cell(pos.Add(gruid.Point{1, 0})).IsPassable() &&
			valid(pos.Add(gruid.Point{0, -1})) && d.Cell(pos.Add(gruid.Point{0, -1})).IsWall() &&
			valid(pos.Add(gruid.Point{0, 1})) && d.Cell(pos.Add(gruid.Point{0, 1})).IsWall())
}

type placement int

const (
	PlacementRandom placement = iota
	PlacementCenter
	PlacementEdge
)

func (dg *dgen) GenRooms(templates []string, n int, pl placement) (ps []gruid.Point, ok bool) {
	ok = true
	for i := 0; i < n; i++ {
		var r *room
		count := 500
		var pos gruid.Point
		var tpl string
		for r == nil && count > 0 {
			count--
			switch pl {
			case PlacementRandom:
				pos = gruid.Point{RandInt(DungeonWidth - 1), RandInt(DungeonHeight - 1)}
			case PlacementCenter:
				pos = gruid.Point{DungeonWidth/2 - 4 + RandInt(5), DungeonHeight/2 - 3 + RandInt(4)}
			case PlacementEdge:
				if RandInt(2) == 0 {
					pos = gruid.Point{RandInt(DungeonWidth / 4), RandInt(DungeonHeight - 1)}
				} else {
					pos = gruid.Point{3*DungeonWidth/4 + RandInt(DungeonWidth/4) - 1, RandInt(DungeonHeight - 1)}
				}
			}
			tpl = templates[RandInt(len(templates))]
			r = dg.NewRoom(pos, tpl)
		}
		if r != nil {
			switch pl {
			case PlacementCenter, PlacementEdge:
				r.special = true
			}
			dg.rooms = append(dg.rooms, r)
			ps = append(ps, pos)
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
		if RandInt(3) == 0 {
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
		if RandInt(3) == 0 {
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
	if RandInt(10) > 0 {
		var c cell
		if RandInt(5) > 1 {
			c = ChasmCell
		} else {
			c = WaterCell
		}
		dg.GenLake(c)
		if RandInt(5) == 0 {
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
	for i := 0; i < 4+RandInt(2); i++ {
		dg.GenBarrel(g)
	}
	dg.AddSpecial(g, ml)
	dg.PR.CCMapAll(newPather(func(pos gruid.Point) bool {
		return valid(pos) && g.Dungeon.Cell(pos).IsPassable()
	}))
	dg.GenMonsters(g)
	dg.PutCavernCells(g)
	if RandInt(2) == 0 {
		dg.GenQueenRock()
	}
}

func (dg *dgen) PutCavernCells(g *game) {
	d := dg.d
	// TODO: improve handling and placement of this
	it := dg.d.Grid.Iterator()
	for it.Next() {
		pos := it.P()
		if terrain(cell(it.Cell())) == GroundCell && !dg.room[pos] && !dg.tunnel[pos] {
			d.SetCell(pos, CavernCell)
		}
	}
}

func (dg *dgen) ClearUnconnected(g *game) {
	d := dg.d
	sp := newPather(func(p gruid.Point) bool { return d.Cell(p).IsPlayerPassable() })
	dg.PR.CCMap(sp, g.Player.Pos)
	mg := rl.MapGen{Grid: dg.d.Grid}
	mg.KeepCC(dg.PR, g.Player.Pos, rl.Cell(WallCell))
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
		if RandInt(3) == 0 {
			ntables++
		} else if RandInt(10) == 0 {
			ntables--
		}
	case RandomWalkTreeCave:
		if RandInt(4) > 0 {
			ntables++
		}
		if RandInt(4) > 0 {
			ntables++
		}
	case RandomSmallWalkCaveUrbanised:
		ntables += 2
		if RandInt(4) > 0 {
			ntables++
		}
	}
	if g.Params.Tables[g.Depth] {
		ntables += 2 + RandInt(2)
	}
	for i := 0; i < ntables; i++ {
		dg.GenTable(g)
	}
	dg.GenLight(g)
	ntrees := 1
	switch ml {
	case AutomataCave:
		if RandInt(4) == 0 {
			ntrees++
		} else if RandInt(8) == 0 {
			ntrees--
		}
	case RandomWalkCave:
		if RandInt(4) > 0 {
			ntrees++
		}
		if RandInt(8) == 0 {
			ntrees++
		}
	case NaturalCave:
		ntrees++
		if RandInt(2) > 0 {
			ntrees++
		}
	case RandomWalkTreeCave, RandomSmallWalkCaveUrbanised:
		if RandInt(2) == 0 {
			ntrees--
		}
	}
	if g.Params.Trees[g.Depth] {
		ntrees += 2 + RandInt(2)
	}
	for i := 0; i < ntrees; i++ {
		dg.GenTree(g)
	}
	nhw := 1
	if RandInt(3) > 0 {
		nhw++
	}
	if g.Params.Holes[g.Depth] {
		nhw += 3 + RandInt(2)
	}
	switch ml {
	case RandomSmallWalkCaveUrbanised:
		if RandInt(4) > 0 {
			nhw++
		}
	}
	dg.PutHoledWalls(g, nhw)
	nwin := 1
	if nhw == 1 {
		nwin++
	}
	if g.Params.Windows[g.Depth] {
		nwin += 4 + RandInt(3)
	}
	switch ml {
	case RandomSmallWalkCaveUrbanised:
		if RandInt(4) > 0 {
			nwin++
		}
	}
	dg.PutWindows(g, nwin)
	if g.Params.Lore[g.Depth] {
		dg.PutLore(g)
	}
}

func (dg *dgen) PutLore(g *game) {
	pos := InvalidPos
	count := 0
	for pos == InvalidPos {
		count++
		if count > 2000 {
			panic("PutLore1")
		}
		pos = dg.rooms[RandInt(len(dg.rooms))].RandomPlace(PlaceItem)
	}
	count = 0
	for {
		count++
		if count > 1000 {
			panic("PutLore2")
		}
		i := RandInt(len(LoreMessages))
		if g.GeneratedLore[i] {
			continue
		}
		g.GeneratedLore[i] = true
		g.Objects.Lore[pos] = i
		g.Objects.Scrolls[pos] = ScrollLore
		g.Dungeon.SetCell(pos, ScrollCell)
		break
	}
}

func (dg *dgen) GenLight(g *game) {
	lights := []gruid.Point{}
	no := 2
	ni := 8
	switch dg.layout {
	case NaturalCave:
		no += RandInt(2)
		ni += RandInt(3)
	case AutomataCave, RandomWalkCave:
		ni += RandInt(4)
	case RandomWalkTreeCave:
		no--
		ni += RandInt(4)
	case RandomSmallWalkCaveUrbanised:
		no--
		no -= RandInt(2)
		ni += 2
		ni += RandInt(4)
	}
	for i := 0; i < no; i++ {
		pos := dg.OutsideGroundCell(g)
		g.Dungeon.SetCell(pos, LightCell)
		lights = append(lights, pos)
	}
	for i := 0; i < ni; i++ {
		pos := dg.rooms[RandInt(len(dg.rooms))].RandomPlaces(PlaceSpecialOrStatic)
		if pos != InvalidPos {
			g.Dungeon.SetCell(pos, LightCell)
			lights = append(lights, pos)
		} else if RandInt(10) > 0 {
			i--
		}
	}
	for _, pos := range lights {
		g.Objects.Lights[pos] = true
	}
	g.ComputeLights()
}

func (dg *dgen) PlayerStartCell(g *game, places []gruid.Point) {
	const far = 30
	r := dg.rooms[len(dg.rooms)-1]
loop:
	for i := len(dg.rooms) - 2; i >= 0; i-- {
		for _, pos := range places {
			if Distance(r.pos, pos) < far && Distance(dg.rooms[i].pos, pos) >= far {
				r = dg.rooms[i]
				dg.rooms[len(dg.rooms)-1], dg.rooms[i] = dg.rooms[i], dg.rooms[len(dg.rooms)-1]
				break loop
			}
		}
	}
	g.Player.Pos = r.RandomPlace(PlacePatrol)
	switch g.Depth {
	case 1, 4:
	default:
		return
	}
	itpos := InvalidPos
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Pos)
	for i := 0; i < len(neighbors); i++ {
		j := RandInt(len(neighbors) - i)
		neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
	}
loopnb:
	for _, npos := range neighbors {
		c := g.Dungeon.Cell(npos)
		if c.IsGround() {
			for _, pl := range r.places {
				if npos == pl.pos {
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
		if count > 1000 {
			panic("GenBanana")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := gruid.Point{x, y}
		c := dg.d.Cell(pos)
		if terrain(c) == GroundCell && !dg.room[pos] {
			dg.d.SetCell(pos, BananaCell)
			g.Objects.Bananas[pos] = true
			break
		}
	}
}

func (dg *dgen) GenPotion(g *game, p potion) {
	count := 0
	pos := InvalidPos
	for pos == InvalidPos {
		count++
		if count > 1000 {
			return
		}
		pos = dg.rooms[RandInt(len(dg.rooms))].RandomPlace(PlaceItem)
	}
	dg.d.SetCell(pos, PotionCell)
	g.Objects.Potions[pos] = p
}

func (dg *dgen) OutsideGroundCell(g *game) gruid.Point {
	count := 0
	for {
		count++
		if count > 1500 {
			panic("OutsideGroundCell")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := gruid.Point{x, y}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		c := dg.d.Cell(pos)
		if terrain(c) == GroundCell && !dg.room[pos] {
			return pos
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
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := gruid.Point{x, y}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		c := dg.d.Cell(pos)
		if terrain(c) == GroundCell && count < 400 || terrain(c) == FoliageCell && count < 350 {
			continue
		}
		if (c.IsGround() || terrain(c) == FoliageCell) && !dg.room[pos] && !dg.d.HasTooManyWallNeighbors(pos) {
			return pos
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
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := gruid.Point{x, y}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		c := dg.d.Cell(pos)
		if terrain(c) == GroundCell && count < 400 {
			continue
		}
		if c.IsGround() && !dg.room[pos] && !dg.d.HasTooManyWallNeighbors(pos) {
			return pos
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
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := gruid.Point{x, y}
		if Distance(pos, g.Player.Pos) < DefaultLOSRange {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		c := dg.d.Cell(pos)
		if terrain(c) == FoliageCell {
			return pos
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
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := gruid.Point{x, y}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		c := dg.d.Cell(pos)
		if !dg.room[pos] && (terrain(c) == FoliageCell || terrain(c) == GroundCell) {
			return pos
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
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := gruid.Point{x, y}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		if Distance(pos, g.Player.Pos) < DefaultLOSRange {
			continue
		}
		c := dg.d.Cell(pos)
		if dg.room[pos] && (terrain(c) == FoliageCell || terrain(c) == GroundCell) {
			return pos
		}
	}
}

func (dg *dgen) GenItem(g *game) {
	plan := g.GenPlan[g.Depth]
	if plan != GenAmulet && plan != GenCloak {
		return
	}
	pos := InvalidPos
	count := 0
	for pos == InvalidPos {
		count++
		if count > 1000 {
			panic("GenItem")
		}
		pos = dg.rooms[RandInt(len(dg.rooms))].RandomPlace(PlaceItem)
	}
	g.Dungeon.SetCell(pos, ItemCell)
	var it item
	switch plan {
	case GenCloak:
		it = g.RandomCloak()
		g.GeneratedCloaks = append(g.GeneratedCloaks, it)
	case GenAmulet:
		it = g.RandomAmulet()
		g.GeneratedAmulets = append(g.GeneratedAmulets, it)
	}
	g.Objects.Items[pos] = it
}

func (dg *dgen) GenBarrierStone(g *game) {
	pos := InvalidPos
	count := 0
	for pos == InvalidPos {
		count++
		if count > 1000 {
			panic("GenBarrierStone")
		}
		pos = dg.rooms[RandInt(len(dg.rooms))].RandomPlaces(PlaceSpecialOrStatic)
	}
	g.Dungeon.SetCell(pos, StoneCell)
	g.Objects.Stones[pos] = SealStone
}

func (dg *dgen) GenMagara(g *game) {
	pos := InvalidPos
	count := 0
	for pos == InvalidPos {
		count++
		if count > 1000 {
			panic("GenMagara")
		}
		pos = dg.rooms[RandInt(len(dg.rooms))].RandomPlace(PlaceItem)
	}
	g.Dungeon.SetCell(pos, MagaraCell)
	mag := g.RandomMagara()
	g.Objects.Magaras[pos] = mag
	g.GeneratedMagaras = append(g.GeneratedMagaras, mag.Kind)
}

func (dg *dgen) GenStairs(g *game, st stair) {
	var ri, pj int
	best := 0
	for i, r := range dg.rooms {
		for j, pl := range r.places {
			score := Distance(pl.pos, g.Player.Pos) + RandInt(20)
			if !pl.used && pl.kind == PlaceSpecialStatic && score > best {
				ri = i
				pj = j
				best = Distance(pl.pos, g.Player.Pos)
			}
		}
	}
	r := dg.rooms[ri]
	r.places[pj].used = true
	r.places[pj].used = true
	pos := r.places[pj].pos
	g.Dungeon.SetCell(pos, StairCell)
	g.Objects.Stairs[pos] = st
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
			if terrain(dg.d.Cell(pl.pos)) == StairCell {
				continue loop
			}
		}
		for j, pl := range r.places {
			score := Distance(pl.pos, g.Player.Pos) + RandInt(20)
			if !pl.used && pl.kind == PlaceSpecialStatic && score > best {
				ri = i
				pj = j
				best = Distance(pl.pos, g.Player.Pos)
			}
		}
	}
	r := dg.rooms[ri]
	r.places[pj].used = true
	r.places[pj].used = true
	pos := r.places[pj].pos
	g.Dungeon.SetCell(pos, FakeStairCell)
	g.Objects.FakeStairs[pos] = true
}

func (dg *dgen) GenBarrel(g *game) {
	pos := InvalidPos
	count := 0
	for pos == InvalidPos {
		count++
		if count > 500 {
			return
		}
		pos = dg.rooms[RandInt(len(dg.rooms))].RandomPlace(PlaceSpecialStatic)
	}
	g.Dungeon.SetCell(pos, BarrelCell)
	g.Objects.Barrels[pos] = true
}

func (dg *dgen) GenTable(g *game) {
	pos := InvalidPos
	count := 0
	for pos == InvalidPos {
		count++
		if count > 500 {
			return
		}
		pos = dg.rooms[RandInt(len(dg.rooms))].RandomPlaces(PlaceSpecialOrStatic)
	}
	g.Dungeon.SetCell(pos, TableCell)
}

func (dg *dgen) GenTree(g *game) {
	pos := dg.OutsideCavernMiddleCell(g)
	if pos != InvalidPos {
		g.Dungeon.SetCell(pos, TreeCell)
	}
}

func (dg *dgen) CaveGroundCell(g *game) gruid.Point {
	count := 0
	for {
		count++
		if count > 1000 {
			panic("CaveGroundCell")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := gruid.Point{x, y}
		c := dg.d.Cell(pos)
		if (terrain(c) == GroundCell || terrain(c) == CavernCell || terrain(c) == QueenRockCell) && !dg.room[pos] {
			return pos
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
	if RandInt(2) == 0 {
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
	if RandInt(2) == 0 {
		// sensing stone less often outside
		instones = append(instones, SensingStone)
	}
	return instones[RandInt(len(instones))]
}

func (dg *dgen) GenStones(g *game) {
	// Magical Stones
	// TODO: move into dungeon generation
	nstones := 3
	switch RandInt(8) {
	case 1, 2, 3, 4, 5:
		nstones++
	case 6, 7:
		nstones += 2
	}
	inroom := 2
	if g.Params.Stones[g.Depth] {
		nstones += 4 + RandInt(3)
		inroom += 2
	}
	if dg.layout == RandomSmallWalkCaveUrbanised {
		inroom++
	}
	for i := 0; i < nstones; i++ {
		pos := InvalidPos
		var st stone
		if i < inroom {
			count := 0
			for pos == InvalidPos {
				count++
				if count > 1500 {
					pos = dg.CaveGroundCell(g)
					break
				}
				pos = dg.rooms[RandInt(len(dg.rooms))].RandomPlace(PlaceStatic)
			}
			st = dg.RandomInStone(g)
		} else {
			pos = dg.CaveGroundCell(g)
			st = dg.RandomOutStone(g)
		}
		g.Objects.Stones[pos] = st
		g.Dungeon.SetCell(pos, StoneCell)
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
	xshift := 10 + RandInt(5)
	yshift := 5 + RandInt(3)
	it := dg.d.Grid.Iterator()
	for it.Next() {
		pos := it.P()
		if pos.X < xshift || pos.Y < yshift || pos.X > DungeonWidth-xshift || pos.Y > DungeonHeight-yshift {
			continue
		}
		c := cell(it.Cell())
		if terrain(c) == WallCell && !dg.room[pos] {
			walls = append(walls, pos)
		}
	}
	count := 0
	var bestpos = walls[RandInt(len(walls))]
	var bestsize int
	d := dg.d
	passable := func(p gruid.Point) func(q gruid.Point) bool {
		return func(q gruid.Point) bool {
			return valid(q) && terrain(dg.d.Cell(q)) == WallCell && !dg.room[q] && Distance(p, q) < 10+RandInt(10)
		}
	}
	for {
		p := walls[RandInt(len(walls))]
		sp := newPather(passable(p))
		size := len(dg.PR.CCMap(sp, p))
		count++
		if Abs(bestsize-90) > Abs(size-90) {
			bestsize = size
			bestpos = p
		}
		if count > 15 || Abs(size-90) < 25 {
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
		pos := idxtopos(i)
		c := dg.d.Cell(pos)
		if terrain(c) == CavernCell {
			cavern = append(cavern, pos)
		}
	}
	if len(cavern) == 0 {
		return
	}
	for i := 0; i < 1+RandInt(2); i++ {
		pos := cavern[RandInt(len(cavern))]
		passable := func(q gruid.Point) bool {
			return valid(q) && terrain(dg.d.Cell(q)) == CavernCell && Distance(q, pos) < 15+RandInt(5)
		}
		cp := newPather(passable)
		for _, p := range dg.PR.CCMap(cp, pos) {
			dg.d.SetCell(p, QueenRockCell)
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

func (d *dungeon) DigBlock(block []gruid.Point) []gruid.Point {
	pos := d.WallCell()
	block = block[:0]
	count := 0
	for {
		count++
		if count > 3000 && count%500 == 0 {
			pos = d.WallCell()
			block = block[:0]
		}
		if count > 10000 {
			panic("DigBlock")
		}
		block = append(block, pos)
		if d.HasFreeNeighbor(pos) {
			break
		}
		pos = RandomNeighbor(pos, false)
		if !valid(pos) {
			block = block[:0]
			pos = d.WallCell()
			continue
		}
		if !valid(pos) {
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
		block = d.DigBlock(block)
		if len(block) == 0 {
			continue loop
		}
		for _, pos := range block {
			if terrain(d.Cell(pos)) != GroundCell {
				d.SetCell(pos, GroundCell)
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
	pos := InvalidPos
	count := 0
loop:
	for pos == InvalidPos {
		count++
		if count > 1000 {
			pos = dg.InsideCell(g)
			break
		}
		for i := 0; i < 20; i++ {
			r := dg.rooms[RandInt(len(dg.rooms)-1)]
			for _, e := range r.places {
				if e.kind == PlaceSpecialStatic {
					pos = r.RandomPlace(pl)
					break
				}
			}
			if pos != InvalidPos && !g.MonsterAt(pos).Exists() {
				break loop
			}
		}
		r := dg.rooms[RandInt(len(dg.rooms)-1)]
		pos = r.RandomPlace(pl)
	}
	bandinfo.Path = append(bandinfo.Path, pos)
	bandinfo.Beh = BehGuard
	return bandinfo
}

func (dg *dgen) BandInfoGuardSpecial(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	pos := InvalidPos
	count := 0
	for _, r := range dg.rooms {
		count++
		if count > 1 {
			log.Print("unavailable special guard position")
			pos = dg.InsideCell(g)
			break
		}
		pos = r.RandomPlace(PlacePatrolSpecial)
		if pos != InvalidPos && !g.MonsterAt(pos).Exists() {
			break
		}
	}
	bandinfo.Path = append(bandinfo.Path, pos)
	bandinfo.Beh = BehGuard
	return bandinfo
}

func (dg *dgen) BandInfoPatrol(g *game, band monsterBand, pl placeKind) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	pos := InvalidPos
	count := 0
	for pos == InvalidPos || g.MonsterAt(pos).Exists() {
		count++
		if count > 4000 {
			pos = dg.InsideCell(g)
			break
		}
		pos = dg.rooms[RandInt(len(dg.rooms)-1)].RandomPlace(pl)
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
	bandinfo.Path = append(bandinfo.Path, pos)
	bandinfo.Path = append(bandinfo.Path, target)
	bandinfo.Beh = BehPatrol
	return bandinfo
}

func (dg *dgen) BandInfoPatrolSpecial(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	pos := InvalidPos
	count := 0
	for _, r := range dg.rooms {
		count++
		if count > 1 {
			log.Print("unavailable special patrol position")
			pos = dg.InsideCell(g)
			break
		}
		pos = r.RandomPlace(PlacePatrolSpecial)
		if pos != InvalidPos && !g.MonsterAt(pos).Exists() {
			break
		}
	}
	target := InvalidPos
	count = 0
	for _, r := range dg.rooms {
		count++
		if count > 1 {
			panic("patrol special")
		}
		target = r.RandomPlace(PlacePatrolSpecial)
		if target != InvalidPos {
			break
		}
	}
	bandinfo.Path = append(bandinfo.Path, pos)
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
			pos := dg.OutsideCell(g)
			if dg.PR.CCMapAt(pos) == dg.PR.CCMapAt(bandinfo.Path[0]) {
				bandinfo.Path = append(bandinfo.Path, pos)
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
			pos := dg.OutsideCell(g)
			if dg.PR.CCMapAt(pos) == dg.PR.CCMapAt(bandinfo.Path[0]) {
				bandinfo.Path = append(bandinfo.Path, pos)
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

func (dg *dgen) PutMonsterBand(g *game, band monsterBand) bool {
	monsters := g.GenBand(band)
	if monsters == nil {
		return false
	}
	awake := RandInt(5) > 0
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
		if RandInt(5) > 0 {
			bdinf = dg.BandInfoPatrolSpecial(g, band)
		} else {
			bdinf = dg.BandInfoGuardSpecial(g, band)
		}
		if !awake && RandInt(2) > 0 {
			awake = true
		}
	case UniqueCrazyImp:
		bdinf = dg.BandInfoOutside(g, band)
		bdinf.Beh = BehCrazyImp
	default:
		bdinf = dg.BandInfoPatrol(g, band, PlacePatrol)
	}
	g.Bands = append(g.Bands, bdinf)
	var pos gruid.Point
	if len(bdinf.Path) == 0 {
		// should not happen now
		pos = g.FreeCellForMonster()
	} else {
		pos = bdinf.Path[0]
		if g.MonsterAt(pos).Exists() {
			log.Printf("already monster at %v mons %v no place for %v", pos, g.MonsterAt(pos).Kind, monsters)
			pos = g.FreeCellForMonster()
		}
	}
	for _, mk := range monsters {
		mons := &monster{Kind: mk}
		if awake {
			mons.State = Wandering
		}
		g.Monsters = append(g.Monsters, mons)
		mons.Init()
		mons.Index = len(g.Monsters) - 1
		mons.Band = len(g.Bands) - 1
		mons.PlaceAtStart(g, pos)
		mons.Target = mons.NextTarget(g)
		pos = g.FreeCellForBandMonster(pos)
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
			switch RandInt(6) {
			case 0:
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneYack}, 2)
			case 1:
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneDog}, 2)
			default:
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneMilfid}, 2)
			}
		case roomCelmists:
			switch RandInt(3) {
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
			if RandInt(3) > 0 {
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneHarpy}, 2)
			} else {
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneSpider}, 2)
			}
		case roomTreeMushrooms:
			if RandInt(3) > 0 {
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneTreeMushroom}, 2)
			} else {
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneHazeCat}, 2)
			}
		case roomMirrorSpecters:
			switch RandInt(6) {
			case 0:
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneAcidMound}, 2)
			case 1:
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneExplosiveNadre}, 2)
			default:
				dg.PutRandomBandN(g, []monsterBand{SpecialLoneMirrorSpecter}, 2)
			}
		case roomShaedra:
			if RandInt(3) > 0 {
				dg.PutRandomBand(g, []monsterBand{SpecialLoneHighGuard})
			} else {
				dg.PutRandomBand(g, []monsterBand{SpecialLoneOricCelmist})
			}
		case roomArtifact:
			switch RandInt(3) {
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
		if RandInt(2) == 0 {
			if RandInt(2) == 0 {
				dg.PutRandomBandN(g, bandsGuard, 5)
			} else {
				dg.PutRandomBandN(g, bandsGuard, 4)
				dg.PutRandomBandN(g, bandsBipeds, 1)
			}
			dg.PutRandomBandN(g, bandsAnimals, 3)
		} else {
			dg.PutRandomBandN(g, bandsGuard, 4)
			if RandInt(5) > 0 {
				dg.PutRandomBandN(g, bandsAnimals, 5)
			} else {
				dg.PutRandomBandN(g, bandsAnimals, 3)
				dg.PutRandomBandN(g, bandsRare, 1)
			}
		}
	case 2:
		// 10-11
		dg.PutRandomBandN(g, bandsGuard, 3)
		switch RandInt(5) {
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
			if RandInt(2) == 0 {
				dg.PutRandomBandN(g, bandFrog, 5)
			} else {
				dg.PutRandomBandN(g, bandYack, 5)
			}
		}
	case 3:
		// 11-12
		dg.PutRandomBandN(g, bandsHighGuard, 2)
		dg.PutRandomBandN(g, bandsGuard, 4)
		switch RandInt(5) {
		case 0, 1:
			// 5
			if RandInt(3) == 0 {
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
		switch RandInt(5) {
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
			if RandInt(2) == 0 {
				dg.PutRandomBandN(g, bandOricCelmist, 4)
			} else {
				dg.PutRandomBandN(g, bandHarmonicCelmist, 4)
			}
			dg.PutRandomBandN(g, bandsPlants, 1)
		}
	case 5:
		// 13-14
		dg.PutRandomBandN(g, bandsHighGuard, 2)
		if RandInt(2) == 0 {
			// 11
			if RandInt(2) == 0 {
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
			if RandInt(2) == 0 {
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
		if RandInt(2) == 0 {
			// 14
			dg.PutRandomBandN(g, bandsGuard, 3)
			dg.PutRandomBandN(g, bandsAnimals, 2)
			dg.PutRandomBandN(g, bandsRare, 3)
			dg.PutRandomBandN(g, bandsBipeds, 1)
			if RandInt(2) == 0 {
				dg.PutRandomBandN(g, bandYackPair, 1)
			} else {
				dg.PutRandomBandN(g, bandWingedMilfidPair, 1)
			}
			dg.PutRandomBandN(g, bandsPlants, 3)
		} else {
			// 16
			dg.PutRandomBandN(g, bandsGuard, 2)
			if RandInt(2) == 0 {
				if RandInt(2) == 0 {
					dg.PutRandomBandN(g, bandYack, 8)
				} else {
					dg.PutRandomBandN(g, bandFrog, 8)
				}
			} else {
				dg.PutRandomBandN(g, bandsRare, 2)
				if RandInt(3) == 0 {
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
		if RandInt(2) == 0 {
			// 18
			dg.PutRandomBandN(g, bandsGuard, 4)
			if RandInt(3) == 0 {
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
			if RandInt(3) == 0 {
				dg.PutRandomBandN(g, bandNadre, 7)
			} else {
				dg.PutRandomBandN(g, bandsAnimals, 5)
				if RandInt(2) == 0 {
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
		if RandInt(2) == 0 {
			// 14
			dg.PutRandomBandN(g, bandsGuard, 5)
			dg.PutRandomBandN(g, bandsRare, 1)
			if RandInt(3) == 0 {
				if RandInt(2) == 0 {
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
		if RandInt(2) == 0 {
			// 18
			dg.PutRandomBandN(g, bandsGuard, 3)
			if RandInt(2) == 0 {
				switch RandInt(4) {
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
			if RandInt(2) == 0 {
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
		if RandInt(2) == 0 {
			// 19
			dg.PutRandomBandN(g, bandsGuard, 7)
			dg.PutRandomBandN(g, bandGuardPair, 1)
			dg.PutRandomBandN(g, bandsRare, 2)
			if RandInt(2) == 0 {
				dg.PutRandomBandN(g, bandsBipeds, 8)
			} else {
				dg.PutRandomBandN(g, bandsBipeds, 4)
				dg.PutRandomBandN(g, bandMirrorSpecter, 4)
			}
		} else {
			// 19
			dg.PutRandomBandN(g, bandGuardPair, 1)
			if RandInt(3) == 0 {
				dg.PutRandomBandN(g, bandsGuard, 4)
				dg.PutRandomBandN(g, bandVampire, 4)
				dg.PutRandomBandN(g, []monsterBand{PairVampire}, 1)
				dg.PutRandomBandN(g, bandsRare, 2)
			} else {
				dg.PutRandomBandN(g, bandsGuard, 6)
				dg.PutRandomBandN(g, bandsBipeds, 3)
				if RandInt(2) == 0 {
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
		if RandInt(2) == 0 {
			// 21
			dg.PutRandomBandN(g, bandsGuard, 5)
			dg.PutRandomBandN(g, bandsRare, 2)
			if RandInt(3) == 0 {
				if RandInt(2) == 0 {
					dg.PutRandomBandN(g, bandOricCelmist, 5)
				} else {
					dg.PutRandomBandN(g, bandHarmonicCelmist, 5)
				}
				dg.PutRandomBandN(g, bandsBipeds, 5)
			} else {
				dg.PutRandomBandN(g, bandsBipeds, 10)
			}
			if RandInt(2) == 0 {
				dg.PutRandomBandN(g, bandVampirePair, 1)
			} else {
				if RandInt(2) == 0 {
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
			if RandInt(2) == 0 {
				dg.PutRandomBandN(g, bandHarmonicCelmistPair, 1)
			} else {
				dg.PutRandomBandN(g, bandNixePair, 1)
			}
			dg.PutRandomBandN(g, bandsAnimals, 1)
			dg.PutRandomBandN(g, bandsPlants, 1)
		}
	}
}
