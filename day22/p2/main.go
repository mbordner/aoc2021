package main

import (
	"aoc2021/common/file"
	"aoc2021/common/geom"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const (
	modeON         = "on"
	modeOFF        = "off"
	initRegionOnly = false
)

var (
	start         = time.Now()
	reInstruction = regexp.MustCompile(`^(on|off)\s+x=(-?\d+)\.\.(-?\d+),y=(-?\d+)\.\.(-?\d+),z=(-?\d+)\.\.(-?\d+)$`)
)

type Instruction struct {
	istr string
	mode string
	min  geom.Point
	max  geom.Point
}

func (i Instruction) String() string {
	return i.istr
}

func (i *Instruction) cap(val int64) *Instruction {
	if i.max.X < -val || i.max.Y < -val || i.max.Z < -val {
		return nil
	}

	if i.min.X > val || i.min.Y > val || i.min.Z > val {
		return nil
	}

	pts := []*geom.Point{&i.min, &i.max}
	for _, pt := range pts {
		if pt.X < -val {
			pt.X = -val
		}
		if pt.X > val {
			pt.X = val
		}
		if pt.Y < -val {
			pt.Y = -val
		}
		if pt.Y > val {
			pt.Y = val
		}
		if pt.Z < -val {
			pt.Z = -val
		}
		if pt.Z > val {
			pt.Z = val
		}
	}
	return i
}

func newInstruction(s string) *Instruction {
	matches := reInstruction.FindAllStringSubmatch(s, -1)

	if matches == nil {
		return nil
	}

	i := new(Instruction)
	i.istr = s
	i.mode = matches[0][1]

	x1, _ := strconv.ParseInt(matches[0][2], 10, 64)
	y1, _ := strconv.ParseInt(matches[0][4], 10, 64)
	z1, _ := strconv.ParseInt(matches[0][6], 10, 64)

	x2, _ := strconv.ParseInt(matches[0][3], 10, 64)
	y2, _ := strconv.ParseInt(matches[0][5], 10, 64)
	z2, _ := strconv.ParseInt(matches[0][7], 10, 64)

	i.min = geom.Point{X: x1, Y: y1, Z: z1}
	i.max = geom.Point{X: x2, Y: y2, Z: z2}

	if initRegionOnly {
		return i.cap(50)
	}
	return i
}

func main() {

	lines := getLines("../data.txt")
	instructions := getInstructions(lines)

	cuboidsON := make(geom.Cuboids, 0, 100)

	for n, i := range instructions {
		fmt.Println(n+1, fmt.Sprintf(" [length of cuboidsON: %d] ", len(cuboidsON)), i)

		o := geom.Cuboid{Min: i.min, Max: i.max}

		if i.mode == modeON {

			cuboidsON = cuboidsON.Merge(o)

		} else if i.mode == modeOFF {

			cuboidsON = cuboidsON.Remove(o)

		}

	}

	fmt.Println("cubes on: ", cuboidsON.PointsCount())
	duration := time.Since(start)
	fmt.Println(duration)

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
