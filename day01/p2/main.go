package main

import (
	"aoc2021/common/file"
	"fmt"
	"strconv"
)

func main() {
	depths := getDepths()
	last := depths[0] + depths[1] + depths[2]
	count := 0
	for i := range depths[1:len(depths)-2] {
		sum := depths[i+1] + depths[i+2] + depths[i+3]
		if sum > last {
			count++
		}
		last = sum
	}
	fmt.Println(count)
}

func getDepths() []int {
	lines, _  := file.GetLines("../p1/data.txt")
	depths := make([]int,0,len(lines))
	for _, line := range lines {
		depth, _ := strconv.Atoi(line)
		depths = append(depths,depth)
	}
	return depths
}