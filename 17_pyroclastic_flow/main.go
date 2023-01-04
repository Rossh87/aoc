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

// take the hash after rock has fallen, combining the following:
// 1. next jet position
// 2. next shape position
// 3. x-offset of the final resting place of the 'start' point for just-settled shape
// with each hash, store:
// 1. rocks settled so far
// 2. pileheight
// once we've found cycle, calculate:
// 1. the height and rocks added by 1 full cycle by subtracting height/settled of current state from height/settled when cycle started
// 2. number of full cycles that can run by dividing the number of rocks remaining to settle by the number of rocks that settle in 1 cycle
// 3. maxHeight += currHeight + (cycleCount * heightPerCycle)
// 4. settledRemaining %= settledPerCycle
// Finally, run sim to completion for settledRemaining
func hash(nextJet, nextShape int, cont chamber.Contour) string {
	return fmt.Sprintf("%d-%d-%s", nextJet, nextShape, cont.Normalized().String())
}

type cacheState struct {
	settled int
	height  int
}

type cache map[string]cacheState

func solveOne(rawDirs string) int {
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

	for int64(settled) < target {
		chamber.AddRock(rocks, js)
		settled++
		hashKey := hash(js.Peek(), rocks.Peek(), chamber.Contour())
		cached, exists := cache[hashKey]

		// if we've never seen this state before, add it to cache and continue
		// NB we're banking here on getting a chache hit before settled grows
		// large enough to overflow a 32-bit int.
		if !exists {
			cache[hashKey] = cacheState{int(settled), chamber.PileHeight()}
			continue
		}

		// if we've found a cycle, fast-forward through as many cycles
		// as possible, then let loop finish the remainder
		heightAddedPerCycle := int64(chamber.PileHeight() - cached.height)
		settledAddedPerCycle := int(settled) - cached.settled
		unsettled := target - int64(settled)
		possibleCycles := unsettled / int64(settledAddedPerCycle)
		heightFromSkippedCycles = possibleCycles * heightAddedPerCycle
		settled += possibleCycles * int64(settledAddedPerCycle)
	}

	return int64(chamber.PileHeight()) + heightFromSkippedCycles
}
