package main

import "testing"

var testInput = `>>><<><>><<<>><>>><<<>>><<<><<<>><>><<>>`

func TestSolveOne(t *testing.T) {
	var want int64 = 3068
	got := solveOne(testInput)

	if got != want {
		t.Fatalf("expected %d, received %d", want, got)
	}
}

func TestSolveTwo(t *testing.T) {
	var want int64 = 1514285714288
	got := solveTwo(testInput, 1e12)

	if got != want {
		t.Fatalf("expected %d, received %d", want, got)
	}
}
