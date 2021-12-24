package main

import (
	"aoc2021/common/file"
	"aoc2021/common/geom"
	"fmt"
	"regexp"
	"strconv"
)

const (
	modeON  = "on"
	modeOFF = "off"
)

var (
	reInstruction = regexp.MustCompile(`(on|off)\s+x=(-?\d+)\.\.(-?\d+),y=(-?\d+)\.\.(-?\d+),z=(-?\d+)\.\.(-?\d+)`)
)

type Instruction struct {
	istr string
	mode string
	min  geom.Pos
	max  geom.Pos
}

func (i Instruction) String() string {
	return i.istr
}

func newInstruction(s string) *Instruction {
	matches := reInstruction.FindAllStringSubmatch(s, -1)

	i := new(Instruction)
	i.istr = s
	i.mode = matches[0][1]

	x1, _ := strconv.ParseInt(matches[0][2], 10, 32)
	y1, _ := strconv.ParseInt(matches[0][4], 10, 32)
	z1, _ := strconv.ParseInt(matches[0][6], 10, 32)

	x2, _ := strconv.ParseInt(matches[0][3], 10, 32)
	y2, _ := strconv.ParseInt(matches[0][5], 10, 32)
	z2, _ := strconv.ParseInt(matches[0][7], 10, 32)

	i.min = geom.Pos{X: int(x1), Y: int(y1), Z: int(z1)}
	i.max = geom.Pos{X: int(x2), Y: int(y2), Z: int(z2)}

	return i
}

func main() {
	lines := getLines("../test3.txt")
	instructions := getInstructions(lines)

	for n, i := range instructions {
		fmt.Println(n, i)

	}

}

func getInstructions(strs []string) []*Instruction {
	instructions := make([]*Instruction, 0, len(strs))

	for _, str := range strs {
		instruction := newInstruction(str)
		if instruction != nil {
			instructions = append(instructions, instruction)
		}
	}

	return instructions
}

func getLines(filename string) []string {
	lines, _ := file.GetLines(filename)
	return lines
}
