package main

import (
	"container/heap"

	"github.com/anaseto/gruid"
)

type Dijkstrer interface {
	Neighbors(gruid.Point) []gruid.Point
	Cost(gruid.Point, gruid.Point) int
}

func Dijkstra(dij Dijkstrer, sources []gruid.Point, maxCost int) nodeMap {
	nodeCache.Index++
	nqs := queueCache[:0]
	nq := &nqs
	heap.Init(nq)
	for _, f := range sources {
		n := nodeCache.get(f)
		n.Open = true
		heap.Push(nq, n)
	}
	for {
		if nq.Len() == 0 {
			return nodeCache
		}
		current := heap.Pop(nq).(*node)
		current.Open = false
		current.Closed = true

		for _, neighbor := range dij.Neighbors(current.Pos) {
			cost := current.Cost + dij.Cost(current.Pos, neighbor)
			neighborNode := nodeCache.get(neighbor)
			if cost < neighborNode.Cost {
				if neighborNode.Open {
					heap.Remove(nq, neighborNode.Index)
				}
				neighborNode.Open = false
				neighborNode.Closed = false
			}
			if !neighborNode.Open && !neighborNode.Closed {
				neighborNode.Cost = cost
				if cost < maxCost {
					neighborNode.Open = true
					neighborNode.Rank = cost
					heap.Push(nq, neighborNode)
				}
			}
		}
	}
}

const unreachable = 9999

// AutoExploreDijkstra is an optimized version of the dijkstra algorithm for
// auto-exploration.
func (g *game) AutoExploreDijkstra(dij Dijkstrer, sources []int) {
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
			if !valid(npos) || d.Cells[nidx].IsWall() { // XXX: IsWall ?
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
