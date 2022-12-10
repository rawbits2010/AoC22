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

	var result int
	result, err := runCode(lines)
	if err != nil {
		fmt.Printf("Error while running part 1 the code: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	fmt.Printf("Result: %d", result)
}

func runCode(lines []string) (int, error) {

	var sumSignalStrength int
	var cyclesOfInterest = []int{20, 60, 100, 140, 180, 220}

	var currCycle int
	var regX int = 1

	var programCounter int = -1 // zero based
	var instruction string
	var instArgs []string
	var instCyclesLeft int
	for { // CPU clock
		currCycle++

		// signal probe
		for _, cycle := range cyclesOfInterest {
			if currCycle == cycle {
				currSignalStrength := currCycle * regX
				sumSignalStrength += currSignalStrength
				fmt.Printf("cycle: %d, regX: %d, signal: %d - sum: %d\n", currCycle, regX, currSignalStrength, sumSignalStrength)
			}
		}

		// setup next instruction
		if instCyclesLeft <= 0 {
			programCounter++

			// end of program
			if programCounter == len(lines) {
				return sumSignalStrength, nil
			}

			tokens := strings.Split(lines[programCounter], " ")
			if len(tokens) == 0 {
				return 0, fmt.Errorf("missing instruction at line number %d", programCounter+1)
			}

			instruction = tokens[0]
			instArgs = tokens[1:]

			switch instruction {
			case "noop":
				instCyclesLeft = 1

			case "addx":
				if len(instArgs) < 1 {
					return 0, fmt.Errorf("not enough arguments for addx in line '%s'", lines[programCounter])
				}
				instCyclesLeft = 2

			default:
				return 0, fmt.Errorf("unknown command in line '%s'", lines[programCounter])
			}

		}
		/*
			if currCycle <= 20 {
				fmt.Printf("cycle: %d, regX: %d - inst: %s, cyclesleft: %d\n", currCycle, regX, lines[programCounter], instCyclesLeft)
			}
		*/
		// process instruction
		instCyclesLeft--
		switch instruction {
		case "noop":
			// chill
			//fmt.Printf("inst: noop, cyclesleft: %d|n", instCyclesLeft)

		case "addx":
			if instCyclesLeft == 0 {
				val, err := strconv.Atoi(instArgs[0])
				if err != nil {
					return 0, fmt.Errorf("invalid argument for addx '%s' in line '%s'", instArgs[0], lines[programCounter])
				}
				//fmt.Printf("inst: addx, arg: %d, cyclesleft: %d\n", val, instCyclesLeft)

				regX += val
			}
		}

	} // CPU clock
}

func getSignalStrength(cycle, register int) int {
	return cycle * register
}
