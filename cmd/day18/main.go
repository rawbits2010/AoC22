package main

import (
	"AoC22/internal/inputhandler"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {

	lines := inputhandler.ReadInput()
	grid, err := create3DGridFrom(lines)
	if err != nil {
		fmt.Printf("Error: couldn't create grid: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	exposedSides := countExposedSides(grid)

	fmt.Printf("Result - Part1: %d", exposedSides)
}

//-----------------------------------------------------------------------------

func countExposedSides(grid [][][]bool) int {

	playField := *NewPlayField(grid)

	var count int
	for colIdx, col := range grid {
		for rowIdx, row := range col {
			for depthIdx, cell := range row {

				// not a stone check
				if !cell {
					continue
				}

				// only check side if exposed - we hit out of bounds in BFS

				// bottom
				if depthIdx != 0 {
					if !grid[colIdx][rowIdx][depthIdx-1] { // only check air
						if _, ok := BFSUntil(playField, colIdx, rowIdx, depthIdx-1); ok {
							count++
						} else {
							// this is an enclosed air pocket
						}
					}
				} else {
					count++
				}

				// top
				if depthIdx != len(grid[colIdx][rowIdx])-1 {
					if !grid[colIdx][rowIdx][depthIdx+1] {
						if _, ok := BFSUntil(playField, colIdx, rowIdx, depthIdx+1); ok {
							count++
						}
					}
				} else {
					count++
				}

				// y sides
				if rowIdx != 0 {
					if !grid[colIdx][rowIdx-1][depthIdx] {
						if _, ok := BFSUntil(playField, colIdx, rowIdx-1, depthIdx); ok {
							count++
						}
					}
				} else {
					count++
				}
				if rowIdx != len(grid[colIdx])-1 {
					if !grid[colIdx][rowIdx+1][depthIdx] {
						if _, ok := BFSUntil(playField, colIdx, rowIdx+1, depthIdx); ok {
							count++
						}
					}
				} else {
					count++
				}

				// x sides
				if colIdx != 0 {
					if !grid[colIdx-1][rowIdx][depthIdx] {
						if _, ok := BFSUntil(playField, colIdx-1, rowIdx, depthIdx); ok {
							count++
						}
					}
				} else {
					count++
				}
				if colIdx != len(grid)-1 {
					if !grid[colIdx+1][rowIdx][depthIdx] {
						if _, ok := BFSUntil(playField, colIdx+1, rowIdx, depthIdx); ok {
							count++
						}
					}
				} else {
					count++
				}
			}
		}
	}

	return count
}

//-BFS-------------------------------------------------------------------------

type PlayField struct {
	Field [][][]bool
}

func NewPlayField(grid [][][]bool) *PlayField {
	return &PlayField{
		Field: grid,
	}
}

func (pf *PlayField) isInBound(x, y, z int) bool {

	if x < 0 || x >= len(pf.Field) {
		return false
	}
	if y < 0 || y >= len(pf.Field[0]) {
		return false
	}
	if z < 0 || z >= len(pf.Field[0][0]) {
		return false
	}
	return true
}

type BFSTile struct {
	X, Y, Z int
	Checked bool
	Prev    *BFSTile
}

// return: destination, found/not
func BFSUntil(playField PlayField, startX, startY, startZ int) (BFSTile, bool) {
	tilesToCheck := make([]BFSTile, 0)
	tilesToCheck = append(tilesToCheck, BFSTile{X: startX, Y: startY, Z: startZ, Prev: nil})

	tilesChecked := make([]BFSTile, 0)

	var currTile BFSTile
	for len(tilesToCheck) > 0 {

		currTile, tilesToCheck = tilesToCheck[0], tilesToCheck[1:]

		if _, ok := isSliceContains(currTile, tilesChecked); ok {
			continue
		}

		tilesChecked = append(tilesChecked, currTile)

		neighbours := [6][3]int{{1, 0, 0}, {-1, 0, 0}, {0, 1, 0}, {0, -1, 0}, {0, 0, 1}, {0, 0, -1}}
		for _, delta := range neighbours {

			posX := currTile.X + delta[0]
			posY := currTile.Y + delta[1]
			posZ := currTile.Z + delta[2]

			// stop if we found an exit
			if !playField.isInBound(posX, posY, posZ) {
				return currTile, true
			}

			// only check if air
			if !playField.Field[posX][posY][posZ] {

				tile := BFSTile{
					X:    posX,
					Y:    posY,
					Z:    posZ,
					Prev: &currTile,
				}

				tilesToCheck = append(tilesToCheck, tile)
			}
		}

	}

	return BFSTile{}, false
}

func isSliceContains(val BFSTile, list []BFSTile) (int, bool) {
	for i := range list {
		if list[i].X == val.X && list[i].Y == val.Y && list[i].Z == val.Z {
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

//-----------------------------------------------------------------------------

func create3DGridFrom(lines []string) ([][][]bool, error) {

	coordsList := make([][]int, 0, 50)

	var dims Dimensions3D
	for lineIdx, line := range lines {

		if len(line) == 0 {
			continue
		}

		coords := strings.Split(line, ",")
		if len(coords) < 3 {
			return nil, fmt.Errorf("too few params for a coord in '%s' at line '%d'", line, lineIdx)
		}

		x, err := strconv.Atoi(coords[0])
		if err != nil {
			return nil, fmt.Errorf("couldn't parse X coord from '%s' at line '%d'", coords[0], lineIdx)
		}

		y, err := strconv.Atoi(coords[1])
		if err != nil {
			return nil, fmt.Errorf("couldn't parse Y coord from '%s' at line '%d'", coords[1], lineIdx)
		}

		z, err := strconv.Atoi(coords[2])
		if err != nil {
			return nil, fmt.Errorf("couldn't parse Z coord from '%s' at line '%d'", coords[2], lineIdx)
		}

		dims.Update(x, y, z)

		coordsList = append(coordsList, []int{x, y, z})
	}

	grid := make([][][]bool, (dims.MaxX-dims.MinX)+1)
	for xIdx := range grid {
		grid[xIdx] = make([][]bool, (dims.MaxY-dims.MinY)+1)
		for yIdx := range grid[xIdx] {
			grid[xIdx][yIdx] = make([]bool, (dims.MaxZ-dims.MinZ)+1)
		}
	}

	for _, coord := range coordsList {
		grid[coord[0]][coord[1]][coord[2]] = true
	}

	return grid, nil
}

type Dimensions3D struct {
	MinX, MaxX int
	MinY, MaxY int
	MinZ, MaxZ int
}

func NewDimensions3D() *Dimensions3D {
	return &Dimensions3D{
		MinX: math.MaxInt,
		MaxX: 0,
		MinY: math.MaxInt,
		MaxY: 0,
		MinZ: math.MaxInt,
		MaxZ: 0,
	}
}

func (d *Dimensions3D) Update(x, y, z int) {
	if d.MinX > x {
		d.MinX = x
	}
	if d.MaxX < x {
		d.MaxX = x
	}
	if d.MinY > y {
		d.MinY = y
	}
	if d.MaxY < y {
		d.MaxY = y
	}
	if d.MinZ > z {
		d.MinZ = z
	}
	if d.MaxZ < z {
		d.MaxZ = z
	}
}
