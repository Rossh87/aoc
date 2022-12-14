package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("./input.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	c := newCPU()

	for scanner.Scan() {
		ln := scanner.Text()

		fields := strings.Fields(ln)

		switch fields[0] {
		case "addx":
			inc, err := strconv.ParseInt(fields[1], 10, 32)

			if err != nil {
				log.Fatal(err)
			}

			c.addx(int(inc))
		case "noop":
			c.noop()
		}
	}

	// cycle := 20
	// sigSum := 0

	// for cycle < 221 {
	// 	// history slice is zero-idx, but cycles are 1-idx
	// 	registerVal := c.history[cycle-1]
	// 	sigStrength := registerVal * cycle
	// 	sigSum += sigStrength
	// 	cycle += 40
	// }

	// fmt.Println(sigSum)
	draw(c.history)
}

type cpu struct {
	history  []int
	register int
}

func newCPU() *cpu {
	h := make([]int, 0)

	return &cpu{h, 1}
}

func (c *cpu) addx(x int) {
	cycles := 2
	for cycles > 0 {
		c.history = append(c.history, c.register)
		cycles--
	}
	c.register += x
}

func (c *cpu) noop() {
	c.history = append(c.history, c.register)
}

func run(s string) int {
	c := newCPU()

	scanner := bufio.NewScanner(strings.NewReader(s))

	for scanner.Scan() {
		ln := scanner.Text()

		fields := strings.Fields(ln)

		switch fields[0] {
		case "addx":
			inc, err := strconv.ParseInt(fields[1], 10, 32)

			if err != nil {
				log.Fatal(err)
			}

			c.addx(int(inc))
		case "noop":
			c.noop()
		}
	}

	cycle := 20
	sigSum := 0

	for cycle < 221 {
		// history slice is zero-idx, but cycles are 1-idx
		registerVal := c.history[cycle-1]
		sigStrength := registerVal * cycle
		sigSum += sigStrength
		cycle += 40
	}

	return sigSum
}

func draw(h []int) {
	ln := []byte{}
	for idx, spritePos := range h {
		pxBeingDrawn := idx % 40

		var char rune

		if pxBeingDrawn >= spritePos-1 && pxBeingDrawn <= spritePos+1 {
			char = '#'
		} else {
			char = '.'
		}

		ln = append(ln, byte(char))

		if pxBeingDrawn == 39 {
			fmt.Println(string(ln))
			ln = []byte{}
		}
	}
}
