package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

func main() {
	r, err := os.Open("./input.txt")

	if err != nil {
		panic(err)
	}

	fmt.Println(solvePartTwo(r))
}

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

type botMaxes map[resource]int

type blueprint struct {
	id          int
	oreBot      oreCost
	clayBot     clayCost
	obsidianBot obsidianCost
	geodeBot    geodeCost
	botMax      botMaxes
}

func parseBlueprint(raw string) blueprint {
	maxes := make(botMaxes)
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
	maxes[ore] = asInts[1]
	bp.clayBot.oreCost = asInts[2]
	maxes[ore] = intMax(asInts[2], maxes[ore])
	bp.obsidianBot.oreCost = asInts[3]
	maxes[ore] = intMax(asInts[3], maxes[ore])
	bp.obsidianBot.clayCost = asInts[4]
	maxes[clay] = asInts[4]
	bp.geodeBot.oreCost = asInts[5]
	maxes[ore] = intMax(asInts[5], maxes[ore])
	bp.geodeBot.obsidianCost = asInts[6]
	maxes[obsidian] = asInts[6]
	bp.botMax = maxes
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

func sortMapKeys(m map[resource]int) []resource {
	out := make([]resource, len(m))
	pos := 0
	for k := range m {
		out[pos] = k
		pos++
	}
	slices.Sort(out)
	return out
}

func (bc botCounts) string() string {
	var b strings.Builder
	sortedKeys := sortMapKeys(bc)
	for _, v := range sortedKeys {
		b.WriteString(fmt.Sprintf("%d-%d", v, bc[v]))
	}

	return b.String()
}

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

func (sp stockPiles) string() string {
	var b strings.Builder
	sortedKeys := sortMapKeys(sp)
	for _, v := range sortedKeys {
		b.WriteString(fmt.Sprintf("%d-%d", v, sp[v]))
	}

	return b.String()
}

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
	runTime int
}

func (state simulationState) string() string {
	return fmt.Sprintf("%d-%s-%s", state.elapsedTime, state.botCounts.string(), state.stockPiles.string())
}

func newSimState(bp blueprint, runTime int) simulationState {
	return simulationState{
		newBotCounts(),
		bp,
		0,
		newStockPiles(),
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
		input.runTime,
	}
}

func (s *simulationState) harvest() {
	for resourceType, botCount := range s.botCounts {
		s.stockPiles[resourceType] += botCount
	}
}

func (s *simulationState) elapseTime(mins int) {
	for mins > 0 {
		mins--
		s.elapsedTime++
		s.harvest()
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

// never build more of a particular kind of bot than can be
// consumed by a single iteration
func (s simulationState) canEnqueue(kind botKind) bool {
	switch kind {
	case ore:
		return s.botCounts[ore] < s.blueprint.botMax[ore]
	case clay:
		return s.botCounts[clay] < s.blueprint.botMax[clay]
	case obsidian:
		return s.botCounts[clay] > 0 && s.botCounts[obsidian] < s.blueprint.botMax[obsidian]
	case geode:
		return s.botCounts[obsidian] > 0
	default:
		panic("bad argument to 'canEnqueue'")
	}
}

func (s *simulationState) addBot(kind botKind) {
	switch kind {
	case ore:
		missingOre := s.blueprint.oreBot.oreCost - s.stockPiles[ore]

		for missingOre > 0 && s.elapsedTime < s.runTime {
			s.elapseTime(1)
			missingOre = s.blueprint.oreBot.oreCost - s.stockPiles[ore]
		}

		// if we hit the end of the simulation time while we're waiting for enough
		// ore to build the next bot, don't actually build the bot
		if s.elapsedTime == s.runTime {
			return
		}

		// otherwise, there's at least 1 iteration remaining before the simulation ends,
		// so use that iteration to debit the cost of a bot, run the harvest for the iteration, and then add
		// the newly-constructed bot to bot totals for this state frame.
		s.stockPiles[ore] -= s.blueprint.oreBot.oreCost
		s.elapseTime(1)
		s.botCounts[ore]++

	case clay:
		missingOre := s.blueprint.clayBot.oreCost - s.stockPiles[ore]

		for missingOre > 0 && s.elapsedTime < s.runTime {
			s.elapseTime(1)
			missingOre = s.blueprint.clayBot.oreCost - s.stockPiles[ore]
		}

		if s.elapsedTime == s.runTime {
			return
		}

		s.stockPiles[ore] -= s.blueprint.clayBot.oreCost
		s.elapseTime(1)
		s.botCounts[clay]++

	case obsidian:
		missingOre := s.blueprint.obsidianBot.oreCost - s.stockPiles[ore]
		missingClay := s.blueprint.obsidianBot.clayCost - s.stockPiles[clay]
		if missingOre > 0 || missingClay > 0 {
			for (missingClay > 0 || missingOre > 0) && s.elapsedTime < s.runTime {
				s.elapseTime(1)
				missingOre = s.blueprint.obsidianBot.oreCost - s.stockPiles[ore]
				missingClay = s.blueprint.obsidianBot.clayCost - s.stockPiles[clay]
			}
		}

		if s.elapsedTime == s.runTime {
			return
		}

		s.stockPiles[ore] -= s.blueprint.obsidianBot.oreCost
		s.stockPiles[clay] -= s.blueprint.obsidianBot.clayCost
		s.elapseTime(1)
		s.botCounts[obsidian]++

	case geode:
		missingOre := s.blueprint.geodeBot.oreCost - s.stockPiles[ore]
		missingObsidian := s.blueprint.geodeBot.obsidianCost - s.stockPiles[obsidian]
		if missingOre > 0 || missingObsidian > 0 {
			for (missingObsidian > 0 || missingOre > 0) && s.elapsedTime < s.runTime {
				s.elapseTime(1)
				missingOre = s.blueprint.geodeBot.oreCost - s.stockPiles[ore]
				missingObsidian = s.blueprint.geodeBot.obsidianCost - s.stockPiles[obsidian]
			}
		}

		if s.elapsedTime == s.runTime {
			return
		}

		s.stockPiles[ore] -= s.blueprint.geodeBot.oreCost
		s.stockPiles[obsidian] -= s.blueprint.geodeBot.obsidianCost
		s.elapseTime(1)
		s.botCounts[geode]++

	default:
		panic("bot enqueueing default case reached in error")
	}
}

// returns an array of all possible states for the beginning of the next minute
func (s *simulationState) nextStates() []simulationState {
	nextStates := []simulationState{}

	if s.sated() {
		newState := copyState(*s)
		// fast-forward to end of simulation
		newState.elapsedTime = s.runTime
		// once we have enough bots online, the theoretical max is the
		// same as the actual max for the given elapsed time.
		newState.stockPiles[geode] = newState.theoreticalMax()
		nextStates = append(nextStates, newState)
		return nextStates
	}

	for botKind := range s.botCounts {
		if s.canEnqueue(botKind) {
			newState := copyState(*s)
			newState.addBot(botKind)
			nextStates = append(nextStates, newState)
		}
	}

	return nextStates
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
	cache := make(map[string]struct{})
	for len(q) > 0 {
		currState := q[0]
		q = q[1:]

		maxOpenedGeodes = intMax(maxOpenedGeodes, currState.stockPiles[geode])

		// don't derive any future states from a current state with no
		// remaining runtime
		if currState.elapsedTime == totalRunTime {
			continue
		}

		// if a proposed state can't possibly outperform some state we've already
		// seen, there's no need to check any of its possible future states
		if maxOpenedGeodes > 0 && currState.theoreticalMax() <= maxOpenedGeodes {
			continue
		}

		cacheKey := currState.string()

		// if we've seen this exact state before, add no future states
		if _, cacheHit := cache[cacheKey]; cacheHit {
			continue
		}

		cache[cacheKey] = struct{}{}

		nextStates := currState.nextStates()

		lo.ForEach(nextStates, func(sim simulationState, idx int) {
			q = append(q, sim)
		})
	}
	return maxOpenedGeodes
}

func solvePartOne(r io.Reader) int {
	scanner := bufio.NewScanner(r)
	qualitySum := 0
	for scanner.Scan() {
		bp := parseBlueprint(scanner.Text())
		qualitySum += getMaxGeodes(24, bp) * bp.id
	}
	return qualitySum
}

func solvePartTwo(r io.Reader) int {
	scanner := bufio.NewScanner(r)
	product := 1
	linesRead := 0
	for scanner.Scan() && linesRead < 3 {
		bp := parseBlueprint(scanner.Text())
		product *= getMaxGeodes(32, bp)
		linesRead++
	}
	return product
}
