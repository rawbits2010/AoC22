package main

import (
	"AoC22/internal/inputhandler"
	"fmt"
	"math/big"
	"os"
	"sort"
)

/*
	var Monkeys = []Monkey{
		{
			Items: []Item{
				{WorryLevel: 99},
				{WorryLevel: 67},
				{WorryLevel: 92},
				{WorryLevel: 61},
				{WorryLevel: 83},
				{WorryLevel: 64},
				{WorryLevel: 98},
			},
			OpFN:            Multiplication,
			OpVal2:          17,
			TestDivision:    3,
			TestResultTrue:  4,
			TestResultFalse: 2,
		},
		{
			Items: []Item{
				{WorryLevel: 78},
				{WorryLevel: 74},
				{WorryLevel: 88},
				{WorryLevel: 89},
				{WorryLevel: 50},
			},
			OpFN:            Multiplication,
			OpVal2:          11,
			TestDivision:    5,
			TestResultTrue:  3,
			TestResultFalse: 5,
		},
		{
			Items: []Item{
				{WorryLevel: 98},
				{WorryLevel: 91},
			},
			OpFN:            Addition,
			OpVal2:          4,
			TestDivision:    2,
			TestResultTrue:  6,
			TestResultFalse: 4,
		},
		{
			Items: []Item{
				{WorryLevel: 59},
				{WorryLevel: 72},
				{WorryLevel: 94},
				{WorryLevel: 91},
				{WorryLevel: 79},
				{WorryLevel: 88},
				{WorryLevel: 94},
				{WorryLevel: 51},
			},
			OpFN:            Power,
			OpVal2:          0,
			TestDivision:    13,
			TestResultTrue:  0,
			TestResultFalse: 5,
		},
		{
			Items: []Item{
				{WorryLevel: 95},
				{WorryLevel: 72},
				{WorryLevel: 78},
			},
			OpFN:            Addition,
			OpVal2:          7,
			TestDivision:    11,
			TestResultTrue:  7,
			TestResultFalse: 6,
		},
		{
			Items: []Item{
				{WorryLevel: 76},
			},
			OpFN:            Addition,
			OpVal2:          8,
			TestDivision:    17,
			TestResultTrue:  0,
			TestResultFalse: 2,
		},
		{
			Items: []Item{
				{WorryLevel: 69},
				{WorryLevel: 60},
				{WorryLevel: 53},
				{WorryLevel: 89},
				{WorryLevel: 71},
				{WorryLevel: 88},
			},
			OpFN:            Addition,
			OpVal2:          5,
			TestDivision:    19,
			TestResultTrue:  7,
			TestResultFalse: 1,
		},
		{
			Items: []Item{
				{WorryLevel: 72},
				{WorryLevel: 54},
				{WorryLevel: 63},
				{WorryLevel: 80},
			},
			OpFN:            Addition,
			OpVal2:          3,
			TestDivision:    7,
			TestResultTrue:  1,
			TestResultFalse: 3,
		},
	}
*/

func CreateTestMonkeyGroup(useRelief bool) []Monkey[int] {
	var temp = make([]Monkey[int], 4)

	temp[0] = *NewMonkey([]int{79, 98}, Multiplication, 19, useRelief, 23, 2, 3)
	temp[1] = *NewMonkey([]int{54, 65, 75, 74}, Addition, 6, useRelief, 19, 2, 0)
	temp[2] = *NewMonkey([]int{79, 60, 97}, Power, 0, useRelief, 13, 1, 3)
	temp[3] = *NewMonkey([]int{74}, Addition, 3, useRelief, 17, 0, 1)

	return temp
}

func CreateTestMonkeyGroupBig(useRelief bool) []Monkey[*big.Int] {
	var temp = make([]Monkey[*big.Int], 4)

	temp[0] = *NewMonkeyBig([]int{79, 98}, MultiplicationBig, 19, useRelief, 23, 2, 3)
	temp[1] = *NewMonkeyBig([]int{54, 65, 75, 74}, AdditionBig, 6, useRelief, 19, 2, 0)
	temp[2] = *NewMonkeyBig([]int{79, 60, 97}, PowerBig, 0, useRelief, 13, 1, 3)
	temp[3] = *NewMonkeyBig([]int{74}, AdditionBig, 3, useRelief, 17, 0, 1)

	return temp
}

func NewMonkey(itemWorryLevels []int, opFN Operation[int], opVal2 int, useRelief bool, testVal2 int, testResTrue int, testResFalse int) *Monkey[int] {
	var temp = Monkey[int]{
		OpFN:            opFN,
		OpVal2:          opVal2,
		UseRelief:       useRelief,
		ReliefFN:        CalcRelief,
		TestFN:          IsDivisable,
		TestVal2:        testVal2,
		TestResultTrue:  testResTrue,
		TestResultFalse: testResFalse,
	}
	temp.Items = make([]Item[int], len(itemWorryLevels))
	for idx, level := range itemWorryLevels {
		temp.Items[idx].WorryLevel = level
	}

	return &temp
}

func NewMonkeyBig(itemWorryLevels []int, opFN Operation[*big.Int], opVal2 int, useRelief bool, testVal2 int, testResTrue int, testResFalse int) *Monkey[*big.Int] {
	var temp = Monkey[*big.Int]{
		OpFN:            opFN,
		OpVal2:          big.NewInt(int64(opVal2)),
		UseRelief:       useRelief,
		ReliefFN:        CalcReliefBig,
		TestFN:          IsDivisableBig,
		TestVal2:        big.NewInt(int64(testVal2)),
		TestResultTrue:  testResTrue,
		TestResultFalse: testResFalse,
	}
	temp.Items = make([]Item[*big.Int], len(itemWorryLevels))
	for idx, level := range itemWorryLevels {
		temp.Items[idx].WorryLevel = big.NewInt(int64(level))
	}

	return &temp
}

//-----------------------------------------------------------------------------

func CreateTestMonkeyGroupModulo(useRelief bool) []Monkey[ModuloInt] {
	/* test
	var temp = make([]Monkey[ModuloInt], 4)

	temp[0] = *NewMonkeyModulo([]int{79, 98}, MultiplicationModulo, 19, useRelief, 23, 2, 3)
	temp[1] = *NewMonkeyModulo([]int{54, 65, 75, 74}, AdditionModulo, 6, useRelief, 19, 2, 0)
	temp[2] = *NewMonkeyModulo([]int{79, 60, 97}, PowerModulo, 0, useRelief, 13, 1, 3)
	temp[3] = *NewMonkeyModulo([]int{74}, AdditionModulo, 3, useRelief, 17, 0, 1)
	*/

	var temp = make([]Monkey[ModuloInt], 8)
	temp[0] = *NewMonkeyModulo([]int{99, 67, 92, 61, 83, 64, 98}, MultiplicationModulo, 17, useRelief, 3, 4, 2)
	temp[1] = *NewMonkeyModulo([]int{78, 74, 88, 89, 50}, MultiplicationModulo, 11, useRelief, 5, 3, 5)
	temp[2] = *NewMonkeyModulo([]int{98, 91}, AdditionModulo, 4, useRelief, 2, 6, 4)
	temp[3] = *NewMonkeyModulo([]int{59, 72, 94, 91, 79, 88, 94, 51}, PowerModulo, 0, useRelief, 13, 0, 5)
	temp[4] = *NewMonkeyModulo([]int{95, 72, 78}, AdditionModulo, 7, useRelief, 11, 7, 6)
	temp[5] = *NewMonkeyModulo([]int{76}, AdditionModulo, 8, useRelief, 17, 0, 2)
	temp[6] = *NewMonkeyModulo([]int{69, 60, 53, 89, 71, 88}, AdditionModulo, 5, useRelief, 19, 7, 1)
	temp[7] = *NewMonkeyModulo([]int{72, 54, 63, 80}, AdditionModulo, 3, useRelief, 7, 1, 3)

	return temp
}

func NewMonkeyModulo(itemWorryLevels []int, opFN Operation[ModuloInt], opVal2 int, useRelief bool, testVal2 int, testResTrue int, testResFalse int) *Monkey[ModuloInt] {
	// least common multiple of (23,19,13,17)
	//mod := 96577//test
	mod := 9699690
	var temp = Monkey[ModuloInt]{
		OpFN:            opFN,
		OpVal2:          *NewModuloInt(opVal2, mod),
		UseRelief:       useRelief,
		ReliefFN:        CalcReliefModulo,
		TestFN:          IsDivisableModulo,
		TestVal2:        *NewModuloInt(testVal2, mod),
		TestResultTrue:  testResTrue,
		TestResultFalse: testResFalse,
	}
	temp.Items = make([]Item[ModuloInt], len(itemWorryLevels))
	for idx, level := range itemWorryLevels {
		temp.Items[idx].WorryLevel = *NewModuloInt(level, mod)
	}

	return &temp
}

type ModuloInt struct {
	val int
	mod int
}

func NewModuloInt(val, m int) *ModuloInt {
	return &ModuloInt{val: val, mod: m}
}

func AdditionModulo(val1, val2 ModuloInt) ModuloInt {
	//fmt.Printf("    increases by %d", val2)
	res := val1.val + val2.val
	if (res > val1.val) == (val2.val > 0) {
		res = res % val1.mod
		return *NewModuloInt(res, val1.mod)
	}
	panic("addition overflow")
}

func MultiplicationModulo(val1, val2 ModuloInt) ModuloInt {
	//fmt.Printf("    multiplied by %d", val2)
	res := val1.val * val2.val
	if (res < 0) == ((val1.val < 0) != (val2.val < 0)) {
		if res/val2.val == val1.val {
			res = res % val1.mod
			return *NewModuloInt(res, val1.mod)
		}
	}
	panic("multiplication overflow")
}

func PowerModulo(val1, val2 ModuloInt) ModuloInt {
	//fmt.Printf("    multiplied by %d", val1)
	return MultiplicationModulo(val1, val1)
}

func IsDivisableModulo(val1, val2 ModuloInt) bool {
	//fmt.Printf("    divided by 3 to %d\n", val1.Int64())
	return val1.val%val2.val == 0
}

// NOTE: is not used in the puzzle and have no idea how to do it
func CalcReliefModulo(val ModuloInt) ModuloInt {
	panic("unimplemeted operation")
}

//-----------------------------------------------------------------------------

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func main() {

	/*
		outputhandler.Initialize()
		defer outputhandler.Reset()

		lines := inputhandler.ReadInput()
	*/
	/*
		// Part 1
		monkeyBusinessLevelPart1, err := startStuffSlingingSimianShenanigans(CreateTestMonkeyGroup(true), 20)
		if err != nil {
			fmt.Printf("Error doing part 1 stuff-slinging simian shenanigans: %v", err)
			os.Exit(int(inputhandler.ErrorCodeProcessing))
		}
	*/
	// Part 2
	//monkeyBusinessLevelPart2, err := startStuffSlingingSimianShenanigans(CreateTestMonkeyGroupBig(false), 10000)
	monkeyBusinessLevelPart2, err := startStuffSlingingSimianShenanigans(CreateTestMonkeyGroupModulo(false), 10000)
	if err != nil {
		fmt.Printf("Error doing part 2 stuff-slinging simian shenanigans: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	fmt.Printf("Result - Part1: %d, Part2: %d", monkeyBusinessLevelPart2, monkeyBusinessLevelPart2)
}

func startStuffSlingingSimianShenanigans[T any](monkeys []Monkey[T], maxRounds int) (int, error) {

	if len(monkeys) < 2 {
		return 0, fmt.Errorf("not enough monkeys for shenanigans '%d'", len(monkeys))
	}

	for round := 1; round <= maxRounds; round++ {
		fmt.Printf("round %d:\n", round)

		for monkeyIdx := range monkeys {
			//fmt.Printf("monkey %d:\n", monkeyIdx)

			for {
				if err := monkeys[monkeyIdx].InspectFirst(); err != nil {
					break // only ErrorOutOfItems possible for now
				}

				if item, toMonkeyIdx, err := monkeys[monkeyIdx].ThrowFirst(); err != nil {
					return 0, fmt.Errorf("monkey '%d' tried to throw with empty hands on round '%d'", monkeyIdx, round) // shouldn't be possible
				} else {
					monkeys[toMonkeyIdx].Catch(item)
				}
			}
		}
		/*
			if round == 1000 {
				fmt.Printf("after round %d:\n", round)
				for monkeyIdx, monkey := range monkeys {
					fmt.Printf("monkey %d: %d\n", monkeyIdx, monkey.Throws)
				}
				fmt.Println()
			}
		*/
		/*
			fmt.Printf("after round %d:\n", round)
			for monkeyIdx, monkey := range monkeys {
				fmt.Printf("monkey %d: ", monkeyIdx)
				for _, item := range monkey.Items {
					fmt.Printf("%v, ", item.WorryLevel)
				}
				fmt.Println()
			}
			fmt.Println()
		*/
	}

	fmt.Printf("after round %d:\n", maxRounds)
	for monkeyIdx, monkey := range monkeys {
		fmt.Printf("monkey %d: %d\n", monkeyIdx, monkey.Throws)
	}
	fmt.Println()

	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].Throws < monkeys[j].Throws
	})

	return monkeys[len(monkeys)-1].Throws * monkeys[len(monkeys)-2].Throws, nil
}

//-----------------------------------------------------------------------------

// Operation is what the inspection does to worry level
type Operation[T any] func(T, T) T

func Addition(val1, val2 int) int {
	//fmt.Printf("    increases by %d", val2)
	res := val1 + val2
	if (res > val1) == (val2 > 0) {
		return res
	}
	panic("addition overflow")
}

func Multiplication(val1, val2 int) int {
	//fmt.Printf("    multiplied by %d", val2)
	res := val1 * val2
	if (res < 0) == ((val1 < 0) != (val2 < 0)) {
		if res/val2 == val1 {
			return res
		}
	}
	panic("multiplication overflow")
}

func Power(val1, val2 int) int {
	//fmt.Printf("    multiplied by %d", val1)
	return val1 * val1
}

func AdditionBig(val1, val2 *big.Int) *big.Int {
	//fmt.Printf("    increases by %d", val2.Int64())
	var res big.Int
	return res.Add(val1, val2)
}

func MultiplicationBig(val1, val2 *big.Int) *big.Int {
	//fmt.Printf("    multiplied by %d", val2.Int64())
	var res big.Int
	return res.Mul(val1, val2)
}

func PowerBig(val1, val2 *big.Int) *big.Int {
	//fmt.Printf("    multiplied by %d", val1.Int64())
	var res big.Int
	return res.Mul(val1, val1)
}

// ModTest tests for divisibility
type ModTest[T any] func(T, T) bool

func IsDivisable(val1, val2 int) bool {
	//fmt.Printf("    divided by 3 to %d\n", val1)
	return val1%val2 == 0
}
func IsDivisableBig(val1, val2 *big.Int) bool {
	//fmt.Printf("    divided by 3 to %d\n", val1.Int64())
	var res big.Int
	return res.Mod(val1, val2).Cmp(big.NewInt(0)) == 0
}

// Relief returns the worry level after relief
type Relief[T any] func(T) T

func CalcRelief(val int) int {
	return val / 3
}

// NOTE: is not used in the puzzle but for completeness sake
func CalcReliefBig(val *big.Int) *big.Int {
	var res big.Int
	return res.Div(val, big.NewInt(3))
}

type Item[T any] struct {
	WorryLevel T
}

var ErrorOutOfItems = fmt.Errorf("out of items")

type Monkey[T any] struct {
	Items           []Item[T]
	OpFN            Operation[T]
	OpVal2          T
	UseRelief       bool
	ReliefFN        Relief[T]
	TestFN          ModTest[T]
	TestVal2        T
	TestResultTrue  int
	TestResultFalse int
	Throws          int
}

func (m *Monkey[T]) InspectFirst() error {

	if len(m.Items) == 0 {
		return ErrorOutOfItems
	}

	oldWorryLevel := m.Items[0].WorryLevel
	//fmt.Printf("  inspects an item with a worry level: %d\n", oldWorryLevel)

	newWorryLevel := m.OpFN(oldWorryLevel, m.OpVal2)
	//fmt.Printf(" to %d\n", newWorryLevel)
	if m.UseRelief {
		m.Items[0].WorryLevel = m.ReliefFN(newWorryLevel)
	} else {
		m.Items[0].WorryLevel = newWorryLevel
	}

	return nil
}

func (m *Monkey[T]) ThrowFirst() (Item[T], int, error) {

	if len(m.Items) == 0 {
		return Item[T]{}, 0, ErrorOutOfItems
	}

	itemToThrow := m.Items[0]
	m.Items = m.Items[1:]

	var throwTo int
	if m.TestFN(itemToThrow.WorryLevel, m.TestVal2) {
		//fmt.Printf("    is divisible by %d\n", m.TestDivision)
		throwTo = m.TestResultTrue
	} else {
		//fmt.Printf("    not divisible by %d\n", m.TestDivision)
		throwTo = m.TestResultFalse
	}

	m.Throws++
	//fmt.Printf("    throw %d to monkey %d\n", itemToThrow.WorryLevel, throwTo)
	return itemToThrow, throwTo, nil
}

func (m *Monkey[T]) Catch(item Item[T]) {
	m.Items = append(m.Items, item)
}
