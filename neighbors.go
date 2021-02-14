package main

import "github.com/anaseto/gruid"

func Neighbors(p gruid.Point, nb []gruid.Point, keep func(gruid.Point) bool) []gruid.Point {
	neighbors := [8]gruid.Point{p.Add(gruid.Point{1, 0}), p.Add(gruid.Point{-1, 0}), p.Add(gruid.Point{0, -1}), p.Add(gruid.Point{0, 1}), p.Add(gruid.Point{1, -1}), p.Add(gruid.Point{-1, -1}), p.Add(gruid.Point{1, 1}), p.Add(gruid.Point{-1, 1})}
	nb = nb[:0]
	for _, q := range neighbors {
		if keep(q) {
			nb = append(nb, q)
		}
	}
	return nb
}

func CardinalNeighbors(p gruid.Point, nb []gruid.Point, keep func(gruid.Point) bool) []gruid.Point {
	neighbors := [4]gruid.Point{p.Add(gruid.Point{1, 0}), p.Add(gruid.Point{-1, 0}), p.Add(gruid.Point{0, -1}), p.Add(gruid.Point{0, 1})}
	nb = nb[:0]
	for _, q := range neighbors {
		if keep(q) {
			nb = append(nb, q)
		}
	}
	return nb
}

func ValidNeighbors(p gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 8)
	nb = Neighbors(p, nb, func(q gruid.Point) bool { return valid(q) })
	return nb
}

func ValidCardinalNeighbors(p gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 4)
	nb = CardinalNeighbors(p, nb, func(q gruid.Point) bool { return valid(q) })
	return nb
}

func (d *dungeon) IsFreeCell(p gruid.Point) bool {
	return valid(p) && d.Cell(p).IsPlayerPassable()
}

func (d *dungeon) NotWallCell(p gruid.Point) bool {
	return valid(p) && !d.Cell(p).IsWall()
}

func (d *dungeon) FreeNeighbors(p gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 8)
	nb = Neighbors(p, nb, d.IsFreeCell)
	return nb
}

func (d *dungeon) CardinalFreeNeighbors(p gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 4)
	nb = CardinalNeighbors(p, nb, d.IsFreeCell)
	return nb
}

func (d *dungeon) CardinalNonWallNeighbors(p gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 4)
	nb = CardinalNeighbors(p, nb, func(q gruid.Point) bool {
		return valid(q) && terrain(d.Cell(q)) != WallCell
	})

	return nb
}

func (d *dungeon) CardinalFlammableNeighbors(p gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 4)
	nb = CardinalNeighbors(p, nb, func(q gruid.Point) bool {
		return valid(q) && d.Cell(q).Flammable()
	})

	return nb
}
