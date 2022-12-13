package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	bytes, err := os.ReadFile("./input.txt")

	if err != nil {
		log.Fatal(err)
	}

	input := string(bytes)

	fmt.Println(coveredByLongTail(input))
}

type point struct {
	x int
	y int
}

func (p point) string() string {
	return fmt.Sprint(p.x) + "-" + fmt.Sprint(p.y)
}

func intAbs(x int) int {
	return int(math.Abs(float64(x)))
}

func tailPos(h point, t point) point {
	dx := intAbs(h.x - t.x)
	dy := intAbs(h.y - t.y)

	// not enough distance to provoke tail move
	if dx < 2 && dy < 2 {
		return t
	}

	// direct horizontal move
	if dy == 0 {
		if h.x < t.x {
			t.x--
			return t
		} else {
			t.x++
		}
		return t
	}

	// direct vertical move
	if dx == 0 {
		if h.y < t.y {
			t.y--
		} else {
			t.y++
		}
		return t
	}

	// horizontal move: close both x and y dist by one
	if h.y < t.y {
		t.y--
	} else {
		t.y++
	}

	if h.x < t.x {
		t.x--
		return t
	} else {
		t.x++
	}

	return t
}

type visited map[string]int

func (v *visited) visit(p point) {
	if _, ok := (*v)[p.string()]; ok {
		(*v)[p.string()]++
		return
	}

	(*v)[p.string()] = 1
}

func (v visited) uniqueVisited() int {
	keyCount := 0

	for range v {
		keyCount++
	}

	return keyCount
}

func coveredByTail(s string) int {
	knot := newKnot()

	v := make(visited)

	v[knot.tail().string()] = 1

	scanner := bufio.NewScanner(strings.NewReader(s))

	for scanner.Scan() {
		line := scanner.Text()

		f := strings.Fields(line)

		moveCount, err := strconv.ParseInt(f[1], 10, 32)

		if err != nil {
			log.Fatal(err)
		}

		switch f[0] {
		case "U":
			for moveCount > 0 {
				t := knot.up()
				v.visit(t)
				moveCount--
			}
		case "D":
			for moveCount > 0 {
				t := knot.down()
				v.visit(t)
				moveCount--
			}
		case "L":
			for moveCount > 0 {
				t := knot.left()
				v.visit(t)
				moveCount--
			}
		case "R":
			for moveCount > 0 {
				t := knot.right()
				v.visit(t)
				moveCount--
			}
		}
	}

	return v.uniqueVisited()
}

func coveredByLongTail(s string) int {
	knot := newLongKnot()

	v := make(visited)

	v[knot.tail().string()] = 1

	scanner := bufio.NewScanner(strings.NewReader(s))

	for scanner.Scan() {
		line := scanner.Text()

		f := strings.Fields(line)

		moveCount, err := strconv.ParseInt(f[1], 10, 32)

		if err != nil {
			log.Fatal(err)
		}

		switch f[0] {
		case "U":
			for moveCount > 0 {
				t := knot.up()
				v.visit(t)
				moveCount--
			}
		case "D":
			for moveCount > 0 {
				t := knot.down()
				v.visit(t)
				moveCount--
			}
		case "L":
			for moveCount > 0 {
				t := knot.left()
				v.visit(t)
				moveCount--
			}
		case "R":
			for moveCount > 0 {
				t := knot.right()
				v.visit(t)
				moveCount--
			}
		}
	}

	return v.uniqueVisited()
}
