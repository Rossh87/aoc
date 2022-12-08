package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCovers(t *testing.T) {
	pairA := pair{3, 5}
	pairB := pair{3, 4}

	got := pairA.covers(pairB)

	if !got {
		t.Fatalf("expected %v, got %v", true, got)
	}

	got = pairB.covers(pairA)

	if got {
		t.Fatalf("expected %v, got %v", false, got)
	}
}

func TestOverlaps(t *testing.T) {
	pairA := pair{2, 4}
	pairB := pair{3, 5}

	got := pairA.overlaps(pairB)

	if !got {
		t.Fatalf("expected %v, got %v", true, got)
	}

	got = pairB.overlaps(pairA)

	if !got {
		t.Fatalf("expected %v, got %v", false, got)
	}
}

func TestGetPair(t *testing.T) {
	given := "2-3"
	want := pair{2, 3}
	got, err := getPair(given)

	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(got, want) {
		t.Fatalf("expected %v, got %v", false, got)
	}
}
