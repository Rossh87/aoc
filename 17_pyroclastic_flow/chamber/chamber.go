package chamber

import (
	"math"
	"strconv"
	"strings"

	"github.com/rossh87/aoc/17_pyroclastic_flow/jets"
	"github.com/rossh87/aoc/17_pyroclastic_flow/point"
	"github.com/rossh87/aoc/17_pyroclastic_flow/rocks"
)

type occupied map[string]struct{}

func (oc occupied) Has(p point.Point) bool {
	key := p.String()
	if _, exists := oc[key]; exists {
		return true
	}
	return false
}

func (oc *occupied) add(p point.Point) {
	key := p.String()
	(*oc)[key] = struct{}{}
}

func (oc *occupied) setOccupied(r rocks.Rock) {
	points := r.List()
	for _, p := range points {
		oc.add(p)
	}
}

type Contour [7]int64

func (c *Contour) Update(r rocks.Rock) {
	// positions in grid are 1-indexed, so we account for that here
	for _, p := range r.List() {
		if p.Y > c[p.X-1] {
			c[p.X-1] = p.Y
		}
	}
}

func (c Contour) String() string {
	var b strings.Builder
	for _, el := range c {
		b.WriteString(strconv.FormatInt(int64(el), 10))
	}
	return b.String()
}

// get relative heights from absolute heights
func (c Contour) Normalized() Contour {
	var min int64 = math.MaxInt64
	for _, el := range c {
		if el < min {
			min = el
		}
	}

	for i := range c {
		c[i] = c[i] - min
	}

	return c
}

type Chamber struct {
	width      int
	occ        occupied
	pileHeight int64
	// Contour stores the highest occupied position in each column. Theoretically,
	// this represents a 'fingerprint' that, when combined with the current shape and
	// jet index, identify a state of the chamber that will repeat itself in a cycle.
	// Gratitude to https://codeberg.org/matta/advent-of-code-go/src/branch/main/2022/day17/main_test.go
	// for the idea of using the 'silhouette' of settled rocks to contribute to
	// the hash, and to https://github.com/hugseverycat/aoc2022/blob/main/day17.py
	// for inspiring the idea of tracking the silhouette as we add rocks rather than doing a BFS
	// or something more complicated.
	contour Contour
}

func New(width int) *Chamber {
	occ := make(occupied)
	contour := [7]int64{0, 0, 0, 0, 0, 0, 0}
	cham := Chamber{width, occ, 0, contour}
	return &cham
}

func (cham *Chamber) PileHeight() int64 {
	return cham.pileHeight
}

func (cham Chamber) Contour() Contour {
	return cham.contour
}

func (cham *Chamber) AddRock(rs *rocks.Rocks, js *jets.Jets) {
	rock := rs.Next(cham.pileHeight)
	for !rock.Stopped() {
		dir := js.Next()
		rock.Move(dir, cham.occ)
	}
	cham.occ.setOccupied(rock)
	cham.contour.Update(rock)
	if rock.MaxHeight() > cham.pileHeight {
		cham.pileHeight = rock.MaxHeight()
	}
}
