package main

import "testing"

var testInput = `>>><<><>><<<>><>>><<<>>><<<><<<>><>><<>>`

func TestSolveOne(t *testing.T) {
	want := 3068
	got := solveOne(testInput)

	if got != want {
		t.Fatalf("expected %d, received %d", want, got)
	}
}
