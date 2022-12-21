package main

import (
	"fmt"
	"io"
	"log"
	"os"
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

	fmt.Println(maxFlow(file))
	end := time.Now()
	fmt.Printf("operation completed in %s\n", end.Sub(start))
}

func maxFlow(r io.Reader) int {
	vs := parseLines(r)

	flow := solve(vs, "AA", 30)

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

func solve(vs valves, start string, simTime int) int {
	q := newSimQueue()

	pather := newPather(vs)

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
