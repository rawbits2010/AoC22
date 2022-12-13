package main

import (
	"AoC22/internal/inputhandler"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
)

func main() {

	lines := inputhandler.ReadInput()

	signal, err := parseSignal(lines)
	if err != nil {
		fmt.Printf("Error parsing signal: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	// Part 1
	//fmt.Println("-----------------------------------------")

	inOrderCount, err := processSignal(signal)
	if err != nil {
		fmt.Printf("Error processing signal Part 1: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	// Part 2
	//fmt.Println("-----------------------------------------")

	decoderKey := findDecoderKey(signal)

	fmt.Printf("Result: Part1: %d, Part2: %d", inOrderCount, decoderKey)
}

func parseSignal(lines []string) ([]interface{}, error) {

	signal := make([]interface{}, 0)
	for _, line := range lines {

		if len(line) == 0 {
			continue
		}

		temp, _, err := parseSlice(0, line, make([]interface{}, 0))
		if err != nil {
			return nil, fmt.Errorf("error while parsing line '%s': %w", line, err)
		}
		//fmt.Println()

		signal = append(signal, temp)
	}

	return signal, nil
}

// assume correct format - no spaces, starts with a '['
func parseSlice(startIdx int, line string, parent []interface{}) ([]interface{}, int, error) {

	var tempVal []byte
	for charIdx := startIdx; charIdx < len(line); charIdx++ {

		switch line[charIdx] {
		case '[':
			//fmt.Println("slice start")

			temp, tempIdx, err := parseSlice(charIdx+1, line, make([]interface{}, 0))
			if err != nil {
				return nil, 0, err
			}
			parent = append(parent, temp)
			charIdx = tempIdx

			if charIdx == len(line)-1 {
				return parent, charIdx, nil
			}

		case ']':
			//fmt.Print("slice end")

			if len(tempVal) > 0 {

				intVal, err := strconv.Atoi(string(tempVal))
				if err != nil {
					return nil, 0, fmt.Errorf("couldn't convert '%s' to int", tempVal)
				}
				parent = append(parent, intVal)

				//fmt.Printf(" - append data %d", intVal)
			}
			//fmt.Println()

			return parent, charIdx, nil

		case ',':
			if len(tempVal) != 0 {

				intVal, err := strconv.Atoi(string(tempVal))
				if err != nil {
					return nil, 0, fmt.Errorf("couldn't convert '%s' to int", tempVal)
				}

				parent = append(parent, intVal)
				//fmt.Printf("append data %d\n", intVal)

				tempVal = []byte{}
				continue
			}

			// could be a slice
			if line[charIdx-1] != ']' {
				return nil, 0, fmt.Errorf("missing value while encountering a ','")
			}

		default:
			tempVal = append(tempVal, line[charIdx])
		}
	}

	return nil, 0, fmt.Errorf("unexpected end of data")
}

//-----------------------------------------------------------------------------

func processSignal(signal []interface{}) (int, error) {

	if len(signal)%2 != 0 {
		return 0, fmt.Errorf("invalid number or signal packets")
	}

	var sumInOrderIdx int
	for packetIdx := 0; packetIdx <= len(signal)-2; packetIdx += 2 {

		order := processPair(signal[packetIdx], signal[packetIdx+1])

		if order != WrongOrder {
			//fmt.Printf("was good: %d\n", packetIdx/2+1)
			sumInOrderIdx += packetIdx/2 + 1
		}

		//fmt.Println()
	}

	return sumInOrderIdx, nil
}

func processPair(left, right interface{}) Order {

	leftKind := reflect.ValueOf(left).Kind()
	rightKind := reflect.ValueOf(right).Kind()

	if leftKind != rightKind { // one is a slice

		if leftKind != reflect.Slice {
			//fmt.Println("slicify left")
			left = slicifyData(left)
			leftKind = reflect.ValueOf(left).Kind()
		}

		if rightKind != reflect.Slice {
			//fmt.Println("slicify right")
			right = slicifyData(right)
			//rightKind = reflect.ValueOf(right).Kind()
		}
	}

	if leftKind == reflect.Slice {
		return compareSlices(left, right)
	} else {
		return compareData(left, right)
	}
}

// left should be longer
// continue on equal if values didn't decide
func compareSlices(left, right interface{}) Order {
	//fmt.Println("compare slice")

	leftValue := reflect.ValueOf(left)
	rightValue := reflect.ValueOf(right)

	var idx int
	for idx = 0; idx < leftValue.Len(); idx++ {

		if idx == rightValue.Len() {
			return WrongOrder
		}

		order := processPair(leftValue.Index(idx).Interface(), rightValue.Index(idx).Interface())
		if order != Undecided {
			return order
		}
	}

	if idx == rightValue.Len() {
		return Undecided
	}

	// still in order
	return RightOrder
}

type Order string

const (
	RightOrder Order = "Right"
	Undecided  Order = "Undecided"
	WrongOrder Order = "Wrong"
)

// NOTE: currently only int is supported
func compareData(left, right interface{}) Order {

	leftValue := reflect.ValueOf(left)
	rightValue := reflect.ValueOf(right)
	//fmt.Printf("comparing: %d, %d\n", leftValue.Int(), rightValue.Int())

	if leftValue.Int() == rightValue.Int() {
		return Undecided
	}
	if leftValue.Int() < rightValue.Int() {
		return RightOrder
	}
	return WrongOrder
}

// returns a slice of the value
func slicifyData(data interface{}) interface{} {
	dataValue := reflect.ValueOf(data)
	temp := reflect.MakeSlice(reflect.SliceOf(dataValue.Type()), 1, 1)
	temp.Index(0).Set(dataValue)
	return temp.Interface()
}

//-----------------------------------------------------------------------------

func findDecoderKey(signal []interface{}) int {

	// adding the divider packets
	// NOTE: could just check the first numbers in the list
	// as these will always be the first of the ones starting with
	// the same numbers...
	dividerPacket1 := [][]int{{2}}
	dividerPacket2 := [][]int{{6}}
	signal = append(signal, dividerPacket1)
	signal = append(signal, dividerPacket2)

	sort.Slice(signal, func(i, j int) bool {
		return processPair(signal[i], signal[j]) == RightOrder
	})

	var dividerIdx1, dividerIdx2 int
	for packetIdx, packet := range signal {
		if processPair(packet, dividerPacket1) == Undecided {
			dividerIdx1 = packetIdx + 1
		}
		if processPair(packet, dividerPacket2) == Undecided {
			dividerIdx2 = packetIdx + 1
		}
	}

	return dividerIdx1 * dividerIdx2
}
