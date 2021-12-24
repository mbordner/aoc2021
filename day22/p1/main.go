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
	capMax  = 50
	capMin  = -50
)

var (
	reInstruction = regexp.MustCompile(`(on|off)\s+x=(-?\d+)\.\.(-?\d+),y=(-?\d+)\.\.(-?\d+),z=(-?\d+)\.\.(-?\d+)`)
)

type ReactorCore map[geom.Pos]bool

func (rc *ReactorCore) on(p geom.Pos) {
	(*rc)[p] = true
}

func (rc *ReactorCore) off(p geom.Pos) {
	if _, exists := (*rc)[p]; exists {
		delete(*rc, p)
	}
}

type Instruction struct {
	mode string
	min  geom.Pos
	max  geom.Pos
}

func cap(val int) int {
	if val < capMin {
		val = capMin
	}
	if val > capMax {
		val = capMax
	}
	return val
}

func newInstruction(s string) *Instruction {
	matches := reInstruction.FindAllStringSubmatch(s, -1)

	i := new(Instruction)
	i.mode = matches[0][1]

	x1, _ := strconv.ParseInt(matches[0][2], 10, 32)
	y1, _ := strconv.ParseInt(matches[0][4], 10, 32)
	z1, _ := strconv.ParseInt(matches[0][6], 10, 32)

	x2, _ := strconv.ParseInt(matches[0][3], 10, 32)
	y2, _ := strconv.ParseInt(matches[0][5], 10, 32)
	z2, _ := strconv.ParseInt(matches[0][7], 10, 32)

	if int(x2) < capMin || int(y2) < capMin || int(z2) < capMin {
		return nil
	}

	if int(x1) > capMax || int(y2) > capMax || int(z2) > capMax {
		return nil
	}

	i.min = geom.Pos{X: cap(int(x1)), Y: cap(int(y1)), Z: cap(int(z1))}
	i.max = geom.Pos{X: cap(int(x2)), Y: cap(int(y2)), Z: cap(int(z2))}

	return i
}

func main() {
	lines := getLines("../data.txt")
	instructions := getInstructions(lines)

	rc := make(ReactorCore)

	for _, i := range instructions {
		bb := geom.BoundingBox{}
		bb.SetExtents(i.min.X, i.min.Y, i.min.Z, i.max.X, i.max.Y, i.max.Z)
		for _, p := range bb.GetPositions() {
			if i.mode == modeON {
				rc.on(p)
			} else {
				rc.off(p)
			}
		}
	}

	fmt.Println(len(rc))
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
