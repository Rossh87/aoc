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

func TestParseTask(t *testing.T) {

}

func TestSolvePartOne(t *testing.T) {
	want := 152
	got := solvePartOne(strings.NewReader(testInput))
	if got != want {
		t.Fatalf("expected %d, got %d", want, got)
	}
}
