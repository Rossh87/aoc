package main

import (
	"strings"
	"testing"
)

func sensorsEqual(a, b sensor) bool {
	return a.x == b.x && a.y == b.y && a.maxDist == b.maxDist
}

func TestParseLine(t *testing.T) {
	given := "Sensor at x=2, y=18: closest beacon is at x=-2, y=15"
	want := newSensor([]int{2, 18, -2, 15})
	got := parseLine(given)

	if !sensorsEqual(*want, *got) {
		t.Fatalf("expected %v, received %v", want, got)
	}
}

var testInput = `Sensor at x=2, y=18: closest beacon is at x=-2, y=15
Sensor at x=9, y=16: closest beacon is at x=10, y=16
Sensor at x=13, y=2: closest beacon is at x=15, y=3
Sensor at x=12, y=14: closest beacon is at x=10, y=16
Sensor at x=10, y=20: closest beacon is at x=10, y=16
Sensor at x=14, y=17: closest beacon is at x=10, y=16
Sensor at x=8, y=7: closest beacon is at x=2, y=10
Sensor at x=2, y=0: closest beacon is at x=2, y=10
Sensor at x=0, y=11: closest beacon is at x=2, y=10
Sensor at x=20, y=14: closest beacon is at x=25, y=17
Sensor at x=17, y=20: closest beacon is at x=21, y=22
Sensor at x=16, y=7: closest beacon is at x=15, y=3
Sensor at x=14, y=3: closest beacon is at x=15, y=3
Sensor at x=20, y=1: closest beacon is at x=15, y=3`

func TestCountSensedInRow(t *testing.T) {
	want := int64(26)

	got := countSensedInRow(10, strings.NewReader(testInput))

	if got != want {
		t.Fatalf("expected %d, received %d", want, got)
	}
}

func TestFindDistressFreq(t *testing.T) {
	want := int64(56000011)

	got := findDistressFreq(strings.NewReader(testInput), 0, 20)

	if got != want {
		t.Fatalf("expected %d, received %d", want, got)
	}
}
