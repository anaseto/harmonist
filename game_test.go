package main

import "testing"

func TestInitLevel(t *testing.T) {
	Testing = true
	for i := 0; i < 50; i++ {
		g := &game{}
		for depth := 0; depth < 11; depth++ {
			g.InitLevel()
			if g.Depth == WinDepth {
				if g.Dungeon.Cell(g.Places.Shaedra).T != StoryCell {
					t.Errorf("Shaedra not there: %+v", g.Places.Shaedra)
				}
			}
			for _, m := range g.Monsters {
				if !g.Dungeon.Cell(m.Pos).IsPassable() {
					t.Errorf("Not free: %+v", m.Pos)
				}
			}
			g.Depth++
		}
	}
}
