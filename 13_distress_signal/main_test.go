package main

import (
	"reflect"
	"strings"
	"testing"
)

type testCase struct {
	left   []any
	right  []any
	expect bool
}

type testCases []testCase

func runCases(tcs testCases, t *testing.T) {
	for _, tc := range tcs {
		got := isOrderedPair(tc.left, tc.right, 0)

		if got != tc.expect {
			t.Fatalf("left: %+v\nright: %+v\nexpected: %t\ngot: %t\n", tc.left, tc.right, tc.expect, got)
		}
	}
}

var case1 = testCase{
	[]any{1, 1, 3, 1, 1},
	[]any{1, 1, 5, 1, 1},
	true,
}

var case2 = testCase{
	[]any{[]any{1}, []any{2, 3, 4}},
	[]any{[]any{1}, 4},
	true,
}

var case3 = testCase{
	[]any{9},
	[]any{[]any{8, 7, 6}},
	false,
}

var case4 = testCase{
	[]any{[]any{4, 4}, 4, 4},
	[]any{[]any{4, 4}, 4, 4, 4},
	true,
}

var case5 = testCase{
	[]any{7, 7, 7, 7},
	[]any{7, 7, 7},
	false,
}

var case6 = testCase{
	make([]any, 0),
	[]any{3},
	true,
}

var case7 = testCase{
	[]any{[]any{[]any{}}},
	[]any{[]any{}},
	false,
}

func TestIsOrderedPair(t *testing.T) {
	cases := testCases{case1, case2, case3, case4, case5, case6, case7}

	runCases(cases, t)
}

func TestStringToSlice(t *testing.T) {
	given := "[1,[2,[3,[4,[5,6,7]]]],8,9]"

	want := []any{1, []any{2, []any{3, []any{4, []any{5, 6, 7}}}}, 8, 9}

	got, _ := stringToSlice([]rune(given), 0)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected %+v, but received %+v", want, got)
	}

	given = "[[1],[2,3,4]]"

	want = []any{[]any{1}, []any{2, 3, 4}}

	got, _ = stringToSlice([]rune(given), 0)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected %+v, but received %+v", want, got)
	}
}

var testInput = `[1,1,3,1,1]
[1,1,5,1,1]

[[1],[2,3,4]]
[[1],4]

[9]
[[8,7,6]]

[[4,4],4,4]
[[4,4],4,4,4]

[7,7,7,7]
[7,7,7]

[]
[3]

[[[]]]
[[]]

[1,[2,[3,[4,[5,6,7]]]],8,9]
[1,[2,[3,[4,[5,6,0]]]],8,9]`

func TestGetOrderedPairSum(t *testing.T) {
	want := 13

	got := getOrderedPairSum(strings.NewReader(testInput))

	if want != got {
		t.Fatalf("expected %d, received  %d", want, got)
	}
}
