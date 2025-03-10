package main

import (
	"strings"
	"testing"
)

var testInput = `        ...#
        .#..
        #...
        ....
...#.......#
........#...
..#....#....
..........#.
        ...#....
        .....#..
        .#......
        ......#.

10R5L5R10L4R5L5`

var testGiven = `  ...#
  #...

1R1L10`

// checked:
// left overflow blocked
func TestBlockedSwitch(t *testing.T) {
	got := solvePartOne(strings.NewReader(testGiven))
	want := 10

	if got != want {
		t.Fatalf("expected %d, received %d", want, got)
	}
}

func TestSolvePartOne(t *testing.T) {
	want := 6032
	got := solvePartOne(strings.NewReader(testInput))
	if got != want {
		t.Fatalf("expected %d, received %d", want, got)
	}
}
