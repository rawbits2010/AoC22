package main

import (
	"AoC22/internal/inputhandler"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Operation string

const (
	Addition       Operation = "+"
	Subtraction    Operation = "-"
	Multiplication Operation = "*"
	Division       Operation = "/"
)

const OpList string = string(Addition + Subtraction + Multiplication + Division)

var ErrorInvalidOperation = fmt.Errorf("invalid operation")
var ErrorUnresolvedOperands = fmt.Errorf("unresolved operands")

type MathOperation struct {
	Val1Ref string
	Val2Ref string
	Op      Operation

	val1     int
	val1Eval bool
	val2     int
	val2Eval bool
}

func NewMathOperation(val1Ref, val2Ref string, op Operation) *MathOperation {
	return &MathOperation{
		Val1Ref:  val1Ref,
		Val2Ref:  val2Ref,
		val1Eval: false,
		val2Eval: false,
		Op:       op,
	}
}

func (m *Monkey) HasFirstOperand() bool {
	return m.Job.val1Eval
}

func (mo *MathOperation) SetFirstOperand(val int) {
	mo.val1 = val
	mo.val1Eval = true
}

func (m *Monkey) HasSecondOperand() bool {
	return m.Job.val2Eval
}
func (mo *MathOperation) SetSecondOperand(val int) {
	mo.val2 = val
	mo.val2Eval = true
}

func (mo *MathOperation) Solve() (int, error) {

	if !mo.val1Eval || !mo.val2Eval {
		return 0, ErrorUnresolvedOperands
	}

	switch mo.Op {
	case Addition:
		return AddInt(mo.val1, mo.val2), nil
	case Subtraction:
		return SubInt(mo.val1, mo.val2), nil
	case Multiplication:
		return MulInt(mo.val1, mo.val2), nil
	case Division:
		return DivInt(mo.val1, mo.val2), nil
	}

	return 0, ErrorInvalidOperation
}

type Monkey struct {
	Name     string
	Job      MathOperation
	Value    int
	hasValue bool
}

func NewMonkey(name string) *Monkey {
	return &Monkey{Name: name}
}

func (m *Monkey) SetValue(val int) {
	m.Value = val
	m.hasValue = true
}

func (m *Monkey) SetOperation(job MathOperation) {
	m.Job = job
	m.hasValue = false
}

func (m *Monkey) HasValue() bool {
	return m.hasValue
}

func (m *Monkey) GetValue() (int, error) {
	if m.hasValue {
		return m.Value, nil
	}
	return 0, fmt.Errorf("unresolved job")
}

//-Math------------------------------------------------------------------------

func AddInt(val1, val2 int) int {
	res := val1 + val2
	if (res > val1) == (val2 > 0) {
		return res
	}
	panic("addition overflow")
}

func SubInt(val1, val2 int) int {
	res := val1 - val2
	if (res < val1) == (val2 > 0) {
		return res
	}
	panic("subtraction overflow")
}

func MulInt(val1, val2 int) int {
	res := val1 * val2
	if (res < 0) == ((val1 < 0) != (val2 < 0)) {
		if res/val2 == val1 {
			return res
		}
	}
	panic("multiplication overflow")
}

func DivInt(val1, val2 int) int {
	if val2 == 0 {
		panic("division by zero")
	}
	res := val1 / val2
	if (res < 0) == ((val1 < 0) != (val2 < 0)) { // sign bit check
		if res*val2 != val1 {
			panic("division has a remainder")
		}
		return res
	}
	panic("division failed")
}

//-Main------------------------------------------------------------------------

func main() {

	lines := inputhandler.ReadInput()

	monkeys, err := parseMonkeys(lines)
	if err != nil {
		fmt.Printf("Error parsing monkeys: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	resultPart1, err := resolveMonkeyEquations(monkeys)
	if err != nil {
		fmt.Printf("Error resolving equations: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	fmt.Printf("Result - Part1: %d", resultPart1)
}

//-----------------------------------------------------------------------------

func parseMonkeys(lines []string) (map[string]*Monkey, error) {

	monkeys := make(map[string]*Monkey)
	for lineIdx, line := range lines {

		if len(line) == 0 {
			continue
		}

		keyval := strings.Split(line, ":")
		if len(keyval) != 2 {
			return nil, fmt.Errorf("invalid key-value format '%s' at line %d", line, lineIdx+1)
		}

		mName := strings.TrimSpace(keyval[0])
		monkey := &Monkey{Name: mName}

		if strings.ContainsAny(keyval[1], OpList) {

			equation := strings.Split(strings.TrimSpace(keyval[1]), " ")
			if len(equation) != 3 {
				return nil, fmt.Errorf("invalid equation '%s' at line %d", keyval[1], lineIdx+1)
			}

			job := NewMathOperation(equation[0], equation[2], Operation(equation[1]))
			monkey.SetOperation(*job)

		} else { // invalid operations will fail here

			val, err := strconv.Atoi(strings.TrimSpace(keyval[1]))
			if err != nil {
				return nil, fmt.Errorf("couldn't parse numeric value from '%s' at line %d", keyval[1], lineIdx+1)
			}

			monkey.SetValue(val)
		}

		if _, found := monkeys[mName]; found {
			return nil, fmt.Errorf("duplicate entry for monkey '%s' found at line %d", mName, lineIdx+1)
		}

		monkeys[mName] = monkey
	}

	return monkeys, nil
}

//-----------------------------------------------------------------------------

// Part 1
func resolveMonkeyEquations(monkeys map[string]*Monkey) (int, error) {

	rootMonkey, found := monkeys["root"]
	if !found {
		return 0, fmt.Errorf("missing root monkey")
	}

	var rootMonkeyValue int
	for {

		val, err := rootMonkey.GetValue()
		if err == nil {
			rootMonkeyValue = val
			break
		}

		for _, thisMonkey := range monkeys {

			// already solved
			if thisMonkey.HasValue() {
				continue
			}

			// try to solve
			result, err := thisMonkey.Job.Solve()
			if err != nil {

				if err != ErrorUnresolvedOperands {
					return 0, fmt.Errorf("error resolving monkey '%s': %w", thisMonkey.Name, err)
				}
				// needs the operands so try to solve them

				if val1Monkey, found := monkeys[thisMonkey.Job.Val1Ref]; !found {
					return 0, fmt.Errorf("unknown monkey '%s' referenced by '%s'", thisMonkey.Job.Val1Ref, thisMonkey.Name)
				} else {
					val1, val1Err := val1Monkey.GetValue()
					if val1Err == nil {
						thisMonkey.Job.SetFirstOperand(val1)
					}
				}

				if val2Monkey, found := monkeys[thisMonkey.Job.Val2Ref]; !found {
					return 0, fmt.Errorf("unknown monkey '%s' referenced by '%s'", thisMonkey.Job.Val2Ref, thisMonkey.Name)
				} else {
					val2, val2Err := val2Monkey.GetValue()
					if val2Err == nil {
						thisMonkey.Job.SetSecondOperand(val2)
					}
				}

				continue
			}

			// done with it
			thisMonkey.SetValue(result)
		}

	}

	return rootMonkeyValue, nil
}

// Part 2
