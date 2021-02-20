package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
)

var Rounds = 40

func (d *dungeon) FreePassableCell() gruid.Point {
	count := 0
	for {
		count++
		if count > maxIterations {
			panic("FreeCell")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		p := gruid.Point{x, y}
		c := d.Cell(p)
		if c.IsPassable() {
			return p
		}
	}
}

func (d *dungeon) connex(pr *paths.PathRange) bool {
	pos := d.FreePassableCell()
	passable := func(p gruid.Point) bool {
		return terrain(d.Cell(p)) != WallCell
	}
	cp := newPather(passable)
	pr.CCMap(cp, pos)
	it := d.Grid.Iterator()
	for it.Next() {
		if cell(it.Cell()).IsPassable() && pr.CCMapAt(it.P()) == -1 {
			return false
		}
	}
	return true
}

func (d *dungeon) String() string {
	b := &bytes.Buffer{}
	it := d.Grid.Iterator()
	for i := 0; it.Next(); i++ {
		c := cell(it.Cell())
		if i > 0 && i%DungeonWidth == 0 {
			fmt.Fprint(b, "\n")
		}
		if terrain(c) == WallCell {
			fmt.Fprint(b, "#")
		} else {
			fmt.Fprint(b, ".")
		}
	}
	return b.String()
}

func TestAutomataCave(t *testing.T) {
	for i := 0; i < Rounds; i++ {
		g := &game{}
		g.initrand()
		g.InitFirstLevel()
		g.InitLevelStructures()
		g.GenRoomTunnels(AutomataCave)
		if !g.Dungeon.connex(g.PR) {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}

func TestRandomWalkCave(t *testing.T) {
	for i := 0; i < Rounds; i++ {
		g := &game{}
		g.initrand()
		g.InitFirstLevel()
		g.InitLevelStructures()
		g.GenRoomTunnels(RandomWalkCave)
		if !g.Dungeon.connex(g.PR) {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}

func TestRandomWalkTreeCave(t *testing.T) {
	for i := 0; i < Rounds; i++ {
		g := &game{}
		g.initrand()
		g.InitFirstLevel()
		g.InitLevelStructures()
		g.GenRoomTunnels(RandomWalkTreeCave)
		if !g.Dungeon.connex(g.PR) {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}

func TestRandomSmallWalkCaveUrbanised(t *testing.T) {
	for i := 0; i < Rounds; i++ {
		g := &game{}
		g.initrand()
		g.InitFirstLevel()
		g.InitLevelStructures()
		g.GenRoomTunnels(RandomSmallWalkCaveUrbanised)
		if !g.Dungeon.connex(g.PR) {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}

func (g *game) initrand() {
	if g.rand == nil {
		g.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
}

func TestNaturalCave(t *testing.T) {
	for i := 0; i < Rounds; i++ {
		g := &game{}
		g.initrand()
		g.InitFirstLevel()
		g.InitLevelStructures()
		g.GenRoomTunnels(NaturalCave)
		if !g.Dungeon.connex(g.PR) {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}
