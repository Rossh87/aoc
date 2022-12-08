package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type pair struct {
	Low  int
	High int
}

func (t pair) covers(toCompare pair) bool {
	return t.Low <= toCompare.Low && t.High >= toCompare.High
}

func (t pair) overlaps(toCompare pair) bool {
	var earlier pair
	var later pair

	if t.Low <= toCompare.Low {
		earlier = t
		later = toCompare
	} else {
		earlier = toCompare
		later = t
	}

	return earlier.High >= later.Low
}

func getPair(input string) (pair, error) {
	empty := pair{0, 0}

	vals := strings.Split(input, "-")

	low, err := strconv.ParseInt(vals[0], 10, 32)

	if err != nil {
		return empty, err
	}

	high, err := strconv.ParseInt(vals[1], 10, 32)

	if err != nil {
		return empty, err
	}

	return pair{int(low), int(high)}, nil
}

func main() {
	file, err := os.Open("./input.txt")

	if err != nil {
		log.Fatal(err)
	}

	count := 0

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		strPairs := strings.Split(line, ",")

		pair1, err := getPair(strPairs[0])

		if err != nil {
			log.Fatal(err)
		}

		pair2, err := getPair(strPairs[1])

		if err != nil {
			log.Fatal(err)
		}

		if pair1.overlaps(pair2) {
			count += 1
		}
	}

	fmt.Println(count)
}
