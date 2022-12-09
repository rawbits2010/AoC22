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

var test = `R 4
U 4
L 3
D 1
R 4
D 1
L 5
R 2`

var fieldColor string
var tailMarkColor string
var startColor string

func main() {

	outputhandler.Initialize()
	defer outputhandler.Reset()

	fieldColor = outputhandler.GetForeground(outputhandler.Gray)
	tailMarkColor = outputhandler.GetForeground(outputhandler.BrightGreen)
	startColor = outputhandler.GetColor(outputhandler.White, outputhandler.BrightRed)

	//lines := strings.Split(test, "\n")
	lines := inputhandler.ReadInput()

	result, err := simulate(lines)
	if err != nil {
		fmt.Printf("Error: simulating movement: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	fmt.Printf("Result: %d", result)
}

func simulate(lines []string) (int, error) {

	bridge := NewRopeBridge()

	for _, line := range lines {

		tokens := strings.Split(line, " ")
		if len(tokens) != 2 {
			return 0, fmt.Errorf("invalid number of inputs in line '%s'", line)
		}

		direction := tokens[0]
		steps, err := strconv.Atoi(tokens[1])
		if err != nil {
			return 0, fmt.Errorf("invalid steps '%s' in line '%s'", tokens[1], line)
		}

		switch direction {
		case "L":
			for i := 1; i <= steps; i++ {
				bridge.MoveLeft()
			}

		case "R":
			for i := 1; i <= steps; i++ {
				bridge.MoveRight()
			}

		case "U":
			for i := 1; i <= steps; i++ {
				bridge.MoveUp()
			}

		case "D":
			for i := 1; i <= steps; i++ {
				bridge.MoveDown()
			}

		default:
			return 0, fmt.Errorf("invalid movement '%s' in line '%s'", tokens[0], line)

		}
	}

	visualizeTailTracks(bridge.TailTrack)

	return len(bridge.TailTrack), nil
}

//-----------------------------------------------------------------------------

type RopeBridge struct {
	Head      Position
	Tail      Position
	TailTrack []Position
}

func NewRopeBridge() *RopeBridge {
	var ropeBridge RopeBridge
	ropeBridge.Head = *NewPosition(0, 0)
	ropeBridge.Tail = *NewPosition(0, 0)
	ropeBridge.TailTrack = []Position{*NewPosition(0, 0)}
	return &ropeBridge
}

func (rb *RopeBridge) MoveLeft() {

	rb.Head.X--

	if isTailTouching(rb.Head, rb.Tail) {
		return
	}

	rb.Tail.X--
	rb.Tail.Y = rb.Head.Y

	rb.updateTailTracking()
}

func (rb *RopeBridge) MoveRight() {

	rb.Head.X++

	if isTailTouching(rb.Head, rb.Tail) {
		return
	}

	rb.Tail.X++
	rb.Tail.Y = rb.Head.Y

	rb.updateTailTracking()
}

func (rb *RopeBridge) MoveUp() {

	rb.Head.Y--

	if isTailTouching(rb.Head, rb.Tail) {
		return
	}

	rb.Tail.Y--
	rb.Tail.X = rb.Head.X

	rb.updateTailTracking()
}

func (rb *RopeBridge) MoveDown() {

	rb.Head.Y++

	if isTailTouching(rb.Head, rb.Tail) {
		return
	}

	rb.Tail.Y++
	rb.Tail.X = rb.Head.X

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

type Position struct {
	X int
	Y int
}

func NewPosition(x, y int) *Position {
	return &Position{X: x, Y: y}
}

func isTailTouching(p1, p2 Position) bool {
	distance := math.Sqrt(math.Abs(float64(p1.X-p2.X)) + math.Abs(float64(p1.Y-p2.Y)))
	if p1.X == p2.X || p1.Y == p2.Y {
		return distance <= 1.1
	}
	return distance < 1.5
}

func visualizeTailTracks(tracks []Position) {

	var minX, maxX, minY, maxY int
	for _, pos := range tracks {
		if pos.X > maxX {
			maxX = pos.X
		}
		if pos.X < minX {
			minX = pos.X
		}
		if pos.Y > maxY {
			maxY = pos.Y
		}
		if pos.Y < minY {
			minY = pos.Y
		}
	}

	offX := -minX
	offY := -minY
	maxX += offX
	maxY += offY
	field := make([][]byte, maxY+1)
	for i := 0; i < maxY+1; i++ {
		row := make([]byte, maxX+1)
		for j := 0; j < maxX+1; j++ {
			row[j] = '.'
		}
		field[i] = row
	}

	for _, pos := range tracks {
		field[pos.Y+offY][pos.X+offX] = '#'
	}
	field[tracks[0].Y+offY][tracks[0].X+offX] = 's'

	for _, row := range field {

		var colorized string = ""
		var currColor string = outputhandler.GetReset()
		var lastRune byte = '\000'

		for _, currRune := range row {

			if currRune != lastRune {
				lastRune = currRune
				switch currRune {
				case '.':
					currColor = fieldColor
				case '#':
					currColor = tailMarkColor
				case 's':
					currColor = startColor
				}
				colorized += currColor
			}
			colorized += string(currRune)
		}

		fmt.Println(colorized)
	}
}
