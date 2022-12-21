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

func (mo *MathOperation) IsSolved() bool {
	return mo.val1Eval && mo.val2Eval
}

func (mo *MathOperation) SetFirstOperand(val int) {
	mo.val1 = val
	mo.val1Eval = true
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

func (mo *MathOperation) UnSolve(result int) (int, int, error) {

	if !mo.val1Eval && !mo.val2Eval {
		return 0, 0, ErrorUnresolvedOperands
	}

	switch mo.Op {
	case Addition:
		if !mo.val1Eval {
			return SubInt(result, mo.val2), mo.val2, nil
		}
		if !mo.val2Eval {
			return mo.val1, SubInt(result, mo.val1), nil
		}
	case Subtraction:
		if !mo.val1Eval {
			return AddInt(result, mo.val2), mo.val2, nil
		}
		if !mo.val2Eval {
			return mo.val1, SubInt(mo.val1, result), nil
		}
	case Multiplication:
		if !mo.val1Eval {
			return DivInt(result, mo.val2), mo.val2, nil
		}
		if !mo.val2Eval {
			return mo.val1, DivInt(result, mo.val1), nil
		}
	case Division:
		if !mo.val1Eval {
			return MulInt(result, mo.val2), mo.val2, nil
		}
		if !mo.val2Eval {
			return mo.val1, DivInt(mo.val1, result), nil
		}
	}

	return 0, 0, ErrorInvalidOperation
}

type Monkey struct {
	Name     string
	Job      MathOperation
	Value    int
	hasValue bool
	hasJob   bool
}

func NewMonkey(name string) *Monkey {
	return &Monkey{
		Name:     name,
		hasValue: false,
		hasJob:   false,
	}
}

func (m *Monkey) SetOperation(job MathOperation) {
	m.Job = job
	m.hasJob = true
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

func (m *Monkey) SetValue(val int) {
	m.Value = val
	m.hasValue = true
}

func (m *Monkey) ClearValue() {
	m.Value = 0
	m.hasValue = false
}

func (m *Monkey) HasJob() bool {
	return m.hasJob
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

	// TODO: ugly, but time...
	monkeys, err = parseMonkeys(lines)
	if err != nil {
		fmt.Printf("Error parsing monkeys: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	resultPart2, err := findMyAnswer(monkeys)
	if err != nil {
		fmt.Printf("Error finding my answer: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}

	fmt.Printf("Result - Part1: %d, Part2: %d", resultPart1, resultPart2)
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

		err = resolvePass(monkeys)
		if err != nil {
			return 0, err
		}
	}

	return rootMonkeyValue, nil
}

func resolvePass(monkeys map[string]*Monkey) error {

	isChanged := false
	for _, thisMonkey := range monkeys {

		// already solved
		if thisMonkey.HasValue() {
			continue
		}

		// try to solve
		result, err := thisMonkey.Job.Solve()
		if err != nil {

			if err != ErrorUnresolvedOperands {
				return fmt.Errorf("error resolving monkey '%s': %w", thisMonkey.Name, err)
			}
			// needs the operands so try to solve them

			if val1Monkey, found := monkeys[thisMonkey.Job.Val1Ref]; !found {
				return fmt.Errorf("unknown monkey '%s' referenced by '%s'", thisMonkey.Job.Val1Ref, thisMonkey.Name)
			} else {
				val1, val1Err := val1Monkey.GetValue()
				if val1Err == nil {
					thisMonkey.Job.SetFirstOperand(val1)
					isChanged = true
				}
			}

			if val2Monkey, found := monkeys[thisMonkey.Job.Val2Ref]; !found {
				return fmt.Errorf("unknown monkey '%s' referenced by '%s'", thisMonkey.Job.Val2Ref, thisMonkey.Name)
			} else {
				val2, val2Err := val2Monkey.GetValue()
				if val2Err == nil {
					thisMonkey.Job.SetSecondOperand(val2)
					isChanged = true
				}
			}

			continue
		}

		// done with it
		thisMonkey.SetValue(result)
		isChanged = true
	}

	if isChanged {
		return nil
	}

	return fmt.Errorf("no advancement happened")
}

// Part 2
func findMyAnswer(monkeys map[string]*Monkey) (int, error) {

	rootMonkey, found := monkeys["root"]
	if !found {
		return 0, fmt.Errorf("missing root monkey")
	}

	myName := "humn"
	myMonkey, found := monkeys[myName]
	if !found {
		return 0, fmt.Errorf("couldn't find monkey '%s", myName)
	}

	// these shouldn't be able to evaluate
	mySide, err := getDependsOn(myMonkey, monkeys)
	if err != nil {
		return 0, fmt.Errorf("failed to find dependents for '%s': %w", myName, err)
	}

	// one ref is the stop condition for below
	var otherSideRootMonkey *Monkey
	if _, found := mySide[rootMonkey.Job.Val1Ref]; !found {
		otherSideRootMonkey, found = monkeys[rootMonkey.Job.Val1Ref]
		if !found {
			return 0, fmt.Errorf("root reference '%s' to unknown monkey", rootMonkey.Job.Val1Ref)
		}
	} else if _, found := mySide[rootMonkey.Job.Val2Ref]; !found {
		otherSideRootMonkey, found = monkeys[rootMonkey.Job.Val2Ref]
		if !found {
			return 0, fmt.Errorf("root reference '%s' to unknown monkey", rootMonkey.Job.Val2Ref)
		}
	} else {
		return 0, fmt.Errorf("both references in root depends on '%s'", myName)
	}

	// this should get the value we need for mySide
	var rootMonkeyValue int
	for {

		val, err := otherSideRootMonkey.GetValue()
		if err == nil {
			rootMonkeyValue = val
			break
		}

		err = resolvePassWithSkips(monkeys, myName)
		if err != nil {
			return 0, err
		}
	}

	// setup for reverse resolve
	var mySideRootMonkey *Monkey
	if val1Monkey, found := mySide[rootMonkey.Job.Val1Ref]; found {
		mySideRootMonkey = val1Monkey
	} else if val2Monkey, found := mySide[rootMonkey.Job.Val2Ref]; found {
		mySideRootMonkey = val2Monkey
	} else {
		return 0, fmt.Errorf("neither of root references depends on '%s'", myName)
	}
	mySideRootMonkey.SetValue(rootMonkeyValue)
	delete(mySide, "root")

	// this should get my value that would let to root's needed value
	myMonkey.hasJob = false
	myMonkey.hasValue = false
	var myValue int
	for {

		val, err := myMonkey.GetValue()
		if err == nil {
			myValue = val
			break
		}

		err = unresolvePass(mySide, monkeys)
		if err != nil {
			return 0, err
		}
	}

	return myValue, nil
}

func resolvePassWithSkips(monkeys map[string]*Monkey, skipName string) error {

	isChanged := false
	for _, thisMonkey := range monkeys {

		if thisMonkey.Name == skipName {
			continue
		}

		// already solved
		if thisMonkey.HasValue() {
			continue
		}

		// try to solve
		result, err := thisMonkey.Job.Solve()
		if err != nil {

			if err != ErrorUnresolvedOperands {
				return fmt.Errorf("error resolving monkey '%s': %w", thisMonkey.Name, err)
			}
			// needs the operands so try to solve them

			if thisMonkey.Job.Val1Ref != skipName {
				if val1Monkey, found := monkeys[thisMonkey.Job.Val1Ref]; !found {
					return fmt.Errorf("unknown monkey '%s' referenced by '%s'", thisMonkey.Job.Val1Ref, thisMonkey.Name)
				} else {
					val1, val1Err := val1Monkey.GetValue()
					if val1Err == nil {
						thisMonkey.Job.SetFirstOperand(val1)
						isChanged = true
					}
				}
			}

			if thisMonkey.Job.Val2Ref != skipName {
				if val2Monkey, found := monkeys[thisMonkey.Job.Val2Ref]; !found {
					return fmt.Errorf("unknown monkey '%s' referenced by '%s'", thisMonkey.Job.Val2Ref, thisMonkey.Name)
				} else {
					val2, val2Err := val2Monkey.GetValue()
					if val2Err == nil {
						thisMonkey.Job.SetSecondOperand(val2)
						isChanged = true
					}
				}
			}

			continue
		}

		// done with it
		thisMonkey.SetValue(result)
		isChanged = true
	}

	if isChanged {
		return nil
	}

	return fmt.Errorf("no advancement happened")
}

func unresolvePass(unresolveMonkeys map[string]*Monkey, monkeys map[string]*Monkey) error {

	isChanged := false
	for _, thisMonkey := range unresolveMonkeys {

		if !thisMonkey.HasJob() {
			continue
		}
		if thisMonkey.Job.IsSolved() {
			continue
		}

		thisVal, err := thisMonkey.GetValue()
		if err != nil {
			continue
		}

		val1, val2, err := thisMonkey.Job.UnSolve(thisVal)
		if err != nil {
			return fmt.Errorf("couldn't unresolve '%s' with result value of '%d': %w", thisMonkey.Name, thisVal, err)
		}
		thisMonkey.Job.SetFirstOperand(val1)
		thisMonkey.Job.SetSecondOperand(val2)

		val1Monkey, found := monkeys[thisMonkey.Job.Val1Ref]
		if !found {
			return fmt.Errorf("unknown reference to monkey '%s' from '%s'", thisMonkey.Job.Val1Ref, thisMonkey.Name)
		}
		val1Monkey.SetValue(val1)

		val2Monkey, found := monkeys[thisMonkey.Job.Val2Ref]
		if !found {
			return fmt.Errorf("unknown reference to monkey '%s' from '%s'", thisMonkey.Job.Val2Ref, thisMonkey.Name)
		}
		val2Monkey.SetValue(val2)

		isChanged = true
	}

	if isChanged {
		return nil
	}

	return fmt.Errorf("no advancement happened")
}

func getDependsOn(startMonkey *Monkey, monkeys map[string]*Monkey) (map[string]*Monkey, error) {

	dependents := make(map[string]*Monkey)
	dependents[startMonkey.Name] = startMonkey

	for {

		// reached root, job is done
		if _, found := dependents["root"]; found {
			break
		}

		// if any of the referenced monkeys referencing a monkey that depends on
		// the start monkey in someway...
		isChanged := false
		for _, thisMonkey := range monkeys {

			if thisMonkey.HasValue() {
				continue
			}

			_, found := dependents[thisMonkey.Job.Val1Ref]
			if found {
				dependents[thisMonkey.Name] = thisMonkey

				isChanged = true
				continue
			}

			_, found = dependents[thisMonkey.Job.Val2Ref]
			if found {
				dependents[thisMonkey.Name] = thisMonkey

				isChanged = true
				continue
			}
		}

		if !isChanged {
			return nil, fmt.Errorf("root is not depending on '%s'", startMonkey.Name)
		}
	}

	return dependents, nil
}
