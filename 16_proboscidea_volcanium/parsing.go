package main

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func extractFlow(s string) int {
	rexp := regexp.MustCompile(`\d+`)
	stringRate := rexp.FindString(s)
	i, err := strconv.Atoi(stringRate)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func decomma(ss *[]string) {
	for i, s := range *ss {
		if i == len(*ss)-1 {
			return
		}

		end := len(s) - 1
		(*ss)[i] = s[:end]
	}
}

func parseLine(line string) valve {
	fields := strings.Fields(line)
	neighbors := fields[9:]
	decomma(&neighbors)
	name := fields[1]
	flow := extractFlow(fields[4])
	return valve{
		name,
		flow,
		false,
		neighbors,
	}
}

func parseLines(r io.Reader) valves {
	out := make(valves)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ln := scanner.Text()
		v := parseLine(ln)
		out[v.name] = v
	}
	return out
}
