package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

// with gratitude to https://github.com/ChrisWojcik/advent-of-code-2022/blob/main/16/1.py

func main() {
	start := time.Now()
	file, err := os.Open("./input.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fmt.Println(maxFlowTwo(file))
	end := time.Now()
	fmt.Printf("operation completed in %s\n", end.Sub(start))
}

func maxFlow(r io.Reader) int {
	vs := parseLines(r)
	p := newPather(vs)
	flow := solve(vs, "AA", 30, *p)
	return flow
}

func maxFlowTwo(r io.Reader) int {
	vs := parseLines(r)
	flow := solveTwo(vs)
	return flow
}

type simState struct {
	timeElapsed  int
	openValves   set[openEvent]
	currentValve string
}

type set[T any] struct {
	vals map[string]T
}

func (s set[T]) has(a string) bool {
	if _, ok := s.vals[a]; ok {
		return true
	}

	return false
}

func (s *set[T]) add(a string, b T) {
	if s.vals == nil {
		s.vals = make(map[string]T)
	}

	s.vals[a] = b
}

func (s set[T]) size() int {
	return len(s.vals)
}

func (s set[T]) copy() set[T] {
	ns := set[T]{}

	for k, v := range s.vals {
		ns.add(k, v)
	}

	return ns
}

type valves map[string]valve

type valve struct {
	name      string
	flow      int
	open      bool
	neighbors []string
}

type openEvent struct {
	flow     int
	openedAt int
}

type simQueue []simState

func (sq *simQueue) push(s simState) {
	*sq = append(*sq, s)
}

func (sq *simQueue) shift() simState {
	shifted := (*sq)[0]
	*sq = (*sq)[1:]
	return shifted
}

func newSimQueue() *simQueue {
	sq := make(simQueue, 0)
	return &sq
}

func sumPressures(opens set[openEvent], simTime int) int {
	flow := 0

	for _, event := range opens.vals {
		runTime := intMax(simTime-event.openedAt, 0)
		valveFlow := event.flow * runTime
		flow += valveFlow
	}

	return flow
}

func intMax(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func solve(vs valves, start string, simTime int, pather pather) int {
	q := newSimQueue()

	openable := []string{}

	for _, v := range vs {
		if v.flow > 0 {
			openable = append(openable, v.name)
		}
	}

	// initial state has the starting valve unopened. This deals with the case
	// of the starting valve having a flow of 0
	initialState := simState{
		timeElapsed:  0,
		openValves:   set[openEvent]{},
		currentValve: start,
	}

	q.push(initialState)

	var maxFlow int

	for len(*q) > 0 {
		currState := q.shift()

		// if we're out of time, OR we've opened every valve,
		// discontinue this 'branch' of state.
		if currState.timeElapsed >= simTime || currState.openValves.size() == len(openable) {
			branchMaxFlow := sumPressures(currState.openValves, simTime)
			maxFlow = intMax(maxFlow, branchMaxFlow)
			continue
		}

		// for every valve yet-to-be-opened, push a state frame representing the
		// state after moving to the valve and opening it.  Note this will cover
		// the current valve if it is unopened; it will have distance of 0.
		for _, nextValve := range openable {
			if !currState.openValves.has(nextValve) {
				// 1 minute per travel unit + 1 minute to actually open
				travelTime := pather.pathLength(currState.currentValve, nextValve)
				openTime := 1
				newElapsed := travelTime + openTime + currState.timeElapsed

				// take a copy for the next state frame we're about to push
				openedValve := vs[nextValve]
				existingOpenEvents := currState.openValves.copy()
				existingOpenEvents.add(openedValve.name, openEvent{
					openedValve.flow,
					newElapsed,
				})

				nextState := simState{
					timeElapsed:  newElapsed,
					openValves:   existingOpenEvents,
					currentValve: nextValve,
				}
				q.push(nextState)
			}
		}
	}

	return maxFlow
}

// TODO: Need a heuristic to do search tree pruning. Without it, this solution takes approx 4mins
// to complete on my local machine.
func solveOneAgent(openable []valve, start string, simTime int, pather pather) int {
	q := newSimQueue()

	// initial state has the starting valve unopened. This deals with the case
	// of the starting valve having a flow of 0
	initialState := simState{
		timeElapsed:  0,
		openValves:   set[openEvent]{},
		currentValve: start,
	}

	q.push(initialState)

	var maxFlow int

	for len(*q) > 0 {
		currState := q.shift()

		// if we're out of time, OR we've opened every valve,
		// discontinue this 'branch' of state.
		if currState.timeElapsed >= simTime || currState.openValves.size() == len(openable) {
			branchMaxFlow := sumPressures(currState.openValves, simTime)
			maxFlow = intMax(maxFlow, branchMaxFlow)
			continue
		}

		// for every valve yet-to-be-opened, push a state frame representing the
		// state after moving to the valve and opening it.  Note this will cover
		// the current valve if it is unopened; it will have distance of 0.
		for _, nextValve := range openable {
			if !currState.openValves.has(nextValve.name) {
				// 1 minute per travel unit + 1 minute to actually open
				travelTime := pather.pathLength(currState.currentValve, nextValve.name)
				openTime := 1
				newElapsed := travelTime + openTime + currState.timeElapsed

				// take a copy for the next state frame we're about to push
				openedValve := nextValve
				existingOpenEvents := currState.openValves.copy()
				existingOpenEvents.add(openedValve.name, openEvent{
					openedValve.flow,
					newElapsed,
				})

				nextState := simState{
					timeElapsed:  newElapsed,
					openValves:   existingOpenEvents,
					currentValve: nextValve.name,
				}
				q.push(nextState)
			}
		}
	}

	return maxFlow
}

func valveSubsets(vs []valve, result *[][]valve, current *[]valve, pos int) {
	max := len(vs)
	if pos >= max {
		// Limit the range of tested subsets to (half of all subsets) +/- (quarter of all subsets).
		// The 25% variance is totally arbitrary; I just picked a number to rule out some combinations.
		// This would fail badly in a scenario where the travel cost to each node is not fixed at 1.
		// E.g. if there were 50 nodes, and 2 of them were VERY far away from the others, the optimal solution
		// would be to sacrifice one agent by sending them to the distant nodes and no other nodes. Since
		// the heuristic here rules out a 48/2 split, it would rule out the solution.
		midpoint := len(vs) / 2
		variance := midpoint / 2
		currLen := len(*current)
		if currLen > midpoint+variance || currLen < midpoint-variance {
			return
		}
		out := make([]valve, len(*current))
		copy(out, *current)
		*result = append(*result, out)
		return
	}

	v := vs[pos]
	*current = append(*current, v)
	valveSubsets(vs, result, current, pos+1)
	*current = (*current)[:len(*current)-1]
	valveSubsets(vs, result, current, pos+1)
}

// Get the complement of the current subset of openable valves
func subsetComplement(openable, subset []valve) []valve {
	ss := set[struct{}]{}
	for _, v := range subset {
		ss.add(v.name, struct{}{})
	}
	out := make([]valve, len(openable)-len(subset))
	if len(out) == 0 {
		return out
	}
	outPos := 0
	for _, v := range openable {
		if ss.has(v.name) {
			continue
		}
		out[outPos] = v
		outPos++
	}
	return out
}

// Since there is no difference between the 2 agents, any run with the same valves in
// *either* the subset *or* its complement will yield the same result.  Since each
// simulation is expensive, we want to avoid running each possible subset more than once.
// By hashing the subset and its complement, we can check whether we've already done
// the computation for a given subset or not, and skip it if we have.
func runHash(subset []valve) string {
	var b strings.Builder
	for _, v := range subset {
		b.WriteString(v.name)
	}
	rns := []rune(b.String())
	sort.Slice(rns, func(i, j int) bool {
		return rns[i] < rns[j]
	})
	return string(rns)
}

// Each agent has identical capabilities, so we should be able to do the following:
// 1. Generate every subset of openable valves (i.e. valves with non-zero flow), and its complement.
// 2. For every subset/complement pair, run 'solve' with 26 minutes of time with the subset, then again with its complement.
// 3. Add the results of each subset complement pair, and compare the result with current maximum.
// 4. If the result is greater than the maximum, update the running answer.
// 5. Return answer.
func solveTwo(allValves valves) int {
	openableValves := []valve{}
	for _, v := range allValves {
		if v.flow > 0 {
			openableValves = append(openableValves, v)
		}
	}
	// the number of subsets for a list of length N is 2^N
	// subsetCount := math.Pow(2, float64(len(openableValves)))
	// subsets := make([][]valve, int(subsetCount))
	subsets := [][]valve{}
	valveSubsets(openableValves, &subsets, &[]valve{}, 0)
	ssCount := len(subsets)
	fmt.Printf("Processing %v subsets\n", ssCount)
	maxFlow := 0
	pather := newPather(allValves)
	seenSubsets := set[struct{}]{}
	for _, ss := range subsets {
		ssHash := runHash(ss)
		if seenSubsets.has(ssHash) {
			continue
		}
		complement := subsetComplement(openableValves, ss)
		compHash := runHash(complement)
		seenSubsets.add(compHash, struct{}{})
		seenSubsets.add(ssHash, struct{}{})
		elephantSum := solveOneAgent(ss, "AA", 26, *pather)
		selfSum := solveOneAgent(complement, "AA", 26, *pather)
		currSum := selfSum + elephantSum
		maxFlow = intMax(maxFlow, currSum)
	}
	return maxFlow
}
