package main

import "testing"

func TestGetDuplicate(t *testing.T) {
	given := []string{"ab", "ba", "ad"}

	want := 'a'

	got := getDuplicate(given)

	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestRunePriority(t *testing.T) {
	given := []rune{'a', 'z', 'A', 'Z'}
	want := []int{1, 26, 27, 52}

	for idx, rn := range given {

		got := runePriority(rn)

		if got != want[idx] {
			t.Fatalf("for input %v, expected %v, got %v", string(rn), want[idx], got)
		}
	}

}
