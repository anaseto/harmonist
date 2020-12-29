package main

import (
	"fmt"
)

func Distance(from, to gruid.Point) int {
	delta := to.Sub(from)
	return Abs(delta.X) + Abs(delta.Y)
}

func MaxCardinalDist(from, to gruid.Point) int {
	delta := to.Sub(from)
	deltaX := Abs(delta.X)
	deltaY := Abs(delta.Y)
	if deltaX > deltaY {
		return deltaX
	}
	return deltaY
}

func DistanceX(from, to gruid.Point) int {
	deltaX := Abs(to.X - from.X)
	return deltaX
}

func DistanceY(from, to gruid.Point) int {
	deltaY := Abs(to.Y - from.Y)
	return deltaY
}

type direction int

const (
	NoDir direction = iota
	E
	ENE
	NE
	NNE
	N
	NNW
	NW
	WNW
	W
	WSW
	SW
	SSW
	S
	SSE
	SE
	ESE
)

func (dir direction) String() (s string) {
	switch dir {
	case NoDir:
		s = ""
	case E:
		s = "E"
	case ENE:
		s = "ENE"
	case NE:
		s = "NE"
	case NNE:
		s = "NNE"
	case N:
		s = "N"
	case NNW:
		s = "NNW"
	case NW:
		s = "NW"
	case WNW:
		s = "WNW"
	case W:
		s = "W"
	case WSW:
		s = "WSW"
	case SW:
		s = "SW"
	case SSW:
		s = "SSW"
	case S:
		s = "S"
	case SSE:
		s = "SSE"
	case SE:
		s = "SE"
	case ESE:
		s = "ESE"
	}
	return s
}

func KeyToDir(k action) (dir direction) {
	switch k {
	case ActionW, ActionRunW:
		dir = W
	case ActionE, ActionRunE:
		dir = E
	case ActionS, ActionRunS:
		dir = S
	case ActionN, ActionRunN:
		dir = N
	}
	return dir
}

func To(dir direction, from gruid.Point) gruid.Point {
	to := from
	switch dir {
	case E, ENE, ESE:
		to = pos.Add(gruid.Point{1, 0})
	case NE:
		to = from.Add(gruid.Point{1, -1})
	case NNE, N, NNW:
		to = from.Add(gruid.Point{0, -1})
	case NW:
		to = from.Add(gruid.Point{-1, -1})
	case WNW, W, WSW:
		to = from.Add(gruid.Point{-1, 0})
	case SW:
		to = from.Add(gruid.Point{-1, 1})
	case SSW, S, SSE:
		to = from.Add(gruid.Point{0, 1})
	case SE:
		to = from.Add(gruid.Point{1, 1})
	}
	return to
}

func Dir(from, to gruid.Point) direction {
	deltaX := Abs(to.X - from.X)
	deltaY := Abs(to.Y - from.Y)
	switch {
	case to.X > from.X && to.Y == from.Y:
		return E
	case to.X > from.X && to.Y < from.Y:
		switch {
		case deltaX > deltaY:
			return ENE
		case deltaX == deltaY:
			return NE
		default:
			return NNE
		}
	case to.X == from.X && to.Y < from.Y:
		return N
	case to.X < from.X && to.Y < from.Y:
		switch {
		case deltaY > deltaX:
			return NNW
		case deltaX == deltaY:
			return NW
		default:
			return WNW
		}
	case to.X < from.X && to.Y == from.Y:
		return W
	case to.X < from.X && to.Y > from.Y:
		switch {
		case deltaX > deltaY:
			return WSW
		case deltaX == deltaY:
			return SW
		default:
			return SSW
		}
	case to.X == from.X && to.Y > from.Y:
		return S
	case to.X > from.X && to.Y > from.Y:
		switch {
		case deltaY > deltaX:
			return SSE
		case deltaX == deltaY:
			return SE
		default:
			return ESE
		}
	default:
		panic(fmt.Sprintf("internal error: invalid gruid.Point:%+v-%+v", to, from))
	}
}

func Parents(pos, from gruid.Point, p []gruid.Point) []gruid.Point {
	switch pos.Dir(from) {
	case E:
		p = append(p, pos.Add(gruid.Point{-1, 0}))
	case ENE:
		p = append(p, pos.Add(gruid.Point{-1, 0}), pos.Add(gruid.Point{-1, 1}))
	case NE:
		p = append(p, pos.Add(gruid.Point{-1, 1}))
	case NNE:
		p = append(p, pos.Add(gruid.Point{0, 1}), pos.Add(gruid.Point{-1, 1}))
	case N:
		p = append(p, pos.Add(gruid.Point{0, 1}))
	case NNW:
		p = append(p, pos.Add(gruid.Point{0, 1}), pos.Add(gruid.Point{1, 1}))
	case NW:
		p = append(p, pos.Add(gruid.Point{1, 1}))
	case WNW:
		p = append(p, pos.Add(gruid.Point{1, 0}), pos.Add(gruid.Point{1, 1}))
	case W:
		p = append(p, pos.Add(gruid.Point{1, 0}))
	case WSW:
		p = append(p, pos.Add(gruid.Point{1, 0}), pos.Add(gruid.Point{1, -1}))
	case SW:
		p = append(p, pos.Add(gruid.Point{1, -1}))
	case SSW:
		p = append(p, pos.Add(gruid.Point{0, -1}), pos.Add(gruid.Point{1, -1}))
	case S:
		p = append(p, pos.Add(gruid.Point{0, -1}))
	case SSE:
		p = append(p, pos.Add(gruid.Point{0, -1}), pos.Add(gruid.Point{-1, -1}))
	case SE:
		p = append(p, pos.Add(gruid.Point{-1, -1}))
	case ESE:
		p = append(p, pos.Add(gruid.Point{-1, 0}), pos.Add(gruid.Point{-1, -1}))
	}
	return p
}

func RandomNeighbor(pos gruid.Point, diag bool) gruid.Point {
	if diag {
		return RandomNeighborDiagonals(pos)
	}
	return RandomNeighborCardinal(pos)
}

func RandomNeighborDiagonals(pos gruid.Point) gruid.Point {
	neighbors := [8]gruid.Point{pos.Add(gruid.Point{1, 0}), pos.Add(gruid.Point{-1, 0}), pos.Add(gruid.Point{0, -1}), pos.Add(gruid.Point{0, 1}), pos.Add(gruid.Point{1, -1}), pos.Add(gruid.Point{-1, -1}), pos.Add(gruid.Point{1, 1}), pos.Add(gruid.Point{-1, 1})}
	var r int
	switch RandInt(8) {
	case 0:
		r = RandInt(len(neighbors[0:4]))
	case 1:
		r = RandInt(len(neighbors[0:2]))
	default:
		r = RandInt(len(neighbors[4:]))
	}
	return neighbors[r]
}

func RandomNeighborCardinal(pos gruid.Point) gruid.Point {
	neighbors := [4]gruid.Point{pos.Add(gruid.Point{1, 0}), pos.Add(gruid.Point{-1, 0}), pos.Add(gruid.Point{0, -1}), pos.Add(gruid.Point{0, 1})}
	var r int
	switch RandInt(4) {
	case 0, 1:
		r = RandInt(len(neighbors[0:2]))
	default:
		r = RandInt(len(neighbors))
	}
	return neighbors[r]
}

func idxtopos(i int) gruid.Point {
	return gruid.Point{i % DungeonWidth, i / DungeonWidth}
}

func idx(pos gruid.Point) int {
	return pos.Y*DungeonWidth + pos.X
}

func valid(pos gruid.Point) bool {
	return pos.Y >= 0 && pos.Y < DungeonHeight && pos.X >= 0 && pos.X < DungeonWidth
}

func (dir direction) InViewCone(from, to gruid.Point) bool {
	if to == from {
		return true
	}
	d := to.Dir(from)
	if d == dir || Distance(from, to) <= 1 {
		return true
	}
	switch dir {
	case E:
		switch d {
		case ESE, ENE, NE, SE:
			return true
		}
	case NE:
		switch d {
		case ENE, NNE, N, E:
			return true
		}
	case N:
		switch d {
		case NNE, NNW, NE, NW:
			return true
		}
	case NW:
		switch d {
		case NNW, WNW, N, W:
			return true
		}
	case W:
		switch d {
		case WNW, WSW, NW, SW:
			return true
		}
	case SW:
		switch d {
		case WSW, SSW, W, S:
			return true
		}
	case S:
		switch d {
		case SSW, SSE, SW, SE:
			return true
		}
	case SE:
		switch d {
		case SSE, ESE, S, E:
			return true
		}
	}
	return false
}

var alternateDirs = []direction{E, NE, N, NW, W, SW, S, SE}

func (dir direction) Left() (d direction) {
	switch dir {
	case E:
		d = NE
	case NE:
		d = N
	case N:
		d = NW
	case NW:
		d = W
	case W:
		d = SW
	case SW:
		d = S
	case S:
		d = SE
	case SE:
		d = E
	default:
		d = alternateDirs[RandInt(len(alternateDirs))]
	}
	return d
}

func (dir direction) Right() (d direction) {
	switch dir {
	case E:
		d = SE
	case NE:
		d = E
	case N:
		d = NE
	case NW:
		d = N
	case W:
		d = NW
	case SW:
		d = W
	case S:
		d = SW
	case SE:
		d = S
	default:
		d = alternateDirs[RandInt(len(alternateDirs))]
	}
	return d
}
