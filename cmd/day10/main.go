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

	// Part 1
	part1ProbeCycles := []int{20, 60, 100, 140, 180, 220}
	part1Probe := NewSignalStrengthProbe(part1ProbeCycles)

	err := runCode(lines, part1Probe)
	if err != nil && !errors.Is(err, ErrorEndOfProgram) {
		fmt.Printf("Error while running part 1 the code: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	fmt.Printf("Result - Part1: %d", part1Probe.sumSignalStrength)
}

func runCode(program []string, probe SignalProber) error {

	// build the PC :) - here, for simplicity
	databus := NewDataBus(program)
	cpu := NewCPU(databus)

	probe.attachProbe(cpu)

	var currCycle int
	for {
		currCycle++

		if probe.needsProbing(currCycle) {
			probe.probe(currCycle)
		}

		if err := cpu.tick(); err != nil {
			return fmt.Errorf("CPU exception: %w", err)
		}

	}
}

//-----------------------------------------------------------------------------

// SignalStrengthProbe is the probe for Part 1
type SignalStrengthProbe struct {
	SignalProbe

	sumSignalStrength int
}

func NewSignalStrengthProbe(probingCycles []int) *SignalStrengthProbe {
	return &SignalStrengthProbe{
		SignalProbe: *NewSignalProbe(probingCycles),
	}
}

func (ssp *SignalStrengthProbe) probe(cycle int) {
	ssp.sumSignalStrength += cycle * ssp.cpu.RegX
}

// SignalProber is the interface for a probe
type SignalProber interface {
	attachProbe(cpu *CPU)
	needsProbing(cycle int) bool
	probe(cycle int)
}

// SignalProbe is a base for an actual signal probing function
type SignalProbe struct {
	cpu *CPU

	probingCycles []int
}

func NewSignalProbe(probingCycles []int) *SignalProbe {
	return &SignalProbe{
		probingCycles: probingCycles,
	}
}

func (sp *SignalProbe) attachProbe(cpu *CPU) {
	sp.cpu = cpu
}

func (sp *SignalProbe) needsProbing(cycle int) bool {
	for _, neededCycle := range sp.probingCycles {
		if cycle == neededCycle {
			return true
		}
	}
	return false
}

//-----------------------------------------------------------------------------

type DataBus struct {
	code []string
}

func NewDataBus(lines []string) *DataBus {
	return &DataBus{code: lines}
}

func (db *DataBus) requestInstruction(index int) (string, error) {
	if index > len(db.code) || index < 0 {
		return "", fmt.Errorf("out of range")
	} else if index == len(db.code) { // for simplicity
		return "", ErrorEndOfProgram
	}
	return db.code[index], nil
}

var ErrorEndOfProgram = fmt.Errorf("end of program")

//-----------------------------------------------------------------------------

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
			return fmt.Errorf("error requesting instruction: %w", err)
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
