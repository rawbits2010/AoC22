package main

import (
	"AoC22/internal/outputhandler"
	"fmt"
	"strconv"
)

// set these with color values from the outputhandler
var fieldColor string
var tailMarkColor string
var bridgeColor string
var startColor string

// VisualizeTailTracks prints out the field containing the tailtrack points to stdout.
// It uses the above specified colors if set.
func VisualizeTailTracks(tracks []Position) {

	offX, offY, maxX, maxY := getDimensions(tracks)
	field := createField(maxX+1, maxY+1)

	for _, pos := range tracks {
		field[pos.Y+offY][pos.X+offX] = '#'
	}
	field[tracks[0].Y+offY][tracks[0].X+offX] = 's'

	colorPrint(field)
}

// VisualizeBridge prints out the field containing the bridge's knots to stdout.
// It uses the above specified colors if set.
func VisualizeBridge(knots []Position) {

	offX, offY, maxX, maxY := getDimensions(knots)
	field := createField(maxX+1, maxY+1)

	field[offY][offX] = 's'
	for idx := len(knots) - 1; idx >= 0; idx-- {
		pos := knots[idx]
		field[pos.Y+offY][pos.X+offX] = strconv.Itoa(idx)[0]
	}
	field[knots[0].Y+offY][knots[0].X+offX] = 'H'

	colorPrint(field)
}

func colorPrint(field [][]byte) {

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
				default:
					currColor = bridgeColor
				}
				colorized += currColor
			}
			colorized += string(currRune)
		}

		fmt.Println(colorized + outputhandler.GetReset())
	}

	fmt.Println()
}

func createField(width, height int) [][]byte {

	field := make([][]byte, height)
	for i := 0; i < height; i++ {
		row := make([]byte, width)
		for j := 0; j < width; j++ {
			row[j] = '.'
		}
		field[i] = row
	}

	return field
}

func getDimensions(posList []Position) (int, int, int, int) {

	var minX, maxX, minY, maxY int
	for _, pos := range posList {
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

	return offX, offY, maxX, maxY
}
