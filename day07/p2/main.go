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

var memo = make(map[int]int)

func expand(cost int) int {
	if c, ok := memo[cost]; ok {
		return c
	}
	c := 0
	for i := 1; i <= cost; i++ {
		c += i
	}
	memo[cost] = c
	return c
}

func getCost(crabs []int, pos int) (cost int) {
	for _, n := range crabs {
		cost +=  expand(int(math.Abs( float64(n) - float64(pos))))
	}
	return
}


func main() {
	crabs := getData()
	min, max, avg := getStats(crabs)
	fmt.Println(min,max,avg)

	cost := getCost(crabs,min)
	pos := min

	for i := min+1; i <= max; i++ {
		c := getCost(crabs,i)
		if c < cost {
			cost = c
			pos = i
		}
		//fmt.Println(">>",i,c)
	}

	fmt.Println(pos,cost)
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
