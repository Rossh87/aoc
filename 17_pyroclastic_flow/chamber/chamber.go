package chamber

import (
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

type Chamber struct {
	width      int
	occ        occupied
	pileHeight int
}

func New(width int) *Chamber {
	occ := make(occupied)
	cham := Chamber{width, occ, 0}
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

func (cham *Chamber) AddRock(rs *rocks.Rocks, js *jets.Jets) {
	rock := rs.Next(cham.pileHeight)
	for !rock.Stopped() {
		dir := js.Next()
		rock.Move(dir, cham.occ)
	}
	cham.occ.setOccupied(rock)
	max := intMax(cham.pileHeight, rock.MaxHeight())
	cham.pileHeight = max
}
