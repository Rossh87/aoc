package rocks

import (
	"github.com/rossh87/aoc/17_pyroclastic_flow/jets"
	"github.com/rossh87/aoc/17_pyroclastic_flow/point"
)

type Shape int

const (
	Minus Shape = iota
	Plus
	Bracket
	Column
	Square
)

func makePoints(start point.Point, deltas []point.Point) ([]point.Point, int64) {
	yMax := start.Y
	out := make([]point.Point, len(deltas)+1)
	out[0] = start
	for i, delta := range deltas {
		dx := delta.X
		dy := delta.Y
		nx := start.X + dx
		ny := start.Y + dy
		out[i+1] = point.Point{X: nx, Y: ny}
		if ny > yMax {
			yMax = ny
		}
	}
	return out, yMax
}

func makeMinus(pileHeight int64) Rock {
	start := point.Point{X: 3, Y: pileHeight + 4}
	deltas := []point.Point{{X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}}
	points, maxHeight := makePoints(start, deltas)
	return Rock{Minus, points, false, maxHeight}
}

func makePlus(pileHeight int64) Rock {
	start := point.Point{X: 3, Y: pileHeight + 5}
	deltas := []point.Point{{X: 1, Y: 0}, {X: 2, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: -1}}
	points, maxHeight := makePoints(start, deltas)
	return Rock{Plus, points, false, maxHeight}
}

func makeBracket(pileHeight int64) Rock {
	start := point.Point{X: 3, Y: pileHeight + 4}
	deltas := []point.Point{{X: 1, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 1}, {X: 2, Y: 2}}
	points, maxHeight := makePoints(start, deltas)
	return Rock{Bracket, points, false, maxHeight}
}

func makeColumn(pileHeight int64) Rock {
	start := point.Point{X: 3, Y: pileHeight + 4}
	deltas := []point.Point{{X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}}
	points, maxHeight := makePoints(start, deltas)
	return Rock{Column, points, false, maxHeight}
}

func makeSquare(pileHeight int64) Rock {
	start := point.Point{X: 3, Y: pileHeight + 4}
	deltas := []point.Point{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 1}}
	points, maxHeight := makePoints(start, deltas)
	return Rock{Square, points, false, maxHeight}
}

type blocked interface {
	Has(p point.Point) bool
}

func none(points []point.Point, testFn func(p point.Point) bool) bool {
	for _, point := range points {
		if testFn(point) {
			return false
		}
	}

	return true
}

func leftIsBlocked(b blocked) func(p point.Point) bool {
	fn := func(p point.Point) bool {

		if p.X == 1 {
			// point is already at left edge of chamber
			return true
		}

		if b.Has(point.Point{X: p.X - 1, Y: p.Y}) {
			// left side of point is in contact with another stopped rock
			return true
		}

		return false
	}

	return fn
}

func rightIsBlocked(b blocked) func(p point.Point) bool {
	fn := func(p point.Point) bool {
		if p.X == 7 {
			// point is already at right edge of chamber
			return true
		}

		if b.Has(point.Point{X: p.X + 1, Y: p.Y}) {
			// right side of point is in contact with another stopped rock
			return true
		}

		return false
	}

	return fn
}

func bottomIsBlocked(b blocked) func(p point.Point) bool {
	fn := func(p point.Point) bool {
		if p.Y == 1 {
			// point is at bottom of chamber
			return true
		}

		if b.Has(point.Point{X: p.X, Y: p.Y - 1}) {
			// bottom of point is in contact with another stopped rock
			return true
		}

		return false
	}

	return fn
}

type Rock struct {
	kind      Shape
	points    []point.Point
	stopped   bool
	maxHeight int64
}

func (r *Rock) Move(dir jets.Push, b blocked) {
	if dir == jets.Left && none(r.points, leftIsBlocked(b)) {
		for i := range r.points {
			r.points[i].X--
		}
	}

	if dir == jets.Right && none(r.points, rightIsBlocked(b)) {
		for i := range r.points {
			r.points[i].X++
		}
	}

	// always try to move down
	if none(r.points, bottomIsBlocked(b)) {
		for i := range r.points {
			r.points[i].Y--
		}
		r.maxHeight--
		return
	}

	// if we can't move down, update rock state
	r.stopped = true
}

func (r Rock) Stopped() bool {
	return r.stopped
}

func (r Rock) List() []point.Point {
	return r.points
}

func (r Rock) MaxHeight() int64 {
	return r.maxHeight
}

type Rocks struct {
	shapes    []Shape
	nextShape int
}

func (rs Rocks) Peek() int {
	return rs.nextShape
}

func (rs *Rocks) Next(pileHeight int64) Rock {
	shape := rs.shapes[rs.nextShape]
	ns := (rs.nextShape + 1) % len(rs.shapes)
	rs.nextShape = ns
	switch shape {
	case Plus:
		return makePlus(pileHeight)
	case Minus:
		return makeMinus(pileHeight)
	case Bracket:
		return makeBracket(pileHeight)
	case Column:
		return makeColumn(pileHeight)
	case Square:
		return makeSquare(pileHeight)
	}
	panic("unrecognized shape")
}

func New() *Rocks {
	shapes := []Shape{Minus, Plus, Bracket, Column, Square}
	rs := Rocks{shapes, 0}
	return &rs
}
