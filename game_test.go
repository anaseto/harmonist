package main

//import "log"
import "testing"

func init() {
	Testing = true
}

func TestInitLevel(t *testing.T) {
	for i := 0; i < 50; i++ {
		g := &game{}
		for depth := 0; depth < MaxDepth; depth++ {
			g.InitLevel()
			if g.Player.P == InvalidPos {
				t.Errorf("Player starting cell is not valid")
			}
			if !terrain(g.Dungeon.Cell(g.Player.P)).IsPlayerPassable() {
				t.Errorf("Player starting cell is not passable: %+v", g.Dungeon.Cell(g.Player.P).ShortDesc(g, g.Player.P))
			}
			if terrain(g.Dungeon.Cell(g.Player.P)) != GroundCell {
				t.Errorf("Player starting cell is not ground: %+v", g.Dungeon.Cell(g.Player.P).ShortDesc(g, g.Player.P))
			}
			if g.Depth == WinDepth {
				if terrain(g.Dungeon.Cell(g.Places.Shaedra)) != StoryCell {
					t.Errorf("Shaedra not there: %+v", g.Places.Shaedra)
				}
				if g.Objects.Story[g.Places.Shaedra] != StoryShaedra {
					t.Errorf("bad Shaedra place: %+v", g.Places.Shaedra)
				}
			}
			if g.Depth == MaxDepth {
				if terrain(g.Dungeon.Cell(g.Places.Artifact)) != StoryCell {
					t.Errorf("Artifact not there: %+v", g.Places.Artifact)
				}
				if g.Objects.Story[g.Places.Artifact] != StoryArtifactSealed {
					t.Errorf("bad Artifact place: %+v", g.Places.Shaedra)
				}
			}
			if g.Depth == MaxDepth || g.Depth == WinDepth {
				if terrain(g.Dungeon.Cell(g.Places.Marevor)) != StoryCell {
					t.Errorf("Marevor not there: %+v", g.Places.Artifact)
				}
				if g.Objects.Story[g.Places.Marevor] != NoStory {
					t.Errorf("bad Marevor place: %+v", g.Places.Shaedra)
				}
			}
			if g.Depth == MaxDepth || g.Depth == WinDepth {
				if terrain(g.Dungeon.Cell(g.Places.Monolith)) != StoryCell {
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
				if !g.Dungeon.Cell(m.P).IsPassable() {
					t.Errorf("Not free: %+v", m.P)
				}
			}
			g.Depth++
		}
	}
}

func BenchmarkLOS(b *testing.B) {
	g := &game{}
	g.InitLevel()
	for i := 0; i < b.N; i++ {
		g.ComputeLOS()
	}
}

func BenchmarkLights(b *testing.B) {
	g := &game{}
	g.InitLevel()
	for i := 0; i < b.N; i++ {
		g.ComputeLights()
	}
}

func BenchmarkInitLevel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := &game{}
		for depth := 0; depth < MaxDepth; depth++ {
			g.InitLevel()
			g.Depth++
		}
	}
}
