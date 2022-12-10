package main

import (
	"AoC22/internal/inputhandler"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	lines := inputhandler.ReadInput()

	var result int
	result, err := runCode(lines)
	if err != nil && !errors.Is(err, ErrorEndOfProgram) {
		fmt.Printf("Error while running part 1 the code: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	fmt.Printf("Result: %d", result)
}

func runCode(lines []string) (int, error) {

	// build the PC :)
	databus := NewDataBus(lines)
	cpu := *NewCPU(databus)

	var sumSignalStrength int
	var cyclesOfInterest = []int{20, 60, 100, 140, 180, 220}

	var currCycle int
	for {
		currCycle++

		// signal probe
		for _, cycle := range cyclesOfInterest {
			if currCycle == cycle {
				currSignalStrength := currCycle * cpu.RegX
				sumSignalStrength += currSignalStrength
				fmt.Printf("cycle: %d, regX: %d, signal: %d - sum: %d\n", currCycle, cpu.RegX, currSignalStrength, sumSignalStrength)
			}
		}

		if err := cpu.tick(); err != nil {
			return sumSignalStrength, fmt.Errorf("CPU exception: %w", err)
		}

	}
}

//-----------------------------------------------------------------------------

type DataBus struct {
	code []string
}

func NewDataBus(lines []string) *DataBus {
	return &DataBus{code: lines}
}

func (db *DataBus) requestInstruction(index int) (string, error) {
	if index >= len(db.code) {
		return "", fmt.Errorf("out of range")
	}
	return db.code[index], nil
}

//-----------------------------------------------------------------------------

var ErrorEndOfProgram = fmt.Errorf("end of program")

type CPU struct {
	RegX           int
	ProgramCounter int // zero based

	Instruction    string
	InstArgs       []string
	InstCyclesLeft int

	databus *DataBus
}

func NewCPU(databus *DataBus) *CPU {
	return &CPU{
		RegX:           1,
		ProgramCounter: -1,
		databus:        databus,
	}
}

func (cpu *CPU) tick() error {

	// setup next instruction
	if cpu.InstCyclesLeft <= 0 {
		cpu.ProgramCounter++

		line, err := cpu.databus.requestInstruction(cpu.ProgramCounter)
		if err != nil {
			return ErrorEndOfProgram // for simplicity
		}

		tokens := strings.Split(line, " ")
		if len(tokens) == 0 {
			return fmt.Errorf("missing instruction at line number %d", cpu.ProgramCounter+1)
		}

		cpu.Instruction = tokens[0]
		cpu.InstArgs = tokens[1:]

		switch cpu.Instruction {
		case "noop":
			cpu.InstCyclesLeft = 1

		case "addx":
			if len(cpu.InstArgs) < 1 {
				return fmt.Errorf("not enough arguments for addx in line '%s'", line)
			}
			cpu.InstCyclesLeft = 2

		default:
			return fmt.Errorf("unknown command in line '%s'", line)
		}

	}
	/*
		if currCycle <= 20 {
			fmt.Printf("cycle: %d, regX: %d - inst: %s, cyclesleft: %d\n", currCycle, regX, lines[programCounter], instCyclesLeft)
		}
	*/
	// process instruction
	cpu.InstCyclesLeft--
	switch cpu.Instruction {
	case "noop":
		// chill
		//fmt.Printf("inst: noop, cyclesleft: %d|n", instCyclesLeft)

	case "addx":
		if cpu.InstCyclesLeft == 0 {
			val, err := strconv.Atoi(cpu.InstArgs[0])
			if err != nil {
				return fmt.Errorf("invalid argument for addx '%s' in line number '%d'", cpu.InstArgs[0], cpu.ProgramCounter+1)
			}
			//fmt.Printf("inst: addx, arg: %d, cyclesleft: %d\n", val, instCyclesLeft)

			cpu.RegX += val
		}
	}

	return nil
}
