package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sync"
)

func main() {
	file, err := os.Open("./input.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fmt.Println(navigateTwo(file))
}

type grid [][]int

// An ICoord represents a point in the 2D grid.
type ICoord interface {
	X() int
	Y() int
}

type coord struct {
	x int
	y int
}

func (c coord) X() int {
	return c.x
}

func (c coord) Y() int {
	return c.y
}

func (c coord) String() string {
	return fmt.Sprintf("%d--%d", c.X(), c.Y())
}

// A cell is a 'coord' with heuristic info added. Cells with a
// LOWER f value have higher priority than cells with HIGHER
// f value.
// f: total heuristic value
// g: movement distance value. Represents the distance of the cell from
// the 'from' cell.
// h: heuristic distance value.  Represents the estimated distance of the cell
// from the 'to' cell.
type cell struct {
	coord
	f      int
	g      int
	h      int
	parent *cell
}

type parsingResult struct {
	g         grid
	source    coord
	dest      coord
	lowPoints []coord
}

type directions [4]coord

func newParsingResult() parsingResult {
	s := coord{}
	d := coord{}
	g := [][]int{}
	l := []coord{}

	return parsingResult{g, s, d, l}
}

func charToHeight(c rune) int {
	return int(c - 'a')
}

func parseInput(input io.Reader) parsingResult {
	scanner := bufio.NewScanner(input)

	res := newParsingResult()
	currRow := 0

	for scanner.Scan() {
		ln := scanner.Text()
		row := make([]int, len(ln))
		rnIdx := 0
		for _, rn := range ln {
			switch rn {
			case 'S':
				row[rnIdx] = 0
				res.source.x = rnIdx
				res.source.y = currRow
				res.lowPoints = append(res.lowPoints, coord{rnIdx, currRow})
			case 'E':
				row[rnIdx] = 25
				res.dest.x = rnIdx
				res.dest.y = currRow
			default:
				height := charToHeight(rn)
				row[rnIdx] = height
				if height == 0 {
					res.lowPoints = append(res.lowPoints, coord{rnIdx, currRow})
				}
			}
			rnIdx++
		}
		res.g = append(res.g, row)
		currRow++
	}

	return res
}

// Because we have 2D grid with 2 axes of movement, this simple
// calculation works to approximate how close a candidate cell is to the
// target cell.
func calcH(target, candidate ICoord) int {
	offSetX := math.Abs(float64(target.X()) - float64(candidate.X()))
	offSetY := math.Abs(float64(target.Y()) - float64(candidate.Y()))
	return int(offSetX) + int(offSetY)
}

// // Set heuristic value for a candidate cell based on movement
// // cost to the candidate cell, and the proximity of the candidate
// // cell to the target cell.
func newNeighbor(current cell, x, y int, dest ICoord) cell {
	n := cell{}
	n.x = x
	n.y = y
	n.g = current.g + 1
	// NB we MUST set X and Y of n before we can call calcH() function
	n.h = calcH(dest, n)
	n.f = n.g + n.h
	n.parent = &current
	return n
}

type priorityQueue []cell

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].f < pq[j].f
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x any) {
	next := x.(cell)
	*pq = append(*pq, next)
}

func (pq *priorityQueue) Pop() any {
	endPos := len((*pq)) - 1
	popped := (*pq)[endPos]
	*pq = (*pq)[:endPos]
	return popped
}

// openQ keeps track of cells that are candidates for our next move choice
type openQ struct {
	q priorityQueue
}

func newOpenQ() *openQ {
	baseQ := priorityQueue{}
	oq := openQ{}
	oq.q = baseQ
	heap.Init(&baseQ)
	return &oq
}

func (oq *openQ) pop() cell {
	baseQ := &oq.q
	el := heap.Pop(baseQ).(cell)
	return el
}

func (oq *openQ) push(c cell) {
	baseQ := &oq.q
	heap.Push(baseQ, c)
}

func (oq openQ) len() int {
	return oq.q.Len()
}

func (oq *openQ) remove(c cell) int {
	removalID := c.String()

	removalIdx := -1

	for idx, el := range oq.q {
		id := el.String()
		if id == removalID {
			removalIdx = idx
			break
		}
	}

	if removalIdx == -1 {
		return removalIdx
	}

	tailPos := len(oq.q) - 1
	tail := oq.q[tailPos]
	oq.q[removalIdx] = tail
	oq.q = oq.q[:tailPos]
	heap.Init(&oq.q)
	return removalIdx
}

// Checks if the open list already has a cell with the SAME
// coordinates as the candidate cell.  If it does, it returns
// the 'f' value of that cell.  Otherwise, it returns -1.
func (oq openQ) has(c cell) int {
	candidateID := c.String()

	for _, el := range oq.q {
		existingID := el.String()

		if existingID == candidateID {
			return el.f
		}
	}

	return -1
}

// Holds identity and heuristic weight of cells that have been processed, i.e.
// that are no longer in consideration for processing.
type closeList map[string]int

func (cl *closeList) add(c cell) {
	cellID := c.String()
	(*cl)[cellID] = c.f
}

// Checks if the closed list already has a cell with the SAME
// coordinates as the candidate cell.  If it does, it returns
// the 'f' value of that cell.  Otherwise, it returns -1.
func (cl *closeList) has(c cell) int {
	candidateID := c.String()

	if existingF, ok := (*cl)[candidateID]; ok {
		return existingF
	}

	return -1
}

func isEligible(neighbor, current ICoord, g grid) bool {
	yMax := len(g) - 1
	xMax := len(g[0]) - 1

	// If candidate is outside the grid, it is ineligible.
	if neighbor.X() < 0 || neighbor.X() > xMax || neighbor.Y() < 0 || neighbor.Y() > yMax {
		return false
	}

	// If candidate is too far above current cell, it is ineligible
	currentHeight := g[current.Y()][current.X()]
	neighborHeight := g[neighbor.Y()][neighbor.X()]
	return neighborHeight-currentHeight <= 1
}

func calcDistance(input parsingResult, start coord) cell {
	// initialize stores
	openQueue := newOpenQ()
	closedList := make(closeList)

	sourceCell := cell{}
	sourceCell.coord = start

	// it costs no movement to get to the starting cell
	sourceCell.g = 0
	sourceCell.h = calcH(input.dest, sourceCell)
	sourceCell.f = sourceCell.g + sourceCell.h

	openQueue.push(sourceCell)

	// utility for checking adjacent cells
	dirs := directions{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	for openQueue.len() > 0 {
		currentCell := openQueue.pop()
		closedList.add(currentCell)

		// Get every viable neighbor
		for _, dir := range dirs {
			neighbX := currentCell.X() + dir.X()
			neighbY := currentCell.Y() + dir.Y()
			neighbor := coord{neighbX, neighbY}

			if isEligible(neighbor, currentCell, input.g) {
				neighborCell := newNeighbor(currentCell, neighbX, neighbY, input.dest)

				if neighborCell.String() == input.dest.String() {
					return neighborCell
				}

				// If there is already a cell for the same coordinate in the closed list,
				//	check its heuristic weight. If the neighbor we're currently considering
				// has a LOWER 'f' value than the cell that has already been closed, REMOVE the
				// cell from the closed list, since the current neighbor is higher-priority.
				closedPriority := closedList.has(neighborCell)
				if closedPriority != -1 {
					if closedPriority > neighborCell.f {
						delete(closedList, neighbor.String())
					} else {
						// if we get here, neighbor has already been processed, we don't need to do it again.
						continue
					}
				}

				// If there is already a cell for the same coordinate in the openQ,
				//	check its heuristic weight. If the neighbor we're currently considering
				// has a LOWER 'f' value than the cell that is already in openQ, REMOVE the
				// cell from openQ, since the current neighbor is higher-priority. Neighborcell
				// will take its place in the open queue at the end of this loop. Otherwise, we
				// don't need to look at this neighbor again.
				existingPriority := openQueue.has(neighborCell)
				if existingPriority != -1 {
					if existingPriority > neighborCell.f {
						openQueue.remove(neighborCell)
					}
				}

				openQueue.push(neighborCell)
			}
		}

		closedList.add(currentCell)
	}

	return cell{}
}

func calcForConcurrent(input parsingResult, start coord, output *[]cell, pos int) {
	res := calcDistance(input, start)

	(*output)[pos] = res
}

func intMin(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// Returns cell representing destination.
func navigate(input io.Reader) cell {
	parsed := parseInput(input)

	return calcDistance(parsed, parsed.source)
}

func navigateTwo(r io.Reader) int {
	parsed := parseInput(r)

	results := make([]cell, len(parsed.lowPoints))

	var wg sync.WaitGroup
	waitCount := len(parsed.lowPoints)
	for waitCount > 0 {
		wg.Add(1)
		waitCount--
	}

	for idx, startCoord := range parsed.lowPoints {
		// careful with closure...
		s := startCoord
		i := idx
		go func() {
			defer wg.Done()
			calcForConcurrent(parsed, s, &results, i)
		}()
	}

	wg.Wait()

	fmt.Println(results)
	var minDist int

	for _, c := range results {
		if c.g == 0 {
			continue
		}
		if minDist == 0 {
			minDist = c.g
		} else {
			minDist = intMin(minDist, c.g)
		}
	}

	return minDist
}
