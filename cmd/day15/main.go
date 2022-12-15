package main

import (
	"AoC22/internal/inputhandler"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

func main() {

	lines := inputhandler.ReadInput()

	sensors, dimensions, err := parseSensorData(lines)
	if err != nil {
		fmt.Printf("Error while parsing data: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	resultPart1 := countNoBeaconPosOnRow(2000000, sensors, dimensions)

	checkArea := Dimensions{
		MinX: 0,
		MaxX: 4000000,
		MinY: 0,
		MaxY: 4000000,
	}
	resultPart2, err := getFreqOfFirstPossibleBeaconPos(sensors, checkArea)
	if err != nil {
		fmt.Printf("Error while processing part 2: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	fmt.Printf("Result - Part1: %d, Part2: %d", resultPart1, resultPart2)
}

//-----------------------------------------------------------------------------

type Dimensions struct {
	MinX, MaxX int
	MinY, MaxY int
}

func NewDimensions() *Dimensions {
	return &Dimensions{
		MinX: math.MaxInt,
		MaxX: 0,
		MinY: math.MaxInt,
		MaxY: 0,
	}
}

func (d *Dimensions) Update(x, y int) {

	if x < d.MinX {
		d.MinX = x
	}
	if x > d.MaxX {
		d.MaxX = x
	}
	if y < d.MinY {
		d.MinY = y
	}
	if y > d.MaxY {
		d.MaxY = y
	}
}

type Position struct {
	X, Y int
}

func (p *Position) DistanceFrom(val Position) int {
	return CalcDistance(*p, val)
}

type Beacon struct {
	Position
}

func NewBeacon(x, y int) *Beacon {
	return &Beacon{
		Position: Position{X: x, Y: y},
	}
}

type Sensor struct {
	Position
	closestBeacon  Beacon
	BeaconDistance int
}

func NewSensor(x, y int) *Sensor {
	return &Sensor{
		Position: Position{X: x, Y: y},
	}
}

func (s *Sensor) SetClosestBeacon(beacon Beacon) {
	s.closestBeacon = beacon
	s.BeaconDistance = s.DistanceFrom(beacon.Position)
}

func parseSensorData(lines []string) ([]Sensor, Dimensions, error) {

	coordsPattern, err := regexp.Compile(`x=(-?\d+)|y=(-?\d+)`)
	if err != nil {
		// shouldn't be possible
		return nil, Dimensions{}, fmt.Errorf("couldn't compile coordinate parser regex")
	}

	dimensions := *NewDimensions()

	sensors := make([]Sensor, len(lines))
	for lineIdx, line := range lines {

		coords := coordsPattern.FindAllStringSubmatch(line, -1)
		if coords == nil || len(coords) < 4 {
			return nil, Dimensions{}, fmt.Errorf("too few coordinates found '%v' at line '%d'", coords, lineIdx)
		}

		senX, err := strconv.Atoi(coords[0][1])
		if err != nil {
			return nil, Dimensions{}, fmt.Errorf("couldn't parse sensor's X coordinate from '%s' at line '%d'", coords[0][0], lineIdx)
		}

		senY, err := strconv.Atoi(coords[1][2])
		if err != nil {
			return nil, Dimensions{}, fmt.Errorf("couldn't parse sensor's Y coordinate from '%s' at line '%d'", coords[1][0], lineIdx)
		}

		beacX, err := strconv.Atoi(coords[2][1])
		if err != nil {
			return nil, Dimensions{}, fmt.Errorf("couldn't parse beacon's X coordinate from '%s' at line '%d'", coords[2][0], lineIdx)
		}

		beacY, err := strconv.Atoi(coords[3][2])
		if err != nil {
			return nil, Dimensions{}, fmt.Errorf("couldn't parse beacon's Y coordinate from '%s' at line '%d'", coords[3][0], lineIdx)
		}

		dimensions.Update(senX, senY)
		dimensions.Update(beacX, beacY)

		sensor := NewSensor(senX, senY)
		beacon := NewBeacon(beacX, beacY)
		sensor.SetClosestBeacon(*beacon)

		sensors = append(sensors, *sensor)
	}

	return sensors, dimensions, nil
}

//-----------------------------------------------------------------------------

// Part 1
func countNoBeaconPosOnRow(row int, sensors []Sensor, dimensions Dimensions) int {
	//fmt.Printf("minX: %d, maxX: %d\n", dimensions.MinX, dimensions.MaxX)

	var count int

	canBeBeacon := func(position Position) bool {

		for _, sensor := range sensors {

			// exclude existing beacons, duh
			if position == sensor.closestBeacon.Position {
				return true
			}

			if sensor.DistanceFrom(position) <= sensor.BeaconDistance {
				return false
			}
		}

		return true
	}

	// middle out algo - if u got the ref :)
	//var line string
	middleIdx := (dimensions.MaxX - dimensions.MinX) / 2

	leftIdx := middleIdx
	for {

		if canBeBeacon(Position{X: leftIdx, Y: row}) {
			//line = "." + line

			if leftIdx < dimensions.MinX {
				break
			}
		} else {
			//line = "#" + line

			count++
		}

		leftIdx--
	}

	rightIdx := middleIdx + 1
	for {

		if canBeBeacon(Position{X: rightIdx, Y: row}) {
			//line = line + "."

			if rightIdx > dimensions.MaxX {
				break
			}
		} else {
			//line = line + "#"

			count++
		}

		rightIdx++
	}
	//fmt.Printf("counted from: %d to: %d\n", leftIdx, rightIdx)

	//fmt.Println(line)

	return count
}

// Part 2
func getFreqOfFirstPossibleBeaconPos(sensors []Sensor, dimensions Dimensions) (int, error) {

	for vIdx := dimensions.MinY; vIdx <= dimensions.MaxY; vIdx++ {
	nextLine:
		for hIdx := dimensions.MinX; hIdx <= dimensions.MaxX; hIdx++ {

			for _, sensor := range sensors {
				if sensor.DistanceFrom(Position{X: hIdx, Y: vIdx}) <= sensor.BeaconDistance {

					// skip this sensor checked area
					hIdx += sensor.BeaconDistance - GetAbs(sensor.Y-vIdx) + (sensor.X - hIdx)

					continue nextLine
				}
			}

			return hIdx*4000000 + vIdx, nil
		}
	}

	return 0, fmt.Errorf("no solution found")
}

//-----------------------------------------------------------------------------

func CalcDistance(pos1, pos2 Position) int {
	return GetAbs(pos1.X-pos2.X) + GetAbs(pos1.Y-pos2.Y)
}

func GetAbs(val int) int {
	if val < 0 {
		return -val
	}
	return val
}
