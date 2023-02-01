package main

import (
	"strings"
	"testing"
)

var testInput = `root: pppw + sjmn
dbpl: 5
cczh: sllz + lgvd
zczc: 2
ptdq: humn - dvpt
dvpt: 3
lfqf: 4
humn: 5
ljgn: 2
sjmn: drzm * dbpl
sllz: 4
pppw: cczh / lfqf
lgvd: ljgn * ptdq
drzm: hmdt - zczc
hmdt: 32`

func TestSolvePartOne(t *testing.T) {
	want := 152
	got := solvePartOne(strings.NewReader(testInput))
	if got != want {
		t.Fatalf("expected %d, got %d", want, got)
	}
}

func TestFindEqualityTarget(t *testing.T) {
	s := make(satisfactionMap)
	ts := parseInput(strings.NewReader(testInput), &s)
	want := 150
	got := findEqualityTarget(ts, &s)
	if got != want {
		t.Fatalf("expected %d, got %d", want, got)
	}
}

func TestUnwindStack(t *testing.T) {
	stack := opStack{
		{'+', 3},
		{'/', 2},
		{'-', 4},
		{'*', 4},
	}
	want := 301
	got := unwindStack(150, stack)
	if got != want {
		t.Fatalf("expected %d, got %d", want, got)
	}
}

func TestSolvePartTwo(t *testing.T) {
	want := 301
	got := solvePartTwo(strings.NewReader(testInput))
	if got != want {
		t.Fatalf("expected %d, got %d", want, got)
	}
}
