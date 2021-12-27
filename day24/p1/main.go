package main

import (
	"aoc2021/common/file"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	numInputSets          = 14
	numInstructionsPerSet = 18
)

var (
	reALUOutput = regexp.MustCompile(`{"w":(-?\d+),"x":(-?\d+),"y":(-?\d+),"z":(-?\d+)}`)
)

type InstructionSets []*InstructionSet
type InstructionSet struct {
	lines []string
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

func getDigits(num int64) ([]int64, error) {
	bytes := fmt.Sprintf("%d", num)
	if strings.Contains(bytes, "0") {
		return nil, errors.New("has zeros")
	}
	vals := make([]int64, len(bytes), len(bytes))
	for i, b := range bytes {
		vals[i] = int64(byte(b) - '0')
	}
	return vals, nil
}

func main() {
	/*
		iss := getInstructionSets("../data.txt")

		num := int64(79996462532205)

		for num > 0 {
			if digits, err := getDigits(num); err == nil {
				output, err := alu.RunALU(iss.combine(), digits)
				if err != nil {
					panic(err)
				}
				matches := reALUOutput.FindAllStringSubmatch(output,-1)
				if matches[0][4] == "0" {
					fmt.Println(digits)
					break
				}
			}
			num--
		}
	*/
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
