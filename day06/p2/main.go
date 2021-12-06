package main

import (
	"aoc2021/common/file"
	"fmt"
	"strconv"
	"strings"
)

func getSum(is []uint64) uint64 {
	count := uint64(0)
	for _, i := range is {
		count += i
	}
	return count
}

func main() {
	data := getData()
	fish := make([]uint64, 9, 9)
	for _, i := range data {
		fish[i]++
	}
	days := 0
	for days < 256 {
		spawning := fish[0]
		for i := 0; i < 8; i++ {
			fish[i] = fish[i+1]
		}
		fish[8] = spawning
		fish[6] += spawning
		days++
	}
	sum := getSum(fish)
	fmt.Println(days, sum)
}

func getData() []int {
	lines, _ := file.GetLines("../data.txt")
	values := strings.Split(lines[0], ",")
	data := make([]int, 0, 100)
	for _, val := range values {
		i, _ := strconv.Atoi(val)
		data = append(data, i)
	}
	return data
}
