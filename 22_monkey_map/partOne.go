package main

import (
	"bufio"
	"io"
)

type row struct {
	id        int
	tileStart int
	tileEnd   int
	walls     map[int]struct{}
}

func parseRow(rawRow string, rowNum int) row {
	countingTiles := false
	ws := make(map[int]struct{})
	ts := 0
	te := 0

	// all ASCII chars, so byte offset should correspond to character offset
	for i, c := range rawRow {
		switch c {
		case ' ':
			if countingTiles {
				countingTiles = false
				// puzzle is 1-indexed
				te = i
				continue
			}
		case '.':
			if !countingTiles {
				countingTiles = true
				ts = i + 1
			}
		case '#':
			if !countingTiles {
				countingTiles = true
				ts = i + 1
			}
			ws[i+1] = struct{}{}
		}
	}

	// if we finished iterating over row characters without hitting any right padding,
	// we still need to set the rightmost bound of the tiles
	if countingTiles {
		te = len(rawRow)
	}

	return row{
		rowNum,
		ts,
		te,
		ws,
	}
}

type board []row

func (b board) isWall(row, col int) bool {
	currRow := b[row-1]
	if _, exists := currRow.walls[col]; exists {
		return true
	}
	return false
}

func (b board) row(row int) row {
	return b[row-1]
}

func (b board) sectionCeil(row, col int) int {
	for row > 1 {
		nextRow := b.row(row - 1)

		if col >= nextRow.tileStart && col <= nextRow.tileEnd {
			row--
		} else {
			break
		}
	}

	return row
}

func (b board) sectionFloor(row, col int) int {
	for row < len(b) {
		nextRow := b.row(row + 1)

		if col >= nextRow.tileStart && col <= nextRow.tileEnd {
			row++
		} else {
			break
		}
	}

	return row
}

type direction int

const (
	right direction = 0
	down  direction = 1
	left  direction = 2
	up    direction = 3
)

type piece struct {
	row  int
	col  int
	face direction
}

func (p *piece) turn(dir rune) {
	switch dir {
	case 'R':
		p.face = (p.face + 1) % 4
	case 'L':
		next := p.face - 1
		if next < 0 {
			next += 4
		}
		p.face = next
	default:
		panic("piece.turn called with unknown argument")
	}
}

func (p *piece) advance(moves int, b board) {
	//   TODO: only run this on vertical moves
	ceil := b.sectionCeil(p.row, p.col)
	floor := b.sectionFloor(p.row, p.col)
Loop:
	for moves > 0 {
		switch p.face {
		case right:
			currRow := b.row(p.row)
			nextCol := p.col + 1
			// wrap around if necessary
			if nextCol > currRow.tileEnd {
				nextCol = currRow.tileStart
			}
			if b.isWall(p.row, nextCol) {
				break Loop
			}
			moves--
			p.col = nextCol

		case left:
			currRow := b.row(p.row)
			nextCol := p.col - 1
			// wrap around if necessary
			if nextCol < currRow.tileStart {
				nextCol = currRow.tileEnd
			}
			if b.isWall(p.row, nextCol) {
				break Loop
			}
			moves--
			p.col = nextCol

		case down:
			nextRow := p.row + 1
			// wrap around if necessary
			if nextRow > len(b) || nextRow > floor {
				nextRow = ceil
			}
			if b.isWall(nextRow, p.col) {
				break Loop
			}
			moves--
			p.row = nextRow

		case up:
			nextRow := p.row - 1
			// wrap around if necessary
			if nextRow < 1 || nextRow < ceil {
				nextRow = floor
			}
			if b.isWall(nextRow, p.col) {
				break Loop
			}
			moves--
			p.row = nextRow
		default:
			panic("unknown piece direction encountered while attempting to advance")
		}
	}
}

func isIntRune(rn rune) bool {
	code := rn - '0'

	return code >= 0 && code < 10
}

func runeToInt(rn rune) int {
	return int(rn - '0')
}

type instructions string

func (ins instructions) run(b board) piece {
	p := piece{
		1,
		b.row(1).tileStart,
		right,
	}

	moves := 0
	for _, rn := range ins {
		if isIntRune(rn) {
			moves *= 10
			nmoves := moves + runeToInt(rn)
			moves = nmoves
		} else {
			p.advance(moves, b)
			moves = 0
			p.turn(rn)
		}
	}
	//   make final moves
	p.advance(moves, b)
	return p
}

func calcAnswer(p piece) int {
	return (p.row * 1000) + (p.col * 4) + int(p.face)
}

func solvePartOne(r io.Reader) int {
	scanner := bufio.NewScanner(r)
	var ins instructions
	board := make(board, 0)
	boardDone := false
	lineNo := 1
	for scanner.Scan() {
		ln := scanner.Text()

		if boardDone {
			ins = instructions(ln)
			break
		}

		if ln == "" {
			boardDone = true
			continue
		}

		board = append(board, parseRow(ln, lineNo))
		lineNo++
	}

	p := ins.run(board)

	return calcAnswer(p)
}
