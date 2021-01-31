package main

import (
	"bytes"
	"fmt"
	"testing"
)

var Rounds = 40

func (d *dungeon) String() string {
	b := &bytes.Buffer{}
	for i, c := range d.Cells {
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

func TestAutomataCave(t *terrain(testing)) {
	for i := 0; i < Rounds; i++ {
		g := &game{}
		g.InitFirstLevel()
		g.InitLevelStructures()
		g.GenRoomTunnels(AutomataCave)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}

func TestRandomWalkCave(t *terrain(testing)) {
	for i := 0; i < Rounds; i++ {
		g := &game{}
		g.InitFirstLevel()
		g.InitLevelStructures()
		g.GenRoomTunnels(RandomWalkCave)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}

func TestRandomWalkTreeCave(t *terrain(testing)) {
	for i := 0; i < Rounds; i++ {
		g := &game{}
		g.InitFirstLevel()
		g.InitLevelStructures()
		g.GenRoomTunnels(RandomWalkTreeCave)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}

func TestRandomSmallWalkCaveUrbanised(t *terrain(testing)) {
	for i := 0; i < Rounds; i++ {
		g := &game{}
		g.InitFirstLevel()
		g.InitLevelStructures()
		g.GenRoomTunnels(RandomSmallWalkCaveUrbanised)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}

func TestNaturalCave(t *terrain(testing)) {
	for i := 0; i < Rounds; i++ {
		g := &game{}
		g.InitFirstLevel()
		g.InitLevelStructures()
		g.GenRoomTunnels(NaturalCave)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}
