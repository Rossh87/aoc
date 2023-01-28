package main

import (
	"strings"
	"testing"
)

var testInput = `1
2
-3
3
-2
0
4`

func TestSolvePartOne(t *testing.T) {
	want := 3
	got := solvePartOne(strings.NewReader(testInput))
	if want != got {
		t.Fatalf("expected %d, got %d", want, got)
	}
}

func TestSolvePartTwo(t *testing.T) {
	want := 1623178306
	got := solvePartTwo(strings.NewReader(testInput))
	if want != got {
		t.Fatalf("expected %d, got %d", want, got)
	}
}
