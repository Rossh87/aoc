package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func getDuplicate(lines []string) rune {

	maps := make([]map[rune]struct{}, 3)

	for pos := range maps {
		maps[pos] = make(map[rune]struct{})
	}

	for lineNumber, line := range lines {
		relevantMap := maps[lineNumber]

		for _, rn := range line {
			if _, ok := relevantMap[rn]; ok {
				continue
			}

			relevantMap[rn] = struct{}{}
		}
	}

	for _, rn := range lines[0] {
		_, inString2 := maps[1][rn]

		_, inString3 := maps[2][rn]

		if inString2 && inString3 {
			return rn
		}
	}

	return '0'
}

func runePriority(rn rune) int {
	if rn > 90 {
		return int(rn) - 96
	}

	return int(rn) - 64 + 26
}

func main() {
	dst := "./input.txt"

	file, err := os.Open(dst)

	if err != nil {
		log.Fatalf("Unable to open %s: %v", dst, err)
	}

	scanner := bufio.NewScanner(file)

	total := 0

	lines := make([]string, 3)
	currLine := 0

	for scanner.Scan() {
		line := scanner.Text()
		linePos := currLine % 3
		lines[linePos] = line

		if linePos == 2 {
			total += runePriority(getDuplicate(lines))
		}

		currLine += 1
	}

	fmt.Println(total)
}
