package main

import "testing"

func TestDir(t *testing.T) {
	type tableTest struct {
		pos gruid.Point
		dir direction
	}
	table := []tableTest{
		{gruid.Point{3, 2}, E},
		{gruid.Point{4, 1}, ENE},
		{gruid.Point{3, 1}, NE},
		{gruid.Point{3, 0}, NNE},
		{gruid.Point{2, 1}, N},
		{gruid.Point{1, 0}, NNW},
		{gruid.Point{1, 1}, NW},
		{gruid.Point{0, 1}, WNW},
		{gruid.Point{1, 2}, W},
		{gruid.Point{0, 3}, WSW},
		{gruid.Point{1, 3}, SW},
		{gruid.Point{1, 4}, SSW},
		{gruid.Point{2, 3}, S},
		{gruid.Point{3, 4}, SSE},
		{gruid.Point{3, 3}, SE},
		{gruid.Point{4, 3}, ESE},
	}
	for _, test := range table {
		if Dir(gruid.Point{2, 2}, test.pos) != test.dir {
			t.Errorf("Bad direction for %+v\n", test)
		}
	}
}
