package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("./input.txt")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(runPartOne(file))
}

type Point interface {
	X() int
	Y() int
}

type Vector []Point

type Rock struct {
	x int
	y int
}

func (s Rock) X() int {
	return s.x
}

func (s Rock) Y() int {
	return s.y
}

type Blocker interface {
	Block(x, y int)
	BlockV(v Vector)
	Blocked(x, y int) bool
	Inbounds(p Point) bool
	Floor() int
}

type Grid struct {
	blocked     map[string]struct{}
	leftBound   int
	rightBound  int
	bottomBound int
}

func (g Grid) Floor() int {
	return g.bottomBound + 2
}

func (g *Grid) Block(x, y int) {
	k := fmt.Sprintf("%d-%d", x, y)
	g.blocked[k] = struct{}{}
}

func (g *Grid) UpdateBounds(x, y int) {
	g.leftBound = intMin(x, g.leftBound)
	g.rightBound = intMax(x, g.rightBound)
	g.bottomBound = intMax(y, g.bottomBound)
}

func intMin(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func intMax(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// NB we only update the bounds when we're populating the initial 'rocks',
// otherwise the floor could infinitely recede!
func (g *Grid) BlockV(v Vector) {
	currPoint := v[0]
	g.Block(currPoint.X(), currPoint.Y())

	for i := 1; i < len(v); i++ {
		nextPoint := v[i]
		offSetX := currPoint.X() - nextPoint.X()
		offSetY := currPoint.Y() - nextPoint.Y()

		// displacement is vertical
		if offSetX == 0 {
			yMin := intMin(currPoint.Y(), nextPoint.Y())
			yMax := intMax(currPoint.Y(), nextPoint.Y())

			for i := yMin; i <= yMax; i++ {
				g.Block(nextPoint.X(), i)
				g.UpdateBounds(nextPoint.X(), i)
			}
		}

		// displacement is horizontal
		if offSetY == 0 {
			xMin := intMin(currPoint.X(), nextPoint.X())
			xMax := intMax(currPoint.X(), nextPoint.X())

			for i := xMin; i <= xMax; i++ {
				g.Block(i, nextPoint.Y())
				g.UpdateBounds(i, nextPoint.Y())
			}
		}

		currPoint = nextPoint
	}
}

func (g Grid) Blocked(x, y int) bool {
	if y == g.Floor() {
		return true
	}

	k := fmt.Sprintf("%d-%d", x, y)

	if _, ok := g.blocked[k]; ok {
		return true
	}

	return false
}

func (g Grid) Inbounds(p Point) bool {
	return p.Y() < g.bottomBound && p.X() > g.leftBound && p.X() < g.rightBound
}

func NewGrid() *Grid {
	blocked := make(map[string]struct{})
	return &Grid{
		blocked,
		500,
		500,
		0,
	}
}

type Sand struct {
	x int
	y int
}

func (s Sand) X() int {
	return s.x
}

func (s Sand) Y() int {
	return s.y
}

func (s Sand) String() string {
	return fmt.Sprintf("%d-%d", s.x, s.y)
}

// Fall moves the sand grain according to its preferred movement.
func (s *Sand) Fall(blocker Blocker) bool {
	dirs := [][]int{{0, 1}, {-1, 1}, {1, 1}}

	for _, dir := range dirs {
		nx := s.x + dir[0]
		ny := s.y + dir[1]

		if !blocker.Blocked(nx, ny) {
			s.x = nx
			s.y = ny
			return false
		}
	}

	return true
}

func NewSand() *Sand {
	return &Sand{500, 0}
}

func countStoppedSand(blocker Blocker) int {
	stoppedCount := 0

outer:
	for {
		grain := NewSand()

		for {
			stopped := grain.Fall(blocker)

			if stopped {
				stoppedCount++
				blocker.Block(grain.X(), grain.Y())
				if grain.X() == 500 && grain.Y() == 0 {
					break outer
				} else {
					break
				}
			}

			// This block only used in part 1
			// if !blocker.Inbounds(grain) {
			// 	break outer
			// }
		}
	}

	return stoppedCount
}

func lineToVector(line string) Vector {
	rexp := regexp.MustCompile(`[\d,]+`)

	stringCoords := rexp.FindAllString(line, -1)

	rocks := make(Vector, len(stringCoords))

	for i, c := range stringCoords {
		coords := strings.Split(c, ",")

		x, err := strconv.Atoi(coords[0])
		if err != nil {
			log.Fatalf("Unable to parse value %v to int", coords[0])
		}

		y, err := strconv.Atoi(coords[1])
		if err != nil {
			log.Fatalf("Unable to parse value %v to int", coords[1])
		}

		rocks[i] = Rock{x, y}
	}

	return rocks
}

func runPartOne(r io.Reader) int {
	scanner := bufio.NewScanner(r)

	grid := NewGrid()

	for scanner.Scan() {
		line := scanner.Text()
		v := lineToVector(line)
		grid.BlockV(v)
	}

	return countStoppedSand(grid)
}
