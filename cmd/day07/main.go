package main

import (
	"AoC22/internal/inputhandler"
	"AoC22/internal/outputhandler"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const treeColor = outputhandler.White
const directoryColor = outputhandler.BrightCyan
const fileColor = outputhandler.BrightGreen
const sizeColor = outputhandler.BrightMagenta

func main() {

	outputhandler.Initialize()
	defer outputhandler.Reset()

	lines := inputhandler.ReadInput()

	rootNode, err := parseFilesystem(lines)
	if err != nil {
		fmt.Printf("Error: while parsing filesystem: %v", err)
		os.Exit(int(inputhandler.ErrorCodeProcessing))
	}
	visualizeFileSystem(rootNode) // had to nerd it, not sorry :)

	updateFolderSizes(rootNode)

	// Part 1
	const part1SizeLimit = 100000
	var sizeAtMost = func(node *Node) bool {
		return node.Size <= part1SizeLimit
	}
	foundDirs := getDirsWithCondition(rootNode, sizeAtMost)

	var part1SumSizes int
	for _, dir := range foundDirs {
		part1SumSizes += dir.Size
	}

	// Part 2
	const totalAvailableSpace = 70000000
	const neededFreeSpace = 30000000
	haveFreeSpace := (totalAvailableSpace - rootNode.Size)
	extraSpaceNeeded := neededFreeSpace - haveFreeSpace

	var part2SizeToDelete int
	if extraSpaceNeeded <= 0 {
		part2SizeToDelete = 0 // already have enough
	} else {

		var sizeGraterThan = func(node *Node) bool {
			return node.Size > extraSpaceNeeded
		}
		foundDirs := getDirsWithCondition(rootNode, sizeGraterThan)

		sort.Slice(foundDirs, func(i int, j int) bool {
			return foundDirs[i].Size < foundDirs[j].Size
		})

		part2SizeToDelete = foundDirs[0].Size
	}

	fmt.Printf("Result - Part1: %d, Part2: %d\n", part1SumSizes, part2SizeToDelete)
}

func parseFilesystem(lines []string) (*Node, error) {

	var rootNode = NewNode(Directory, "/", 0)
	var currNode = rootNode // just to be safe

	var currCommand string
	var currCmdArgs []string
	for _, line := range lines {

		// we have a prompt
		var isPrompt bool
		if line[:1] == "$" {
			isPrompt = true

			tokens := strings.Split(line, " ")
			if len(tokens) < 2 {
				return nil, fmt.Errorf("invalid command '%s'", line)
			}
			currCommand = tokens[1]
			currCmdArgs = tokens[2:]
		}

		// process command for this line
		switch currCommand {
		case "cd":

			if len(currCmdArgs) < 1 {
				return nil, fmt.Errorf("too few arguments for 'cd' in line '%s'", line)
			}

			switch currCmdArgs[0] {
			case "/":
				currNode = rootNode

			case "..":
				currNode = currNode.Parent

			default:
				var found = false
				for idx, node := range currNode.ChildNodes {
					if currCmdArgs[0] == node.Name {
						currNode = currNode.ChildNodes[idx]
						found = true
						break
					}
				}

				if !found {
					return nil, fmt.Errorf("invalid child directory referenced '%s' in line '%s'", currCmdArgs[0], line)
				}
			}

		case "ls":
			if !isPrompt { // we need only the data

				tokens := strings.Split(line, " ")
				if len(tokens) < 2 {
					return nil, fmt.Errorf("too few elements for directory listing in line '%s'", line)
				}

				var node *Node
				if tokens[0] == "dir" {
					node = NewNode(Directory, tokens[1], 0)
				} else {
					size, err := strconv.Atoi(tokens[0])
					if err != nil {
						return nil, fmt.Errorf("couldn't parse file size in line '%s'", line)
					}
					node = NewNode(File, tokens[1], size)
				}

				currNode.AddNode(node)
			}

		default:
			return nil, fmt.Errorf("unknown command '%s' in line '%s'", currCommand, line)
		}
	}

	return rootNode, nil
}

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[37m"
	ColorWhite  = "\033[97m"
)

func visualizeFileSystem(node *Node) {
	visualizeFileSystem_recurse(node, make([]bool, 0))
	fmt.Print(ColorReset)
}

func visualizeFileSystem_recurse(node *Node, isLastList []bool) {

	if len(isLastList) > 0 {
		fmt.Print(outputhandler.GetForeground(treeColor))
		for _, isLast := range isLastList[:len(isLastList)-1] {
			if isLast {
				fmt.Printf("    ")
			} else {
				fmt.Printf("│   ")
			}
		}
		if isLastList[len(isLastList)-1] {
			fmt.Printf("└── ")
		} else {
			fmt.Printf("├── ")
		}
	}

	switch node.Type {
	case File:
		fmt.Printf("%s%s %s%d\n", outputhandler.GetForeground(fileColor), node.Name, outputhandler.GetForeground(sizeColor), node.Size)

	case Directory:

		fmt.Print(outputhandler.GetForeground(directoryColor))
		fmt.Printf("%s\n", node.Name)

		for childIdx := range node.ChildNodes {
			visualizeFileSystem_recurse(node.ChildNodes[childIdx], append(isLastList, (childIdx == len(node.ChildNodes)-1)))
		}
	}
}

func updateFolderSizes(node *Node) int {

	var size int
	for idx, chn := range node.ChildNodes {
		if chn.Type == Directory {
			size += updateFolderSizes(node.ChildNodes[idx])
		} else {

			size += chn.Size

		}
	}
	node.Size = size
	return size
}

type conditionFunc func(*Node) bool

func getDirsWithCondition(node *Node, condition conditionFunc) []*Node {
	return getDirsWithCondition_recurse(node, condition, make([]*Node, 0))
}

func getDirsWithCondition_recurse(node *Node, condition conditionFunc, dirs []*Node) []*Node {

	for childIdx, child := range node.ChildNodes {
		if child.Type == Directory {
			dirs = getDirsWithCondition_recurse(node.ChildNodes[childIdx], condition, dirs)
		}
	}

	if condition(node) {
		return append(dirs, node)
	}

	return dirs
}

//-----------------------------------------------------------------------------

type NodeType int

const (
	Directory NodeType = 1
	File      NodeType = 2
)

type Node struct {
	Type       NodeType
	Name       string
	Size       int
	ChildNodes []*Node
	Parent     *Node
}

func NewNode(nodeType NodeType, name string, size int) *Node {
	return &Node{
		Type:       nodeType,
		Name:       name,
		Size:       size,
		ChildNodes: make([]*Node, 0, 10),
		Parent:     nil,
	}
}

func (n *Node) AddNode(node *Node) {
	node.Parent = n
	n.ChildNodes = append(n.ChildNodes, node)
}
