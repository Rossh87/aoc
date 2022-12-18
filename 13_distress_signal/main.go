package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"sort"
)

func main() {
	file, _ := os.Open("./input.txt")

	res := getDividerProduct(file)

	fmt.Println(res)
}

type pairStatus int

const (
	valid pairStatus = iota
	invalid
	inconclusive
)

func isOrderedPair(first, second []any, idx int) pairStatus {
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
		// if we get to this point, everything up to idx was equal, so we don't have an answer
		if l == nil && r == nil {
			return inconclusive
		}

		// if left runs out first, we're done
		if l == nil {
			return valid
		}

		// otherwise it's out-of-order
		return invalid
	}

	switch lType := l.(type) {
	case int:
		switch rType := r.(type) {
		// int int
		case int:
			if lType > rType {
				return invalid
			}

			if lType < rType {
				return valid
			}

			return isOrderedPair(first, second, idx+1)
		// int []any
		case []any:
			// make copies--if we mutate slices in-place, it messes up part 2
			double := make([]any, len(first))
			copy(double, first)
			double[idx] = []any{lType}
			return isOrderedPair(double, second, idx)
		default:
			log.Fatalf("received unexpected types at idx %d in %v", idx, second)
		}

	case []any:
		switch rType := r.(type) {
		// []any int
		case int:
			// make copies--if we mutate slices in-place, it messes up part 2
			double := make([]any, len(second))
			copy(double, second)
			double[idx] = []any{rType}
			return isOrderedPair(first, double, idx)
		// []any []any
		case []any:
			subResult := isOrderedPair(lType, rType, 0)

			if subResult == inconclusive {
				return isOrderedPair(first, second, idx+1)
			}

			return subResult
		default:
			log.Fatalf("received unexpected types at idx %d in %v", idx, second)
		}
	default:
		log.Fatalf("received unexpected types at idx %d in %v", idx, first)
	}

	return invalid
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
		if isOrderedPair(pair.left, pair.right, 0) == valid {
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

type Packet []any

type Packets []Packet

func packPackets(r io.Reader, pax *Packets) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		s, _ := stringToSlice([]rune(line), 0)
		*pax = append(*pax, s)
	}
}

func sortPackets(r io.Reader) Packets {
	pax := Packets{[]any{[]any{2}}, []any{[]any{6}}}

	packPackets(r, &pax)

	sort.Slice(pax, func(i, j int) bool {
		return isOrderedPair(pax[i], pax[j], 0) == valid
	})

	return pax
}

func getDividerProduct(r io.Reader) int {
	sorted := sortPackets(r)

	var ans [2]int

	for i, el := range sorted {
		// remember to account for 1-indexing
		if reflect.DeepEqual(Packet{[]any{2}}, el) {
			ans[0] = i + 1
		}

		if reflect.DeepEqual(Packet{[]any{6}}, el) {
			ans[1] = i + 1
		}
	}

	return ans[0] * ans[1]
}
