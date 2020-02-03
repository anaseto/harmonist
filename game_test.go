package main

import "testing"

func TestInitLevel(t *testing.T) {
	Testing = true
	for i := 0; i < 50; i++ {
		g := &game{}
		for depth := 0; depth < MaxDepth; depth++ {
			g.InitLevel()
			if g.Player.Pos == InvalidPos {
				t.Errorf("Player starting cell is not valid")
			}
			if !g.Dungeon.Cell(g.Player.Pos).T.IsPlayerPassable() {
				t.Errorf("Player starting cell is not passable: %+v", g.Dungeon.Cell(g.Player.Pos).ShortDesc(g, g.Player.Pos))
			}
			if g.Dungeon.Cell(g.Player.Pos).T != GroundCell {
				t.Errorf("Player starting cell is not ground: %+v", g.Dungeon.Cell(g.Player.Pos).ShortDesc(g, g.Player.Pos))
			}
			if g.Depth == WinDepth {
				if g.Dungeon.Cell(g.Places.Shaedra).T != StoryCell {
					t.Errorf("Shaedra not there: %+v", g.Places.Shaedra)
				}
				if g.Objects.Story[g.Places.Shaedra] != StoryShaedra {
					t.Errorf("bad Shaedra place: %+v", g.Places.Shaedra)
				}
			}
			if g.Depth == MaxDepth {
				if g.Dungeon.Cell(g.Places.Artifact).T != StoryCell {
					t.Errorf("Artifact not there: %+v", g.Places.Artifact)
				}
				if g.Objects.Story[g.Places.Artifact] != StoryArtifactSealed {
					t.Errorf("bad Artifact place: %+v", g.Places.Shaedra)
				}
			}
			if g.Depth == MaxDepth || g.Depth == WinDepth {
				if g.Dungeon.Cell(g.Places.Marevor).T != StoryCell {
					t.Errorf("Marevor not there: %+v", g.Places.Artifact)
				}
				if g.Objects.Story[g.Places.Marevor] != NoStory {
					t.Errorf("bad Marevor place: %+v", g.Places.Shaedra)
				}
			}
			if g.Depth == MaxDepth || g.Depth == WinDepth {
				if g.Dungeon.Cell(g.Places.Monolith).T != StoryCell {
					t.Errorf("Monolith not there: %+v", g.Places.Artifact)
				}
				if g.Objects.Story[g.Places.Monolith] != NoStory {
					t.Errorf("bad Monolith place: %+v", g.Places.Shaedra)
				}
			}
			if g.Depth != WinDepth {
				if len(g.Objects.Magaras) != 1 {
					t.Errorf("bad number of magaras: %+v", g.Objects.Magaras)
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
