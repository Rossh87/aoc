package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/samber/lo"
)

type oreCost struct {
	oreCost int
}
type clayCost struct {
	oreCost int
}
type obsidianCost struct {
	oreCost  int
	clayCost int
}
type geodeCost struct {
	oreCost      int
	obsidianCost int
}

type blueprint struct {
	id          int
	oreBot      oreCost
	clayBot     clayCost
	obsidianBot obsidianCost
	geodeBot    geodeCost
}

// type blueprints []blueprint

func parseBlueprint(raw string) blueprint {
	exp := regexp.MustCompile(`\d+`)
	matches := exp.FindAllString(raw, -1)
	if len(matches) != 7 {
		panic(fmt.Sprintf("string %s contained inappropriate number of digits:\n expected 7, but received %d", raw, len(matches)))
	}
	asInts := lo.Map(matches, func(item string, idx int) int {
		v, err := strconv.ParseInt(item, 10, 32)
		if err != nil {
			panic(err)
		}
		return int(v)
	})
	bp := blueprint{}
	bp.id = asInts[0]
	bp.oreBot.oreCost = asInts[1]
	bp.clayBot.oreCost = asInts[2]
	bp.obsidianBot.oreCost = asInts[3]
	bp.obsidianBot.clayCost = asInts[4]
	bp.geodeBot.oreCost = asInts[5]
	bp.geodeBot.obsidianCost = asInts[6]
	return bp
}

type resource int

const (
	ore resource = iota
	clay
	obsidian
	geode
	none
)

type botKind = resource

type botCounts map[resource]int

// NB we start with 1 ore bot
func newBotCounts() botCounts {
	bc := make(botCounts)
	bc[ore] = 1
	bc[clay] = 0
	bc[obsidian] = 0
	bc[geode] = 0
	return bc
}

type stockPiles map[resource]int

func newStockPiles() stockPiles {
	sp := make(stockPiles)
	sp[ore] = 0
	sp[clay] = 0
	sp[obsidian] = 0
	sp[geode] = 0
	return sp
}

type simulationState struct {
	botCounts
	blueprint
	elapsedTime int
	stockPiles
	enqueued resource
	runTime  int
}

func newSimState(bp blueprint, runTime int) simulationState {
	return simulationState{
		newBotCounts(),
		bp,
		0,
		newStockPiles(),
		none,
		runTime,
	}
}

func copyState(input simulationState) simulationState {
	countCopy := make(botCounts)
	stockPilesCopy := make(stockPiles)
	for k, v := range input.botCounts {
		countCopy[k] = v
	}
	for k, v := range input.stockPiles {
		stockPilesCopy[k] = v
	}
	return simulationState{
		countCopy,
		input.blueprint,
		input.elapsedTime,
		stockPilesCopy,
		input.enqueued,
		input.runTime,
	}
}

func (s *simulationState) harvest() {
	for resourceType, botCount := range s.botCounts {
		s.stockPiles[resourceType] += botCount
	}
}

func (s *simulationState) constructBot() {
	if s.enqueued == none {
		return
	}
	botType := s.enqueued
	s.enqueued = none
	s.botCounts[botType]++
}

func (s *simulationState) elapseTime(mins int) {
	// we complete a harvest phase and a build phase
	// every minute that elapses
	for mins > 0 {
		mins--
		s.elapsedTime++
		s.harvest()
		s.constructBot()
	}
}

func (s simulationState) theoreticalMax() int {
	timeRemaining := s.runTime - s.elapsedTime
	thMax := s.stockPiles[geode]
	thMax += s.botCounts[geode] * timeRemaining
	thMax += (timeRemaining * (timeRemaining + 1)) / 2
	return thMax
}

// once the sim is sated, theoretical max becomes actual max
func (s simulationState) sated() bool {
	return s.botCounts[ore] == s.blueprint.geodeBot.oreCost &&
		s.botCounts[obsidian] == s.blueprint.geodeBot.obsidianCost
}

// TODO: bad modeling here
func (s simulationState) canEnqueue(kind botKind) bool {
	switch kind {
	case ore:
		return s.botCounts[ore] < 5
	case clay:
		return s.botCounts[clay] < 5
	case obsidian:
		return s.botCounts[clay] > 0 && s.botCounts[obsidian] < 5
	case geode:
		return s.botCounts[obsidian] > 0
	default:
		panic("bad argument to 'canEnqueue'")
	}
}

func (s *simulationState) addBot(kind botKind) bool {
	switch kind {
	case ore:
		missingOre := s.blueprint.oreBot.oreCost - s.stockPiles[ore]

		for missingOre > 0 {
			s.elapseTime(1)
			missingOre = s.blueprint.oreBot.oreCost - s.stockPiles[ore]
		}

		if s.elapsedTime > s.runTime {
			return false
		}

		s.stockPiles[ore] -= s.blueprint.oreBot.oreCost
		s.enqueued = ore
		return true

	case clay:
		missingOre := s.blueprint.clayBot.oreCost - s.stockPiles[ore]

		for missingOre > 0 {
			s.elapseTime(1)
			missingOre = s.blueprint.clayBot.oreCost - s.stockPiles[ore]
		}

		if s.elapsedTime > s.runTime {
			return false
		}

		s.stockPiles[ore] -= s.blueprint.clayBot.oreCost
		s.enqueued = clay
		return true

	case obsidian:
		missingOre := s.blueprint.obsidianBot.oreCost - s.stockPiles[ore]
		missingClay := s.blueprint.obsidianBot.clayCost - s.stockPiles[clay]
		if missingOre > 0 || missingClay > 0 {
			for missingClay > 0 || missingOre > 0 {
				s.elapseTime(1)
				missingOre = s.blueprint.obsidianBot.oreCost - s.stockPiles[ore]
				missingClay = s.blueprint.obsidianBot.clayCost - s.stockPiles[clay]
			}
		}
		if s.elapsedTime > s.runTime {
			return false
		}
		s.stockPiles[ore] -= s.blueprint.obsidianBot.oreCost
		s.stockPiles[clay] -= s.blueprint.obsidianBot.clayCost
		s.enqueued = obsidian
		return true

	case geode:
		missingOre := s.blueprint.geodeBot.oreCost - s.stockPiles[ore]
		missingObsidian := s.blueprint.geodeBot.obsidianCost - s.stockPiles[obsidian]
		if missingOre > 0 || missingObsidian > 0 {
			for missingObsidian > 0 || missingOre > 0 {
				s.elapseTime(1)
				missingOre = s.blueprint.geodeBot.oreCost - s.stockPiles[ore]
				missingObsidian = s.blueprint.geodeBot.obsidianCost - s.stockPiles[obsidian]
			}
		}
		if s.elapsedTime > s.runTime {
			return false
		}
		s.stockPiles[ore] -= s.blueprint.geodeBot.oreCost
		s.stockPiles[obsidian] -= s.blueprint.geodeBot.obsidianCost
		s.enqueued = geode
		return true
	default:
		panic("bot enqueueing default case reached in error")
	}
}

// 1. don't add a state frame for building any robots if current income-per-round is enough
// to build 1 geode robot per turn
// 2. track best theoretical max seen for given time elapsed. For each possible next frame, calculate
// the theoretical max it could produce.  If it is less than the saved theoretical max, abandon the frame,
// i.e. don't add it to the queue.

// returns an array of all possible states for the beginning of the next minute
func (s *simulationState) nextStates() []simulationState {
	// this adds any bots enqueued by previous calls to this function
	s.elapseTime(1)
	next := []simulationState{}

	// once we're sated, we can calculate the max geodes programatically and skip a bunch of frames by
	// pushing a frame with all time elapsed
	if s.sated() {
		newState := copyState(*s)
		newState.elapsedTime = s.runTime
		newState.stockPiles[geode] = newState.theoreticalMax()
		next = append(next, newState)
		return next
	}

	// next = append(next, copyState(*s))

	for botKind := range s.botCounts {
		if s.canEnqueue(botKind) {
			newState := copyState(*s)
			newState.addBot(botKind)
			next = append(next, newState)
		}
	}

	return next
}

func intMax(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func getMaxGeodes(totalRunTime int, bp blueprint) int {
	baseState := newSimState(bp, totalRunTime)
	q := []simulationState{baseState}
	maxOpenedGeodes := 0
	cache := map[int]int{}
	loopcount := 0
	for len(q) > 0 {
		loopcount++
		currState := q[0]
		q = q[1:]

		if currState.elapsedTime == totalRunTime {
			maxOpenedGeodes = intMax(maxOpenedGeodes, currState.stockPiles[geode])
			continue
		}

		// generating next states handles all updates due to elapsed time
		nextStates := currState.nextStates()

		lo.ForEach(nextStates, func(sim simulationState, idx int) {
			cachedMax, exists := cache[sim.elapsedTime]
			thMax := sim.theoreticalMax()
			if exists && thMax < cachedMax {
				return
			}
			cache[sim.elapsedTime] = thMax
			q = append(q, sim)
		})
	}
	return maxOpenedGeodes
}

func main() {
	fmt.Println("hello")
}

func solve() int {
	return 42
}
