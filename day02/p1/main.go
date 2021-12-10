package main

import (
	"aoc2021/common/file"
	"fmt"
	"regexp"
	"strconv"
)

var (
	reInstruction = regexp.MustCompile(`^(\w+)\s+(\d+)$`)
)

type instruction struct {
	dir string
	amt int
}

type pos struct {
	x int
	y int
	z int
}

func main() {
	sub := pos{}
	instructions := getInstructions("../data.txt")
	for _, instruction := range instructions {
		switch instruction.dir {
		case "forward":
			sub.x += instruction.amt
		case "down":
			sub.z += instruction.amt
		case "up":
			sub.z -= instruction.amt
		}
	}

	fmt.Println(sub.z * sub.x)
}

func getInstructions(filename string) []instruction {
	instructions := make([]instruction, 0, 100)
	lines, _ := file.GetLines(filename)
	for _, line := range lines {
		matches := reInstruction.FindStringSubmatch(line)
		dir := matches[1]
		amt, _ := strconv.Atoi(matches[2])
		instructions = append(instructions, instruction{
			dir: dir,
			amt: amt,
		})
	}
	return instructions
}
