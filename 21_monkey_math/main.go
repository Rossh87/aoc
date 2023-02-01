package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("./input.txt")

	if err != nil {
		panic(err)
	}

	fmt.Println(solvePartTwo(file))
}

type operator rune

func (op operator) inverse(knownInputPos int) operator {
	switch op {
	case '+':
		return '-'
	case '-':
		// unknown - known = needed
		if knownInputPos == 2 {
			return '+'
		}
		// known - unknown = needed
		// s for 'special subtraction', need particular equation to
		// invert subtraction in this direction
		return 's'
	case '*':
		return '/'
	case '/':
		// uknown / known = needed
		if knownInputPos == 2 {
			return '*'
		}
		// known / unknown = needed
		// d for 'special division', need particular equation to
		// invert division in this direction
		return 'd'
	default:
		panic(op)
	}
}

func toOperator(input string) operator {
	if len(input) > 1 {
		panic("received operator base string of length greater than 1")
	}

	return operator([]rune(input)[0])
}

const (
	add      operator = '+'
	subtract operator = '-'
	multiply operator = '*'
	divide   operator = '/'
	none     operator = '0'
)

type taskID string

type operand struct {
	order int
	value int
}

type task struct {
	id       taskID
	done     bool
	inDegree int
	operator operator
	result   int
	requires map[taskID]operand
}

func (t task) operate() int {
	if !t.done {
		panic("operating on unready task")
	}
	operands := make([]operand, len(t.requires))
	pos := 0
	for _, op := range t.requires {
		operands[pos] = op
		pos++
	}
	if len(operands) != 2 {
		panic("wrong number of operands")
	}

	sort.Slice(operands, func(i, j int) bool {
		return operands[i].order < operands[j].order
	})

	switch t.operator {
	case '+':
		return operands[0].value + operands[1].value
	case '-':
		return operands[0].value - operands[1].value
	case '*':
		return operands[0].value * operands[1].value
	case '/':
		return operands[0].value / operands[1].value
	}
	return 0
}

type tasks map[taskID]task

func newTasks() *tasks {
	t := make(tasks)
	return &t
}

func extractId(input string) string {
	return input[:len(input)-1]
}

type satisfactionMap map[taskID]taskID

func (s satisfactionMap) addRelation(satisfier, satisfied taskID) {
	s[satisfier] = satisfied
}

func (ts *tasks) add(rawTask string, s *satisfactionMap) {
	split := strings.Split(rawTask, " ")
	currTaskId := taskID(extractId(split[0]))
	var t task
	// doesn't depend on anything else
	if len(split) == 2 {
		result, err := strconv.ParseInt(split[1], 10, 32)

		if err != nil {
			panic(err)
		}

		t = task{
			currTaskId,
			true,
			0,
			none,
			int(result),
			nil,
		}

		(*ts)[currTaskId] = t
		return
	}
	op := split[2]
	dep1 := taskID(split[1])
	dep2 := taskID(split[3])

	t = task{
		currTaskId,
		false,
		2,
		toOperator(op),
		0,
		map[taskID]operand{
			dep1: {
				1,
				0,
			},
			dep2: {
				2,
				0,
			},
		},
	}
	(*ts)[currTaskId] = t

	s.addRelation(dep1, currTaskId)
	s.addRelation(dep2, currTaskId)
}

func parseInput(r io.Reader, s *satisfactionMap) *tasks {
	scanner := bufio.NewScanner(r)
	ts := newTasks()
	for scanner.Scan() {
		ln := scanner.Text()
		ts.add(ln, s)
	}
	return ts
}

func solvePartOne(r io.Reader) int {
	s := make(satisfactionMap)
	ts := parseInput(r, &s)

	q := []task{}

	for _, t := range *ts {
		if t.done {
			q = append(q, t)
		}
	}

	for len(q) > 0 {
		completedTask := q[0]
		q = q[1:]

		if completedTask.id == "root" {
			return completedTask.operate()
		}

		idToUpdate := s[completedTask.id]

		target, updateTargetExists := (*ts)[idToUpdate]

		if !updateTargetExists {
			panic("trying to update a task that doesn't exist")
		}

		target.inDegree--
		operand := target.requires[completedTask.id]
		operand.value = completedTask.result
		target.requires[completedTask.id] = operand
		if target.inDegree == 0 {
			target.done = true
			target.result = target.operate()
			q = append(q, target)
		}
		(*ts)[idToUpdate] = target
	}

	panic("completed queue without finding a result")
}

func findEqualityTarget(ts *tasks, s *satisfactionMap) int {
	q := []task{}

	for _, t := range *ts {
		// we don't want to start into the tree that includes
		// unknown "humn" value
		if t.id == "humn" {
			continue
		}

		if t.done {
			q = append(q, t)
		}
	}

	for len(q) > 0 {
		completedTask := q[0]
		q = q[1:]

		idToUpdate := (*s)[completedTask.id]

		target, updateTargetExists := (*ts)[idToUpdate]

		if !updateTargetExists {
			panic("trying to update a task that doesn't exist")
		}

		// the first time we update root, we have one of the
		// inputs to root, so we can be done
		if target.id == "root" {
			return completedTask.result
		}

		target.inDegree--
		operand := target.requires[completedTask.id]
		operand.value = completedTask.result
		target.requires[completedTask.id] = operand
		if target.inDegree == 0 {
			target.done = true
			target.result = target.operate()
			q = append(q, target)
		}
		(*ts)[idToUpdate] = target
	}

	panic("completed queue without finding a result")
}

type stackOp struct {
	op    operator
	value int
}

type opStack []stackOp

func buildOpStack(ts *tasks, s *satisfactionMap) opStack {
	stack := make(opStack, 0)
	// we start with the first node that depends on "humn"
	current := (*s)["humn"]

	for current != "root" {
		currNode := (*ts)[current]
		// Choose whichever ancestor node is NOT on the path upwards
		// from "root" to "humn".  That ancestor gives us the known
		// value for the equation that must equal the equality target
		// of the "root" node.
		for parentId, parentVal := range currNode.requires {
			if parentId != "humn" && (*ts)[parentId].done {
				sop := stackOp{
					currNode.operator.inverse(parentVal.order),
					parentVal.value,
				}
				stack = append(stack, sop)
			}
		}
		current = (*s)[current]
	}

	return stack
}

// once we know what the final value needs to be,
// we work through the stack of operations backwards, performing
// each operation we encounter.  In essence, we are re-winding
// a DFS from "humn" to "root"
func unwindStack(value int, stack opStack) int {
	l := len(stack)

	for i := l - 1; i >= 0; i-- {
		op := stack[i]
		switch op.op {
		// op.value is the known input to the equation on the way down to "root"
		case '+':
			value += op.value
		case '-':
			value -= op.value
		case '*':
			value *= op.value
		case '/':
			value /= op.value
		case 's':
			value = -1 * (value - op.value)
		case 'd':
			value = op.value / value
		default:
			panic("found unmatched operator while unwinding opstack")
		}
	}

	return value
}

// answer for actual input is 3327575724809
func solvePartTwo(r io.Reader) int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	s := make(satisfactionMap)
	ts := parseInput(r, &s)
	eqTarget := findEqualityTarget(ts, &s)
	stack := buildOpStack(ts, &s)
	return unwindStack(eqTarget, stack)
}
