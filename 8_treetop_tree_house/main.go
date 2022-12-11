package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, _ := os.Open("./input.txt")
	scanner := bufio.NewScanner(file)
	g := grid{}

	for scanner.Scan() {
		row := scanner.Text()

		rowToGrid(row, &g)
	}

	fmt.Println(countScenicScore(g))
}

type grid [][]int

type maxEntry struct {
	topMax    int
	leftMax   int
	bottomMax int
	rightMax  int
}

type dp [][]maxEntry

func rowToGrid(s string, g *grid) {
	newRow := []int{}

	for _, rn := range s {
		v := int(rn - '0')

		newRow = append(newRow, v)
	}

	*g = append(*g, newRow)
}

func intMax(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// leftMax[row][col] = largest value that exists TO THE LEFT OF grid[row][col]
func setLeftMax(row, col int, g grid, idx *dp) {
	if col == 0 {
		(*idx)[row][col].leftMax = -1
		return
	}

	lNeighbor := (*idx)[row][col-1]

	lNeighborVal := g[row][col-1]

	(*idx)[row][col].leftMax = intMax(lNeighborVal, lNeighbor.leftMax)
}

// leftMax[row][col] = largest value that exists ABOVE grid[row][col]
func setTopMax(row, col int, g grid, idx *dp) {
	if row == 0 {
		(*idx)[row][col].topMax = -1
		return
	}

	tNeighbor := (*idx)[row-1][col]

	tNeighborVal := g[row-1][col]

	(*idx)[row][col].topMax = intMax(tNeighborVal, tNeighbor.topMax)
}

// leftMax[row][col] = largest value that exists ABOVE grid[row][col]
func setRightMax(row, col int, g grid, idx *dp) {
	cols := len(g[0])

	if col == cols-1 {
		(*idx)[row][col].rightMax = -1
		return
	}

	rNeighbor := (*idx)[row][col+1]

	rNeighborVal := g[row][col+1]

	(*idx)[row][col].rightMax = intMax(rNeighborVal, rNeighbor.rightMax)
}

// leftMax[row][col] = largest value that exists ABOVE grid[row][col]
func setBottomMax(row, col int, g grid, idx *dp) {
	rows := len(g)

	if row == rows-1 {
		(*idx)[row][col].bottomMax = -1
		return
	}

	bNeighbor := (*idx)[row+1][col]

	bNeighborVal := g[row+1][col]

	(*idx)[row][col].bottomMax = intMax(bNeighborVal, bNeighbor.bottomMax)
}

func countVisible(grid [][]int) int {
	rows := len(grid)
	cols := len(grid[0])

	// populate idx
	idx := make(dp, rows)

	for i := 0; i < rows; i++ {
		entries := make([]maxEntry, cols)

		for j := 0; j < cols; j++ {
			entries[j] = maxEntry{}
		}

		idx[i] = entries
	}

	// set left and top maxes
	for row, rowSlice := range grid {
		for col := range rowSlice {
			setTopMax(row, col, grid, &idx)
			setLeftMax(row, col, grid, &idx)
		}
	}

	// set right and bottom maxes
	for row := rows - 1; row >= 0; row-- {
		for col := cols - 1; col >= 0; col-- {
			setBottomMax(row, col, grid, &idx)
			setRightMax(row, col, grid, &idx)
		}
	}

	visible := 0

	for row, rowSlice := range grid {
		for col, height := range rowSlice {
			maxes := idx[row][col]

			if maxes.bottomMax < height || maxes.topMax < height || maxes.rightMax < height || maxes.leftMax < height {
				visible++
			}
		}
	}

	return visible
}

func oneGridScore(row, col int, g grid) int {
	rows := len(g)
	cols := len(g[0])
	currHeight := g[row][col]

	lTrees := 0
	rTrees := 0
	tTrees := 0
	bTrees := 0

	// top
	ptr := row - 1

	for ptr >= 0 {
		tTrees++

		if g[ptr][col] >= currHeight {
			break
		}

		ptr--
	}

	// bottom
	ptr = row + 1

	for ptr <= rows-1 {
		bTrees++

		if g[ptr][col] >= currHeight {
			break
		}

		ptr++
	}

	// left
	ptr = col - 1

	for ptr >= 0 {
		lTrees++

		if g[row][ptr] >= currHeight {
			break
		}

		ptr--
	}

	// right
	ptr = col + 1

	for ptr <= cols-1 {
		rTrees++

		if g[row][ptr] >= currHeight {
			break
		}

		ptr++
	}

	return lTrees * rTrees * bTrees * tTrees
}

func countScenicScore(g grid) int {
	maxScore := 0

	for row, rowSlice := range g {
		for col := range rowSlice {
			maxScore = intMax(maxScore, oneGridScore(row, col, g))
		}
	}

	return maxScore
}
