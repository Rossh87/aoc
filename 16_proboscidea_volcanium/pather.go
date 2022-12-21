package main

import (
	"container/heap"
)

type node struct {
	name string
	dist int
}

type minHeap struct {
	seen map[string]struct{}
	q    []node
}

func (h minHeap) Len() int {
	return len(h.q)
}

func (h minHeap) Less(i, j int) bool {
	return h.q[i].dist < h.q[j].dist
}

func (h minHeap) Swap(i, j int) {
	h.q[i], h.q[j] = h.q[j], h.q[i]
}

func (h *minHeap) Push(x any) {
	item := x.(node)
	(*h).q = append((*h).q, item)
}

func (h *minHeap) Pop() any {
	end := len(h.q) - 1
	popped := h.q[end]
	h.q = h.q[:end]
	return popped
}

func (h *minHeap) Mark(name string) {
	h.seen[name] = struct{}{}
}

func (h minHeap) Seen(name string) bool {
	if _, ok := h.seen[name]; ok {
		return true
	}

	return false
}

func newMinHeap() *minHeap {
	h := minHeap{
		make(map[string]struct{}),
		[]node{},
	}
	return &h
}

// map[name] = distance from currValve to valve <name>
func djikstra(currValve string, vs valves) map[string]int {
	distances := make(map[string]int)
	mh := newMinHeap()
	heap.Init(mh)
	heap.Push(mh, node{currValve, 0})

	for mh.Len() > 0 {
		currNode := heap.Pop(mh).(node)
		if mh.Seen(currNode.name) {
			continue
		}
		mh.Mark(currNode.name)
		distances[currNode.name] = currNode.dist
		valve := vs[currNode.name]
		for _, neighborName := range valve.neighbors {
			heap.Push(mh, node{neighborName, currNode.dist + 1})
		}
	}
	return distances
}

type distances map[string]map[string]int

type pather struct {
	paths distances
	vs    valves
}

// todo: extremely lazy to not populate inverse path with chart
func (p *pather) pathLength(from, to string) int {
	if p.paths == nil {
		p.paths = make(distances)
	}

	if distChart, exists := p.paths[from]; exists {
		return distChart[to]
	}

	distChart := djikstra(from, p.vs)
	p.paths[from] = distChart
	return distChart[to]
}

func newPather(vs valves) *pather {
	ds := make(map[string]map[string]int)
	pather := pather{
		ds,
		vs,
	}
	return &pather
}
