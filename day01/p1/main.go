package main

import (
	"aoc2021/common/file"
	"fmt"
	"strconv"
)

func main() {
	depths := getDepths()
	last := depths[0]
	count := 0
	for _, depth := range depths[1:] {
		if depth > last {
			count++
		}
		last = depth
	}
	fmt.Println(count)
}

func getDepths() []int {
	lines, _ := file.GetLines("./data.txt")
	depths := make([]int, 0, len(lines))
	for _, line := range lines {
		depth, _ := strconv.Atoi(line)
		depths = append(depths, depth)
	}
	return depths
}
