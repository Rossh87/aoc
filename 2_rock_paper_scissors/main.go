package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type choice string

const (
	rock     choice = "rock"
	paper    choice = "paper"
	scissors choice = "scissors"
)

// map a string character to a 'move' the player has chosen
type choiceMap = map[string]choice

// map the 'hand' a player has chosen to a base value the player receives just for making that choice
type choiceScores map[choice]int

type outcome string

const (
	win  outcome = "win"
	loss outcome = "loss"
	draw outcome = "draw"
)

// map the outcome of a round to the bonus value player will receive
type resultScores map[outcome]int

// map cipher string character to an outcome suggested by the strategy
type outcomeMap map[string]outcome

// get from opponent choice to player's choice by way of desired outcome
type playerChoiceMap map[choice]map[outcome]choice

func scoreFromOutcome(choicePair string, cm choiceMap, om outcomeMap, cs choiceScores, rs resultScores) int {
	pcm := playerChoiceMap{
		rock: map[outcome]choice{
			win:  paper,
			loss: scissors,
			draw: rock,
		},
		paper: map[outcome]choice{
			win:  scissors,
			loss: rock,
			draw: paper,
		},
		scissors: map[outcome]choice{
			win:  rock,
			loss: paper,
			draw: scissors,
		},
	}

	pair := strings.Split(choicePair, " ")

	oppChoice := cm[pair[0]]

	desiredOutcome := om[pair[1]]

	outcomeScore := rs[desiredOutcome]

	playerChoice := pcm[oppChoice][desiredOutcome]

	playerChoiceScore := cs[playerChoice]

	return outcomeScore + playerChoiceScore
}

func main() {
	cm := choiceMap{
		"A": rock,
		"B": paper,
		"C": scissors,
	}

	cs := choiceScores{
		rock:     1,
		paper:    2,
		scissors: 3,
	}

	rs := resultScores{
		"win":  6,
		"loss": 0,
		"draw": 3,
	}

	om := outcomeMap{
		"X": "loss",
		"Y": "draw",
		"Z": "win",
	}

	file, err := os.Open("./input.txt")

	if err != nil {
		log.Fatalf("failed opening file: %v", err)
	}

	scanner := bufio.NewScanner(file)

	score := 0

	for scanner.Scan() {
		pair := scanner.Text()

		roundScore := scoreFromOutcome(pair, cm, om, cs, rs)

		if err != nil {
			log.Fatalf("invalid round result: %v", err)
		}

		score += roundScore
	}

	fmt.Println(score)
}
