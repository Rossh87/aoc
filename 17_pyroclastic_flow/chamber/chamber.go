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

type Contour [7]int

func (c *Contour) Update(r rocks.Rock) {
	// positions in grid are 1-indexed, so we account for that here
	for _, p := range r.List() {
		c[p.X-1] = intMax(c[p.X-1], p.Y)
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
	min := math.MaxInt
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
	pileHeight int
	// contour stores the uppermost occupied position in each column. Theoretically,
	// this represents a 'fingerprint' that, when combined with the current shape and
	// jet index, form the hash of a cycle that will repeat
	contour Contour
}

func New(width int) *Chamber {
	occ := make(occupied)
	contour := [7]int{-1, -1, -1, -1, -1, -1, -1}
	cham := Chamber{width, occ, 0, contour}
	return &cham
}

func intMax(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func (cham *Chamber) PileHeight() int {
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
	max := intMax(cham.pileHeight, rock.MaxHeight())
	cham.pileHeight = max
}
