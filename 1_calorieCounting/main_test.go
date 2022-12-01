package main

import (
	"fmt"
	"strings"
	"testing"
)

type testCases struct {
	inputs   []input
	expected []int
}

type input struct {
	s        string
	aggCount int
}

var t1 = `
1

1

2

3
`

var t2 = `

1

2


`

func TestMaxCals(t *testing.T) {
	inputs := []input{{t1, 2}, {t2, 1}}
	expected := []int{5, 2}
	tcs := testCases{inputs, expected}

	for idx, input := range tcs.inputs {
		result, err := maxCals(strings.NewReader(input.s), input.aggCount)

		if err != nil {
			t.Fatalf("%v", err)
		}

		if result != tcs.expected[idx] {
			t.Fatalf(fmt.Sprintf("expected %d, got %d", tcs.expected[idx], result))
		}
	}
}
