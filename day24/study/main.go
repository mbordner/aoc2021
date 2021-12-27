package main

import (
	"aoc2021/common/file"
	"aoc2021/day24/alu"
	"fmt"
	"regexp"
	"strings"
)

const (
	numInputSets          = 14
	numInstructionsPerSet = 18
)

var (
	reDigit    = regexp.MustCompile(`-?\d+`)
	reVariable = regexp.MustCompile(`[w-z]`)
)

type inputSets []*inputSet

type inputSet struct {
	lines  []string
	tokens [][]string
}

func newInputSet() *inputSet {
	is := new(inputSet)
	is.lines = make([]string, 0, numInstructionsPerSet)
	is.tokens = make([][]string, 0, numInstructionsPerSet)
	return is
}

func (is *inputSet) addLine(line string) {
	is.lines = append(is.lines, line)
	tokens := strings.Split(line, " ")
	is.tokens = append(is.tokens, tokens)
}

func (iss *inputSets) compare() {
	for i := 0; i < len(*iss)-1; i++ {
		isI := (*iss)[i]
		isJ := (*iss)[i+1]
		fmt.Println(fmt.Sprintf("------- comparing is %d to %d", i+1, i+2))
		for n := 0; n < numInstructionsPerSet; n++ {
			if isI.lines[n] != isJ.lines[n] {
				fmt.Println(i+1, " inst ", n+1, ": ", isI.lines[n])
				fmt.Println(i+2, " inst ", n+1, ": ", isJ.lines[n])
			}
		}
	}
}

func main() {
	lines, _ := file.GetLines("../data.txt")

	iss := make(inputSets, 0, 14)

	if numInputSets*numInstructionsPerSet == len(lines) {
		l := 0
		for i := 0; i < numInputSets; i++ {
			is := newInputSet()
			for j := 0; j < numInstructionsPerSet; j++ {
				is.addLine(lines[l])
				l++
			}
			iss = append(iss, is)
		}
	}

	iss.compare()

	fmt.Println("==============")

	for a := int64(1); a <= 9; a++ {
		inputs := []int64{a}
		_, _, _, _, output, err := alu.RunALU(iss[0].lines, inputs)
		if err != nil {
			break
		}
		fmt.Println(inputs, output)
	}

}
