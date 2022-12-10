package main

import (
	"AoC22/internal/inputhandler"
	"AoC22/internal/outputhandler"
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
		fmt.Printf("Error while running part 1 code: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	// Part 2
	part2ProbeCycles := []int{40, 80, 120, 160, 200, 240}
	part2Probe := NewDisplaySignalProbe(part2ProbeCycles)

	err = runCode(lines, part2Probe)
	if err != nil && !errors.Is(err, ErrorEndOfProgram) {
		fmt.Printf("Error while running part 2 code: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	fmt.Printf("Result - Part1: %d, Part 2:\n", part1Probe.SumSignalStrength)
	vizualizeDisplaySignalProbe(part2Probe)
}

func vizualizeDisplaySignalProbe(probe *DisplaySignalProbe) {

	outputhandler.Initialize()
	defer outputhandler.Reset()

	litPixelColor := outputhandler.GetColor(outputhandler.White, outputhandler.BrightGreen)
	unlitPixelColor := outputhandler.GetForeground(outputhandler.Gray)

	for _, row := range probe.Display {

		var colorized string = ""
		var currColor string = outputhandler.GetReset()
		var lastRune rune = '\000'

		for _, currRune := range row {

			if currRune != lastRune {
				lastRune = currRune
				switch currRune {
				case '.':
					currColor = unlitPixelColor
				case '#':
					currColor = litPixelColor
				}
				colorized += currColor
			}
			colorized += string(currRune)
		}

		fmt.Println(colorized + outputhandler.GetReset())
	}

}

func runCode(program []string, probe SignalProber) error {

	// build the PC :) - here, for simplicity
	databus := NewDataBus(program)
	cpu := NewCPU(databus)
	gpu := NewGPU(cpu)

	probe.AttachProbe(cpu, gpu)

	var currCycle int
	for {
		currCycle++

		gpu.Tick() // just like OpenGL, the show must go on :)

		if probe.NeedsProbing(currCycle) {
			probe.Probe(currCycle)
		}

		if err := cpu.Tick(); err != nil {
			return fmt.Errorf("CPU exception: %w", err)
		}

	}
}

//-----------------------------------------------------------------------------

// DisplaySignalProbe is the probe for Part 2
type DisplaySignalProbe struct {
	SignalProbe

	Display []string
}

func NewDisplaySignalProbe(probingCycles []int) *DisplaySignalProbe {
	return &DisplaySignalProbe{
		SignalProbe: *NewSignalProbe(probingCycles),
		Display:     make([]string, 0, len(probingCycles)),
	}
}

func (dp *DisplaySignalProbe) Probe(cycle int) {
	dp.Display = append(dp.Display, string(dp.gpu.currScanline[:]))
}

// SignalStrengthProbe is the probe for Part 1
type SignalStrengthProbe struct {
	SignalProbe

	SumSignalStrength int
}

func NewSignalStrengthProbe(probingCycles []int) *SignalStrengthProbe {
	return &SignalStrengthProbe{
		SignalProbe: *NewSignalProbe(probingCycles),
	}
}

func (ssp *SignalStrengthProbe) Probe(cycle int) {
	ssp.SumSignalStrength += cycle * ssp.cpu.RegX
}

// SignalProber is the interface for a probe
type SignalProber interface {
	AttachProbe(cpu *CPU, gpu *GPU)
	NeedsProbing(cycle int) bool
	Probe(cycle int)
}

// SignalProbe is a base for an actual signal probing function
type SignalProbe struct {
	cpu *CPU
	gpu *GPU

	probingCycles []int
}

func NewSignalProbe(probingCycles []int) *SignalProbe {
	return &SignalProbe{
		probingCycles: probingCycles,
	}
}

func (sp *SignalProbe) AttachProbe(cpu *CPU, gpu *GPU) {
	sp.cpu = cpu
	sp.gpu = gpu
}

func (sp *SignalProbe) NeedsProbing(cycle int) bool {
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

func (db *DataBus) RequestInstruction(index int) (string, error) {
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

func (cpu *CPU) Tick() error {

	// setup next instruction
	if cpu.InstCyclesLeft <= 0 {
		cpu.ProgramCounter++

		line, err := cpu.databus.RequestInstruction(cpu.ProgramCounter)
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

	// process instruction
	cpu.InstCyclesLeft--
	switch cpu.Instruction {
	case "noop":
		// chill

	case "addx":
		if cpu.InstCyclesLeft == 0 {
			val, err := strconv.Atoi(cpu.InstArgs[0])
			if err != nil {
				return fmt.Errorf("invalid argument for addx '%s' in line number '%d'", cpu.InstArgs[0], cpu.ProgramCounter+1)
			}

			cpu.RegX += val
		}
	}

	return nil
}

//-----------------------------------------------------------------------------

const scanlineLength int = 40
const spriteLength int = 3

type GPU struct {
	currScanline [scanlineLength]byte
	currPixelIdx int

	cpu *CPU
}

func NewGPU(cpu *CPU) *GPU {
	return &GPU{
		currPixelIdx: -1,
		cpu:          cpu,
	}
}

func (gpu *GPU) Tick() {
	gpu.currPixelIdx++

	// new scanline
	if gpu.currPixelIdx == scanlineLength {
		for idx := range gpu.currScanline {
			gpu.currScanline[idx] = ' '
		}
		gpu.currPixelIdx = 0
	}

	// register X is the middle of the sprite
	var spriteLeftPos = gpu.cpu.RegX - spriteLength/2
	if gpu.currPixelIdx >= spriteLeftPos && gpu.currPixelIdx < spriteLeftPos+spriteLength {
		gpu.currScanline[gpu.currPixelIdx] = '#'
	} else {
		gpu.currScanline[gpu.currPixelIdx] = '.'
	}
}
