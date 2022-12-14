package main

import (
	"AoC22/internal/inputhandler"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {

	lines := inputhandler.ReadInput()

	rockPaths, dimensions, err := parseScan(lines)
	if err != nil {
		fmt.Printf("Error while parsing data: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	caveSlice := NewCaveSlice(dimensions, rockPaths)
	caveSlice.Simulate(500 - dimensions.MinX)

	restedSandCount := caveSlice.countRested()

	fmt.Printf("Result - Part1: %d", restedSandCount)
}

//-----------------------------------------------------------------------------

type Position struct {
	X, Y int
}

type Dimensions struct {
	MinX, MaxX int
	MinY, MaxY int
}

func parseScan(scanData []string) ([][]Position, Dimensions, error) {

	rockPaths := make([][]Position, 0)

	var minX, minY = math.MaxInt, math.MaxInt
	var maxX, maxY int
	for rockPathIdx, rockPathLine := range scanData {

		rockCoords := strings.Split(rockPathLine, " -> ")
		if len(rockCoords) < 2 {
			return nil, Dimensions{}, fmt.Errorf("too few coordinates in line number '%d'", rockPathIdx+1)
		}

		rockPath := make([]Position, 0)

		for _, coordStr := range rockCoords {

			coord := strings.Split(coordStr, ",")
			if len(coord) < 2 {
				return nil, Dimensions{}, fmt.Errorf("invalid coordinates '%s' in line number '%d'", coordStr, rockPathIdx+1)
			}

			coordX, err := strconv.Atoi(coord[0])
			if err != nil {
				return nil, Dimensions{}, fmt.Errorf("failed to convert X coordinate to int in '%s' at line number '%d'", coordStr, rockPathIdx+1)
			}
			coordY, err := strconv.Atoi(coord[1])
			if err != nil {
				return nil, Dimensions{}, fmt.Errorf("failed to convert Y coordinate to int in '%s' at line number '%d'", coordStr, rockPathIdx+1)
			}

			pathCoord := Position{coordX, coordY}
			rockPath = append(rockPath, pathCoord)

			if pathCoord.X < minX {
				minX = pathCoord.X
			}
			if pathCoord.X > maxX {
				maxX = pathCoord.X
			}
			if pathCoord.Y < minY {
				minY = pathCoord.Y
			}
			if pathCoord.Y > maxY {
				maxY = pathCoord.Y
			}

		}

		rockPaths = append(rockPaths, rockPath)
	}

	return rockPaths, Dimensions{MinX: minX, MaxX: maxX, MinY: minY, MaxY: maxY}, nil
}

type CellType byte

const (
	Air        CellType = '.'
	Rock       CellType = '#'
	SandStatic CellType = 'o'
	SandMoving CellType = '~'
)

type CaveSlice struct {
	Field      [][]CellType
	Dimensions Dimensions
}

func NewCaveSlice(dimensions Dimensions, rockPaths [][]Position) *CaveSlice {

	caveSlice := CaveSlice{}

	caveSlice.Dimensions.MinX = 0
	caveSlice.Dimensions.MaxX = dimensions.MaxX - dimensions.MinX
	caveSlice.Dimensions.MinY = 0
	caveSlice.Dimensions.MaxY = dimensions.MaxY // leave Y so there is space for simulation

	caveSlice.Field = make([][]CellType, caveSlice.Dimensions.MaxY+1)
	for vIdx := range caveSlice.Field {
		caveSlice.Field[vIdx] = make([]CellType, caveSlice.Dimensions.MaxX+1)

		for hIdx := range caveSlice.Field[vIdx] {
			caveSlice.Field[vIdx][hIdx] = Air
		}
	}

	for _, rockPath := range rockPaths {
		for pathIdx := 0; pathIdx < len(rockPath)-1; pathIdx++ {

			startX, startY, endX, endY := GetMinMax(rockPath[pathIdx], rockPath[pathIdx+1])
			for x := startX - dimensions.MinX; x <= endX-dimensions.MinX; x++ {
				for y := startY; y <= endY; y++ {

					caveSlice.Field[y][x] = Rock
				}
			}
		}
	}

	return &caveSlice
}

func (cs *CaveSlice) Simulate(dropInPos int) {

	var iterCount int
	for {
		iterCount++

		// seed
		cs.Field[0][dropInPos] = SandMoving

		simDone := cs.doAnIteration()
		//visualize(cs)

		if simDone /*|| iterCount == 15*/ {
			break
		}
	}

	visualize(cs)
}

// returns true when simulation done (sand fall out to the void)
func (cs *CaveSlice) doAnIteration() bool {

	simDone := false
	for vIdx := cs.Dimensions.MaxY; vIdx >= cs.Dimensions.MinY; vIdx-- {
		for hIdx := cs.Dimensions.MinX; hIdx <= cs.Dimensions.MaxX; hIdx++ {

			if cs.Field[vIdx][hIdx] == SandMoving {

				cs.Field[vIdx][hIdx] = Air

				// check bottom
				if vIdx+1 > cs.Dimensions.MaxY {
					simDone = true
					continue
				}
				if cs.Field[vIdx+1][hIdx] == Air {
					cs.Field[vIdx+1][hIdx] = SandMoving
					continue
				}

				// bottom left
				if hIdx-1 < cs.Dimensions.MinX {
					simDone = true
					continue
				}
				if cs.Field[vIdx+1][hIdx-1] == Air {
					cs.Field[vIdx+1][hIdx-1] = SandMoving
					continue
				}

				// bottom right
				if hIdx+1 > cs.Dimensions.MaxX {
					simDone = true
					continue
				}
				if cs.Field[vIdx+1][hIdx+1] == Air {
					cs.Field[vIdx+1][hIdx+1] = SandMoving
					continue
				}

				// comes to rest
				cs.Field[vIdx][hIdx] = SandStatic

			}
		}
	}

	return simDone
}

func (cs *CaveSlice) countRested() int {

	var counter int
	for _, line := range cs.Field {
		for _, cell := range line {
			if cell == SandStatic {
				counter++
			}
		}
	}

	return counter
}

// returns smallerXY, biggerXY
func GetMinMax(val1 Position, val2 Position) (int, int, int, int) {

	var minX, maxX int
	if val1.X < val2.X {
		minX = val1.X
		maxX = val2.X
	} else {
		minX = val2.X
		maxX = val1.X
	}

	var minY, maxY int
	if val1.Y < val2.Y {
		minY = val1.Y
		maxY = val2.Y
	} else {
		minY = val2.Y
		maxY = val1.Y
	}

	return minX, minY, maxX, maxY
}

//-----------------------------------------------------------------------------

func visualize(caveSlice *CaveSlice) {

	for vIdx := 0; vIdx <= caveSlice.Dimensions.MaxY; vIdx++ {
		fmt.Println(string(caveSlice.Field[vIdx]))
	}

	fmt.Println()
}
