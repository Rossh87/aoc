package main

import (
	"strings"
	"testing"
)

var testInput = `Sabqponm
abcryxxl
accszExk
acctuvwj
abdefghi`

func TestNavigate(t *testing.T) {
	want := 31
	got := navigate(strings.NewReader(testInput)).g

	if want != got {
		t.Fatalf("expected %d, but got %d", want, got)
	}
}
