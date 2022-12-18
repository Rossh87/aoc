package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestLineToVector(t *testing.T) {
	given := "498,4 -> 498,6 -> 496,6"
	want := Vector{Rock{498, 4}, Rock{498, 6}, Rock{496, 6}}
	got := lineToVector(given)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected %v, received %v", want, got)
	}
}

var testInput = `498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9`

func TestRun(t *testing.T) {
	want := 93

	got := runPartOne(strings.NewReader(testInput))

	if got != want {
		t.Fatalf("expected %v, received %v", want, got)
	}
}
