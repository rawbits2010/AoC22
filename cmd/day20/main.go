package main

import (
	"AoC22/internal/inputhandler"
	"fmt"
	"os"
	"strconv"
)

type Coord struct {
	Value int
	Idx   int
}

func MulInt(val1, val2 int) int {
	res := val1 * val2
	if (res < 0) == ((val1 < 0) != (val2 < 0)) {
		if res/val2 == val1 {
			return res
		}
	}
	panic("multiplication overflow")
}

func main() {

	lines := inputhandler.ReadInput()

	coordList, err := parseCoords(lines)
	if err != nil {
		fmt.Printf("Error parsing input: %v\n", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	// Part 1

	mixedCoords := mix(append([]Coord{}, coordList...))

	zeroIdx, found := getIdxOf(0, mixedCoords)
	if !found {
		fmt.Printf("Error: couldn't find starting index '0' in '%v'\n", mixedCoords)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	resultPart1 := sumCoords([]int{1000, 2000, 3000}, zeroIdx, mixedCoords)

	// Part 2

	moddedCoordList := make([]Coord, len(coordList))
	for coordIdx, coord := range coordList {
		moddedCoordList[coordIdx] = Coord{Value: MulInt(coord.Value, 811589153), Idx: coord.Idx}
	}

	mixedCoords = moddedCoordList
	for idx := 0; idx < 10; idx++ {
		mixedCoords = mix(mixedCoords)
	}

	zeroIdx, found = getIdxOf(0, mixedCoords)
	if !found {
		fmt.Printf("Error: couldn't find starting index '0' in '%v'\n", mixedCoords)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	resultPart2 := sumCoords([]int{1000, 2000, 3000}, zeroIdx, mixedCoords)

	fmt.Printf("Result - Part1: %d, Part2: %d", resultPart1, resultPart2)
}

//-----------------------------------------------------------------------------

func parseCoords(lines []string) ([]Coord, error) {

	coords := make([]Coord, 0, len(lines))
	for lineIdx, line := range lines {

		if len(line) == 0 {
			continue
		}

		val, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("couldn't convert '%s' at line '%d'", line, lineIdx)
		}

		coords = append(coords, Coord{Value: val, Idx: lineIdx})
	}

	return coords, nil
}

//-----------------------------------------------------------------------------

func sumCoords(posList []int, zeroIdx int, mixedCoords []Coord) int {

	var result int
	for _, pos := range posList {
		cordAtPos := mixedCoords[(pos+zeroIdx)%len(mixedCoords)].Value
		//fmt.Printf("%dth - %d\n", pos, cordAtPos)

		result += cordAtPos
	}

	return result
}

func mix(coords []Coord) []Coord {

	for currCoordIdx := 0; currCoordIdx < len(coords); currCoordIdx++ {

		var newCoords []Coord
		for idx, coord := range coords {
			if coord.Idx != currCoordIdx {
				continue
			}

			ringCoords := append(coords[:idx], coords[idx+1:]...)

			var moveToIdx int
			if coord.Value > 0 {
				moveToIdx = (idx + (coord.Value % len(ringCoords))) % len(ringCoords)
			} else if coord.Value < 0 {
				moveToIdx = len(ringCoords) - (((len(ringCoords) - idx) - (coord.Value % len(ringCoords))) % len(ringCoords))
			} else {
				moveToIdx = idx
			}

			newCoords = append([]Coord{}, ringCoords[:moveToIdx]...)
			newCoords = append(newCoords, coord)
			newCoords = append(newCoords, ringCoords[moveToIdx:]...)

			//fmt.Printf("moving: '%d' - %v\n", coord.Value, newCoords)
			break
		}

		coords = newCoords
	}

	return coords
}

func getIdxOf(coord int, coordList []Coord) (int, bool) {
	for idx, currCoord := range coordList {
		if currCoord.Value == coord {
			return idx, true
		}
	}
	return 0, false
}
