package main

import (
	"strings"
	"testing"
)

var testInput = `2,2,2
1,2,2
3,2,2
2,1,2
2,3,2
2,2,1
2,2,3
2,2,4
2,2,6
1,2,5
3,2,5
2,1,5
2,3,5`

func TestCountArea(t *testing.T) {
	want := 64
	got := countArea(strings.NewReader(testInput))
	if got != want {
		t.Fatalf("expected %d, but received %d", want, got)
	}
}

func TestCountAreaPartTwo(t *testing.T) {
	want := 58
	got := countAreaPartTwo(strings.NewReader(testInput))
	if got != want {
		t.Fatalf("expected %d, but received %d", want, got)
	}
}
