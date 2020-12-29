package main

import "github.com/anaseto/gruid"

func (pos gruid.Point) Neighbors(nb []gruid.Point, keep func(gruid.Point) bool) []gruid.Point {
	neighbors := [8]gruid.Point{pos.Add(gruid.Point{1, 0}), pos.Add(gruid.Point{-1, 0}), pos.Add(gruid.Point{0, -1}), pos.Add(gruid.Point{0, 1}), pos.Add(gruid.Point{1, -1}), pos.Add(gruid.Point{-1, -1}), pos.Add(gruid.Point{1, 1}), pos.Add(gruid.Point{-1, 1})}
	nb = nb[:0]
	for _, npos := range neighbors {
		if keep(npos) {
			nb = append(nb, npos)
		}
	}
	return nb
}

func (pos gruid.Point) CardinalNeighbors(nb []gruid.Point, keep func(gruid.Point) bool) []gruid.Point {
	neighbors := [4]gruid.Point{pos.Add(gruid.Point{1, 0}), pos.Add(gruid.Point{-1, 0}), pos.Add(gruid.Point{0, -1}), pos.Add(gruid.Point{0, 1})}
	nb = nb[:0]
	for _, npos := range neighbors {
		if keep(npos) {
			nb = append(nb, npos)
		}
	}
	return nb
}

func (pos gruid.Point) OutsideNeighbors() []gruid.Point {
	nb := make([]gruid.Point, 0, 8)
	nb = pos.Neighbors(nb, func(npos gruid.Point) bool {
		return !valid(npos)
	})
	return nb
}

func (pos gruid.Point) ValidNeighbors() []gruid.Point {
	nb := make([]gruid.Point, 0, 8)
	nb = pos.Neighbors(nb, func(npos gruid.Point) bool { return valid(npos) })
	return nb
}

func (pos gruid.Point) ValidCardinalNeighbors() []gruid.Point {
	nb := make([]gruid.Point, 0, 4)
	nb = pos.CardinalNeighbors(nb, func(npos gruid.Point) bool { return valid(npos) })
	return nb
}

func (d *dungeon) IsFreeCell(pos gruid.Point) bool {
	return valid(pos) && d.Cell(pos).T.IsPlayerPassable()
}

func (d *dungeon) NotWallCell(pos gruid.Point) bool {
	return valid(pos) && !d.Cell(pos).IsWall()
}

func (d *dungeon) FreeNeighbors(pos gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 8)
	nb = pos.Neighbors(nb, d.IsFreeCell)
	return nb
}

func (d *dungeon) CardinalFreeNeighbors(pos gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 4)
	nb = pos.CardinalNeighbors(nb, d.IsFreeCell)
	return nb
}

func (d *dungeon) CardinalNonWallNeighbors(pos gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 4)
	nb = pos.CardinalNeighbors(nb, func(npos gruid.Point) bool {
		return valid(npos) && d.Cell(npos).T != WallCell
	})
	return nb
}

func (d *dungeon) CardinalFlammableNeighbors(pos gruid.Point) []gruid.Point {
	nb := make([]gruid.Point, 0, 4)
	nb = pos.CardinalNeighbors(nb, func(npos gruid.Point) bool {
		return valid(npos) && d.Cell(npos).Flammable()
	})
	return nb
}
