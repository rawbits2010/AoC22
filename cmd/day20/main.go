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

func main() {

	lines := inputhandler.ReadInput()

	coordList, err := parseCoords(lines)
	if err != nil {
		fmt.Printf("Error parsing input: %v\n", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	mixedCoords := mix(coordList)
	zeroIdx, found := getIdxOf(0, mixedCoords)
	if !found {
		fmt.Printf("Error: couldn't find starting index '0' in '%v'\n", mixedCoords)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	result := sumCoords([]int{1000, 2000, 3000}, zeroIdx, mixedCoords)

	fmt.Printf("Result: %d", result)
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
		result += mixedCoords[(pos+zeroIdx)%len(mixedCoords)].Value
	}

	return result
}

func mix(coords []Coord) []Coord {

	for currCoordIdx := 0; currCoordIdx < len(coords); currCoordIdx++ {

		newCoords := make([]Coord, 0)
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

			newCoords = append(newCoords, ringCoords[:moveToIdx]...)
			newCoords = append(newCoords, coord)
			newCoords = append(newCoords, ringCoords[moveToIdx:]...)

			//fmt.Printf("moving: '%d' - %v\n", coord.Value, newCoords)
			break
		}

		coords = newCoords[:]
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
