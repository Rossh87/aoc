package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("./input.txt")

	if err != nil {
		panic(err)
	}

	fmt.Println(solvePartOne(file))
}
