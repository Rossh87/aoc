package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, _ := os.Open("./input.txt")

	res := getOrderedPairSum(file)

	fmt.Println(res)
}

func isOrderedPair(first, second []any, idx int) bool {
	var l any

	if idx > len(first)-1 {
		l = nil
	} else {
		l = first[idx]
	}

	var r any

	if idx > len(second)-1 {
		r = nil
	} else {
		r = second[idx]
	}

	// we've run out of elements...
	if l == nil || r == nil {
		// everything preceding was true if we got here, so we're done
		if l == nil && r == nil {
			return true
		}

		// if left runs out first, we're done
		if l == nil {
			return true
		}

		// otherwise it's out-of-order
		return false
	}

	switch lType := l.(type) {
	case int:
		switch rType := r.(type) {
		// int int
		case int:
			if lType > rType {
				return false
			}
			return isOrderedPair(first, second, idx+1)
		// int []any
		case []any:
			first[idx] = []any{lType}
			return isOrderedPair(first, second, idx)
		default:
			log.Fatalf("received unexpected types at idx %d in %v", idx, second)
		}

	case []any:
		switch rType := r.(type) {
		// []int int
		case int:
			second[idx] = []any{rType}
			return isOrderedPair(first, second, idx)
		// []int []any
		case []any:
			return isOrderedPair(lType, rType, 0)
		default:
			log.Fatalf("received unexpected types at idx %d in %v", idx, second)
		}
	default:
		log.Fatalf("received unexpected types at idx %d in %v", idx, first)
	}

	return false
}

type Pair struct {
	leftSet  bool
	rightSet bool
	left     []any
	right    []any
}

type Pairs []Pair

func getOrderedPairIndices(p Pairs) []int {
	result := []int{}

	for idx, pair := range p {
		if isOrderedPair(pair.left, pair.right, 0) {
			result = append(result, idx)
		}
	}

	return result
}

func stringToSlice(rns []rune, start int) ([]any, int) {
	wrapper := []any{}

	valSet := false
	currVal := 0

	for i := start + 1; i < len(rns); i++ {
		rn := rns[i]

		if rn == '[' {
			inner, skipTo := stringToSlice(rns, i)
			wrapper = append(wrapper, inner)
			i = skipTo
			continue
		}

		if rn == ']' {
			if valSet {
				wrapper = append(wrapper, currVal)
			}
			return wrapper, i
		}

		if rn == ',' {
			if valSet {
				wrapper = append(wrapper, currVal)
				currVal = 0
				valSet = false
			}
			continue
		}

		valSet = true
		v := int(rn - '0')
		currVal *= 10
		currVal += v
	}

	return wrapper, len(rns)
}

func getOrderedPairSum(r io.Reader) int {
	scanner := bufio.NewScanner(r)

	pairs := Pairs{}

	pair := Pair{}

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			pairs = append(pairs, pair)
			pair = Pair{}
			continue
		}

		list, _ := stringToSlice([]rune(line), 0)

		if pair.leftSet {
			pair.right = list
			pair.rightSet = true
		} else {
			pair.left = list
			pair.leftSet = true
		}
	}

	// append final pair when parsing is done
	pairs = append(pairs, pair)

	indices := getOrderedPairIndices(pairs)

	Sum := 0

	for _, v := range indices {
		// add one since input is 1-indexed
		Sum += v + 1
	}

	return Sum
}
