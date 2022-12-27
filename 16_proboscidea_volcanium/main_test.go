package main

import (
	"strings"
	"testing"
)

func neighborsEqual(a, b []string) bool {
	for i, el := range b {
		if a[i] != el {
			return false
		}
	}

	return true
}

func valveEqual(a, b valve) bool {
	return a.name == b.name &&
		a.flow == b.flow && a.open == b.open &&
		neighborsEqual(a.neighbors, b.neighbors)
}

func TestParseLine(t *testing.T) {
	ln := "Valve WT has flow rate=10; tunnels lead to valves BD, FQ"
	want := valve{
		"WT",
		10,
		false,
		[]string{"BD", "FQ"},
	}
	got := parseLine(ln)

	if !valveEqual(want, got) {
		t.Fatalf("expected %+v, received %+v", want, got)
	}
}

var testInput = `Valve AA has flow rate=0; tunnels lead to valves DD, II, BB
Valve BB has flow rate=13; tunnels lead to valves CC, AA
Valve CC has flow rate=2; tunnels lead to valves DD, BB
Valve DD has flow rate=20; tunnels lead to valves CC, AA, EE
Valve EE has flow rate=3; tunnels lead to valves FF, DD
Valve FF has flow rate=0; tunnels lead to valves EE, GG
Valve GG has flow rate=0; tunnels lead to valves FF, HH
Valve HH has flow rate=22; tunnel leads to valve GG
Valve II has flow rate=0; tunnels lead to valves AA, JJ
Valve JJ has flow rate=21; tunnel leads to valve II`

func TestMaxFlow(t *testing.T) {
	want := 1651
	got := maxFlow(strings.NewReader(testInput))
	if got != want {
		t.Fatalf("expected %+v, received %+v", want, got)
	}
}

func TestMaxFlowPartTwo(t *testing.T) {
	want := 1707
	got := maxFlowTwo(strings.NewReader(testInput))
	if got != want {
		t.Fatalf("expected %+v, received %+v", want, got)
	}
}
