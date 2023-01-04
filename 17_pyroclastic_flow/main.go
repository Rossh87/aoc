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

	ans := solveOne(string(bytes))
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
func solveOne(rawDirs string) int {
	js := jets.New(rawDirs)
	rocks := rocks.New()
	chamber := chamber.New(7)
	stoppedCount := 0

	for stoppedCount < 2022 {
		chamber.AddRock(rocks, js)
		stoppedCount++
	}

	return chamber.PileHeight()
}
