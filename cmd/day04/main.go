package main

import (
	"AoC22/internal/inputhandler"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	lines := inputhandler.ReadInput()

	overlapsPart1, err := countOverlapse(lines, isFullRangeOverlap)
	if err != nil {
		fmt.Printf("Error: while processing part 1: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	overlapsPart2, err := countOverlapse(lines, isPartialOverlap)
	if err != nil {
		fmt.Printf("Error: while processing part 2: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	fmt.Printf("Overlaps - Part1: %d, Part2: %d", overlapsPart1, overlapsPart2)
}

type overlapCheckFN func(Range, Range) bool

// Part 1
func isFullRangeOverlap(assignment1, assignment2 Range) bool {
	if assignment1.Overlapse(assignment2) {
		return true
	}
	if assignment2.Overlapse(assignment1) {
		return true
	}
	return false
}

// Part 2
func isPartialOverlap(assignment1, assignment2 Range) bool {

	if assignment1.Contains(assignment2.Min) || assignment1.Contains(assignment2.Max) {
		return true
	}
	if assignment2.Contains(assignment1.Min) || assignment2.Contains(assignment1.Max) {
		return true
	}
	return false
}

//-Common----------------------------------------------------------------------

type Range struct {
	Min, Max int
}

func NewRange(min, max int) *Range {
	return &Range{Min: min, Max: max}
}

func (r Range) Overlapse(other Range) bool {
	return r.Min <= other.Min && other.Max <= r.Max
}

func (r Range) Contains(val int) bool {
	return val >= r.Min && val <= r.Max
}

func countOverlapse(lines []string, overlapCheck overlapCheckFN) (int, error) {

	var overlapCount int
	for _, line := range lines {

		assignments := strings.Split(line, ",")
		if len(assignments) != 2 {
			return 0, fmt.Errorf("invalid assigment count '%d' in line '%s'", len(assignments), line)
		}

		ass1, err := getAssignemtRange(assignments[0])
		if err != nil {
			return 0, fmt.Errorf("invalid assignment '%s' in line '%s'", assignments[0], line)
		}

		ass2, err := getAssignemtRange(assignments[1])
		if err != nil {
			return 0, fmt.Errorf("invalid assignment '%s' in line '%s'", assignments[1], line)
		}

		if overlapCheck(*ass1, *ass2) {
			overlapCount++
		}
	}

	return overlapCount, nil
}

func getAssignemtRange(assignment string) (*Range, error) {

	ranges := strings.Split(assignment, "-")
	if len(ranges) != 2 {
		return nil, fmt.Errorf("invalid assignment range '%s'", assignment)
	}

	val1, err := strconv.Atoi(ranges[0])
	if err != nil {
		return nil, fmt.Errorf("error converting range limit '%s' in assignment '%s'", ranges[0], assignment)
	}

	val2, err := strconv.Atoi(ranges[1])
	if err != nil {
		return nil, fmt.Errorf("error converting range limit '%s' in assignment '%s'", ranges[0], assignment)
	}

	// noone stated the format of the ranges
	if val1 < val2 {
		return NewRange(val1, val2), nil
	}
	return NewRange(val2, val1), nil
}
