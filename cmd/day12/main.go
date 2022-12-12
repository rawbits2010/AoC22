package main

import (
	"AoC22/internal/inputhandler"
	"AoC22/internal/outputhandler"
	"fmt"
	"math"
	"sort"
)

var pathColor string
var mapColor string

func visualizePath(steps []Location, playfield PlayField) {
	linesHeightMap := make([]string, playfield.Height)
	linesSteps := make([]string, playfield.Height)
	for vIdx := range linesSteps {

		var lineHeightMap string
		lineSteps := make([]byte, playfield.Width)

		var lastPosIsOnPath int = -1
		for hIdx := range lineSteps {

			var isOnPath = 0
			var currStepIdx int
			var currStep Location
			for stepIdx, step := range steps {
				if step.x == hIdx && step.y == vIdx {
					currStepIdx = stepIdx
					currStep = step
					isOnPath = 1
					break
				}
			}

			// on height map
			currRune := byte(playfield.getHeightAt(hIdx, vIdx))
			if isOnPath != lastPosIsOnPath {
				lastPosIsOnPath = isOnPath
				if isOnPath == 1 {
					lineHeightMap += pathColor
				} else {
					lineHeightMap += mapColor
				}
			}
			lineHeightMap += string(currRune)

			// steps took
			if isOnPath == 0 {
				lineSteps[hIdx] = '.'
			} else {
				if currStepIdx < len(steps)-1 {
					if currStep.x < steps[currStepIdx+1].x && currStep.y == steps[currStepIdx+1].y {
						lineSteps[hIdx] = '>'
					} else {
						if currStep.x > steps[currStepIdx+1].x && currStep.y == steps[currStepIdx+1].y {
							lineSteps[hIdx] = '<'
						} else {
							if currStep.x == steps[currStepIdx+1].x && currStep.y < steps[currStepIdx+1].y {
								lineSteps[hIdx] = 'v'
							} else {
								if currStep.x == steps[currStepIdx+1].x && currStep.y > steps[currStepIdx+1].y {
									lineSteps[hIdx] = '^'
								} else {
									lineSteps[hIdx] = 'O' // shouldn't be possible
								}
							}
						}
					}

				} else {
					lineSteps[hIdx] = 'E'
				}
			}
		}
		linesHeightMap[vIdx] = string(lineHeightMap)
		linesSteps[vIdx] = string(lineSteps)
	}

	for _, line := range linesHeightMap {
		fmt.Println(line)
	}
	fmt.Println()
	for _, line := range linesSteps {
		fmt.Println(line)
	}

}

func main() {

	outputhandler.Initialize()
	defer outputhandler.Reset()
	pathColor = outputhandler.GetForeground(outputhandler.BrightGreen)
	mapColor = outputhandler.GetReset()

	lines := inputhandler.ReadInput()
	playField, start, goal := parseInput(lines)

	// Part 1
	stepsPart1, found := pathFind(start, goal, *playField)
	if !found {
		panic("no solution")
	}
	//ReverseSlice(stepsPart1)
	//visualizePath(stepsPart1, *playField)

	// Part 2
	stepsListPart2 := make([]int, 0)
	for vIdx, heightMapLine := range playField.heightMap {
		for hIdx, height := range heightMapLine {

			if height != int('a') {
				continue
			}

			steps, found := pathFind(Location{x: hIdx, y: vIdx}, Location{x: goal.x, y: goal.y}, *playField)
			if !found {
				//panic("no solution")
				continue
			}

			//ReverseSlice(steps)
			//visualizePath(steps, *playField)

			stepsListPart2 = append(stepsListPart2, len(steps))
		}
	}
	sort.Ints(stepsListPart2)

	fmt.Printf("Result - Part1: %d, Part2: %d", len(stepsPart1), stepsListPart2[0])
}

func parseInput(lines []string) (*PlayField, Location, Location) {
	heightMap := make([][]int, 0, len(lines))
	var startPos Location
	var goalPos Location
	for vIdx, line := range lines {
		row := make([]int, len(line))
		for hIdx, val := range line {
			switch val {
			case 'S':
				startPos = Location{x: hIdx, y: vIdx}
				val = 'a'
			case 'E':
				goalPos = Location{x: hIdx, y: vIdx}
				val = 'z'
			}
			row[hIdx] = int(val)
		}
		heightMap = append(heightMap, row)
	}

	return NewPlayField(heightMap), startPos, goalPos
}

type Position struct {
	X, Y int
}

type PlayField struct {
	heightMap [][]int
	Width     int
	Height    int
}

func NewPlayField(heightMap [][]int) *PlayField {
	return &PlayField{
		heightMap: heightMap,
		Width:     len(heightMap[0]),
		Height:    len(heightMap),
	}
}

func (pf *PlayField) getHeightAt(x, y int) int {
	return pf.heightMap[y][x]
}

func filterMovableTiles(tiles []Location, currLocation Location, playField PlayField) []Location {

	var temp []Location
	for i := 0; i < len(tiles); i++ {

		if playField.getHeightAt(tiles[i].x, tiles[i].y) > playField.getHeightAt(currLocation.x, currLocation.y)+1 {
			continue
		}

		temp = append(temp, tiles[i])
	}

	return temp
}

//-A*--------------------------------------------------------------------------

func pathFind(unitLocation, targetLocation Location, playfield PlayField) ([]Location, bool) {

	var tilesToCheck []Location
	var tilesChecked []Location

	tilesToCheck = append(tilesToCheck, unitLocation)

	for {
		if len(tilesToCheck) <= 0 {
			break
		}

		sort.Slice(tilesToCheck, func(i, j int) bool {
			return tilesToCheck[i].aStarVals.getF() < tilesToCheck[j].aStarVals.getF()
		})
		var currentTile = tilesToCheck[0]

		if currentTile.x == targetLocation.x && currentTile.y == targetLocation.y {
			var currTile = &currentTile
			var result []Location
			for {
				if currTile.x == unitLocation.x && currTile.y == unitLocation.y {
					break
				}
				result = append(result, *currTile)

				currTile = currTile.aStarVals.prev
			}

			return result, true
		}

		tilesToCheck = tilesToCheck[1:]
		tilesChecked = append(tilesChecked, currentTile)

		var tilesAround = getTilesAround(currentTile, playfield)
		tilesAround = filterMovableTiles(tilesAround, currentTile, playfield)
		for idx := range tilesAround {

			if _, ok := isSliceContains(tilesAround[idx], tilesChecked); ok {
				continue
			}

			var curr *Location
			if i, ok := isSliceContains(tilesAround[idx], tilesToCheck); ok {
				curr = &tilesToCheck[i]
			} else {
				tilesToCheck = append(tilesToCheck, tilesAround[idx])
				curr = &tilesToCheck[len(tilesToCheck)-1]
			}

			tg := currentTile.aStarVals.g + calcDistance(currentTile, tilesAround[idx])
			if tg > tilesAround[idx].aStarVals.g {
				curr.aStarVals.g = tg
				curr.aStarVals.h = calcDistance(targetLocation, tilesAround[idx])

				curr.aStarVals.prev = &currentTile
			}

		}

	}

	return []Location{}, false
}

type Location struct {
	x int
	y int

	aStarVals AStar
}

type AStar struct {
	g float64
	h float64

	prev *Location
}

func (a AStar) getF() float64 {
	return a.g + a.h
}

func getTilesAround(location Location, playfield PlayField) []Location {

	var temp []Location

	if location.x > 0 {
		temp = append(temp, Location{x: location.x - 1, y: location.y})
	}

	if location.x < playfield.Width-1 {
		temp = append(temp, Location{x: location.x + 1, y: location.y})
	}

	if location.y > 0 {
		temp = append(temp, Location{x: location.x, y: location.y - 1})
	}

	if location.y < playfield.Height-1 {
		temp = append(temp, Location{x: location.x, y: location.y + 1})
	}

	return temp
}

func calcDistance(unitLocation Location, targetLocation Location) float64 {
	dX := unitLocation.x - targetLocation.x
	dY := unitLocation.y - targetLocation.y
	return math.Sqrt(float64(dX*dX + dY*dY))
}

//-Utils-----------------------------------------------------------------------

func isSliceContains(val Location, list []Location) (int, bool) {
	for i := range list {
		if list[i].x == val.x && list[i].y == val.y {
			return i, true
		}
	}
	return 0, false
}

func ReverseSlice[T comparable](s []T) {
	sort.SliceStable(s, func(i, j int) bool {
		return i > j
	})
}
