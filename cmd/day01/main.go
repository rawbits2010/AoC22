package main

import (
	"AoC22/internal/inputhandler"
	"fmt"
	"os"
	"strconv"
)

func main() {

	lines := inputhandler.ReadInput()

	maxPart1, err := CalcPart1Calories(lines)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	maxPart2, err := CalcPart2Calories(lines)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	fmt.Printf("MaxCalories - Part1: %d, Part2: %d\n", maxPart1, maxPart2)
}

func CalcPart1Calories(lines []string) (int64, error) {

	var max, curr int64
	for _, line := range lines {

		if len(line) == 0 {

			if max < curr {
				max = curr
			}

			curr = 0

			continue
		}

		value, err := strconv.ParseInt(line, 10, 0)
		if err != nil {
			return 0, fmt.Errorf("invalid value in data (%s)", line)
		}

		curr += value
	}

	if max < curr {
		max = curr
	}

	return max, nil
}

func CalcPart2Calories(lines []string) (int64, error) {

	var max = make([]int64, 3)
	var curr int64
	for _, line := range lines {

		if len(line) == 0 {
			for idx, _ := range max {
				if max[idx] <= curr {
					max[idx], curr = curr, max[idx]
				}
			}

			curr = 0

			continue
		}

		value, err := strconv.ParseInt(line, 10, 0)
		if err != nil {
			return 0, fmt.Errorf("invalid value in data (%s)", line)
		}

		curr += value
	}

	for idx, _ := range max {
		if max[idx] <= curr {
			max[idx], curr = curr, max[idx]
		}
	}

	return max[2] + max[1] + max[0], nil
}
