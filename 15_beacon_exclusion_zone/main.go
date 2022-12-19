package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/samber/lo"
)

func main() {
	file, err := os.Open("./input.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Println(findDistressFreq(file, 0, 4e6))
}

type Point interface {
	Stringer
	X() int
	Y() int
}

type Stringer fmt.Stringer

func intAbs(a, b int) int {
	diff := a - b

	if diff < 0 {
		return diff * -1
	}

	return diff
}

func manhattanDist(p1, p2 Point) int {
	dx := intAbs(p1.X(), p2.X())
	dy := intAbs(p1.Y(), p2.Y())
	return dx + dy
}

type Sensor interface {
	Point
	Sense(p Point) bool
	Reach() int
}

type sensor struct {
	maxDist int
	x       int
	y       int
	beacon  Point
}

func (s sensor) Reach() int {
	return s.maxDist
}

func (s sensor) X() int {
	return s.x
}

func (s sensor) Y() int {
	return s.y
}

func (s sensor) String() string {
	return fmt.Sprintf("%d-%d", s.x, s.y)
}

func (s sensor) Sense(p Point) bool {
	distToCandidate := manhattanDist(s, p)
	return distToCandidate <= s.maxDist
}

func newSensor(vs []int) *sensor {
	s := sensor{}
	b := beacon{vs[2], vs[3]}
	s.x = vs[0]
	s.y = vs[1]
	s.maxDist = manhattanDist(s, b)
	s.beacon = b
	return &s
}

type beacon struct {
	x int
	y int
}

func (b beacon) X() int {
	return b.x
}

func (b beacon) Y() int {
	return b.y
}

func (b beacon) String() string {
	return fmt.Sprintf("%d-%d", b.x, b.y)
}

func parseLine(line string) *sensor {
	rexp := regexp.MustCompile(`(?:x|y)=(-?\d+)`)
	vals := rexp.FindAllStringSubmatch(line, -1)
	flatVals := lo.Map(vals, func(item []string, idx int) int {
		sVal := item[1]
		intVal, err := strconv.Atoi(sVal)
		if err != nil {
			log.Fatalf("failed parsing value %v to integer", sVal)
		}
		return intVal
	})
	s := newSensor(flatVals)
	return s
}

func intMin(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func intMax(a, b int) int {
	if a > b {
		return a
	}

	return b
}

type StringSet interface {
	Has(s Stringer) bool
	Set(s Stringer)
	Delete(s Stringer)
}

type stringSet struct {
	m map[string]struct{}
}

func (se stringSet) Has(s Stringer) bool {
	k := s.String()
	if _, ok := se.m[k]; ok {
		return true
	}
	return false
}

func (se *stringSet) Set(s Stringer) {
	k := s.String()
	se.m[k] = struct{}{}
}

func (se *stringSet) Delete(s Stringer) {
	delete((*se).m, s.String())
}

func newstringSet() *stringSet {
	s := stringSet{}
	s.m = make(map[string]struct{})
	return &s
}

type gridPoint struct {
	x int
	y int
}

func (g gridPoint) X() int {
	return g.x
}

func (g gridPoint) Y() int {
	return g.y
}

func (g gridPoint) String() string {
	return fmt.Sprintf("%d-%d", g.x, g.y)
}

func countSensedInRow(rowNumber int, r io.Reader) int64 {
	scanner := bufio.NewScanner(r)
	var xMin *int
	var xMax *int
	sensors := []*sensor{}
	beacons := newstringSet()
	var maxDist int
	for scanner.Scan() {
		line := scanner.Text()
		// populate a sensor object
		sensor := parseLine(line)
		maxDist = intMax(maxDist, sensor.maxDist)

		// add sensor to list of sensors
		sensors = append(sensors, sensor)

		// update grid bounds based on sensor range
		lowX := sensor.X() - sensor.maxDist
		highX := sensor.X() + sensor.maxDist

		if xMin == nil {
			xMin = &lowX
		} else {
			*xMin = intMin(lowX, *xMin)
		}

		if xMax == nil {
			xMax = &highX
		} else {
			*xMax = intMax(highX, *xMax)
		}

		// add sensed beacon to list of beacons so we don't count that
		// position as a cell that CANNOT have a beacon.
		beacons.Set(sensor.beacon)
	}

	sensed := newstringSet()
	var sensedCount int64
	for i := *xMin; i <= *xMax; i++ {
		point := gridPoint{i, rowNumber}
		for _, sensor := range sensors {
			if sensor.Sense(point) && !sensed.Has(point) && !beacons.Has(point) {
				sensed.Set(point)
				sensedCount++
			}
		}
	}
	return sensedCount
}

func checkDistances(currIdx int, ss []Sensor, point Point) bool {
	for i, s := range ss {
		if i == currIdx {
			continue
		}

		if manhattanDist(point, s) <= s.Reach() {
			return false
		}
	}

	return true
}

func inbounds(p Point, floor, ceil int) bool {
	return p.X() >= floor && p.X() <= ceil && p.Y() >= floor && p.Y() <= ceil
}

// We can test if a point is in range of a sensor by comparing the manhattan distance
// from the sensor to point. Any point whose x and y are inbounds, AND that is not
// in range of any sensor, must be the answer. However, a square grid with 4e6
// cells per side is too big to manually search.
// Since there is exactly one solution inside the search space, we can infer
// that the answer MUST be adjacent to the area covered by one of the sensors.
// So we can reduce the search space by ONLY checking points that are adjacent
// to the 'frontier' of each sensor. Any cell that is on the frontier of one sensor,
// AND is out-of-range for all the other sensors, is the answer.
func checkSensors(ss []Sensor, idx int, floor, ceil int) *gridPoint {
	// If this block runs, there is no solution for the given input.
	if idx >= len(ss) {
		return nil
	}

	sensor := ss[idx]

	// Imagine that all the points covered by a given sensor is a clock.
	// 'boundingVal' is the length of the clock hand.
	// Check every point that touches the circumference
	// of the circle covered by the clock hand in one revolution.
	boundingVal := sensor.Reach() + 1

	for xOffset := boundingVal; xOffset >= boundingVal*-1; xOffset-- {
		xPos := sensor.X() + xOffset
		// yOffset changes inversely with xPos: as xPos gets closer
		// to the center of the clock, yOffset gets bigger.
		yOffset := boundingVal - intAbs(sensor.X(), xPos)
		// for each value of xPos, check a point on the top half of
		// the circumference, and its reflection over the x-axis on the
		// bottom half of the circumference.
		up := gridPoint{xPos, sensor.Y() + yOffset}
		down := gridPoint{xPos, sensor.Y() - yOffset}

		// as soon as we find a perimeter cell that is within bounds, and is out-of-reach
		// for all sensors, we have the solution.
		if inbounds(up, floor, ceil) && checkDistances(idx, ss, up) {
			return &up
		}

		if inbounds(down, floor, ceil) && checkDistances(idx, ss, down) {
			return &down
		}
	}

	// otherwise, check the perimeter cells of the next sensor
	return checkSensors(ss, idx+1, floor, ceil)
}

func calcFreq(p Point) int64 {
	return int64(4e6)*int64(p.X()) + int64(p.Y())
}

func findDistressFreq(r io.Reader, floor, ceil int) int64 {
	scanner := bufio.NewScanner(r)
	sensors := []Sensor{}
	for scanner.Scan() {
		line := scanner.Text()
		sensor := parseLine(line)
		sensors = append(sensors, *sensor)
	}
	ans := checkSensors(sensors, 0, floor, ceil)

	if ans == nil {
		log.Fatal("Got nil as answer")
	}

	return calcFreq(*ans)
}
