package main

import (
	"AoC22/internal/inputhandler"
	"fmt"
)

func main() {

	lines := inputhandler.ReadInput()

	sopMarkerEndIdxPart1, err := findSOPMarkerEndIndex(lines[0], 4)
	if err != nil {
		fmt.Printf("Error: while searching for %d long marker: %v", 4, err)
	}

	sopMarkerEndIdxPart2, err := findSOPMarkerEndIndex(lines[0], 14)
	if err != nil {
		fmt.Printf("Error: while searching for %d long marker: %v", 14, err)
	}

	fmt.Printf("SOP marker end index - Part1: %d, Part2 %d", sopMarkerEndIdxPart1, sopMarkerEndIdxPart2)
}

func findSOPMarkerEndIndex(signal string, markerSize int) (int, error) {

	markerBuff := make([]rune, markerSize)
	currBufferIdx := 0
	for charIdx, char := range signal {

		// check if in last 4
		for markerIdx, markerChar := range markerBuff {
			if char == markerChar {
				// reset and copy the rest from the match
				tempBuff := make([]rune, markerSize)
				j := 0
				for i := markerIdx + 1; i < markerSize; i++ {
					if markerBuff[i] == 0 {
						continue
					}
					tempBuff[j] = markerBuff[i]
					j++
				}
				markerBuff = tempBuff
				currBufferIdx = j
			}
		}

		// new one, insert
		markerBuff[currBufferIdx] = char
		currBufferIdx++

		// found it!
		if currBufferIdx == markerSize {
			return charIdx + 1, nil // zero based
		}
	}

	return 0, fmt.Errorf("marker not found")
}
