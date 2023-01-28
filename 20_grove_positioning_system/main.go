package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(solvePartTwo(file))
}

type node struct {
	v    int
	next *node
	prev *node
}

func (n *node) forward() {
	ftmp := n.next
	btmp := n.prev

	n.next = n.next.next
	n.prev = ftmp

	n.next.prev = n

	ftmp.next = n
	ftmp.prev = btmp

	btmp.next = ftmp
}

func (n *node) back() {
	ftmp := n.next
	btmp := n.prev

	n.prev = n.prev.prev
	n.next = btmp

	n.prev.next = n

	btmp.next = ftmp
	btmp.prev = n

	ftmp.prev = btmp
}

func (n *node) move(count, circuitLen int) {
	moves := int(math.Abs(float64(count))) % circuitLen
	for moves > 0 {
		moves--
		if count > 0 {
			n.forward()
			continue
		}
		n.back()
	}
}

type index map[int]*node

type nodeList []*node

type dll struct {
	nodeList
	index
}

func newDLL() *dll {
	out := dll{
		make(nodeList, 0),
		make(index),
	}

	return &out
}

func (d *dll) nodeFromValue(value int) *node {
	return d.index[value]
}

func (d *dll) append(val int, mapFn func(v int) int) {
	newNodeRef := &node{mapFn(val), nil, nil}
	// if list is empty, the new node points to itself to maintain the 'circleness'
	// of the structure
	if len(d.nodeList) == 0 {
		newNodeRef.prev = newNodeRef
		newNodeRef.next = newNodeRef
	} else {
		head := d.nodeList[0]
		tail := d.nodeList[len(d.nodeList)-1]
		tail.next = newNodeRef
		newNodeRef.prev = tail
		newNodeRef.next = head
		head.prev = newNodeRef
	}

	d.nodeList = append(d.nodeList, newNodeRef)
	d.index[val] = newNodeRef
}

func (d *dll) mix(rounds int) {
	for rounds > 0 {
		for _, node := range d.nodeList {
			node.move(node.v, len(d.nodeList)-1)
		}
		rounds--
	}
}

func (d *dll) sum(start int, sumPositions []int) int {
	currNode := d.nodeFromValue(start)
	sum := 0
	currPos := 0
	targetPos := 0
	for targetPos < len(sumPositions) {
		currNode = currNode.next
		currPos++
		if currPos == sumPositions[targetPos] {
			sum += currNode.v
			targetPos++
		}
	}
	return sum
}

func parseInput(r io.Reader, mapFn func(v int) int) *dll {
	scanner := bufio.NewScanner(r)
	dl := newDLL()
	for scanner.Scan() {
		ln := scanner.Text()
		v, err := strconv.ParseInt(ln, 10, 32)
		if err != nil {
			panic(err)
		}
		dl.append(int(v), mapFn)
	}
	return dl
}

func id[T any](v T) T {
	return v
}

func solvePartOne(r io.Reader) int {
	dl := parseInput(r, id[int])
	dl.mix(1)
	return dl.sum(0, []int{1000, 2000, 3000})
}

func solvePartTwo(r io.Reader) int {
	mapInput := func(v int) int {
		return v * 811589153
	}
	dl := parseInput(r, mapInput)
	dl.mix(10)
	return dl.sum(0, []int{1000, 2000, 3000})
}
