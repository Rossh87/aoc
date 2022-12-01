package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	r, err := os.Open("./input.txt")

	if err != nil {
		log.Fatalf("%v", err)
	}

	v, err := maxCals(r, 3)

	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println(v)
}

type PriorityQ []int

func (pq PriorityQ) Len() int {
	return len(pq)
}

func (pq PriorityQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQ) Push(a any) {
	*pq = append(*pq, a.(int))
}

func (pq *PriorityQ) Pop() any {
	old := *pq
	n := len(old)
	popped := old[n-1]
	*pq = old[0 : n-1]
	return popped
}

func (pq PriorityQ) Less(i, j int) bool {
	return pq[i] < pq[j]
}

func maxCals(r io.Reader, aggCount int) (int, error) {
	scanner := bufio.NewScanner(r)

	h := &PriorityQ{}

	heap.Init(h)

	insertItem := func(item int) {
		heap.Push(h, item)

		if h.Len() > aggCount {
			heap.Pop(h)
		}
	}

	currSum := 0

	for scanner.Scan() {
		text := scanner.Text()

		if text == "" {
			insertItem(currSum)
			// fmt.Println(*h)
			currSum = 0
			continue
		}

		val, err := strconv.ParseInt(text, 10, 32)

		if err != nil {
			return 0, fmt.Errorf("failed to convert line value %v to int: %v", text, err)
		}

		currSum += int(val)
	}

	insertItem(currSum)

	aggSum := 0

	for h.Len() > 0 {
		aggSum += heap.Pop(h).(int)
	}

	return aggSum, nil
}
