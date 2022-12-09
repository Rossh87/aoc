package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.ReadFile("./input.txt")

	if err != nil {
		log.Fatal(err)
	}

	ans := findMarker(string(file), 14)

	fmt.Println(ans)
}

func allUnique(m map[rune]int) bool {
	for _, v := range m {
		if v > 1 {
			return false
		}
	}

	return true
}

type tracker struct {
	counts map[rune]int
	s      string
}

// add returns true if all runes including the added rune
// are unique
func (t *tracker) add(char rune) bool {
	t.s += string(char)

	if _, ok := t.counts[char]; ok {
		t.counts[char] += 1
		return false
	}

	t.counts[char] = 1

	return allUnique(t.counts)
}

func (t *tracker) shift() {
	head := t.s[0]

	t.s = t.s[1:]

	newCount := t.counts[rune(head)] - 1

	if newCount == 0 {
		delete(t.counts, rune(head))
		return
	}

	t.counts[rune(head)] = newCount
}

func findMarker(buf string, markerLength int) int {
	trk := tracker{make(map[rune]int), ""}

	rns := []rune(buf)

	for i := 0; i < markerLength-1; i++ {
		rn := rns[i]
		trk.add(rn)
	}

	for i := 3; i < len(rns); i++ {
		if trk.add(rns[i]) {
			// answer is expected to be 1-indexed
			return i + 1
		}

		trk.shift()
	}

	return -1
}
