package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	r, err := os.Open("./input.txt")

	if err != nil {
		panic(err)
	}

	// fmt.Println(countArea(r))
	fmt.Println(countAreaPartTwo(r))
}

type delta [3]int
type deltas [6]delta
type coord string
type coords map[coord]struct{}
type negativeCoords map[coord]bool
type bounds struct {
	minX int
	minY int
	minZ int
	maxX int
	maxY int
	maxZ int
}

func newBounds() *bounds {
	b := bounds{
		-1,
		-1,
		-1,
		-1,
		-1,
		-1,
	}
	return &b
}

func (b *bounds) update(c coord) {
	vs := c.Slice()
	for i, v := range vs {
		switch i {
		case 0:
			if v < b.minX || b.minX < 0 {
				b.minX = v
			}

			if v > b.maxX || b.maxX < 0 {
				b.maxX = v
			}
		case 1:
			if v < b.minY || b.minY < 0 {
				b.minY = v
			}

			if v > b.maxY || b.maxY < 0 {
				b.maxY = v
			}
		case 2:
			if v < b.minZ || b.minZ < 0 {
				b.minZ = v
			}

			if v > b.maxZ || b.maxZ < 0 {
				b.maxZ = v
			}
		default:
			panic("bounds received coordinate with too many values")
		}
	}
}

func (b bounds) outOfBounds(c coord) bool {
	for i, v := range c.Slice() {
		switch i {
		case 0:
			if v <= b.minX || v >= b.maxX {
				return true
			}
		case 1:
			if v <= b.minY || v >= b.maxY {
				return true
			}
		case 2:
			if v <= b.minZ || v >= b.maxZ {
				return true
			}
		default:
			panic("bounds received coordinate with too many values")
		}
	}

	return false
}

var ds deltas = [6]delta{{0, 0, 1}, {0, 0, -1}, {0, 1, 0}, {0, -1, 0}, {1, 0, 0}, {-1, 0, 0}}

func (c coord) String() string {
	return string(c)
}

func (c coord) Slice() []int {
	out := make([]int, 3)
	vs := strings.Split(c.String(), ",")
	for i, stringVal := range vs {
		intVal, err := strconv.ParseInt(stringVal, 10, 32)
		if err != nil {
			panic(err)
		}
		out[i] = int(intVal)
	}
	return out
}

func (c coord) AddDelta(d delta) coord {
	a := strings.Split(c.String(), ",")
	var res = [3]int{}
	for i, s := range a {
		v, err := strconv.ParseInt(s, 0, 32)

		if err != nil {
			panic(err)
		}

		res[i] = int(v)
	}

	for i, v := range res {
		res[i] = v + d[i]
	}

	return coord(fmt.Sprintf("%d,%d,%d", res[0], res[1], res[2]))
}

func countExposedSides(current coord, ds deltas, all coords) int {
	res := 0

	for _, d := range ds {
		toCheck := current.AddDelta(d)

		// if cell adjacent to current face is occupied, the current
		// face is not exposed
		if _, exists := all[toCheck]; exists {
			continue
		}

		res++
	}

	return res
}

func countArea(r io.Reader) int {
	allCoords := make(coords)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		ln := scanner.Text()
		allCoords[coord(ln)] = struct{}{}
	}

	res := 0

	for c := range allCoords {
		res += countExposedSides(c, ds, allCoords)
	}

	return res
}

// get list of all negative coords that cause an edge to be counted as surface area.
// calc min/max bounds for x, y, and z coords
// check each negative coord to see if it belongs to an internal negative space:
// grab a coord and BFS, adding all adjacent negative spaces to q. With each popped
// coord, check if it is outside the bounds of grid.  If so, the negative space can 'escape',
// and we can return early.  Otherwise, track popped coord so we don't re-add it to the q.
// If BFS returns true, remove all coords in the negative space from the list that can trigger
// a surface area addition.
// Finally, walk all input coords again, using logic from part 1 with the updated list of
// 'triggering' spaces.
func getNegativeCoords(occupied coords, ds deltas) negativeCoords {
	res := make(negativeCoords)

	for c := range occupied {
		for _, d := range ds {
			toCheck := c.AddDelta(d)
			_, isOccupied := occupied[toCheck]
			_, seen := res[toCheck]
			if !isOccupied && !seen {
				res[toCheck] = true
			}
		}
	}

	return res
}

func isEnclosed(start coord, checked coords, occupied coords, b bounds, ds deltas) (bool, coords) {
	enclosed := true
	// tracks cells in current negative space
	currentVoid := make(coords)
	q := []coord{start}
	currentVoid[start] = struct{}{}
	checked[start] = struct{}{}

	for len(q) > 0 {
		current := q[0]
		q = q[1:]
		for _, d := range ds {
			toCheck := current.AddDelta(d)

			// if the space we're checking is occupied by lava, don't push it
			if _, isOccupied := occupied[toCheck]; isOccupied {
				continue
			}

			// if we've already checked space during this BFS, don't push it
			if _, seen := currentVoid[toCheck]; seen {
				continue
			}

			// If space to check is empty and at outer bounds of shape, the negative
			// space it belongs to 'escapes', i.e. is not fully enclosed.
			if b.outOfBounds(toCheck) {
				enclosed = false
				continue
			}

			currentVoid[toCheck] = struct{}{}
			checked[toCheck] = struct{}{}
			q = append(q, toCheck)
		}
	}

	return enclosed, currentVoid
}

func countAreaFromNegatives(occupied coords, outerNegatives negativeCoords, ds deltas) int {
	count := 0
	for c := range occupied {
		for _, d := range ds {
			toCheck := c.AddDelta(d)

			if shouldCount, exists := outerNegatives[toCheck]; shouldCount && exists {
				count++
			}
		}
	}
	return count
}

func countAreaPartTwo(r io.Reader) int {
	occupied := make(coords)
	scanner := bufio.NewScanner(r)
	b := newBounds()

	for scanner.Scan() {
		c := coord(scanner.Text())
		occupied[c] = struct{}{}
		b.update(c)
	}

	allNegatives := getNegativeCoords(occupied, ds)

	checked := make(coords)

	for coord, stillValid := range allNegatives {
		if !stillValid {
			continue
		}

		if _, seen := checked[coord]; seen {
			continue
		}

		if shouldUpdate, toInvalidate := isEnclosed(coord, checked, occupied, *b, ds); shouldUpdate {
			for coordToInvalidate := range toInvalidate {
				if _, exists := allNegatives[coordToInvalidate]; exists {
					allNegatives[coordToInvalidate] = false
				}
			}
		}
	}

	return countAreaFromNegatives(occupied, allNegatives, ds)
}
