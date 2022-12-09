package stringtable

import (
	"testing"
)

func TestParseTableRow(t *testing.T) {
	given := "    [A] [B]"

	want := []byte{'x', 'A', 'B'}

	got := ParseTableRow(given)

	for idx, byte := range got {
		if byte != want[idx] {
			t.Fatalf("expected %v at position %v, but received %v", want[idx], idx, byte)
		}
	}

	given = "[A] [B] [C]"

	want = []byte{'A', 'B', 'C'}

	got = ParseTableRow(given)

	for idx, byte := range got {
		if byte != want[idx] {
			t.Fatalf("expected %v at position %v, but received %v", want[idx], idx, byte)
		}
	}
}
