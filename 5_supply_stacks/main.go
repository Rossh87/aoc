package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/rossh87/aoc/5_supply_stacks/stacks"
	"github.com/rossh87/aoc/5_supply_stacks/stringtable"
)

func main() {
	file, err := os.Open("./input.txt")

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	var s stacks.Stacks

	for scanner.Scan() {
		line := scanner.Text()

		if lineNumber == 0 {
			s = *stacks.New(line)
		}

		lineNumber++

		switch {
		case line == "":
			s.Reverse()
			continue

		case line[1] == '1':
			continue

		case strings.HasPrefix(line, "move"):
			fields := strings.Fields(line)
			src, _ := strconv.ParseInt(fields[3], 10, 32)
			src -= 1
			dest, _ := strconv.ParseInt(fields[5], 10, 32)
			dest -= 1
			count, _ := strconv.ParseInt(fields[1], 10, 32)
			s.MoveMultiple(int(count), int(src), int(dest))
			continue

		default:
			s.AddRow(stringtable.ParseTableRow(line))
		}

	}

	res := ""

	for _, stack := range s {
		res += string(stack[len(stack)-1])
	}

	fmt.Println(res)
}
