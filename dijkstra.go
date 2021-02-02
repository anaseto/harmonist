package main

import (
	"github.com/anaseto/gruid/paths"
)

const unreachable = 9999

// AutoExploreDijkstra is an optimized version of the dijkstra algorithm for
// auto-exploration.
func (g *game) AutoExploreDijkstra(dij paths.Dijkstra, sources []int) {
	d := g.Dungeon
	dmap := DijkstraMapCache[:]
	var visited [DungeonNCells]bool
	var queue [DungeonNCells]int
	var qstart, qend int
	for i := 0; i < DungeonNCells; i++ {
		dmap[i] = unreachable
	}
	for _, s := range sources {
		dmap[s] = 0
		queue[qend] = s
		qend++
		visited[s] = true
	}
	for qstart < qend {
		cidx := queue[qstart]
		qstart++
		cpos := idxtopos(cidx)
		for _, npos := range dij.Neighbors(cpos) {
			nidx := idx(npos)
			if !valid(npos) || d.Cell(npos).IsWall() { // XXX: IsWall ?
				continue
			}
			if !visited[nidx] {
				queue[qend] = nidx
				qend++
				visited[nidx] = true
				dmap[nidx] = 1 + dmap[cidx]
			}
		}
	}
}
