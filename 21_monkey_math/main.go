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

	fmt.Println(solvePartOne(file))
}

type operator rune

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

type satisfactionMap map[taskID][]taskID

func (s satisfactionMap) addRelation(satisfier, satisfied taskID) {
	if _, exists := s[satisfier]; exists {
		s[satisfied] = append(s[satisfier], satisfied)
		return
	}

	s[satisfier] = []taskID{satisfied}
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

func solvePartOne(r io.Reader) int {
	scanner := bufio.NewScanner(r)
	ts := newTasks()
	s := make(satisfactionMap)
	for scanner.Scan() {
		ln := scanner.Text()
		ts.add(ln, &s)
	}

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
			return completedTask.result
		}

		updates := s[completedTask.id]

		for _, idToUpdate := range updates {
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
	}

	panic("completed queue without finding a result")
}

func solvePartTwo(r io.Reader) int {
	scanner := bufio.NewScanner(r)
	ts := newTasks()
	s := make(satisfactionMap)
	for scanner.Scan() {
		ln := scanner.Text()
		ts.add(ln, &s)
	}

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
			return completedTask.result
		}

		updates := s[completedTask.id]

		for _, idToUpdate := range updates {
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
	}

	panic("completed queue without finding a result")
}
