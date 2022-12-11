package main

import (
	"bufio"
	"strings"
	"testing"
)

var testInput string = `30373
25512
65332
33549
35390
`

func TestCountVisible(t *testing.T) {
	g := grid{}

	scanner := bufio.NewScanner(strings.NewReader(testInput))

	for scanner.Scan() {
		row := scanner.Text()

		rowToGrid(row, &g)
	}

	want := 21

	got := countVisible(g)

	if want != got {
		t.Fatalf("expected %d, but received %d", want, got)
	}
}

func TestCountScenicScore(t *testing.T) {
	g := grid{}

	scanner := bufio.NewScanner(strings.NewReader(testInput))

	for scanner.Scan() {
		row := scanner.Text()

		rowToGrid(row, &g)
	}

	want := 8

	got := countScenicScore(g)

	if want != got {
		t.Fatalf("expected %d, but received %d", want, got)
	}
}
