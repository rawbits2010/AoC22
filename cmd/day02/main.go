package main

import (
	"AoC22/internal/inputhandler"
	"fmt"
	"os"
	"strings"
)

func main() {

	lines := inputhandler.ReadInput()

	scorePart1, err := CalcPart1Score(lines)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	scorePart2, err := CalcPart2Score(lines)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	fmt.Printf("Scores - Part1: %d, Part2: %d\n", scorePart1, scorePart2)
}

//-Part1-----------------------------------------------------------------------

func CalcPart1Score(lines []string) (int, error) {

	var score int
	for _, line := range lines {

		elfHand, myHand, err := ReadHands(line)
		if err != nil {
			return 0, fmt.Errorf("invalid input in line '%s': %w", line, err)
		}
		score += HandPoints[myHand]

		outcome, err := ResolveOutcome(elfHand, myHand)
		if err != nil {
			return 0, fmt.Errorf("couldn't resolve outcome for elf hand '%s' and my hand '%s': %w", elfHand, myHand, err)
		}
		score += OutcomePoints[outcome]
	}

	return score, nil
}

func ReadHands(line string) (ElfHand, MyHand, error) {
	hands := strings.Fields(line)

	var elfHand ElfHand = ElfInvalid
	for _, hand := range ValidElfHands {
		if string(hand) == hands[0] {
			elfHand = hand
		}
	}
	if elfHand == ElfInvalid {
		return ElfInvalid, MyInvalid, ErrorInvalidHand
	}

	var myHand MyHand = MyInvalid
	for _, hand := range ValidMyHands {
		if string(hand) == hands[1] {
			myHand = hand
		}
	}
	if myHand == MyInvalid {
		return ElfInvalid, MyInvalid, ErrorInvalidHand
	}

	return elfHand, myHand, nil
}

func ResolveOutcome(elfHand ElfHand, myHand MyHand) (Outcome, error) {

	for _, outcome := range ValidOutcomes {
		if HandForOutcome[elfHand][outcome] == myHand {
			return outcome, nil
		}
	}

	return OutcomeInvalid, ErrorInvalidOutcome
}

//-Part2-----------------------------------------------------------------------

func CalcPart2Score(lines []string) (int, error) {

	var score int
	for _, line := range lines {

		elfHand, expectedOutcome, err := ReadRoundPlan(line)
		if err != nil {
			return 0, fmt.Errorf("invalid input in line '%s': %w", line, err)
		}

		outcomePos := OutcomePosition[expectedOutcome]
		expectedHand := HandForOutcome[elfHand][outcomePos]
		score += HandPoints[expectedHand]

		outcome, err := ResolveOutcome(elfHand, expectedHand)
		if err != nil {
			return 0, fmt.Errorf("couldn't resolve outcome for elf hand '%s' and expected hand '%s': %w", elfHand, expectedHand, err)
		}
		score += OutcomePoints[outcome]
	}

	return score, nil
}

func ReadRoundPlan(line string) (ElfHand, ExpectedOutcome, error) {
	plan := strings.Fields(line)

	var elfHand ElfHand = ElfInvalid
	for _, hand := range ValidElfHands {
		if string(hand) == plan[0] {
			elfHand = hand
		}
	}
	if elfHand == ElfInvalid {
		return ElfInvalid, ExpectedInvalid, ErrorInvalidHand
	}

	var expectedOutcome ExpectedOutcome = ExpectedInvalid
	for _, outcome := range ValidExpectedOutcomes {
		if string(outcome) == plan[1] {
			expectedOutcome = outcome
		}
	}
	if expectedOutcome == ExpectedInvalid {
		return ElfInvalid, ExpectedInvalid, ErrorInvalidOutcome
	}

	return elfHand, expectedOutcome, nil
}

//-----------------------------------------------------------------------------

var ErrorInvalidHand = fmt.Errorf("invalid hand")
var ErrorInvalidOutcome = fmt.Errorf("invalid outcome")

type ElfHand string

const (
	ElfRock     ElfHand = "A"
	ElfPaper    ElfHand = "B"
	ElfScissors ElfHand = "C"
	ElfInvalid  ElfHand = ""
)

var ValidElfHands = []ElfHand{ElfRock, ElfPaper, ElfScissors}

type MyHand string

const (
	MyRock     MyHand = "X"
	MyPaper    MyHand = "Y"
	MyScissors MyHand = "Z"
	MyInvalid  MyHand = ""
)

var ValidMyHands = []MyHand{MyRock, MyPaper, MyScissors}

var HandPoints = map[MyHand]int{
	MyRock:     1,
	MyPaper:    2,
	MyScissors: 3,
}

type Outcome int

const (
	Loose          Outcome = 0
	Draw           Outcome = 1
	Win            Outcome = 2
	OutcomeInvalid Outcome = -1
)

var ValidOutcomes = []Outcome{Loose, Draw, Win}

var OutcomePoints = map[Outcome]int{
	Loose: 0,
	Draw:  3,
	Win:   6,
}

var HandForOutcome = map[ElfHand][]MyHand{
	ElfRock:     {MyScissors, MyRock, MyPaper},
	ElfPaper:    {MyRock, MyPaper, MyScissors},
	ElfScissors: {MyPaper, MyScissors, MyRock},
}

//
// Part 2

type ExpectedOutcome string

const (
	ExpectedLoose   ExpectedOutcome = "X"
	ExpectedDraw    ExpectedOutcome = "Y"
	ExpectedWin     ExpectedOutcome = "Z"
	ExpectedInvalid ExpectedOutcome = ""
)

var ValidExpectedOutcomes = []ExpectedOutcome{ExpectedLoose, ExpectedDraw, ExpectedWin}

var OutcomePosition = map[ExpectedOutcome]Outcome{
	ExpectedLoose: Loose,
	ExpectedDraw:  Draw,
	ExpectedWin:   Win,
}
