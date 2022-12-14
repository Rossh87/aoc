package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	file, _ := os.Open("./input.txt")

	fmt.Println(calcMonkeyBiz(file, 10000))
}

func populateMonkey(m *monkey, data string) {
	data = strings.Trim(data, " ")
	nums := strings.Split(data, ", ")

	for _, num := range nums {
		v, err := strconv.ParseInt(num, 10, 32)

		if err != nil {
			log.Fatal(err)
		}

		m.items = append(m.items, int(v))
	}
}

func stringToOp(s string) operation {
	switch s {
	case "+":
		return add
	case "*":
		return multiply
	default:
		return "*"
	}
}

type operation string

const (
	add      operation = "+"
	multiply operation = "*"
	square   operation = "^"
)

type monkey struct {
	inspectionCount int
	items           []int
	operator        operation
	operand         int
	testDivisor     int
	onTrueDest      int
	onFalseDest     int
}

// booooo.  Terrible implementation...
func (m *monkey) inspect(lcd int) (item, destination int) {
	if len(m.items) <= 0 {
		return -1, -1
	}

	(*m).inspectionCount++

	el := (m.items)[0]

	// 1. apply operation
	switch m.operator {
	case add:
		el += m.operand
	case multiply:
		el *= m.operand
	case square:
		el *= el
	}

	if lcd == -1 {
		el = el / 3
	} else {
		el = el % lcd
	}

	// 3. test and throw
	var dest int

	if el%m.testDivisor == 0 {
		dest = m.onTrueDest
	} else {
		dest = m.onFalseDest
	}

	m.items = m.items[1:]

	return el, dest
}

type monkeys []monkey

func (ms *monkeys) playRound() {
	lowestCommonDenom := 1

	for _, m := range *ms {
		lowestCommonDenom *= m.testDivisor
	}

	for idx := range *ms {
		monkey := &(*ms)[idx]
		for len(monkey.items) > 0 {
			item, dest := monkey.inspect(lowestCommonDenom)
			recipient := &(*ms)[dest]
			recipient.items = append(recipient.items, item)
		}
	}
}

func setupMonkeys(r io.Reader) monkeys {
	scanner := bufio.NewScanner(r)

	ms := monkeys{}
	currMonkey := monkey{}

	for scanner.Scan() {
		ln := scanner.Text()
		ln = strings.Trim(ln, " \t")

		switch {
		case ln == "":
			ms = append(ms, currMonkey)
			currMonkey = monkey{}
			continue

		case strings.HasPrefix(ln, "Starting"):
			populateMonkey(&currMonkey, strings.Split(ln, ":")[1])

		case strings.HasPrefix(ln, "Operation"):
			fields := strings.Fields(ln)
			currMonkey.operator = stringToOp(fields[len(fields)-2])
			rawOperand := fields[len(fields)-1]

			if rawOperand == "old" {
				currMonkey.operator = square
				currMonkey.operand = 1
			} else {
				operand, err := strconv.ParseInt(fields[len(fields)-1], 10, 32)
				if err != nil {
					log.Fatal(err)
				}
				currMonkey.operand = int(operand)
			}

		case strings.HasPrefix(ln, "Test"):
			fields := strings.Fields(ln)
			divisor, err := strconv.ParseInt(fields[len(fields)-1], 10, 32)
			if err != nil {
				log.Fatal(err)
			}
			currMonkey.testDivisor = int(divisor)

			// get 'true' dest
			scanner.Scan()
			ln = scanner.Text()
			fields = strings.Fields(ln)
			trueDest, err := strconv.ParseInt(fields[len(fields)-1], 10, 32)
			if err != nil {
				log.Fatal(err)
			}
			currMonkey.onTrueDest = int(trueDest)

			// get 'false' dest
			scanner.Scan()
			ln = scanner.Text()
			fields = strings.Fields(ln)
			falseDest, err := strconv.ParseInt(fields[len(fields)-1], 10, 32)
			if err != nil {
				log.Fatal(err)
			}
			currMonkey.onFalseDest = int(falseDest)
		}
	}
	ms = append(ms, currMonkey)
	return ms
}

func calcMonkeyBiz(r io.Reader, rounds int) int {
	ms := setupMonkeys(r)

	for i := 0; i < rounds; i++ {
		ms.playRound()
	}

	sort.Slice(ms, func(i, j int) bool {
		return ms[i].inspectionCount < ms[j].inspectionCount
	})

	return ms[len(ms)-1].inspectionCount * ms[len(ms)-2].inspectionCount
}
