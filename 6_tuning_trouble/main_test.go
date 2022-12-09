package main

import "testing"

func TestFindMarker(t *testing.T) {
	given := "bvwbjplbgvbhsrlpgdmjqwftvncz"
	want := 5
	got := findMarker(given, 4)

	if got != want {
		t.Fatalf("expected %d, received %d", want, got)
	}

	given = "nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg"
	want = 10
	got = findMarker(given, 4)

	if got != want {
		t.Fatalf("expected %d, received %d", want, got)
	}

	given = "zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw"
	want = 11
	got = findMarker(given, 4)

	if got != want {
		t.Fatalf("expected %d, received %d", want, got)
	}

	given = "mjqjpqmgbljsphdztnvjfqwrcgsmlb"
	want = 19
	got = findMarker(given, 14)

	if got != want {
		t.Fatalf("expected %d, received %d", want, got)
	}

	given = "nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg"
	want = 29
	got = findMarker(given, 14)

	if got != want {
		t.Fatalf("expected %d, received %d", want, got)
	}
}
