package stacks

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	given := "    [A] [B]"

	got := New(given)

	if len(*got) != 3 {
		t.Fatalf("expected 4 stacks, received %d", len(*got))
	}

	given = "[A]     [B]"

	got = New(given)

	if len(*got) != 3 {
		t.Fatalf("expected 4 stacks, received %d", len(*got))
	}
}

func TestStacksReverse(t *testing.T) {
	s := Stacks{{'A', 'B'}, {'C', 'D'}}

	want := [][]byte{{'B', 'A'}, {'D', 'C'}}

	s.Reverse()

	for idx, expected := range want {
		for resIdx, res := range expected {
			if s[idx][resIdx] != res {
				fmt.Println(s)
				t.Fatal("reverse failed")
			}
		}
	}
}

func TestAddRow(t *testing.T) {
	s := Stacks{{'A', 'B'}, {'C', 'D'}}

	given := []byte{'J', 'x'}

	want := [][]byte{{'A', 'B', 'J'}, {'C', 'D'}}

	s.AddRow(given)

	for outerIdx, wantedSlice := range want {
		for innerIdx, wanted := range wantedSlice {
			if s[outerIdx][innerIdx] != wanted {
				fmt.Println(s)
				t.Fatal("addrow failed")
			}
		}
	}
}

func TestMove(t *testing.T) {
	s := Stacks{{'A', 'B'}, {'C', 'D'}}

	s.Move(2, 1, 0)

	t.Logf("%v", s)

	want := [][]byte{{'A', 'B', 'D', 'C'}, {}}

	for outerIdx, wantedSlice := range want {
		for innerIdx, wanted := range wantedSlice {
			if s[outerIdx][innerIdx] != wanted {
				fmt.Println(s)
				t.Fatal("move failed")
			}
		}
	}
}

func TestMoveMultipl(t *testing.T) {
	s := Stacks{{'A', 'B'}, {'C', 'D'}}

	s.MoveMultiple(2, 1, 0)

	t.Logf("%v", s)

	want := [][]byte{{'A', 'B', 'C', 'D'}, {}}

	for outerIdx, wantedSlice := range want {
		for innerIdx, wanted := range wantedSlice {
			if s[outerIdx][innerIdx] != wanted {
				fmt.Println(s)
				t.Fatal("move failed")
			}
		}
	}
}
