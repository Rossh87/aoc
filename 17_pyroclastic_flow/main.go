package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rossh87/aoc/17_pyroclastic_flow/chamber"
	"github.com/rossh87/aoc/17_pyroclastic_flow/jets"
	"github.com/rossh87/aoc/17_pyroclastic_flow/rocks"
)

func main() {
	bytes, err := os.ReadFile("./input.txt")

	if err != nil {
		log.Fatal(err)
	}

	ans := solveTwo(string(bytes), 1e12)
	fmt.Println(ans)
}

func hash(nextJet, nextShape int, cont chamber.Contour) string {
	str := cont.Normalized().String()
	return fmt.Sprintf("%d-%d-%s", nextJet, nextShape, str)
}

type cacheState struct {
	settled int64
	height  int64
}

type cache map[string]cacheState

func solveOne(rawDirs string) int64 {
	js := jets.New(rawDirs)
	rocks := rocks.New()
	chamber := chamber.New(7)
	settled := 0

	for settled < 2022 {
		chamber.AddRock(rocks, js)
		settled++
	}

	return chamber.PileHeight()
}

func solveTwo(rawDirs string, target int64) int64 {
	js := jets.New(rawDirs)
	rocks := rocks.New()
	chamber := chamber.New(7)
	cache := make(cache)
	var settled int64 = 0
	var heightFromSkippedCycles int64 = 0

	for settled < target {
		chamber.AddRock(rocks, js)
		settled++
		hashKey := hash(js.Peek(), rocks.Peek(), chamber.Contour())
		cached, exists := cache[hashKey]

		// if we've fast-forwarded once, we don't want to
		// do it again
		if heightFromSkippedCycles > 0 {
			continue
		}

		// if we've never seen this state before, add it to cache and continue
		if !exists {
			cache[hashKey] = cacheState{settled, chamber.PileHeight()}
			continue
		}

		// if we've found a cycle, fast-forward through as many cycles
		// as possible, then let loop finish the remainder
		heightAddedPerCycle := chamber.PileHeight() - cached.height
		settledAddedPerCycle := settled - cached.settled
		unsettled := target - settled
		possibleCycles := unsettled / settledAddedPerCycle
		heightFromSkippedCycles = possibleCycles * heightAddedPerCycle
		settled += possibleCycles * settledAddedPerCycle
	}

	return chamber.PileHeight() + heightFromSkippedCycles
}
