package main

import (
	"aoc2021/common/datastructure"
	"aoc2021/common/file"
	"aoc2021/day24/alu"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	numInputSets          = 14
	numInstructionsPerSet = 18
)

type InstructionSets []*InstructionSet
type InstructionSet struct {
	lines []string
}

func (is *InstructionSet) getTokens(lineNums []int) [][]string {
	tokens := make([][]string, 0, len(lineNums))
	for _, lineNum := range lineNums {
		tokens = append(tokens, strings.Split(is.lines[lineNum], " "))
	}
	return tokens
}

func NewInstructionSet() *InstructionSet {
	is := new(InstructionSet)
	is.lines = make([]string, 0, numInstructionsPerSet)
	return is
}

func (is *InstructionSet) addLine(line string) {
	is.lines = append(is.lines, line)
}

func (iss InstructionSets) combine() []string {
	instructions := make([]string, 0, len(iss)*numInstructionsPerSet)
	for i := 0; i < len(iss); i++ {
		for j := 0; j < numInstructionsPerSet; j++ {
			instructions = append(instructions, iss[i].lines[j])
		}
	}
	return instructions
}

func getDigits(bytes string) ([]int64, error) {
	if strings.Contains(bytes, "0") {
		return nil, errors.New("has zeros")
	}
	vals := make([]int64, len(bytes), len(bytes))
	for i, b := range bytes {
		vals[i] = int64(byte(b) - '0')
	}
	return vals, nil
}

type stackVals struct {
	index int
	addY  int
}

func main() {

	iss := getInstructionSets("../data.txt")

	addX := make([]int, numInputSets, numInputSets)
	divVal := make([]int, numInputSets, numInputSets)
	addY := make([]int, numInputSets, numInputSets)

	for i := range iss {
		for j, tokens := range iss[i].getTokens([]int{4, 5, 15}) {
			val, _ := strconv.ParseInt(tokens[2], 10, 32)
			switch j {
			case 0:
				divVal[i] = int(val)
			case 1:
				addX[i] = int(val)
			case 2:
				addY[i] = int(val)
			}
		}
	}

	s := datastructure.NewStack(14)
	largestSolutionStr := make([]byte, numInputSets, numInputSets)
	largestSolutionInts := make([]int64, numInputSets, numInputSets)

	smallestSolutionStr := make([]byte, numInputSets, numInputSets)
	smallestSolutionInts := make([]int64, numInputSets, numInputSets)

	for i, addY := range addY {
		if divVal[i] == 1 {
			s.Push(stackVals{index: i, addY: addY})
		} else if divVal[i] == 26 {

			prevVals := s.Pop().(stackVals)

			var w0Largest, w1Largest, w0Smallest, w1Smallest int

			addY0 := prevVals.addY
			prevIndex := prevVals.index

			addX1 := addX[i]

		outer1:
			for j := 9; j > 0; j-- {
				for k := 9; k > 0; k-- {
					if j+addY0 == k-addX1 {
						w0Largest = j
						w1Largest = k
						break outer1
					}
				}
			}

			if w0Largest+w1Largest == 0 {
				panic("not supposed to happen, can't find a solution")
			}

			largestSolutionStr[prevIndex] = byte(w0Largest) + '0'
			largestSolutionStr[i] = byte(w1Largest) + '0'

			largestSolutionInts[prevIndex] = int64(w0Largest)
			largestSolutionInts[i] = int64(w1Largest)

		outer2:
			for j := 1; j <= 9; j++ {
				for k := 1; k <= 9; k++ {
					if j+addY0 == k-addX1 {
						w0Smallest = j
						w1Smallest = k
						break outer2
					}
				}
			}

			if w0Smallest+w1Smallest == 0 {
				panic("not supposed to happen, can't find a solution")
			}

			smallestSolutionStr[prevIndex] = byte(w0Smallest) + '0'
			smallestSolutionStr[i] = byte(w1Smallest) + '0'

			smallestSolutionInts[prevIndex] = int64(w0Smallest)
			smallestSolutionInts[i] = int64(w1Smallest)

		} else {
			panic("shouldn't happen")
		}
	}

	_, _, _, z, output, err := alu.RunALU(iss.combine(), largestSolutionInts)

	if z == 0 && err == nil {
		fmt.Println("found largest solution: ", string(largestSolutionStr))
	} else {
		if err != nil {
			fmt.Println("error checking largest solution: ", err)
		} else {
			fmt.Println("largest solution does not check out")
			fmt.Println("final ALU state: ", output)
		}
	}

	_, _, _, z, output, err = alu.RunALU(iss.combine(), smallestSolutionInts)

	if z == 0 && err == nil {
		fmt.Println("found smallest solution: ", string(smallestSolutionStr))
	} else {
		if err != nil {
			fmt.Println("error checking smallest solution: ", err)
		} else {
			fmt.Println("smallest solution does not check out")
			fmt.Println("final ALU state: ", output)
		}
	}

}

func getInstructionSets(filename string) InstructionSets {
	lines, _ := file.GetLines(filename)
	iss := make(InstructionSets, 0, numInputSets)

	l := 0
	for i := 0; i < numInputSets; i++ {
		is := NewInstructionSet()
		for j := 0; j < numInstructionsPerSet; j++ {
			is.addLine(lines[l])
			l++
		}
		iss = append(iss, is)
	}

	return iss
}
