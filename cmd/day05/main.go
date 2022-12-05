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

	topBoxesPart1, err := processInput(lines, false)
	if err != nil {
		fmt.Printf("Error: while processing input: %v\n", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	topBoxesPart2, err := processInput(lines, true)
	if err != nil {
		fmt.Printf("Error: while processing input: %v\n", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	fmt.Printf("Top boxes - Part1: %s, Part2: %s", topBoxesPart1, topBoxesPart2)
}

//-Common----------------------------------------------------------------------

func processInput(lines []string, canDoMultiple bool) (string, error) {

	// build the supply stacks
	buildMode := true

	stackCount := int(math.Ceil(float64(len(lines[0])) / 4))
	supply := NewSupplyStacks(stackCount, canDoMultiple)

	for _, line := range lines {

		if buildMode {

			// hopefully this means moves are next
			if len(line) == 0 {

				// idx 0 should be bottom
				for idx := range supply.CargoStacks {
					supply.CargoStacks[idx].Reverse()
				}

				buildMode = false // got the supplys, start rearranging
				continue
			}

			// this line has at least one box
			if strings.Contains(line, "[") {

				for i := 0; i < stackCount; i++ {
					box := rune(line[1+i*4])
					if box == ' ' {
						continue
					}
					supply.CargoStacks[i+1].AddBoxes([]rune{rune(box)}) // +1 'cause 1st is padding
				}
			}

			continue
		}

		arrangementStep := strings.Split(line, " ")
		if len(arrangementStep) != 6 {
			return "", fmt.Errorf("invalid arrangement element count '%d' in line '%s'", len(arrangementStep), line)
		}

		moveCount, err := strconv.Atoi(arrangementStep[1])
		if err != nil {
			return "", fmt.Errorf("invalid move count '%s' in line '%s'", arrangementStep[1], line)
		}

		moveFrom, err := strconv.Atoi(arrangementStep[3])
		if err != nil {
			return "", fmt.Errorf("invalid stack reference '%s' in line '%s'", arrangementStep[3], line)
		}

		moveTo, err := strconv.Atoi(arrangementStep[5])
		if err != nil {
			return "", fmt.Errorf("invalid stack reference '%s' in line '%s'", arrangementStep[5], line)
		}

		move := *NewMove(moveCount, moveFrom, moveTo)

		supply.Rearrange(move)
	}

	topBoxes := supply.ReadTopBoxes()

	return string(topBoxes), nil
}

type Move struct {
	Count int
	From  int
	To    int
}

func NewMove(count, from, to int) *Move {
	return &Move{
		Count: count,
		From:  from,
		To:    to,
	}
}

type Stack struct {
	CargoBoxes []rune // 0 idx is bottom
}

func NewStack() *Stack {
	return &Stack{
		CargoBoxes: make([]rune, 0, 40), // TODO: optimize this magic number
	}
}

func (s *Stack) AddBoxes(boxes []rune) {
	s.CargoBoxes = append(s.CargoBoxes, boxes...)
}

func (s *Stack) RemoveBoxes(quantity int) ([]rune, error) {

	if len(s.CargoBoxes) < quantity {
		return nil, fmt.Errorf("too many boxes requested - '%d' out of '%d'", quantity, len(s.CargoBoxes))
	}

	boxes := s.CargoBoxes[len(s.CargoBoxes)-quantity:]

	s.CargoBoxes = s.CargoBoxes[:len(s.CargoBoxes)-quantity]

	return boxes, nil
}

func (s *Stack) Reverse() {
	reverseSplice(s.CargoBoxes)
}

func (s *Stack) Size() int {
	return len(s.CargoBoxes)
}

func (s *Stack) PeekTop() rune {
	return s.CargoBoxes[len(s.CargoBoxes)-1:][0]
}

type SupplyStacks struct {
	CargoStacks   []Stack
	canDoMultiple bool
}

func NewSupplyStacks(count int, doMultiple bool) *SupplyStacks {
	return &SupplyStacks{
		CargoStacks:   make([]Stack, count+1), // +1 so there is no +1/-1 shenanigans everywhere
		canDoMultiple: doMultiple,
	}
}

func (ss *SupplyStacks) Rearrange(move Move) error {

	if move.From < 1 || move.From > len(ss.CargoStacks) {
		return fmt.Errorf("moving from invalid stack '%d' out of '%d'", move.From, len(ss.CargoStacks))
	}

	if move.To < 1 || move.To > len(ss.CargoStacks) {
		return fmt.Errorf("moving to invalid stack '%d' out of '%d'", move.To, len(ss.CargoStacks))
	}

	stack, err := ss.CargoStacks[move.From].RemoveBoxes(move.Count)
	if err != nil {
		return fmt.Errorf("couldn't remove from stack '%d': %w", move.From, err)
	}

	if !ss.canDoMultiple {
		reverseSplice(stack) // actually can only move 1 at a time so technically order would reverse
	}

	ss.CargoStacks[move.To].AddBoxes(stack)

	return nil
}

func (ss SupplyStacks) ReadTopBoxes() []rune {

	var topBoxes = make([]rune, len(ss.CargoStacks))
	for idx := range ss.CargoStacks {

		// idk what if there is no box in a stack - assume it's a space?
		if ss.CargoStacks[idx].Size() == 0 {
			topBoxes[idx] = ' '
		} else {
			topBoxes[idx] = ss.CargoStacks[idx].PeekTop()
		}

	}

	return topBoxes[1:] // don't forget idx 0 is an extra
}

func reverseSplice[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
