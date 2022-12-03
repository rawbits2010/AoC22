package main

import (
	"AoC22/internal/inputhandler"
	"fmt"
	"os"
)

func main() {

	lines := inputhandler.ReadInput()

	resultPart1, err := calcPart1Result(lines)
	if err != nil {
		fmt.Printf("Error processing part 1: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	resultPart2, err := calcPart2Result(lines)
	if err != nil {
		fmt.Printf("Error processing part 2: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	fmt.Printf("Result - Part1: %d, Part2: %d", resultPart1, resultPart2)
}

func calcPart2Result(lines []string) (int, error) {

	if len(lines)%3 != 0 {
		return 0, fmt.Errorf("invalid line count '%d'", len(lines))
	}

	var result int
	for groupIdx := 0; groupIdx < len(lines); groupIdx += 3 {

		groupSacks := lines[groupIdx : groupIdx+3]

		groupChecklists := make([][]bool, 3)
		for elfIdx, sacks := range groupSacks {

			checklist, err := checklistItems(sacks)
			if err != nil {
				return 0, fmt.Errorf("error processing line '%s': %w", sacks, err)
			}

			groupChecklists[elfIdx] = checklist
		}

		var badge int = 0
		for itemIdx := 0; itemIdx < 53; itemIdx++ {
			if groupChecklists[0][itemIdx] && groupChecklists[1][itemIdx] && groupChecklists[2][itemIdx] {
				badge = itemIdx
				break // it's stated that there can be only one
			}
		}
		if badge == 0 {
			return 0, fmt.Errorf("couldn't find bagde for lines '%s;%s;%s'", groupSacks[0], groupSacks[1], groupSacks[2])
		}

		result += badge
	}

	return result, nil
}

func calcPart1Result(lines []string) (int, error) {
	var result int
	for _, line := range lines {

		if len(line)%2 != 0 {
			return 0, fmt.Errorf("invalid line length '%d' for line '%s'", len(line), line)
		}

		sackCount := len(line) / 2
		comp1 := line[:sackCount]
		comp2 := line[sackCount:]

		checklist, err := checklistItems(comp1)
		if err != nil {
			return 0, fmt.Errorf("error processing line '%s': %w", line, err)
		}

		var mistake int = 0
		for _, char := range comp2 {

			intVal, err := getItemPriority(char) // runes are Unicode but for our purpose UTF-8 behaves like ASCII
			if err != nil {
				return 0, fmt.Errorf("invalid character '%s' in line '%s'", string(char), line)
			}

			if checklist[intVal] {
				mistake = intVal
				break // it's stated that there can be only one
			}
		}

		if mistake == 0 {
			return 0, fmt.Errorf("couldn't find mistake n line '%s'", line)
		}

		result += mistake
	}

	return result, nil
}

func getItemPriority(item rune) (int, error) {

	switch intVal := int(item); {
	case intVal >= int('a'):
		return intVal - 96, nil
	case intVal <= int('Z'):
		return intVal - 64 + 26, nil
	}

	return 0, fmt.Errorf("invalid character '%s'", string(item))
}

func checklistItems(compartment string) ([]bool, error) {

	checklist := make([]bool, 53) // +1 so no need for -1 indexing everywhere
	for _, char := range compartment {

		intVal, err := getItemPriority(char) // runes are Unicode but for our purpose UTF-8 behaves like ASCII
		if err != nil {
			return nil, fmt.Errorf("invalid character '%s' in compartment '%s'", string(char), compartment)
		}

		checklist[intVal] = true
	}

	return checklist, nil
}
