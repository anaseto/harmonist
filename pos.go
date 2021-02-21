package main

import (
	"github.com/anaseto/gruid"
)

func distance(from, to gruid.Point) int {
	delta := to.Sub(from)
	return abs(delta.X) + abs(delta.Y)
}

func distanceChebyshev(from, to gruid.Point) int {
	delta := to.Sub(from)
	deltaX := abs(delta.X)
	deltaY := abs(delta.Y)
	if deltaX > deltaY {
		return deltaX
	}
	return deltaY
}

// ZP is the zero value for gruid.Point.
var ZP gruid.Point = gruid.Point{}

func dirString(dir gruid.Point) (s string) {
	switch dir {
	case ZP:
		s = ""
	case gruid.Point{1, 0}:
		s = "E"
	case gruid.Point{1, -1}:
		s = "NE"
	case gruid.Point{0, -1}:
		s = "N"
	case gruid.Point{-1, -1}:
		s = "NW"
	case gruid.Point{-1, 0}:
		s = "W"
	case gruid.Point{-1, 1}:
		s = "SW"
	case gruid.Point{0, 1}:
		s = "S"
	case gruid.Point{1, 1}:
		s = "SE"
	}
	return s
}

func keyToDir(k action) (p gruid.Point) {
	switch k {
	case ActionW, ActionRunW:
		p = gruid.Point{-1, 0}
	case ActionE, ActionRunE:
		p = gruid.Point{1, 0}
	case ActionS, ActionRunS:
		p = gruid.Point{0, 1}
	case ActionN, ActionRunN:
		p = gruid.Point{0, -1}
	}
	return p
}

func sign(n int) int {
	var i int
	switch {
	case n > 0:
		i = 1
	case n < 0:
		i = -1
	}
	return i
}

// dirnorm returns a normalized direction between two points, so that
// directions that aren't cardinal nor diagonal are transformed into the
// cardinal part (this corresponds to pruned intermediate nodes in diagonal
// jump).
func dirnorm(p, q gruid.Point) gruid.Point {
	dir := q.Sub(p)
	dx := abs(dir.X)
	dy := abs(dir.Y)
	dir = gruid.Point{sign(dir.X), sign(dir.Y)}
	switch {
	case dx == dy:
	case dx > dy:
		dir.Y = 0
	default:
		dir.X = 0
	}
	return dir
}

func idxtopos(i int) gruid.Point {
	return gruid.Point{i % DungeonWidth, i / DungeonWidth}
}

func idx(p gruid.Point) int {
	return p.Y*DungeonWidth + p.X
}

func valid(p gruid.Point) bool {
	return p.Y >= 0 && p.Y < DungeonHeight && p.X >= 0 && p.X < DungeonWidth
}

func inViewCone(dir, from, to gruid.Point) bool {
	if to == from || distance(from, to) <= 1 {
		return true
	}
	d := dirnorm(from, to)
	return d == dir || leftDir(d) == dir || rightDir(d) == dir
}

func leftDir(dir gruid.Point) gruid.Point {
	switch {
	case dir.X == 0 || dir.Y == 0:
		return left(dir, dir)
	default:
		return gruid.Point{(dir.Y + dir.X) / 2, (dir.Y - dir.X) / 2}
	}
}

func rightDir(dir gruid.Point) gruid.Point {
	switch {
	case dir.X == 0 || dir.Y == 0:
		return right(dir, dir)
	default:
		return gruid.Point{(dir.X - dir.Y) / 2, (dir.Y + dir.X) / 2}
	}
}

func right(p gruid.Point, dir gruid.Point) gruid.Point {
	return gruid.Point{p.X - dir.Y, p.Y + dir.X}
}

func left(p gruid.Point, dir gruid.Point) gruid.Point {
	return gruid.Point{p.X + dir.Y, p.Y - dir.X}
}
