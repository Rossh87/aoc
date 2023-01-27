package main

import "testing"

// var testInput = `Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 14 clay. Each geode robot costs 2 ore and 7 obsidian.
// Blueprint 2: Each ore robot costs 2 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 8 clay. Each geode robot costs 3 ore and 12 obsidian.`

func TestParseBlueprint(t *testing.T) {
	given := "Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 14 clay. Each geode robot costs 2 ore and 7 obsidian."
	want := blueprint{
		1,
		oreCost{4},
		clayCost{2},
		obsidianCost{3, 14},
		geodeCost{2, 7},
	}
	got := parseBlueprint(given)
	if got != want {
		t.Fatalf("expected %d, received %d", want, got)
	}
}

func TestGetMaxGeodes(t *testing.T) {
	given := "Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 14 clay. Each geode robot costs 2 ore and 7 obsidian."
	want := 9
	got := getMaxGeodes(24, parseBlueprint(given))
	if got != want {
		t.Fatalf("expected %d, received %d", want, got)
	}
}
func TestSolve(t *testing.T) {
	want := 33
	got := solve()
	if got != want {
		t.Fatalf("expected %d, received %d", want, got)
	}
}
