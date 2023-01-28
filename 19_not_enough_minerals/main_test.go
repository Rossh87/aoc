package main

import (
	"strings"
	"testing"
)

var testInput = `Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 14 clay. Each geode robot costs 2 ore and 7 obsidian.
Blueprint 2: Each ore robot costs 2 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 8 clay. Each geode robot costs 3 ore and 12 obsidian.`

func TestGetMaxGeodes(t *testing.T) {
	given := "Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 14 clay. Each geode robot costs 2 ore and 7 obsidian."
	want := 9
	got := getMaxGeodes(24, parseBlueprint(given))
	if got != want {
		t.Fatalf("expected %d, received %d", want, got)
	}
}

func TestGetMaxGeodesTwo(t *testing.T) {
	given := "Blueprint 2: Each ore robot costs 2 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 8 clay. Each geode robot costs 3 ore and 12 obsidian."
	want := 62
	got := getMaxGeodes(32, parseBlueprint(given))
	if got != want {
		t.Fatalf("expected %d, received %d", want, got)
	}
}
func TestSolve(t *testing.T) {
	want := 33
	got := solvePartOne(strings.NewReader(testInput))
	if got != want {
		t.Fatalf("expected %d, received %d", want, got)
	}
}
