package main

import (
	"AoC22/internal/inputhandler"
	"AoC22/internal/outputhandler"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {

	outputhandler.Initialize()
	defer outputhandler.Reset()

	fieldColor = outputhandler.GetForeground(outputhandler.Gray)
	tailMarkColor = outputhandler.GetForeground(outputhandler.BrightGreen)
	bridgeColor = outputhandler.GetForeground(outputhandler.BrightGreen)
	startColor = outputhandler.GetColor(outputhandler.White, outputhandler.BrightRed)

	lines := inputhandler.ReadInput()

	tailTrackCountPart1, err := simulate(lines, 2)
	if err != nil {
		fmt.Printf("Error: simulating movement: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	tailTrackCountPart2, err := simulate(lines, 10)
	if err != nil {
		fmt.Printf("Error: simulating movement: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	fmt.Printf("Result - Part1: %d, Part2: %d", tailTrackCountPart1, tailTrackCountPart2)
}

func simulate(lines []string, knots int) (int, error) {

	if knots < 2 {
		return 0, fmt.Errorf("invalid number of knots '%d' need at least 2", knots)
	}

	bridge := NewRopeBridge(knots)

	for lineIdx, line := range lines {
		_ = lineIdx

		tokens := strings.Split(line, " ")
		if len(tokens) != 2 {
			return 0, fmt.Errorf("invalid number of inputs in line '%s'", line)
		}

		direction := tokens[0]
		steps, err := strconv.Atoi(tokens[1])
		if err != nil {
			return 0, fmt.Errorf("invalid steps '%s' in line '%s'", tokens[1], line)
		}

		for i := 1; i <= steps; i++ {

			switch direction {
			case "L":
				bridge.MoveLeft()

			case "R":
				bridge.MoveRight()

			case "U":
				bridge.MoveUp()

			case "D":
				bridge.MoveDown()

			default:
				return 0, fmt.Errorf("invalid movement '%s' in line '%s'", tokens[0], line)

			}

			//VisualizeBridge(bridge.Knots)
		}
	}

	//VisualizeTailTracks(bridge.TailTrack)

	return len(bridge.TailTrack), nil
}

//-----------------------------------------------------------------------------

type RopeBridge struct {
	Knots        []Position
	Head         *Position
	Tail         *Position
	RopeSections []RopeSection
	TailTrack    []Position
}

func NewRopeBridge(knots int) *RopeBridge {

	var ropeBridge RopeBridge

	ropeBridge.Knots = make([]Position, knots) // inited to 0,0 by default
	ropeBridge.RopeSections = make([]RopeSection, 0, knots-1)
	for i := 0; i < knots-1; i++ {
		ropeBridge.RopeSections = append(ropeBridge.RopeSections, *NewRopeSection(&ropeBridge.Knots[i], &ropeBridge.Knots[i+1]))
	}

	ropeBridge.Head = &ropeBridge.Knots[0]
	ropeBridge.Tail = &ropeBridge.Knots[len(ropeBridge.Knots)-1]

	ropeBridge.TailTrack = []Position{*NewPosition(0, 0)}

	return &ropeBridge
}

func (rb *RopeBridge) MoveLeft() {

	rb.Head.X--

	rb.updateChain()
}

func (rb *RopeBridge) MoveRight() {

	rb.Head.X++

	rb.updateChain()
}

func (rb *RopeBridge) MoveUp() {

	rb.Head.Y--

	rb.updateChain()
}

func (rb *RopeBridge) MoveDown() {

	rb.Head.Y++

	rb.updateChain()
}

func (rb *RopeBridge) updateChain() {

	for _, section := range rb.RopeSections {
		section.Update()
	}

	rb.updateTailTracking()
}

func (rb *RopeBridge) updateTailTracking() {

	for _, pos := range rb.TailTrack {
		if pos.X == rb.Tail.X && pos.Y == rb.Tail.Y {
			return
		}
	}

	rb.TailTrack = append(rb.TailTrack, *NewPosition(rb.Tail.X, rb.Tail.Y))
}

type DragDirection string

const (
	Left      DragDirection = "L"
	Right     DragDirection = "R"
	Up        DragDirection = "U"
	Down      DragDirection = "D"
	UpLeft    DragDirection = "UL"
	UpRight   DragDirection = "UR"
	DownLeft  DragDirection = "DL"
	DownRight DragDirection = "DR"
	NoDrag    DragDirection = ""
)

type RopeSection struct {
	FirstKnot  *Position
	SecondKnot *Position
}

func NewRopeSection(first *Position, second *Position) *RopeSection {
	return &RopeSection{FirstKnot: first, SecondKnot: second}
}

func (s *RopeSection) Update() {

	if s.isTouching() {
		return
	}

	drag := s.dragDirection()
	switch drag {
	case Left:
		s.SecondKnot.X--
		s.SecondKnot.Y = s.FirstKnot.Y

	case Right:
		s.SecondKnot.X++
		s.SecondKnot.Y = s.FirstKnot.Y

	case Up:
		s.SecondKnot.Y--
		s.SecondKnot.X = s.FirstKnot.X

	case Down:
		s.SecondKnot.Y++
		s.SecondKnot.X = s.FirstKnot.X

	case UpLeft:
		s.SecondKnot.Y--
		s.SecondKnot.X--

	case UpRight:
		s.SecondKnot.Y--
		s.SecondKnot.X++

	case DownLeft:
		s.SecondKnot.Y++
		s.SecondKnot.X--

	case DownRight:
		s.SecondKnot.Y++
		s.SecondKnot.X++

	}
}

func (s *RopeSection) dragDirection() DragDirection {

	// larger step determines the direction
	hOff := s.FirstKnot.X - s.SecondKnot.X
	vOff := s.FirstKnot.Y - s.SecondKnot.Y

	if math.Abs(float64(hOff)) > math.Abs(float64(vOff)) {
		if hOff > 0 {
			return Right
		} else if hOff < 0 {
			return Left
		}
	} else if math.Abs(float64(hOff)) < math.Abs(float64(vOff)) {
		if vOff > 0 {
			return Down
		} else if vOff < 0 {
			return Up
		}
	} else {
		if hOff > 0 {
			if vOff > 0 {
				return DownRight
			} else if vOff < 0 {
				return UpRight
			}
		} else if hOff < 0 {
			if vOff > 0 {
				return DownLeft
			} else if vOff < 0 {
				return UpLeft
			}
		}
	}

	return NoDrag
}

func (s *RopeSection) isTouching() bool {

	distance := math.Sqrt(math.Abs(float64(s.FirstKnot.X-s.SecondKnot.X)) + math.Abs(float64(s.FirstKnot.Y-s.SecondKnot.Y)))

	// neighbour
	if s.FirstKnot.X == s.SecondKnot.X || s.FirstKnot.Y == s.SecondKnot.Y {
		return distance <= 1.1
	}

	// diagonal
	return distance < 1.5
}

type Position struct {
	X int
	Y int
}

func NewPosition(x, y int) *Position {
	return &Position{X: x, Y: y}
}
