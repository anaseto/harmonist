package main

import "github.com/anaseto/gruid"

func (g *game) cardinalNeighbors(p gruid.Point) []gruid.Point {
	return g.nbs.Cardinal(p, valid)
}

func (g *game) playerPassableNeighbors(p gruid.Point) []gruid.Point {
	d := g.Dungeon
	return g.nbs.All(p, func(q gruid.Point) bool {
		return valid(q) && d.Cell(q).IsPlayerPassable()
	})
}

func (g *game) nonWallNeighbors(p gruid.Point) []gruid.Point {
	return g.nbs.Cardinal(p, func(q gruid.Point) bool {
		return valid(q) && terrain(g.Dungeon.Cell(q)) != WallCell
	})
}

func (g *game) flammableNeighbors(p gruid.Point) []gruid.Point {
	return g.nbs.Cardinal(p, func(q gruid.Point) bool {
		return valid(q) && g.Dungeon.Cell(q).Flammable()
	})
}

func randomNeighbor(p gruid.Point) gruid.Point {
	switch RandInt(6) {
	case 0, 1:
		return p.Add(gruid.Point{1, 0})
	case 2, 3:
		return p.Add(gruid.Point{-1, 0})
	case 4:
		return p.Add(gruid.Point{0, 1})
	default:
		return p.Add(gruid.Point{0, -1})
	}
}
