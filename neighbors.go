package main

import "github.com/anaseto/gruid"

func Neighbors(pos gruid.Point, nb []gruid.Point, keep func(gruid.Point) bool) []gruid.Point {
	neighbors := [8]gruid.Point{pos.Add(gruid.Point{1, 0}), pos.Add(gruid.Point{-1, 0}), pos.Add(gruid.Point{0, -1}), pos.Add(gruid.Point{0, 1}), pos.Add(gruid.Point{1, -1}), pos.Add(gruid.Point{-1, -1}), pos.Add(gruid.Point{1, 1}), pos.Add(gruid.Point{-1, 1})}
	nb = nb[:0]
	for _, npos := range neighbors {
		if keep(npos) {
			nb = append(nb, npos)
		}
	}
	return nb
}

func CardinalNeighbors(pos gruid.Point, nb []gruid.Point, keep func(gruid.Point) bool) []gruid.Point {
	neighbors := [4]gruid.Point{pos.Add(gruid.Point{1, 0}), pos.Add(gruid.Point{-1, 0}), pos.Add(gruid.Point{0, -1}), pos.Add(gruid.Point{0, 1})}
	nb = nb[:0]
	for _, npos := range neighbors {
		if keep(npos) {
			nb = append(nb, npos)
		}
	}
	return nb
}

func ValidNeighbors(pos gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 8)
	nb = Neighbors(pos, nb, func(npos gruid.Point) bool { return valid(npos) })
	return nb
}

func ValidCardinalNeighbors(pos gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 4)
	nb = CardinalNeighbors(pos, nb, func(npos gruid.Point) bool { return valid(npos) })
	return nb
}

func (d *dungeon) IsFreeCell(pos gruid.Point) bool {
	return valid(pos) && d.Cell(pos).IsPlayerPassable()
}

func (d *dungeon) NotWallCell(pos gruid.Point) bool {
	return valid(pos) && !d.Cell(pos).IsWall()
}

func (d *dungeon) FreeNeighbors(pos gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 8)
	nb = Neighbors(pos, nb, d.IsFreeCell)
	return nb
}

func (d *dungeon) CardinalFreeNeighbors(pos gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 4)
	nb = CardinalNeighbors(pos, nb, d.IsFreeCell)
	return nb
}

func (d *dungeon) CardinalNonWallNeighbors(pos gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 4)
	nb = CardinalNeighbors(pos, nb, func(npos gruid.Point) bool {
		return valid(npos) && terrain(d.Cell(npos)) != WallCell
	})

	return nb
}

func (d *dungeon) CardinalFlammableNeighbors(pos gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 4)
	nb = CardinalNeighbors(pos, nb, func(npos gruid.Point) bool {
		return valid(npos) && d.Cell(npos).Flammable()
	})

	return nb
}
