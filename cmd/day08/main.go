package main

import (
	"AoC22/internal/inputhandler"
	"fmt"
)

func main() {

	lines := inputhandler.ReadInput()

	var visibleCount int
	var highestScenicScore int

	var forestWidth = len(lines[0])
	var forestHeight = len(lines)

	// check trees 1-by-1
	for hIdx := 0; hIdx < forestWidth; hIdx++ {
		for vIdx := 0; vIdx < forestHeight; vIdx++ {

			visible, scenicScore := checkTree(hIdx, vIdx, lines)

			// Part 1
			if visible {
				visibleCount++
			}

			// Part 2
			if scenicScore > highestScenicScore {
				highestScenicScore = scenicScore
			}

		}
	}

	fmt.Printf("Result - Part1: %d, Part2: %d\n", visibleCount, highestScenicScore)
}

func checkTree(hIdx, vIdx int, forest []string) (bool, int) {

	var isVisible = false
	var sumScenicScore int

	visibility, scenicScore := checkLeftSide(hIdx, vIdx, forest)
	if visibility {
		isVisible = true
	}
	sumScenicScore = scenicScore

	visibility, scenicScore = checkRightSide(hIdx, vIdx, forest)
	if visibility {
		isVisible = true
	}
	sumScenicScore *= scenicScore

	visibility, scenicScore = checkUpSide(hIdx, vIdx, forest)
	if visibility {
		isVisible = true
	}
	sumScenicScore *= scenicScore

	visibility, scenicScore = checkDownSide(hIdx, vIdx, forest)
	if visibility {
		isVisible = true
	}
	sumScenicScore *= scenicScore

	return isVisible, sumScenicScore
}

func checkLeftSide(hIdx, vIdx int, forest []string) (bool, int) {
	if hIdx == 0 {
		return true, 0
	} else {
		for i := hIdx - 1; i >= 0; i-- {
			if forest[vIdx][i] >= forest[vIdx][hIdx] {
				return false, hIdx - i
			}
		}
	}
	return true, hIdx
}

func checkRightSide(hIdx, vIdx int, forest []string) (bool, int) {
	if hIdx == len(forest[0])-1 {
		return true, 0
	} else {
		for i := hIdx + 1; i < len(forest[vIdx]); i++ {
			if forest[vIdx][i] >= forest[vIdx][hIdx] {
				return false, i - hIdx
			}
		}
	}
	return true, len(forest[vIdx]) - 1 - hIdx
}

func checkUpSide(hIdx, vIdx int, forest []string) (bool, int) {
	if vIdx == 0 {
		return true, 0
	} else {
		for i := vIdx - 1; i >= 0; i-- {
			if forest[i][hIdx] >= forest[vIdx][hIdx] {
				return false, vIdx - i
			}
		}
	}
	return true, vIdx
}

func checkDownSide(hIdx, vIdx int, forest []string) (bool, int) {
	if vIdx == len(forest)-1 {
		return true, 0
	} else {
		for i := vIdx + 1; i < len(forest); i++ {
			if forest[i][hIdx] >= forest[vIdx][hIdx] {
				return false, i - vIdx
			}
		}
	}
	return true, len(forest) - 1 - vIdx
}
