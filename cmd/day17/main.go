package main

import (
	"AoC22/internal/inputhandler"
	"fmt"
	"math"
	"strconv"
)

// Shapes
type Shape [][]rune

var (
	FallenHero = Shape{
		{'#', '#', '#', '#'},
	}
	SiameseTweewee = Shape{
		{'.', '#', '.'},
		{'#', '#', '#'},
		{'.', '#', '.'},
	}
	Ricky = Shape{
		{'.', '.', '#'},
		{'.', '.', '#'},
		{'#', '#', '#'},
	}
	Hero = Shape{
		{'#'},
		{'#'},
		{'#'},
		{'#'},
	}
	Smashboy = Shape{
		{'#', '#'},
		{'#', '#'},
	}

	ShapeList = []Shape{
		FallenHero, SiameseTweewee, Ricky, Hero, Smashboy,
	}
)

type CellType rune

const (
	Air          CellType = '.'
	StaticBlock  CellType = '#'
	FallingBlock CellType = '@'
)

type JetDirection rune

const (
	Left  JetDirection = '<'
	Right JetDirection = '>'
)

const (
	shapeStartVGap = +3
	shapeStartHGap = +2
)

//-----------------------------------------------------------------------------

func main() {

	lines := inputhandler.ReadInput()
	jets := lines[0]
	chamber := *NewVerticalChamber()
	highestPointPart1 := letTheBlocksFall(jets, ShapeList, chamber, 2022, true)

	chamber = *NewVerticalChamber()
	highestPointPart2 := letTheBlocksFall(jets, ShapeList, chamber, 1000000000000, true)

	fmt.Printf("Result - Part1: %d, Part2: %d", highestPointPart1, highestPointPart2)
}

//-----------------------------------------------------------------------------

type MoveResult string

const (
	WallHit  MoveResult = "Hit the wall"
	BlockHit MoveResult = "Hit other blocks"
	Success  MoveResult = "Success"
)

type VerticalChamber struct {
	Field         [][]byte // 0 idx is top
	Width, Height int

	CurrShape          Shape
	CurrShapeBottomIdx int
	CurrShapeLeftIdx   int
}

func NewVerticalChamber() *VerticalChamber {

	chamber := VerticalChamber{}

	chamber.Height = 30
	chamber.Width = 7

	chamber.Field = make([][]byte, chamber.Height)
	for vIdx := range chamber.Field {
		chamber.Field[vIdx] = make([]byte, chamber.Width)
		for hIdx := range chamber.Field[vIdx] {
			chamber.Field[vIdx][hIdx] = byte(Air)
		}
	}

	chamber.CurrShape = nil
	chamber.CurrShapeBottomIdx = -1
	chamber.CurrShapeLeftIdx = -1

	return &chamber
}

// position, not index!
func (vch *VerticalChamber) SpawnAt(vPos, hPos int, shape Shape) bool {

	if vch.CheckCollisionAt(hPos, vPos, vch.CurrShape) {
		return false
	}

	vch.CurrShape = shape
	vch.CurrShapeBottomIdx = vPos
	vch.CurrShapeLeftIdx = hPos

	chamberVIdx := vch.convertVPosToIdx(vch.CurrShapeBottomIdx + len(shape))
	chamberHIdx := vch.CurrShapeLeftIdx

	for shapeVIdx, shapeLine := range shape {
		for shapeHIdx, shapeBlock := range shapeLine {
			if shapeBlock == rune(StaticBlock) {
				vch.Field[chamberVIdx+shapeVIdx][chamberHIdx+shapeHIdx] = byte(FallingBlock)
			}
		}
	}

	return true
}

func (vch *VerticalChamber) MoveLeft() MoveResult {

	newLeftIdx := vch.CurrShapeLeftIdx - 1

	// reached the wall?
	if newLeftIdx < 0 {
		return WallHit
	}

	// hit a static block?
	if vch.CheckCollisionAt(newLeftIdx, vch.CurrShapeBottomIdx, vch.CurrShape) {
		return BlockHit
	}

	// copy
	for shapeVIdx := range vch.CurrShape {
		chamberVIdx := vch.convertVIdxToIdx(vch.CurrShapeBottomIdx + shapeVIdx)

		for cellIdx := 0; cellIdx < vch.Width; cellIdx++ {

			// clear old
			if vch.Field[chamberVIdx][cellIdx] == byte(FallingBlock) {
				vch.Field[chamberVIdx][cellIdx] = byte(Air)
			}

			// copy to new pos
			if cellIdx+1 == vch.Width {
				continue
			}
			if vch.Field[chamberVIdx][cellIdx+1] == byte(FallingBlock) {
				vch.Field[chamberVIdx][cellIdx] = byte(FallingBlock)
			}
		}
	}

	vch.CurrShapeLeftIdx = newLeftIdx
	return Success
}

func (vch *VerticalChamber) MoveRight() MoveResult {

	newLeftIdx := vch.CurrShapeLeftIdx + 1

	// reached the wall?
	if newLeftIdx+len(vch.CurrShape[0]) > vch.Width {
		return WallHit
	}

	// hit a static block?
	if vch.CheckCollisionAt(newLeftIdx, vch.CurrShapeBottomIdx, vch.CurrShape) {
		return BlockHit
	}

	// copy
	for shapeVIdx := range vch.CurrShape {
		chamberVIdx := vch.convertVIdxToIdx(vch.CurrShapeBottomIdx + shapeVIdx)

		for cellIdx := vch.Width - 1; cellIdx >= 0; cellIdx-- {

			// clear old
			if vch.Field[chamberVIdx][cellIdx] == byte(FallingBlock) {
				vch.Field[chamberVIdx][cellIdx] = byte(Air)
			}

			// copy to new pos
			if cellIdx-1 < 0 {
				continue
			}
			if vch.Field[chamberVIdx][cellIdx-1] == byte(FallingBlock) {
				vch.Field[chamberVIdx][cellIdx] = byte(FallingBlock)
			}
		}
	}

	vch.CurrShapeLeftIdx = newLeftIdx
	return Success
}

func (vch *VerticalChamber) MoveDown() MoveResult {

	newBottomIdx := vch.CurrShapeBottomIdx - 1

	// hit chamber bottom
	if newBottomIdx < 0 {
		return WallHit
	}

	// hit a static block?
	if vch.CheckCollisionAt(vch.CurrShapeLeftIdx, newBottomIdx, vch.CurrShape) {
		return BlockHit
	}

	// do fall
	chamberVIdxStart := vch.convertVIdxToIdx(newBottomIdx)
	chamberVIdxEnd := vch.convertVIdxToIdx(newBottomIdx + len(vch.CurrShape))
	for chamberVIdx := chamberVIdxStart; chamberVIdx >= chamberVIdxEnd; chamberVIdx-- {

		for shapeHIdx := range vch.CurrShape[0] {
			chamberHIdx := vch.CurrShapeLeftIdx + shapeHIdx

			// clear old
			if vch.Field[chamberVIdx][chamberHIdx] == byte(FallingBlock) {
				vch.Field[chamberVIdx][chamberHIdx] = byte(Air)
			}

			if chamberVIdx-1 == vch.Height {
				continue
			}
			if vch.Field[chamberVIdx-1][chamberHIdx] == byte(FallingBlock) {
				vch.Field[chamberVIdx][chamberHIdx] = byte(FallingBlock)
			}
		}
	}

	vch.CurrShapeBottomIdx = newBottomIdx
	return Success
}

func (vch *VerticalChamber) Solidify() {

	chamberVIdx := vch.convertVPosToIdx(vch.CurrShapeBottomIdx + len(vch.CurrShape))
	chamberHIdx := vch.CurrShapeLeftIdx

	for shapeVIdx, shapeLine := range vch.CurrShape {
		for shapeHIdx, shapeBlock := range shapeLine {
			if shapeBlock == rune(StaticBlock) {
				vch.Field[(chamberVIdx + shapeVIdx)][chamberHIdx+shapeHIdx] = byte(StaticBlock)
			}
		}
	}
}

func (vch *VerticalChamber) CheckCollisionAt(shapeLeftIdx int, shapeBottomIdx int, shape Shape) bool {

	for shapeVIdx := range shape {
		chamberVIdx := vch.convertVIdxToIdx((shapeBottomIdx + (len(shape) - 1)) - shapeVIdx)

		for shapeHIdx := range shape[shapeVIdx] {
			chamberHIdx := shapeLeftIdx + shapeHIdx

			// until finds a statick block in shape
			if shape[shapeVIdx][shapeHIdx] != rune(StaticBlock) {
				continue
			}

			// if collides
			if vch.Field[chamberVIdx][chamberHIdx] == byte(StaticBlock) {
				return true
			}
		}
	}

	return false
}

// position, not index!
func (vch *VerticalChamber) FindHighestBlock() int {

	for vIdx, chamberLine := range vch.Field {
		for _, cell := range chamberLine {
			if cell != byte(Air) {
				return vch.Height - vIdx
			}
		}
	}

	return 0 // nothing in field
}

func (vch *VerticalChamber) ExtendChamber() {

	extendBy := 30

	chamberExtension := make([][]byte, extendBy)
	for vIdx := range chamberExtension {
		chamberExtension[vIdx] = make([]byte, vch.Width)
		for hIdx := range chamberExtension[vIdx] {
			chamberExtension[vIdx][hIdx] = byte(Air)
		}
	}
	chamberExtension = append(chamberExtension, vch.Field...)
	vch.Height += extendBy

	vch.Field = chamberExtension
}

// position -> chamber index!
func (vch *VerticalChamber) convertVPosToIdx(vPos int) int {
	return vch.Height - vPos
}

// index -> chamber index!
func (vch *VerticalChamber) convertVIdxToIdx(vIdx int) int {
	return (vch.Height - 1) - vIdx
}

//-Part 2 extensions-----------------------------------------------------------

func (vch *VerticalChamber) GetFreeHeightMap() []int {

	heightMap := make([]int, 7)
	for idx := range heightMap {
		heightMap[idx] = math.MaxInt
	}

	highestBlockIdx := vch.convertVPosToIdx(vch.FindHighestBlock())
	for hIdx := 0; hIdx < vch.Width; hIdx++ {
		for vIdx := highestBlockIdx; vIdx < vch.Height; vIdx++ {
			if vch.Field[vIdx][hIdx] == byte(StaticBlock) {
				heightMap[hIdx] = vIdx - highestBlockIdx
				break
			}
		}
	}

	return heightMap
}

type Hash struct {
	heightMap    string
	nextShapeIdx int
	nextJetIdx   int
}

type State struct {
	shapesFallen int
	highestPoint int
}

//-----------------------------------------------------------------------------

func letTheBlocksFall(jets string, shapeList []Shape, chamber VerticalChamber, maxShapesToFall int, withPatternDetection bool) int {

	highestPoint := 0

	currShapeIdx := 0
	currJetIdx := 0

	chamberState := make(map[Hash]State) // fallen shape count
	var shapesFallenOffset int = 0
	var highestPointOffset int = 0

	shapesFallen := 0
	haveFallingRock := false
	for {

		// add new block
		if !haveFallingRock {

			// find highest point
			highestPoint = chamber.FindHighestBlock()

			// reached the limit
			if shapesFallen+shapesFallenOffset == maxShapesToFall {
				break
			}

			// Part 2
			if withPatternDetection {

				var heightMap string
				for _, height := range chamber.GetFreeHeightMap() {
					heightMap += strconv.Itoa(height) + ","
				}
				hash := Hash{
					heightMap:    heightMap,
					nextShapeIdx: currShapeIdx,
					nextJetIdx:   currJetIdx,
				}

				state, ok := chamberState[hash]
				if ok { // found the pattern

					// fast forward
					periodShapeCount := shapesFallen - state.shapesFallen
					skippingPeriods := (maxShapesToFall - shapesFallen) / periodShapeCount
					shapesFallenOffset = skippingPeriods * periodShapeCount
					periodHighestPoint := highestPoint - state.highestPoint
					highestPointOffset = skippingPeriods * periodHighestPoint

					withPatternDetection = false // done with this
				} else {
					chamberState[hash] = State{shapesFallen: shapesFallen, highestPoint: highestPoint}
				}
			}

			// extend chamber sim area if needed
			neededHeight := highestPoint + (shapeStartVGap + 1) + len(shapeList[currShapeIdx])
			if chamber.Height <= neededHeight {
				chamber.ExtendChamber()
			}

			// add block
			shapeBottomPos := highestPoint + shapeStartVGap
			shapeLeftPos := shapeStartHGap
			chamber.SpawnAt(shapeBottomPos, shapeLeftPos, shapeList[currShapeIdx])

			haveFallingRock = true

			//fmt.Println("new shape added")
			//visualize(chamber.Field, shapesFallen)
		}

		// apply jet
		switch jets[currJetIdx] {
		case byte(Left):
			moveResult := chamber.MoveLeft()
			_ = moveResult
			//fmt.Printf("move shape to the left - %s\n", moveResult)

		case byte(Right):
			moveResult := chamber.MoveRight()
			_ = moveResult
			//fmt.Printf("move shape to the right - %s\n", moveResult)
		}
		//visualize(chamber.Field, shapesFallen)
		currJetIdx = (currJetIdx + 1) % len(jets)

		// fall check
		moveResult := chamber.MoveDown()
		//fmt.Printf("move shape down - %s\n", moveResult)
		if moveResult != Success {

			chamber.Solidify()
			haveFallingRock = false
			shapesFallen++
			fmt.Printf("%.2f%%\n", float64(shapesFallen+shapesFallenOffset)/float64(maxShapesToFall)*100)

			currShapeIdx = (currShapeIdx + 1) % len(ShapeList)
		}
		//visualize(chamber.Field, shapesFallen)
	}

	//visualize(chamber.Field, shapesFallen)
	return highestPoint + highestPointOffset
}

func visualize(lines [][]byte, blockNum int) {
	/*
		if blockNum != maxShapesFallen-1 {
			return
		}
	*/
	for lineIdx, line := range lines {
		fmt.Println(string(line), len(lines)-lineIdx)
	}
	fmt.Println()
}
