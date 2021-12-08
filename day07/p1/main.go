package main

import (
	"aoc2021/common/file"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func getStats(nums []int) (min, max, avg int) {
	min = nums[0]
	max = 0
	sum := 0
	for _, n := range nums {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
		sum += n
	}
	avg = int(float32(sum) / float32(len(nums)))
	return
}

func getCost(crabs []int, pos int) (cost int) {
	for _, n := range crabs {
		cost +=  int(math.Abs( float64(n) - float64(pos)))
	}
	return
}


func main() {
	crabs := getData()
	min, max, avg := getStats(crabs)
	fmt.Println(min,max,avg)

	cost := getCost(crabs,min)

	for i := min+1; i <= max; i++ {
		c := getCost(crabs,i)
		if c < cost {
			cost = c
		}
	}

	fmt.Println(cost)
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
