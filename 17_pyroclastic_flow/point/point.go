package point

import "fmt"

type Point struct {
	X int
	Y int64
}

func (p Point) String() string {
	return fmt.Sprintf("%d-%d", p.X, p.Y)
}
