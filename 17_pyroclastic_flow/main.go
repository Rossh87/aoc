package main

import "fmt"

func main() {
	fmt.Println("hello")
}

// storing each occupied space will take less room than generating the entire grid
// floor can be height 0
// we need to constantly track the max height of the simulation as rocks come to rest
// encode shapes as structs that can generate/update themselves
